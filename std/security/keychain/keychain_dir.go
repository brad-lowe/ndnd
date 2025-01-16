package keychain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	enc "github.com/named-data/ndnd/std/encoding"
	"github.com/named-data/ndnd/std/log"
	"github.com/named-data/ndnd/std/ndn"
	"github.com/named-data/ndnd/std/security"
)

const EXT_KEY = ".key"
const EXT_CERT = ".cert"

// KeyChainDir is a directory-based keychain.
type KeyChainDir struct {
	wmut sync.Mutex
	mem  ndn.KeyChain
	path string
}

// NewKeyChainDir creates a new in-memory keychain.
func NewKeyChainDir(path string, pubStore ndn.Store) (ndn.KeyChain, error) {
	kc := &KeyChainDir{
		wmut: sync.Mutex{},
		mem:  NewKeyChainMem(pubStore),
		path: path,
	}

	// Populate keychain from disk
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), EXT_KEY) &&
			!strings.HasSuffix(entry.Name(), EXT_CERT) {
			continue
		}

		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(path, entry.Name())
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Warn(kc, "Failed to read keychain entry", "file", filename, "error", err)
			continue
		}

		err = InsertFile(kc.mem, content)
		if err != nil {
			log.Error(kc, "Failed to insert keychain entries", "file", filename, "error", err)
		}
	}

	return kc, nil
}

func (kc *KeyChainDir) String() string {
	return fmt.Sprintf("KeyChainDir (%s)", kc.path)
}

func (kc *KeyChainDir) GetIdentities() []ndn.Identity {
	return kc.mem.GetIdentities()
}

func (kc *KeyChainDir) GetIdentity(name enc.Name) ndn.Identity {
	return kc.mem.GetIdentity(name)
}

func (kc *KeyChainDir) InsertKey(signer ndn.Signer) error {
	err := kc.mem.InsertKey(signer)
	if err != nil {
		return err
	}

	secret, err := EncodeSecret(signer)
	if err != nil {
		return err
	}

	return kc.writeFile(secret.Join(), EXT_KEY)
}

func (kc *KeyChainDir) InsertCert(wire []byte) error {
	err := kc.mem.InsertCert(wire)
	if err != nil {
		return err
	}

	return kc.writeFile(wire, EXT_CERT)
}

func (kc *KeyChainDir) writeFile(wire []byte, ext string) error {
	hash := sha256.Sum256(wire)
	filename := hex.EncodeToString(hash[:])
	path := filepath.Join(kc.path, filename+ext)

	str, err := security.PemEncode(wire)
	if err != nil {
		return err
	}

	kc.wmut.Lock()
	defer kc.wmut.Unlock()

	return os.WriteFile(path, str, 0644)
}
