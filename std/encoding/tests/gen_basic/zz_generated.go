// Code generated by ndn tlv codegen DO NOT EDIT.
package gen_basic

import (
	"encoding/binary"
	"io"
	"strings"
	"time"

	enc "github.com/named-data/ndnd/std/encoding"
)

type FakeMetaInfoEncoder struct {
	length uint
}

type FakeMetaInfoParsingContext struct {
}

func (encoder *FakeMetaInfoEncoder) Init(value *FakeMetaInfo) {

	l := uint(0)
	l += 1
	l += uint(1 + enc.Nat(value.Number).EncodingLength())
	l += 1
	l += uint(1 + enc.Nat(uint64(value.Time/time.Millisecond)).EncodingLength())
	if value.Binary != nil {
		l += 1
		l += uint(enc.TLNum(len(value.Binary)).EncodingLength())
		l += uint(len(value.Binary))
	}
	encoder.length = l

}

func (context *FakeMetaInfoParsingContext) Init() {

}

func (encoder *FakeMetaInfoEncoder) EncodeInto(value *FakeMetaInfo, buf []byte) {

	pos := uint(0)

	buf[pos] = byte(24)
	pos += 1

	buf[pos] = byte(enc.Nat(value.Number).EncodeInto(buf[pos+1:]))
	pos += uint(1 + buf[pos])
	buf[pos] = byte(25)
	pos += 1

	buf[pos] = byte(enc.Nat(uint64(value.Time / time.Millisecond)).EncodeInto(buf[pos+1:]))
	pos += uint(1 + buf[pos])
	if value.Binary != nil {
		buf[pos] = byte(26)
		pos += 1
		pos += uint(enc.TLNum(len(value.Binary)).EncodeInto(buf[pos:]))
		copy(buf[pos:], value.Binary)
		pos += uint(len(value.Binary))
	}
}

func (encoder *FakeMetaInfoEncoder) Encode(value *FakeMetaInfo) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *FakeMetaInfoParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*FakeMetaInfo, error) {

	var handled_Number bool = false
	var handled_Time bool = false
	var handled_Binary bool = false

	progress := -1
	_ = progress

	value := &FakeMetaInfo{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 24:
				if true {
					handled = true
					handled_Number = true
					value.Number = uint64(0)
					{
						for i := 0; i < int(l); i++ {
							x := byte(0)
							x, err = reader.ReadByte()
							if err != nil {
								if err == io.EOF {
									err = io.ErrUnexpectedEOF
								}
								break
							}
							value.Number = uint64(value.Number<<8) | uint64(x)
						}
					}
				}
			case 25:
				if true {
					handled = true
					handled_Time = true
					{
						timeInt := uint64(0)
						timeInt = uint64(0)
						{
							for i := 0; i < int(l); i++ {
								x := byte(0)
								x, err = reader.ReadByte()
								if err != nil {
									if err == io.EOF {
										err = io.ErrUnexpectedEOF
									}
									break
								}
								timeInt = uint64(timeInt<<8) | uint64(x)
							}
						}
						value.Time = time.Duration(timeInt) * time.Millisecond
					}
				}
			case 26:
				if true {
					handled = true
					handled_Binary = true
					value.Binary = make([]byte, l)
					_, err = reader.ReadFull(value.Binary)
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Number && err == nil {
		err = enc.ErrSkipRequired{Name: "Number", TypeNum: 24}
	}
	if !handled_Time && err == nil {
		err = enc.ErrSkipRequired{Name: "Time", TypeNum: 25}
	}
	if !handled_Binary && err == nil {
		value.Binary = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *FakeMetaInfo) Encode() enc.Wire {
	encoder := FakeMetaInfoEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *FakeMetaInfo) Bytes() []byte {
	return value.Encode().Join()
}

func ParseFakeMetaInfo(reader enc.WireView, ignoreCritical bool) (*FakeMetaInfo, error) {
	context := FakeMetaInfoParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type OptFieldEncoder struct {
	length uint
}

type OptFieldParsingContext struct {
}

func (encoder *OptFieldEncoder) Init(value *OptField) {

	l := uint(0)
	if optval, ok := value.Number.Get(); ok {
		l += 1
		l += uint(1 + enc.Nat(optval).EncodingLength())
	}
	if optval, ok := value.Time.Get(); ok {
		l += 1
		l += uint(1 + enc.Nat(uint64(optval/time.Millisecond)).EncodingLength())
	}
	if value.Binary != nil {
		l += 1
		l += uint(enc.TLNum(len(value.Binary)).EncodingLength())
		l += uint(len(value.Binary))
	}
	if value.Bool {
		l += 1
		l += 1
	}
	encoder.length = l

}

func (context *OptFieldParsingContext) Init() {

}

func (encoder *OptFieldEncoder) EncodeInto(value *OptField, buf []byte) {

	pos := uint(0)

	if optval, ok := value.Number.Get(); ok {
		buf[pos] = byte(24)
		pos += 1

		buf[pos] = byte(enc.Nat(optval).EncodeInto(buf[pos+1:]))
		pos += uint(1 + buf[pos])

	}
	if optval, ok := value.Time.Get(); ok {
		buf[pos] = byte(25)
		pos += 1

		buf[pos] = byte(enc.Nat(uint64(optval / time.Millisecond)).EncodeInto(buf[pos+1:]))
		pos += uint(1 + buf[pos])

	}
	if value.Binary != nil {
		buf[pos] = byte(26)
		pos += 1
		pos += uint(enc.TLNum(len(value.Binary)).EncodeInto(buf[pos:]))
		copy(buf[pos:], value.Binary)
		pos += uint(len(value.Binary))
	}
	if value.Bool {
		buf[pos] = byte(48)
		pos += 1
		buf[pos] = byte(0)
		pos += 1
	}
}

func (encoder *OptFieldEncoder) Encode(value *OptField) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *OptFieldParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*OptField, error) {

	var handled_Number bool = false
	var handled_Time bool = false
	var handled_Binary bool = false
	var handled_Bool bool = false

	progress := -1
	_ = progress

	value := &OptField{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 24:
				if true {
					handled = true
					handled_Number = true
					{
						optval := uint64(0)
						optval = uint64(0)
						{
							for i := 0; i < int(l); i++ {
								x := byte(0)
								x, err = reader.ReadByte()
								if err != nil {
									if err == io.EOF {
										err = io.ErrUnexpectedEOF
									}
									break
								}
								optval = uint64(optval<<8) | uint64(x)
							}
						}
						value.Number.Set(optval)
					}
				}
			case 25:
				if true {
					handled = true
					handled_Time = true
					{
						timeInt := uint64(0)
						timeInt = uint64(0)
						{
							for i := 0; i < int(l); i++ {
								x := byte(0)
								x, err = reader.ReadByte()
								if err != nil {
									if err == io.EOF {
										err = io.ErrUnexpectedEOF
									}
									break
								}
								timeInt = uint64(timeInt<<8) | uint64(x)
							}
						}
						optval := time.Duration(timeInt) * time.Millisecond
						value.Time.Set(optval)
					}
				}
			case 26:
				if true {
					handled = true
					handled_Binary = true
					value.Binary = make([]byte, l)
					_, err = reader.ReadFull(value.Binary)
				}
			case 48:
				if true {
					handled = true
					handled_Bool = true
					value.Bool = true
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Number && err == nil {
		value.Number.Unset()
	}
	if !handled_Time && err == nil {
		value.Time.Unset()
	}
	if !handled_Binary && err == nil {
		value.Binary = nil
	}
	if !handled_Bool && err == nil {
		value.Bool = false
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *OptField) Encode() enc.Wire {
	encoder := OptFieldEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *OptField) Bytes() []byte {
	return value.Encode().Join()
}

func ParseOptField(reader enc.WireView, ignoreCritical bool) (*OptField, error) {
	context := OptFieldParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type WireNameFieldEncoder struct {
	length uint

	Wire_length uint
	Name_length uint
}

type WireNameFieldParsingContext struct {
}

func (encoder *WireNameFieldEncoder) Init(value *WireNameField) {
	if value.Wire != nil {
		encoder.Wire_length = 0
		for _, c := range value.Wire {
			encoder.Wire_length += uint(len(c))
		}
	}
	if value.Name != nil {
		encoder.Name_length = 0
		for _, c := range value.Name {
			encoder.Name_length += uint(c.EncodingLength())
		}
	}

	l := uint(0)
	if value.Wire != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Wire_length).EncodingLength())
		l += encoder.Wire_length
	}
	if value.Name != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Name_length).EncodingLength())
		l += encoder.Name_length
	}
	encoder.length = l

}

func (context *WireNameFieldParsingContext) Init() {

}

func (encoder *WireNameFieldEncoder) EncodeInto(value *WireNameField, buf []byte) {

	pos := uint(0)

	if value.Wire != nil {
		buf[pos] = byte(1)
		pos += 1
		pos += uint(enc.TLNum(encoder.Wire_length).EncodeInto(buf[pos:]))
		for _, w := range value.Wire {
			copy(buf[pos:], w)
			pos += uint(len(w))
		}
	}
	if value.Name != nil {
		buf[pos] = byte(2)
		pos += 1
		pos += uint(enc.TLNum(encoder.Name_length).EncodeInto(buf[pos:]))
		for _, c := range value.Name {
			pos += uint(c.EncodeInto(buf[pos:]))
		}
	}
}

func (encoder *WireNameFieldEncoder) Encode(value *WireNameField) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *WireNameFieldParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*WireNameField, error) {

	var handled_Wire bool = false
	var handled_Name bool = false

	progress := -1
	_ = progress

	value := &WireNameField{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 1:
				if true {
					handled = true
					handled_Wire = true
					value.Wire, err = reader.ReadWire(int(l))
				}
			case 2:
				if true {
					handled = true
					handled_Name = true
					delegate := reader.Delegate(int(l))
					value.Name, err = delegate.ReadName()
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Wire && err == nil {
		value.Wire = nil
	}
	if !handled_Name && err == nil {
		value.Name = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *WireNameField) Encode() enc.Wire {
	encoder := WireNameFieldEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *WireNameField) Bytes() []byte {
	return value.Encode().Join()
}

func ParseWireNameField(reader enc.WireView, ignoreCritical bool) (*WireNameField, error) {
	context := WireNameFieldParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type MarkersEncoder struct {
	length uint

	startMarker     int
	startMarker_pos int
	Wire_length     uint
	argument        int
	Name_length     uint
	endMarker       int
	endMarker_pos   int
}

type MarkersParsingContext struct {
	startMarker int

	argument int

	endMarker int
}

func (encoder *MarkersEncoder) Init(value *Markers) {

	if value.Wire != nil {
		encoder.Wire_length = 0
		for _, c := range value.Wire {
			encoder.Wire_length += uint(len(c))
		}
	}

	if value.Name != nil {
		encoder.Name_length = 0
		for _, c := range value.Name {
			encoder.Name_length += uint(c.EncodingLength())
		}
	}

	l := uint(0)
	encoder.startMarker = int(l)
	if value.Wire != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Wire_length).EncodingLength())
		l += encoder.Wire_length
	}

	if value.Name != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Name_length).EncodingLength())
		l += encoder.Name_length
	}
	encoder.endMarker = int(l)
	encoder.length = l

}

func (context *MarkersParsingContext) Init() {

}

func (encoder *MarkersEncoder) EncodeInto(value *Markers, buf []byte) {

	pos := uint(0)

	encoder.startMarker_pos = int(pos)
	if value.Wire != nil {
		buf[pos] = byte(1)
		pos += 1
		pos += uint(enc.TLNum(encoder.Wire_length).EncodeInto(buf[pos:]))
		for _, w := range value.Wire {
			copy(buf[pos:], w)
			pos += uint(len(w))
		}
	}

	if value.Name != nil {
		buf[pos] = byte(2)
		pos += 1
		pos += uint(enc.TLNum(encoder.Name_length).EncodeInto(buf[pos:]))
		for _, c := range value.Name {
			pos += uint(c.EncodeInto(buf[pos:]))
		}
	}
	encoder.endMarker_pos = int(pos)
}

func (encoder *MarkersEncoder) Encode(value *Markers) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *MarkersParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*Markers, error) {

	var handled_startMarker bool = false
	var handled_Wire bool = false
	var handled_argument bool = false
	var handled_Name bool = false
	var handled_endMarker bool = false

	progress := -1
	_ = progress

	value := &Markers{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		for handled := false; !handled && progress < 5; progress++ {
			switch typ {
			case 1:
				if progress+1 == 1 {
					handled = true
					handled_Wire = true
					value.Wire, err = reader.ReadWire(int(l))
				}
			case 2:
				if progress+1 == 3 {
					handled = true
					handled_Name = true
					delegate := reader.Delegate(int(l))
					value.Name, err = delegate.ReadName()
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
				switch progress {
				case 0 - 1:
					handled_startMarker = true
					context.startMarker = int(startPos)
				case 1 - 1:
					handled_Wire = true
					value.Wire = nil
				case 2 - 1:
					handled_argument = true
					// base - skip
				case 3 - 1:
					handled_Name = true
					value.Name = nil
				case 4 - 1:
					handled_endMarker = true
					context.endMarker = int(startPos)
				}
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_startMarker && err == nil {
		context.startMarker = int(startPos)
	}
	if !handled_Wire && err == nil {
		value.Wire = nil
	}
	if !handled_argument && err == nil {
		// base - skip
	}
	if !handled_Name && err == nil {
		value.Name = nil
	}
	if !handled_endMarker && err == nil {
		context.endMarker = int(startPos)
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

type NoCopyStructEncoder struct {
	length uint

	wirePlan []uint

	Wire1_length uint

	Wire2_length uint
}

type NoCopyStructParsingContext struct {
}

func (encoder *NoCopyStructEncoder) Init(value *NoCopyStruct) {
	if value.Wire1 != nil {
		encoder.Wire1_length = 0
		for _, c := range value.Wire1 {
			encoder.Wire1_length += uint(len(c))
		}
	}

	if value.Wire2 != nil {
		encoder.Wire2_length = 0
		for _, c := range value.Wire2 {
			encoder.Wire2_length += uint(len(c))
		}
	}

	l := uint(0)
	if value.Wire1 != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Wire1_length).EncodingLength())
		l += encoder.Wire1_length
	}
	l += 1
	l += uint(1 + enc.Nat(value.Number).EncodingLength())
	if value.Wire2 != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Wire2_length).EncodingLength())
		l += encoder.Wire2_length
	}
	encoder.length = l

	wirePlan := make([]uint, 0, 8)
	l = uint(0)
	if value.Wire1 != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Wire1_length).EncodingLength())
		wirePlan = append(wirePlan, l)
		l = 0
		for range value.Wire1 {
			wirePlan = append(wirePlan, l)
			l = 0
		}
	}
	l += 1
	l += uint(1 + enc.Nat(value.Number).EncodingLength())
	if value.Wire2 != nil {
		l += 1
		l += uint(enc.TLNum(encoder.Wire2_length).EncodingLength())
		wirePlan = append(wirePlan, l)
		l = 0
		for range value.Wire2 {
			wirePlan = append(wirePlan, l)
			l = 0
		}
	}
	if l > 0 {
		wirePlan = append(wirePlan, l)
	}
	encoder.wirePlan = wirePlan
}

func (context *NoCopyStructParsingContext) Init() {

}

func (encoder *NoCopyStructEncoder) EncodeInto(value *NoCopyStruct, wire enc.Wire) {

	wireIdx := 0
	buf := wire[wireIdx]

	pos := uint(0)

	if value.Wire1 != nil {
		buf[pos] = byte(1)
		pos += 1
		pos += uint(enc.TLNum(encoder.Wire1_length).EncodeInto(buf[pos:]))
		wireIdx++
		pos = 0
		if wireIdx < len(wire) {
			buf = wire[wireIdx]
		} else {
			buf = nil
		}
		for _, w := range value.Wire1 {
			wire[wireIdx] = w
			wireIdx++
			pos = 0
			if wireIdx < len(wire) {
				buf = wire[wireIdx]
			} else {
				buf = nil
			}
		}
	}
	buf[pos] = byte(2)
	pos += 1

	buf[pos] = byte(enc.Nat(value.Number).EncodeInto(buf[pos+1:]))
	pos += uint(1 + buf[pos])
	if value.Wire2 != nil {
		buf[pos] = byte(3)
		pos += 1
		pos += uint(enc.TLNum(encoder.Wire2_length).EncodeInto(buf[pos:]))
		wireIdx++
		pos = 0
		if wireIdx < len(wire) {
			buf = wire[wireIdx]
		} else {
			buf = nil
		}
		for _, w := range value.Wire2 {
			wire[wireIdx] = w
			wireIdx++
			pos = 0
			if wireIdx < len(wire) {
				buf = wire[wireIdx]
			} else {
				buf = nil
			}
		}
	}
}

func (encoder *NoCopyStructEncoder) Encode(value *NoCopyStruct) enc.Wire {
	total := uint(0)
	for _, l := range encoder.wirePlan {
		total += l
	}
	content := make([]byte, total)

	wire := make(enc.Wire, len(encoder.wirePlan))
	for i, l := range encoder.wirePlan {
		if l > 0 {
			wire[i] = content[:l]
			content = content[l:]
		}
	}
	encoder.EncodeInto(value, wire)

	return wire
}

func (context *NoCopyStructParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*NoCopyStruct, error) {

	var handled_Wire1 bool = false
	var handled_Number bool = false
	var handled_Wire2 bool = false

	progress := -1
	_ = progress

	value := &NoCopyStruct{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 1:
				if true {
					handled = true
					handled_Wire1 = true
					value.Wire1, err = reader.ReadWire(int(l))
				}
			case 2:
				if true {
					handled = true
					handled_Number = true
					value.Number = uint64(0)
					{
						for i := 0; i < int(l); i++ {
							x := byte(0)
							x, err = reader.ReadByte()
							if err != nil {
								if err == io.EOF {
									err = io.ErrUnexpectedEOF
								}
								break
							}
							value.Number = uint64(value.Number<<8) | uint64(x)
						}
					}
				}
			case 3:
				if true {
					handled = true
					handled_Wire2 = true
					value.Wire2, err = reader.ReadWire(int(l))
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Wire1 && err == nil {
		value.Wire1 = nil
	}
	if !handled_Number && err == nil {
		err = enc.ErrSkipRequired{Name: "Number", TypeNum: 2}
	}
	if !handled_Wire2 && err == nil {
		value.Wire2 = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *NoCopyStruct) Encode() enc.Wire {
	encoder := NoCopyStructEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *NoCopyStruct) Bytes() []byte {
	return value.Encode().Join()
}

func ParseNoCopyStruct(reader enc.WireView, ignoreCritical bool) (*NoCopyStruct, error) {
	context := NoCopyStructParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type StrFieldEncoder struct {
	length uint
}

type StrFieldParsingContext struct {
}

func (encoder *StrFieldEncoder) Init(value *StrField) {

	l := uint(0)
	l += 1
	l += uint(enc.TLNum(len(value.Str1)).EncodingLength())
	l += uint(len(value.Str1))
	if optval, ok := value.Str2.Get(); ok {
		l += 1
		l += uint(enc.TLNum(len(optval)).EncodingLength())
		l += uint(len(optval))
	}
	encoder.length = l

}

func (context *StrFieldParsingContext) Init() {

}

func (encoder *StrFieldEncoder) EncodeInto(value *StrField, buf []byte) {

	pos := uint(0)

	buf[pos] = byte(1)
	pos += 1
	pos += uint(enc.TLNum(len(value.Str1)).EncodeInto(buf[pos:]))
	copy(buf[pos:], value.Str1)
	pos += uint(len(value.Str1))
	if optval, ok := value.Str2.Get(); ok {
		buf[pos] = byte(2)
		pos += 1
		pos += uint(enc.TLNum(len(optval)).EncodeInto(buf[pos:]))
		copy(buf[pos:], optval)
		pos += uint(len(optval))
	}
}

func (encoder *StrFieldEncoder) Encode(value *StrField) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *StrFieldParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*StrField, error) {

	var handled_Str1 bool = false
	var handled_Str2 bool = false

	progress := -1
	_ = progress

	value := &StrField{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 1:
				if true {
					handled = true
					handled_Str1 = true
					{
						var builder strings.Builder
						_, err = reader.CopyN(&builder, int(l))
						if err == nil {
							value.Str1 = builder.String()
						}
					}
				}
			case 2:
				if true {
					handled = true
					handled_Str2 = true
					{
						var builder strings.Builder
						_, err = reader.CopyN(&builder, int(l))
						if err == nil {
							value.Str2.Set(builder.String())
						}
					}
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Str1 && err == nil {
		err = enc.ErrSkipRequired{Name: "Str1", TypeNum: 1}
	}
	if !handled_Str2 && err == nil {
		value.Str2.Unset()
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *StrField) Encode() enc.Wire {
	encoder := StrFieldEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *StrField) Bytes() []byte {
	return value.Encode().Join()
}

func ParseStrField(reader enc.WireView, ignoreCritical bool) (*StrField, error) {
	context := StrFieldParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}

type FixedUintFieldEncoder struct {
	length uint
}

type FixedUintFieldParsingContext struct {
}

func (encoder *FixedUintFieldEncoder) Init(value *FixedUintField) {

	l := uint(0)
	l += 1
	l += 1 + 1
	if value.U32.IsSet() {
		l += 1
		l += 1 + 4
	}
	if value.U64.IsSet() {
		l += 1
		l += 1 + 8
	}
	if value.BytePtr != nil {
		l += 1
		l += 2
	}
	encoder.length = l

}

func (context *FixedUintFieldParsingContext) Init() {

}

func (encoder *FixedUintFieldEncoder) EncodeInto(value *FixedUintField, buf []byte) {

	pos := uint(0)

	buf[pos] = byte(1)
	pos += 1
	buf[pos] = 1
	buf[pos+1] = byte(value.Byte)
	pos += 2
	if optval, ok := value.U32.Get(); ok {
		buf[pos] = byte(2)
		pos += 1
		buf[pos] = 4
		binary.BigEndian.PutUint32(buf[pos+1:], uint32(optval))
		pos += 5
	}
	if optval, ok := value.U64.Get(); ok {
		buf[pos] = byte(3)
		pos += 1
		buf[pos] = 8
		binary.BigEndian.PutUint64(buf[pos+1:], uint64(optval))
		pos += 9
	}
	if value.BytePtr != nil {
		buf[pos] = byte(4)
		pos += 1
		buf[pos] = 1
		buf[pos+1] = byte(*value.BytePtr)
		pos += 2
	}
}

func (encoder *FixedUintFieldEncoder) Encode(value *FixedUintField) enc.Wire {

	wire := make(enc.Wire, 1)
	wire[0] = make([]byte, encoder.length)
	buf := wire[0]
	encoder.EncodeInto(value, buf)

	return wire
}

func (context *FixedUintFieldParsingContext) Parse(reader enc.WireView, ignoreCritical bool) (*FixedUintField, error) {

	var handled_Byte bool = false
	var handled_U32 bool = false
	var handled_U64 bool = false
	var handled_BytePtr bool = false

	progress := -1
	_ = progress

	value := &FixedUintField{}
	var err error
	var startPos int
	for {
		startPos = reader.Pos()
		if startPos >= reader.Length() {
			break
		}
		typ := enc.TLNum(0)
		l := enc.TLNum(0)
		typ, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}
		l, err = reader.ReadTLNum()
		if err != nil {
			return nil, enc.ErrFailToParse{TypeNum: 0, Err: err}
		}

		err = nil
		if handled := false; true {
			switch typ {
			case 1:
				if true {
					handled = true
					handled_Byte = true
					value.Byte, err = reader.ReadByte()
					if err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
				}
			case 2:
				if true {
					handled = true
					handled_U32 = true
					{
						optval := uint32(0)
						optval = uint32(0)
						{
							for i := 0; i < int(l); i++ {
								x := byte(0)
								x, err = reader.ReadByte()
								if err != nil {
									if err == io.EOF {
										err = io.ErrUnexpectedEOF
									}
									break
								}
								optval = uint32(optval<<8) | uint32(x)
							}
						}
						value.U32.Set(optval)
					}
				}
			case 3:
				if true {
					handled = true
					handled_U64 = true
					{
						optval := uint64(0)
						optval = uint64(0)
						{
							for i := 0; i < int(l); i++ {
								x := byte(0)
								x, err = reader.ReadByte()
								if err != nil {
									if err == io.EOF {
										err = io.ErrUnexpectedEOF
									}
									break
								}
								optval = uint64(optval<<8) | uint64(x)
							}
						}
						value.U64.Set(optval)
					}
				}
			case 4:
				if true {
					handled = true
					handled_BytePtr = true
					{
						buf, err := reader.ReadBuf(1)
						if err == io.EOF {
							err = io.ErrUnexpectedEOF
						}
						value.BytePtr = &buf[0]
					}
				}
			default:
				if !ignoreCritical && ((typ <= 31) || ((typ & 1) == 1)) {
					return nil, enc.ErrUnrecognizedField{TypeNum: typ}
				}
				handled = true
				err = reader.Skip(int(l))
			}
			if err == nil && !handled {
			}
			if err != nil {
				return nil, enc.ErrFailToParse{TypeNum: typ, Err: err}
			}
		}
	}

	startPos = reader.Pos()
	err = nil

	if !handled_Byte && err == nil {
		err = enc.ErrSkipRequired{Name: "Byte", TypeNum: 1}
	}
	if !handled_U32 && err == nil {
		value.U32.Unset()
	}
	if !handled_U64 && err == nil {
		value.U64.Unset()
	}
	if !handled_BytePtr && err == nil {
		value.BytePtr = nil
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (value *FixedUintField) Encode() enc.Wire {
	encoder := FixedUintFieldEncoder{}
	encoder.Init(value)
	return encoder.Encode(value)
}

func (value *FixedUintField) Bytes() []byte {
	return value.Encode().Join()
}

func ParseFixedUintField(reader enc.WireView, ignoreCritical bool) (*FixedUintField, error) {
	context := FixedUintFieldParsingContext{}
	context.Init()
	return context.Parse(reader, ignoreCritical)
}
