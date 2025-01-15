package sec

import (
	"crypto/elliptic"
	"fmt"
	"os"
	"strconv"

	enc "github.com/named-data/ndnd/std/encoding"
	"github.com/named-data/ndnd/std/ndn"
	"github.com/named-data/ndnd/std/security"
	"github.com/named-data/ndnd/std/security/crypto"
	"github.com/named-data/ndnd/std/security/keychain"
)

func keygen(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: sec keygen <identity> <key-type> [params]\n")
		fmt.Fprintf(os.Stderr, "  key-type: rsa|ecc|ed25519\n")
		fmt.Fprintf(os.Stderr, "Example: %s sec keygen /ndn/alice ed25519\n", os.Args[0])
		return
	}

	identity := args[1]
	keyType := args[2]
	args = args[3:]

	name, err := enc.NameFromStr(identity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid identity: %s\n", identity)
		os.Exit(1)
	}

	var signer ndn.Signer
	switch keyType {
	case "rsa":
		signer = keygenRsa(args, name)
	case "ed25519":
		signer = keygenEd25519(args, name)
	case "ecc":
		signer = keygecEcc(args, name)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported key type: %s\n", keyType)
		os.Exit(1)
	}

	secret, err := keychain.EncodeSecret(signer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode secret key: %s\n", err)
		os.Exit(1)
	}

	out, err := security.TxtFrom(secret.Join())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to convert secret key to text: %s\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(out)
}

func keygenRsa(args []string, name enc.Name) ndn.Signer {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: keygen rsa <key-size>\n")
		os.Exit(1)
	}

	keySize, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid key size: %s\n", args[0])
		os.Exit(1)
	}

	signer, err := crypto.KeygenRsa(name, keySize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate RSA key: %s\n", err)
		os.Exit(1)
	}

	return signer
}

func keygenEd25519(_ []string, name enc.Name) ndn.Signer {
	signer, err := crypto.KeygenEd25519(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate Ed25519 key: %s\n", err)
	}
	return signer
}

func keygecEcc(args []string, name enc.Name) ndn.Signer {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: keygen ecc <curve>\n")
		fmt.Fprintf(os.Stderr, "Supported curves: secp256r1, secp384r1, secp521r1\n")
		os.Exit(1)
	}

	var curve elliptic.Curve
	switch args[0] {
	case "secp256r1":
		curve = elliptic.P256()
	case "secp384r1":
		curve = elliptic.P384()
	case "secp521r1":
		curve = elliptic.P521()
	default:
		fmt.Fprintf(os.Stderr, "Unsupported curve: %s\n", args[0])
	}

	signer, err := crypto.KeygenEcc(name, curve)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate EC key: %s\n", err)
	}
	return signer
}
