package keychain_test

import (
	"encoding/base64"
	"os"
	"testing"
	"time"

	enc "github.com/named-data/ndnd/std/encoding"
	"github.com/named-data/ndnd/std/ndn"
	spec "github.com/named-data/ndnd/std/ndn/spec_2022"
	"github.com/named-data/ndnd/std/object"
	sec "github.com/named-data/ndnd/std/security"
	"github.com/named-data/ndnd/std/security/keychain"
	sig "github.com/named-data/ndnd/std/security/signer"
	"github.com/named-data/ndnd/std/utils"
	"github.com/stretchr/testify/require"
)

const CERT_ROOT = `
Bv0BSwcjCANuZG4IA0tFWQgIJ8SyKp97gScIA25kbjYIAAABgHX6c7QUCRgBAhkE
ADbugBVbMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEPuDnW4oq0mULLT8PDXh0
zuBg+0SJ1yPC85jylUU+hgxX9fDNyjlykLrvb1D6IQRJWJHMKWe6TJKPUhGgOT65
8hZyGwEDHBYHFAgDbmRuCANLRVkICCfEsiqfe4En/QD9Jv0A/g8yMDIyMDQyOVQx
NTM5NTD9AP8PMjAyNjEyMzFUMjM1OTU5/QECKf0CACX9AgEIZnVsbG5hbWX9AgIV
TkROIFRlc3RiZWQgUm9vdCAyMjA0F0gwRgIhAPYUOjNakdfDGh5j9dcCGOz+Ie1M
qoAEsjM9PEUEWbnqAiEApu0rg9GAK1LNExjLYAF6qVgpWQgU+atPn63Gtuubqyg=
`

const KEY_ALICE = `
-----BEGIN NDN KEY-----
Name: /ndn/alice/KEY/cK%1D%A4%E1%5B%91%CF
SigType: Ed25519

BsoHGwgDbmRuCAVhbGljZQgDS0VZCAhjSx2k4VuRzxQDGAEJFUC64F62YK0/v5z4
fjONZO7Y4PNqy7FiDnar33uVO71FLK6Vp8GrPCkEhuODl6GBv2nUuovtO9KtHW11
8apSS093FiIbAQUcHQcbCANuZG4IBWFsaWNlCANLRVkICGNLHaThW5HPF0Cw3Oh7
I2jmBBxop1bIPXq292TfltVwhdbB3/yUXkKcg3BYbY6vcAhNNqrG2B+G/iHvKGsy
DpvDtnlEN72hIeIP
-----END NDN KEY-----
`

var CERT_ROOT_NAME, _ = enc.NameFromStr("/ndn/KEY/%27%C4%B2%2A%9F%7B%81%27/ndn/v=1651246789556")
var KEY_ALICE_NAME, _ = enc.NameFromStr("/ndn/alice/KEY/cK%1D%A4%E1%5B%91%CF")

func signCert(t *testing.T, signer ndn.Signer) []byte {
	certData, _, _ := spec.Spec{}.ReadData(enc.NewWireReader(
		utils.WithoutErr(sig.MarshalSecret(signer))))
	cert, err := sec.SignCert(sec.SignCertArgs{
		Signer:    signer,
		Data:      certData,
		IssuerId:  enc.NewStringComponent(enc.TypeGenericNameComponent, "Test"),
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),
	})
	require.NoError(t, err)
	return cert.Join()
}

func TestKeyChainMem(t *testing.T) {
	utils.SetTestingT(t)

	store := object.NewMemoryStore()
	kc := keychain.NewKeyChainMem(store)

	// Insert a key
	idName1, _ := enc.NameFromStr("/my/test/identity")
	signer11 := utils.WithoutErr(sig.KeygenEd25519(sec.MakeKeyName(idName1)))
	require.NoError(t, kc.InsertKey(signer11))

	// Check key in keychain
	identity1 := kc.GetIdentity(idName1)
	require.NotNil(t, identity1)
	require.Equal(t, idName1, identity1.Name())
	require.Len(t, identity1.Keys(), 1)
	require.Equal(t, signer11, identity1.Keys()[0].Signer())

	// Insert another key for the same identity
	signer12 := utils.WithoutErr(sig.KeygenEd25519(sec.MakeKeyName(idName1)))
	require.NoError(t, kc.InsertKey(signer12))
	require.Len(t, identity1.Keys(), 2)

	// Generate cert11 for first signer
	// Make sure signer is the default signer
	cert111 := signCert(t, signer11)
	require.NoError(t, kc.InsertCert(cert111))
	require.Equal(t, signer11, identity1.Keys()[0].Signer())

	// Generate newer cert for second signer
	time.Sleep(5 * time.Millisecond) // new version
	cert121 := signCert(t, signer12)
	require.NoError(t, kc.InsertCert(cert121))

	// Check if the default signer changes to the newer signer
	key12 := identity1.Keys()[0]
	require.Equal(t, signer12, key12.Signer())
	require.Len(t, key12.UniqueCerts(), 1)

	// Insert another cert for second signer
	cert122 := signCert(t, signer12)
	require.NoError(t, kc.InsertCert(cert122))
	require.Len(t, key12.UniqueCerts(), 1) // same issuer

	// Lookup non-existing identity
	idName2, _ := enc.NameFromStr("/my/test/identity2")
	require.Nil(t, kc.GetIdentity(idName2))

	// Insert key for identity2
	signer21 := utils.WithoutErr(sig.KeygenEd25519(sec.MakeKeyName(idName2)))
	require.NoError(t, kc.InsertKey(signer21))
	identity2 := kc.GetIdentity(idName2)
	require.NotNil(t, identity2)
	require.Len(t, identity2.Keys(), 1)
	require.Equal(t, signer21, identity2.Keys()[0].Signer())

	// Insert cert for another key for identity2 before key
	signer22 := utils.WithoutErr(sig.KeygenEd25519(sec.MakeKeyName(idName2)))
	cert22 := signCert(t, signer22)
	require.NoError(t, kc.InsertCert(cert22))
	require.Len(t, identity2.Keys(), 1)
	require.NoError(t, kc.InsertKey(signer22))
	require.Len(t, identity2.Keys(), 2)
	require.Equal(t, signer22, identity2.Keys()[0].Signer()) // has cert

	// Insert invalid key
	signerInvalid := sig.NewSha256Signer()
	require.Error(t, kc.InsertKey(signerInvalid))

	// Insert a certificate.
	certRoot, _ := base64.StdEncoding.DecodeString(CERT_ROOT)
	require.NoError(t, kc.InsertCert(certRoot))

	// Check certificate in store
	data, err := store.Get(CERT_ROOT_NAME, false)
	require.NoError(t, err)
	require.Equal(t, certRoot, data)
}

func TestKeyChainDir(t *testing.T) {
	utils.SetTestingT(t)

	store := object.NewMemoryStore()

	// Create a temporary directory
	dirname := "./ndn-test-keychain"
	require.NoError(t, os.RemoveAll(dirname))
	defer os.RemoveAll(dirname)
	require.NoError(t, os.Mkdir(dirname, 0755))

	// Write root cert (raw)
	rootCert, _ := base64.StdEncoding.DecodeString(CERT_ROOT)
	require.NoError(t, os.WriteFile(dirname+"/root.cert", rootCert, 0644))

	// Write Alice key (text)
	require.NoError(t, os.WriteFile(dirname+"/alice.key", []byte(KEY_ALICE), 0644))

	// Create a keychain
	kc, err := keychain.NewKeyChainDir(dirname, store)
	require.NoError(t, err)

	// Check root cert
	data, err := store.Get(CERT_ROOT_NAME, false)
	require.NoError(t, err)
	require.Equal(t, rootCert, data)

	// Check Alice key
	identity := kc.GetIdentity(KEY_ALICE_NAME[:2])
	require.NotNil(t, identity)
	require.Len(t, identity.Keys(), 1)
	require.Equal(t, identity.Keys()[0].KeyName(), KEY_ALICE_NAME)

	// Check Alice key is not in store
	data, _ = store.Get(KEY_ALICE_NAME, false)
	require.Nil(t, data)
}
