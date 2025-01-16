package security

import (
	"encoding/pem"
	"errors"

	enc "github.com/named-data/ndnd/std/encoding"
	"github.com/named-data/ndnd/std/log"
	"github.com/named-data/ndnd/std/ndn"

	spec "github.com/named-data/ndnd/std/ndn/spec_2022"
)

const PEM_TYPE_CERT = "NDN Certificate"
const PEM_TYPE_SECRET = "NDN Key"

// PemEncode converts an NDN data to a text representation following RFC 7468.
func PemEncode(raw []byte) ([]byte, error) {
	data, _, err := spec.Spec{}.ReadData(enc.NewBufferReader(raw))
	if err != nil {
		return nil, err
	}

	if data.ContentType() == nil {
		return nil, ndn.ErrInvalidValue{Item: "content type"}
	}

	if data.Signature() == nil {
		return nil, ndn.ErrInvalidValue{Item: "signature"}
	}

	// Explanatory text before the block
	headers := map[string]string{
		"Name": data.Name().String(),
	}

	var pemType string
	switch *data.ContentType() {
	case ndn.ContentTypeKey:
		pemType = PEM_TYPE_CERT
	case ndn.ContentTypeSecret:
		pemType = PEM_TYPE_SECRET
	default:
		return nil, errors.New("unsupported content type")
	}

	if nb, na := data.Signature().Validity(); nb != nil && na != nil {
		headers["Validity"] = nb.String() + " - " + na.String()
	}

	switch data.Signature().SigType() {
	case ndn.SignatureDigestSha256:
		headers["SigType"] = "Digest-SHA256"
	case ndn.SignatureSha256WithRsa:
		headers["SigType"] = "RSA-SHA256"
	case ndn.SignatureSha256WithEcdsa:
		headers["SigType"] = "ECDSA-SHA256"
	case ndn.SignatureHmacWithSha256:
		headers["SigType"] = "HMAC-SHA256"
	case ndn.SignatureEd25519:
		headers["SigType"] = "Ed25519"
	default:
		headers["SigType"] = "Unknown"
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:    pemType,
		Headers: headers,
		Bytes:   raw,
	}), nil
}

// PemDecode converts a text representation of an NDN data.
func PemDecode(str []byte) [][]byte {
	ret := make([][]byte, 0)

	for {
		block, rest := pem.Decode(str)
		if block == nil {
			break
		}
		str = rest

		if block.Type != PEM_TYPE_CERT && block.Type != PEM_TYPE_SECRET {
			log.Warn(nil, "Unsupported PEM type", "type", block.Type)
			continue
		}

		ret = append(ret, block.Bytes)
	}

	return ret
}
