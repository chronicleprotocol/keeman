package rand

import (
	"bytes"
	"encoding/binary"
	"math/rand"
)

func SeededRandBytesGen(seedBytes []byte, len int) (func() []byte, error) {
	var seed int64
	buf := bytes.NewBuffer(seedBytes)
	if err := binary.Read(buf, binary.BigEndian, &seed); err != nil {
		return nil, err
	}
	r := rand.New(rand.NewSource(seed))

	return func() []byte {
		rb := make([]byte, len)
		r.Read(rb)
		return rb
	}, nil
}
