// Package basic gives a default implementation of the Engine interface.
// It only connects to local forwarding node via Unix socket.
package basic

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	enc "github.com/named-data/ndnd/std/encoding"
	"github.com/named-data/ndnd/std/engine/face"
	"github.com/named-data/ndnd/std/log"
	"github.com/named-data/ndnd/std/ndn"
	mgmt "github.com/named-data/ndnd/std/ndn/mgmt_2022"
	spec "github.com/named-data/ndnd/std/ndn/spec_2022"
	"github.com/named-data/ndnd/std/utils"
)

const DefaultInterestLife = 4 * time.Second
const TimeoutMargin = 10 * time.Millisecond

type fibEntry = ndn.InterestHandler

type pendInt struct {
	callback    ndn.ExpressCallbackFunc
	deadline    time.Time
	canBePrefix bool
	// mustBeFresh is actually not useful, since Freshness is decided by the cache, not us.
	mustBeFresh   bool
	impSha256     []byte
	timeoutCancel func() error
}

type pitEntry = []*pendInt

type Engine struct {
	face  face.Face
	timer ndn.Timer

	// fib contains the registered Interest handlers.
	fib *NameTrie[fibEntry]

	// pit contains pending outgoing Interests.
	pit *NameTrie[pitEntry]

	// Since there is only one main coroutine, no need for RW locks.
	fibLock sync.Mutex
	pitLock sync.Mutex

	// mgmtConf is the configuration for the management protocol.
	mgmtConf *mgmt.MgmtConfig

	// cmdChecker is used to validate NFD management packets.
	cmdChecker ndn.SigChecker
}

func (e *Engine) String() string {
	return "basic-engine"
}

func (e *Engine) EngineTrait() ndn.Engine {
	return e
}

func (*Engine) Spec() ndn.Spec {
	return spec.Spec{}
}

func (e *Engine) Timer() ndn.Timer {
	return e.timer
}

func (e *Engine) AttachHandler(prefix enc.Name, handler ndn.InterestHandler) error {
	e.fibLock.Lock()
	defer e.fibLock.Unlock()
	n := e.fib.MatchAlways(prefix)
	if n.Value() != nil {
		return ndn.ErrMultipleHandlers
	}
	n.SetValue(handler)
	return nil
}

func (e *Engine) DetachHandler(prefix enc.Name) error {
	e.fibLock.Lock()
	defer e.fibLock.Unlock()

	n := e.fib.ExactMatch(prefix)
	if n == nil {
		return ndn.ErrInvalidValue{Item: "prefix", Value: prefix}
	}
	n.SetValue(nil)
	n.Prune()
	return nil
}

func (e *Engine) onPacket(frame []byte) error {
	// Copy received buffer from face so face can reuse it
	frameCopy := make([]byte, len(frame))
	copy(frameCopy, frame)
	reader := enc.NewFastBufReader(frameCopy)

	var nackReason uint64 = spec.NackReasonNone
	var pitToken []byte = nil
	var incomingFaceId enc.Optional[uint64]
	var raw enc.Wire = nil

	if hasLogTrace() {
		wire := reader.Range(0, reader.Length())
		log.Trace(e, "Received packet bytes", "wire", hex.EncodeToString(wire.Join()))
	}

	// Parse the outer packet - could be either L2 or L3
	pkt, ctx, err := spec.ReadPacket(reader)
	if err != nil {
		// Recoverable error. Should continue.
		log.Error(e, "Failed to parse packet", "err", err)
		return nil
	}

	// Now, exactly one of Interest, Data, LpPacket is not nil
	// First check LpPacket, and do further parse.
	if pkt.LpPacket != nil {
		lpPkt := pkt.LpPacket
		if lpPkt.FragIndex.IsSet() || lpPkt.FragCount.IsSet() {
			log.Warn(e, "Fragmented LpPackets are not supported - DROP")
			return nil
		}

		// Parse the inner packet.
		raw = pkt.LpPacket.Fragment
		if len(raw) == 1 {
			pkt, ctx, err = spec.ReadPacket(enc.NewFastBufReader(raw[0]))
		} else {
			pkt, ctx, err = spec.ReadPacket(enc.NewFastReader(raw))
		}

		// Make sure there is an inner packet.
		if err != nil || (pkt.Data == nil) == (pkt.Interest == nil) {
			if hasLogTrace() {
				wire := reader.Range(0, reader.Length())
				log.Trace(e, "Failed to parse packet bytes", "wire", hex.EncodeToString(wire.Join()))
			}

			// Recoverable error. Should continue.
			log.Error(e, "Failed to parse packet in LpPacket", "err", err)
			return nil
		}

		// Set parameters
		if lpPkt.Nack != nil {
			nackReason = lpPkt.Nack.Reason
		}
		pitToken = lpPkt.PitToken
		incomingFaceId = lpPkt.IncomingFaceId
	} else {
		raw = reader.Range(0, reader.Length())
	}

	// Now pkt is either Data or Interest (including Nack).
	if nackReason != spec.NackReasonNone {
		if pkt.Interest == nil {
			log.Error(e, "Nack received for non-Interest", "reason", nackReason)
			return nil
		}
		log.Trace(e, "Nack received", "reason", nackReason, "name", pkt.Interest.Name())
		e.onNack(pkt.Interest.NameV, nackReason)
	} else if pkt.Interest != nil {
		log.Trace(e, "Interest received", "name", pkt.Interest.Name())
		e.onInterest(ndn.InterestHandlerArgs{
			Interest:       pkt.Interest,
			RawInterest:    raw,
			SigCovered:     ctx.Interest_context.SigCovered(),
			PitToken:       pitToken,
			IncomingFaceId: incomingFaceId.Ptr(),
		})
	} else if pkt.Data != nil {
		log.Trace(e, "Data received", "name", pkt.Data.Name())
		// PitToken is not used for now
		e.onData(pkt.Data, ctx.Data_context.SigCovered(), raw, pitToken)
	} else {
		panic("[BUG] unexpected packet type") // checked above
	}

	return nil
}

func (e *Engine) onInterest(args ndn.InterestHandlerArgs) {
	name := args.Interest.Name()

	// Compute deadline
	args.Deadline = e.timer.Now()
	if args.Interest.Lifetime() != nil {
		args.Deadline = args.Deadline.Add(*args.Interest.Lifetime())
	} else {
		args.Deadline = args.Deadline.Add(DefaultInterestLife)
	}

	// Match node
	handler := func() ndn.InterestHandler {
		e.fibLock.Lock()
		defer e.fibLock.Unlock()
		n := e.fib.PrefixMatch(name)

		// If we have the prefix-free condition, we can return the value here
		// directly. But we need longest prefix match now.
		// return n.Value()

		for n != nil && n.Value() == nil {
			n = n.Parent()
		}
		if n != nil {
			return n.Value()
		}
		return nil
	}()
	if handler == nil {
		log.Warn(e, "No handler for interest", "name", name)
		return
	}

	// The reply callback function
	args.Reply = func(encodedData enc.Wire) error {
		now := e.timer.Now()
		if args.Deadline.Before(now) {
			log.Warn(e, "Deadline exceeded - DROP", "name", name)
			return ndn.ErrDeadlineExceed
		}
		if !e.face.IsRunning() {
			log.Error(e, "Cannot send through a closed face - DROP", "name", name)
			return ndn.ErrFaceDown
		}
		if args.PitToken != nil {
			lpPkt := &spec.Packet{
				LpPacket: &spec.LpPacket{
					PitToken: args.PitToken,
					Fragment: encodedData,
				},
			}
			encoder := spec.PacketEncoder{}
			encoder.Init(lpPkt)
			wire := encoder.Encode(lpPkt)
			if wire == nil {
				return ndn.ErrFailedToEncode
			}
			return e.face.Send(wire)
		} else {
			return e.face.Send(encodedData)
		}
	}

	// Call the handler. The handler should create goroutine to avoid blocking.
	// Do not `go` here because if Data is ready at hand, creating a go routine may be slower. Not tested though.
	handler(args)
}

func (e *Engine) onDataMatch(pkt *spec.Data, raw enc.Wire) pitEntry {
	e.pitLock.Lock()
	defer e.pitLock.Unlock()

	n := e.pit.PrefixMatch(pkt.NameV)
	if n == nil {
		log.Warn(e, "Received data for an unknown interest - DROP", "name", pkt.Name())
		return nil
	}

	ret := make(pitEntry, 0, 4)
	for cur := n; cur != nil; cur = cur.Parent() {
		entries := cur.Value()
		for i := 0; i < len(entries); i++ {
			entry := entries[i]

			// we don't check MustBeFresh, as it is the job of the cache/forwarder.
			// check CanBePrefix
			if cur.Depth() < len(pkt.NameV) && !entry.canBePrefix {
				continue
			}

			// check ImplicitDigest256
			if entry.impSha256 != nil {
				h := sha256.New()
				for _, buf := range raw {
					h.Write(buf)
				}
				digest := h.Sum(nil)
				if !bytes.Equal(entry.impSha256, digest) {
					continue
				}
			}

			// pop entry
			entries[i] = entries[len(entries)-1]
			entries = entries[:len(entries)-1]
			i-- // recheck the current index
			ret = append(ret, entry)
		}
		cur.SetValue(entries)
	}

	n.PruneIf(func(lst []*pendInt) bool { return len(lst) == 0 })

	return ret
}

func (e *Engine) onData(pkt *spec.Data, sigCovered enc.Wire, raw enc.Wire, pitToken []byte) {
	for _, entry := range e.onDataMatch(pkt, raw) {
		entry.timeoutCancel()
		if entry.callback == nil {
			panic("[BUG] PIT has empty entry")
		}

		entry.callback(ndn.ExpressCallbackArgs{
			Result:     ndn.InterestResultData,
			Data:       pkt,
			RawData:    raw,
			SigCovered: sigCovered,
			NackReason: spec.NackReasonNone,
		})
	}
}

func (e *Engine) onNack(name enc.Name, reason uint64) {
	entries := func() []*pendInt {
		e.pitLock.Lock()
		defer e.pitLock.Unlock()

		n := e.pit.ExactMatch(name)
		if n == nil {
			log.Warn(e, "Received Nack for an unknown interest - DROP", "name", name)
			return nil
		}

		ret := n.Value()
		n.SetValue(nil)
		n.Prune()
		return ret
	}()

	for _, entry := range entries {
		entry.timeoutCancel()

		if entry.callback == nil {
			panic("[BUG] PIT has empty entry")
		}

		entry.callback(ndn.ExpressCallbackArgs{
			Result:     ndn.InterestResultNack,
			NackReason: reason,
		})
	}
}

func (e *Engine) onError(err error) error {
	log.Error(e, "Error on face", "err", err, "face", e.face)
	// TODO: Handle Interest cancellation
	return err
}

func (e *Engine) Start() error {
	if e.face.IsRunning() {
		return errors.New("face is already running")
	}

	e.face.OnPacket(e.onPacket)
	e.face.OnError(e.onError)

	err := e.face.Open()
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) Stop() error {
	if !e.face.IsRunning() {
		return errors.New("face is not running")
	}
	return e.face.Close()
}

func (e *Engine) IsRunning() bool {
	return e.face.IsRunning()
}

func (e *Engine) onExpressTimeout(n *NameTrie[pitEntry]) {
	now := e.timer.Now()

	expired := func() []*pendInt {
		e.pitLock.Lock()
		defer e.pitLock.Unlock()

		ret := make([]*pendInt, 0, 4)
		entries := n.Value()
		for i := 0; i < len(entries); i++ {
			entry := entries[i]
			if entry.deadline.After(now) {
				continue
			}

			// pop entry
			entries[i] = entries[len(entries)-1]
			entries = entries[:len(entries)-1]
			i-- // recheck the current index
			ret = append(ret, entry)
		}

		n.SetValue(entries)
		n.PruneIf(func(lst []*pendInt) bool { return len(lst) == 0 })

		return ret
	}()

	for _, entry := range expired {
		if entry.callback == nil {
			panic("[BUG] PIT has empty entry")
		}

		entry.callback(ndn.ExpressCallbackArgs{
			Result:     ndn.InterestResultTimeout,
			NackReason: spec.NackReasonNone,
		})
	}
}

func (e *Engine) Express(interest *ndn.EncodedInterest, callback ndn.ExpressCallbackFunc) error {
	var impSha256 []byte = nil

	finalName := interest.FinalName
	nodeName := interest.FinalName

	if callback == nil {
		callback = func(ndn.ExpressCallbackArgs) {}
	}

	// Handle implicit digest
	if len(finalName) <= 0 {
		return ndn.ErrInvalidValue{Item: "finalName", Value: finalName}
	}
	lastComp := finalName[len(finalName)-1]
	if lastComp.Typ == enc.TypeImplicitSha256DigestComponent {
		impSha256 = lastComp.Val
		nodeName = finalName[:len(finalName)-1]
	}

	// Handle deadline
	lifetime := DefaultInterestLife
	if interest.Config.Lifetime != nil {
		lifetime = *interest.Config.Lifetime
	}
	deadline := e.timer.Now().Add(lifetime)

	// Inject interest into PIT
	func() {
		e.pitLock.Lock()
		defer e.pitLock.Unlock()

		n := e.pit.MatchAlways(nodeName)
		entry := &pendInt{
			callback:    callback,
			deadline:    deadline,
			canBePrefix: interest.Config.CanBePrefix,
			mustBeFresh: interest.Config.MustBeFresh,
			impSha256:   impSha256,
			timeoutCancel: e.timer.Schedule(lifetime+TimeoutMargin, func() {
				e.onExpressTimeout(n)
			}),
		}
		n.SetValue(append(n.Value(), entry))
	}()

	// Wrap the interest in link packet if needed
	wire := interest.Wire
	if interest.Config.NextHopId != nil {
		lpPkt := &spec.Packet{
			LpPacket: &spec.LpPacket{
				Fragment:      wire,
				NextHopFaceId: enc.OptionPtr(interest.Config.NextHopId),
			},
		}
		encoder := spec.PacketEncoder{}
		encoder.Init(lpPkt)
		wire = encoder.Encode(lpPkt)
	}

	// Send interest to face
	err := e.face.Send(wire)
	if err != nil {
		log.Error(e, "Failed to send interest", "err", err)
	}

	log.Trace(e, "Interest sent", "name", finalName)
	return err
}

func (e *Engine) ExecMgmtCmd(module string, cmd string, args any) (any, error) {
	cmdArgs, ok := args.(*mgmt.ControlArgs)
	if !ok {
		return nil, ndn.ErrInvalidValue{Item: "args", Value: args}
	}

	intCfg := &ndn.InterestConfig{
		Lifetime:    utils.IdPtr(1 * time.Second),
		Nonce:       utils.ConvertNonce(e.timer.Nonce()),
		MustBeFresh: true,

		// Signed interest shenanigans (NFD wants this)
		SigNonce: e.timer.Nonce(),
		SigTime:  utils.IdPtr(time.Duration(e.timer.Now().UnixMilli()) * time.Millisecond),
	}
	interest, err := e.mgmtConf.MakeCmd(module, cmd, cmdArgs, intCfg)
	if err != nil {
		return nil, err
	}

	type mgmtResp struct {
		err error
		val *mgmt.ControlResponse
	}
	respCh := make(chan *mgmtResp)

	err = e.Express(interest, func(args ndn.ExpressCallbackArgs) {
		resp := &mgmtResp{}
		defer func() {
			respCh <- resp
			close(respCh)
		}()

		if args.Result == ndn.InterestResultNack {
			resp.err = fmt.Errorf("nack received: %v", args.NackReason)
		} else if args.Result == ndn.InterestResultTimeout {
			resp.err = ndn.ErrDeadlineExceed
		} else if args.Result == ndn.InterestResultData {
			data := args.Data
			valid := e.cmdChecker(data.Name(), args.SigCovered, data.Signature())
			if !valid {
				resp.err = fmt.Errorf("command signature is not valid")
			} else {
				ret, err := mgmt.ParseControlResponse(enc.NewFastReader(data.Content()), true)
				if err != nil {
					resp.err = err
				} else {
					resp.val = ret
					if ret.Val != nil {
						if ret.Val.StatusCode == 200 {
							return
						} else {
							resp.err = fmt.Errorf("command failed due to error %d: %s",
								ret.Val.StatusCode, ret.Val.StatusText)
						}
					} else {
						resp.err = fmt.Errorf("improper response")
					}
				}
			}
		} else {
			resp.err = fmt.Errorf("unknown result: %v", args.Result)
		}
	})
	if err != nil {
		return nil, err
	}

	resp := <-respCh
	return resp.val, resp.err
}

func (e *Engine) RegisterRoute(prefix enc.Name) error {
	_, err := e.ExecMgmtCmd("rib", "register", &mgmt.ControlArgs{Name: prefix})
	if err != nil {
		log.Error(e, "Failed to register prefix", "err", err, "name", prefix)
		return err
	} else {
		log.Info(e, "Prefix registered", "name", prefix)
	}
	return nil
}

func (e *Engine) UnregisterRoute(prefix enc.Name) error {
	_, err := e.ExecMgmtCmd("rib", "unregister", &mgmt.ControlArgs{Name: prefix})
	if err != nil {
		log.Error(e, "Failed to unregister prefix", "err", err, "name", prefix)
		return err
	} else {
		log.Info(e, "Prefix unregistered", "name", prefix)
	}
	return nil
}

func NewEngine(face face.Face, timer ndn.Timer, cmdSigner ndn.Signer, cmdChecker ndn.SigChecker) *Engine {
	if face == nil || timer == nil || cmdSigner == nil || cmdChecker == nil {
		return nil
	}
	mgmtCfg := mgmt.NewConfig(face.IsLocal(), cmdSigner, spec.Spec{})
	return &Engine{
		face:       face,
		timer:      timer,
		mgmtConf:   mgmtCfg,
		cmdChecker: cmdChecker,
		fib:        NewNameTrie[fibEntry](),
		pit:        NewNameTrie[pitEntry](),
		fibLock:    sync.Mutex{},
		pitLock:    sync.Mutex{},
	}
}

func hasLogTrace() bool {
	return log.Default().Level() <= log.LevelTrace
}
