package eth

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

func NewKeyWithID(privateKey *ecdsa.PrivateKey) (*keystore.Key, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}, nil
}
