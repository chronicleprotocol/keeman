package ssb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/ssbc/go-ssb"

	"github.com/chronicleprotocol/keeman/pkg/rand"
)

type Caps struct {
	Shs  []byte
	Sign []byte
}

func (c Caps) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Shs    string `json:"shs"`
		Sign   string `json:"sign,omitempty"`
		Invite string `json:"invite,omitempty"`
	}{
		Shs:  base64.StdEncoding.EncodeToString(c.Shs),
		Sign: base64.StdEncoding.EncodeToString(c.Sign),
	})
}

func NewCaps(seed []byte) (*Caps, error) {
	randBytes, err := rand.SeededRandBytesGen(seed, 32)
	if err != nil {
		return nil, err
	}
	return &Caps{
		Shs:  randBytes(),
		Sign: randBytes(),
	}, nil
}

type Secret struct{ ssb.KeyPair }

func (s Secret) MarshalJSON() ([]byte, error) {
	b := new(bytes.Buffer)
	if err := ssb.EncodeKeyPairAsJSON(s.KeyPair, b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func NewSecret(b []byte) (*Secret, error) {
	kp, err := ssb.NewKeyPair(bytes.NewReader(b), "ed25519")
	if err != nil {
		return nil, err
	}
	return &Secret{KeyPair: kp}, nil
}
