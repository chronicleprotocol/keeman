//  Copyright (C) 2020 Maker Ecosystem Growth Holdings, INC.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as
//  published by the Free Software Foundation, either version 3 of the
//  License, or (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ssb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"go.cryptoscope.co/ssb"

	"github.com/chronicleprotocol/keeman/rand"
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
