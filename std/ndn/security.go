package ndn

import (
	"time"

	enc "github.com/named-data/ndnd/std/encoding"
)

// Signature is the abstract of the signature of a packet.
// Some of the fields are invalid for Data or Interest.
type Signature interface {
	SigType() SigType
	KeyName() enc.Name
	SigNonce() []byte
	SigTime() *time.Time
	SigSeqNum() *uint64
	Validity() (notBefore, notAfter *time.Time)
	SigValue() []byte
}

// Signer is the interface of a NDN packet signer.
type Signer interface {
	// SigInfo returns the configuration of the signature.
	Type() SigType
	// KeyName returns the key name of the signer.
	KeyName() enc.Name
	// KeyLocator returns the key locator name of the signer.
	KeyLocator() enc.Name
	// EstimateSize gives the approximate size of the signature in bytes.
	EstimateSize() uint
	// Sign computes the signature value of a wire.
	Sign(enc.Wire) ([]byte, error)
	// Public returns the public key of the signer or nil.
	Public() ([]byte, error)
}

// SigChecker is a basic function to check the signature of a packet.
// In NTSchema, policies&sub-trees are supposed to be used for validation;
// SigChecker is only designed for low-level engine.
// Create a go routine for time consuming jobs.
type SigChecker func(name enc.Name, sigCovered enc.Wire, sig Signature) bool

// KeyChain is the interface of a keychain.
// Note that Keychains are not thread-safe, and the owner should provide a lock.
type KeyChain interface {
	// String provides the log identifier of the keychain.
	String() string
	// GetIdentities returns all identities in the keychain.
	GetIdentities() []KeyChainIdentity
	// GetIdentity returns the identity by full name.
	GetIdentity(enc.Name) KeyChainIdentity
	// InsertKey inserts a key to the keychain.
	InsertKey(Signer) error
	// InsertCert inserts a certificate to the keychain.
	InsertCert([]byte) error
}

// KeyChainIdentity is the interface of a signing identity.
type KeyChainIdentity interface {
	// Name returns the full name of the identity.
	Name() enc.Name
	// Keys returns all keys of the identity.
	Keys() []KeyChainKey
}

// KeyChainKey is the interface of a key in the keychain.
type KeyChainKey interface {
	// KeyName returns the key name of the key.
	KeyName() enc.Name
	// Signer returns the signer of the key.
	Signer() Signer
	// UniqueCerts returns a list of unique cert names in the key.
	// The version number is always set to zero.
	UniqueCerts() []enc.Name
}
