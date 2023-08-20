package eth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"io"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

func NewKeyWithID(privateKey *ecdsa.PrivateKey) (*keystore.Key, error) {
	return NewKeyWithIDFromReader(privateKey, rand.Reader)
}
func NewKeyWithIDFromReader(privateKey *ecdsa.PrivateKey, r io.Reader) (*keystore.Key, error) {
	id, err := uuid.NewRandomFromReader(r)
	if err != nil {
		return nil, err
	}
	return &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}, nil
}
