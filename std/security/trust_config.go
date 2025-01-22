package security

import (
	"fmt"
	"sync"

	enc "github.com/named-data/ndnd/std/encoding"
	"github.com/named-data/ndnd/std/log"
	"github.com/named-data/ndnd/std/ndn"
	spec "github.com/named-data/ndnd/std/ndn/spec_2022"
)

// TrustConfig is the configuration of the trust module.
type TrustConfig struct {
	// KeyChain is the keychain.
	KeyChain ndn.KeyChain
	// Schema is the trust schema.
	Schema ndn.TrustSchema
	// Roots are the full names of the trust anchors.
	Roots []enc.Name

	// mutex is the lock for keychain.
	mutex sync.RWMutex
}

// NewTrustConfig creates a new TrustConfig.
// ALl roots must be full names and already present in the keychain.
func NewTrustConfig(keyChain ndn.KeyChain, schema ndn.TrustSchema, roots []enc.Name) (*TrustConfig, error) {
	// Check if we have some roots
	if len(roots) == 0 {
		return nil, fmt.Errorf("no trust anchors provided")
	}

	// Check if all roots are present in the keychain
	for _, root := range roots {
		if cert, _ := keyChain.Store().Get(root, false); cert == nil {
			return nil, fmt.Errorf("trust anchor not found in keychain: %s", root)
		}
	}

	return &TrustConfig{
		KeyChain: keyChain,
		Schema:   schema,
		Roots:    roots,
	}, nil
}

// Suggest suggests a signer for a given name.
func (tc *TrustConfig) Suggest(name enc.Name) ndn.Signer {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return tc.Schema.Suggest(name, tc.KeyChain)
}

// ValidateArgs are the arguments for the Validate function.
type ValidateArgs struct {
	Data     ndn.Data
	Fetch    func(enc.Name, *ndn.InterestConfig, func(ndn.Data, []byte, error))
	Callback func(bool, error)

	cert  ndn.Data
	depth int
}

// Validate validates a Data packet using a fetch API.
func (tc *TrustConfig) Validate(args ValidateArgs) {
	if args.Data == nil {
		args.Callback(false, fmt.Errorf("data is nil"))
		return
	}

	if args.depth == 0 {
		args.depth = 32
	} else if args.depth <= 1 {
		args.Callback(false, fmt.Errorf("max depth reached"))
		return
	} else {
		args.depth--
	}

	if args.cert != nil {
		// Check if the data claims to be a root certificate.
		// This breaks the recursion for validation.
		if args.Data.Name().Equal(args.cert.Name()) {
			for _, root := range tc.Roots {
				if args.Data.Name().Equal(root) {
					args.Callback(true, nil)
					return
				}
			}
			args.Callback(false, fmt.Errorf("data claims to be a trust anchor: %s", args.Data.Name()))
			return
		}

		// Check schema if the key is allowed
		if !tc.Schema.Check(args.Data.Name(), args.cert.Name()) {
			args.Callback(false, fmt.Errorf("key is not allowed: %s signed by %s", args.Data.Name(), args.cert.Name()))
			return
		}

		// TODO: validate signature

		// Recursively validate the certificate
		tc.Validate(ValidateArgs{
			Data:     args.cert,
			Fetch:    args.Fetch,
			Callback: args.Callback,
			depth:    args.depth,
			cert:     nil,
		})
		return
	}

	// Get the certificate using the key locator
	signature := args.Data.Signature()
	if signature == nil {
		args.Callback(false, fmt.Errorf("signature is nil"))
		return
	}

	// Get the key locator
	keyLocator := signature.KeyName()
	if keyLocator == nil {
		args.Callback(false, fmt.Errorf("key locator is nil"))
		return
	}

	// Detect if this is a self-signed certificate, and automatically pick the cert
	// as itself to verify in this case.
	if args.Data.ContentType() != nil && *args.Data.ContentType() == ndn.ContentTypeKey && keyLocator.IsPrefix(args.Data.Name()) {
		args.cert = args.Data
		tc.Validate(args)
		return
	}

	// Attempt to get cert from store
	certBytes, err := tc.KeyChain.Store().Get(keyLocator, true)
	if err != nil {
		log.Error(nil, "Failed to get certificate from store", "error", err)
		args.Callback(false, err)
		return // store is likely broken
	}
	if len(certBytes) > 0 {
		// Attempt to parse the certificate
		args.cert, _, err = spec.Spec{}.ReadData(enc.NewBufferReader(certBytes))
		if err != nil {
			log.Error(nil, "Failed to parse certificate in store", "error", err)
			args.cert = nil
		}
	}

	// If not found, attempt to fetch cert from network
	if args.cert == nil {
		args.Fetch(keyLocator, &ndn.InterestConfig{
			CanBePrefix: true,
			MustBeFresh: true,
		}, func(cert ndn.Data, wire []byte, err error) {
			if err != nil {
				args.Callback(false, err)
				return // failed to fetch cert
			}
			args.cert = cert

			// Monkey patch the callback to store the cert in keychain
			// if the validation passes.
			origCallback := args.Callback
			args.Callback = func(valid bool, err error) {
				if valid && err == nil {
					tc.mutex.Lock()
					err := tc.KeyChain.InsertCert(wire)
					tc.mutex.Unlock()
					if err != nil {
						log.Error(tc.KeyChain, "Failed to insert certificate", "error", err)
					}
				}
				origCallback(valid, err) // continue validation
			}

			tc.Validate(args)
		})
	} else {
		tc.Validate(args)
	}
}
