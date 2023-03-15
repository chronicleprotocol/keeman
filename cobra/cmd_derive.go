package cobra

import (
	"bytes"
	"crypto/ecdsa"
	rand2 "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/spf13/cobra"

	"github.com/chronicleprotocol/keeman/eth"
	"github.com/chronicleprotocol/keeman/rand"
	"github.com/chronicleprotocol/keeman/ssb"
	"github.com/chronicleprotocol/keeman/tor"
)

const (
	Ethereum                  = "m/44'/60'/0'/0"
	EthereumClassic           = "m/44'/61'/0'/0"
	EthereumTestnetRopsten    = "m/44'/1'/0'/0"
	EthereumLedger            = "m/44'/60'/0'"
	EthereumClassicLedger     = "m/44'/60'/160720'/0"
	EthereumLedgerLive        = "m/44'/60'"
	EthereumClassicLedgerLive = "m/44'/61'"
	RSKMainnet                = "m/44'/137'/0'/0"
	Expanse                   = "m/44'/40'/0'/0"
	Ubiq                      = "m/44'/108'/0'/0"
	Ellaism                   = "m/44'/163'/0'/0"
	EtherGem                  = "m/44'/1987'/0'/0"
	Callisto                  = "m/44'/820'/0'/0"
	EthereumSocial            = "m/44'/1128'/0'/0"
	Musicoin                  = "m/44'/184'/0'/0"
	EOSClassic                = "m/44'/2018'/0'/0"
	Akroma                    = "m/44'/200625'/0'/0"
	EtherSocialNetwork        = "m/44'/31102'/0'/0"
	PIRL                      = "m/44'/164'/0'/0"
	GoChain                   = "m/44'/6060'/0'/0"
	Ether                     = "m/44'/1313114'/0'/0"
	Atheios                   = "m/44'/1620'/0'/0"
	TomoChain                 = "m/44'/889'/0'/0"
	MixBlockchain             = "m/44'/76'/0'/0"
	Iolite                    = "m/44'/1171337'/0'/0"
	ThunderCore               = "m/44'/1001'/0'/0"
)

// Paths:
//   m/<env=[0,1,...]>'/<purpose>/<role>/<idx>
//
// Key purpose (prefixes):
//   eth:       m/0'/0
//   libp2p:    m/0'/1
//   caps:      m/0'/2
//   onion:     m/0'/3
//
// Node roles:
//   eth:     0
//   boot:    1
//   feed:    2
//   feed_lb: 3
//   bb:      4
//   relay:   5
//   spectre: 6
//   ghost:   7
//   monitor: 8
//   lair:    9

func NewDerive(opts *Options) *cobra.Command {
	var prefix, password, format string
	var lineNum int
	cmd := &cobra.Command{
		Use:     "derive [--prefix path] [--format " + strings.Join(FormatList, "|") + "] [--password] path...",
		Aliases: []string{"der", "d"},
		Short:   "Derive a key pair from the provided mnemonic phrase.",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = []string{"0"}
			}
			if lineNum < 1 {
				return fmt.Errorf("line number must be greater than 0")
			}
			mnemonic, err := lineFromFile(opts.InputFile, lineNum-1)
			if err != nil {
				return err
			}
			wallet, err := hdwallet.NewFromMnemonic(mnemonic)
			if err != nil {
				return err
			}
			if prefix != "" && !strings.HasSuffix(prefix, "/") {
				prefix += "/"
			}
			for _, arg := range args {
				dp, err := accounts.ParseDerivationPath(prefix + arg)
				if err != nil {
					return err
				}
				log.Println(dp.String())
				acc, err := wallet.Derive(dp, false)
				if err != nil {
					return err
				}
				log.Println(acc.Address.String())
				privateKey, err := wallet.PrivateKey(acc)
				if err != nil {
					return err
				}
				b, err := formattedBytes(format, privateKey, password)
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(
		&lineNum,
		"line",
		"l",
		1,
		"which seed line to take from the input file",
	)
	cmd.Flags().StringVarP(
		&prefix,
		"prefix",
		"p",
		"",
		"derivation path prefix",
	)
	cmd.Flags().StringVarP(
		&password,
		"password",
		"w",
		"",
		"encryption password",
	)
	cmd.Flags().StringVarP(
		&format,
		"format",
		"f",
		FormatEth,
		"output format",
	)
	return cmd
}

const (
	FormatEth        = "eth"
	FormatEthPlain   = "eth-plain"
	FormatEthStatic  = "eth-static"
	FormatLibP2P     = "libp2p"
	FormatSSB        = "ssb"
	FormatSSBCaps    = "caps"
	FormatSSBSHS     = "shs"
	FormatBytes32    = "b32"
	FormatPrivHex    = "privhex"
	FormatOnionV3    = "onion"
	FormatOnionV3Adr = "onion-adr"
	FormatOnionV3Pub = "onion-pub"
	FormatOnionV3Sec = "onion-sec"
)

var FormatList = []string{
	FormatEth,
	FormatEthPlain,
	FormatEthStatic,
	FormatLibP2P,
	FormatSSB,
	FormatSSBCaps,
	FormatSSBSHS,
	FormatBytes32,
	FormatPrivHex,
	FormatOnionV3,
	FormatOnionV3Adr,
	FormatOnionV3Pub,
	FormatOnionV3Sec,
}

func formattedBytes(format string, privateKey *ecdsa.PrivateKey, password string) ([]byte, error) {
	switch format {
	case FormatLibP2P:
		randBytes, err := seededRandBytesFunc(privateKey, 32)
		if err != nil {
			return nil, err
		}
		return hexEncodeBytes(randBytes()), nil
	case FormatBytes32, FormatSSBSHS:
		randBytes, err := seededRandBytesFunc(privateKey, 32)
		if err != nil {
			return nil, err
		}
		return b64Encode(randBytes()), nil
	case FormatSSB:
		o, err := ssb.NewSecret(crypto.FromECDSA(privateKey))
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	case FormatSSBCaps:
		o, err := ssb.NewCaps(crypto.FromECDSA(privateKey))
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	case FormatEth:
		o, err := eth.NewKeyWithID(privateKey)
		if err != nil {
			return nil, err
		}
		return keystore.EncryptKey(
			o,
			password,
			keystore.StandardScryptN,
			keystore.StandardScryptP,
		)
	case FormatEthStatic:
		o, err := eth.NewKeyWithID(privateKey)
		if err != nil {
			return nil, err
		}
		t := rand2.Reader
		defer func() { rand2.Reader = t }()
		bytesFunc, err := seededRandBytesFunc(privateKey, 48)
		if err != nil {
			return nil, err
		}
		rand2.Reader = bytes.NewReader(bytesFunc())
		return keystore.EncryptKey(
			o,
			password,
			keystore.StandardScryptN,
			keystore.StandardScryptP,
		)
	case FormatEthPlain:
		o, err := eth.NewKeyWithID(privateKey)
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	case FormatPrivHex:
		return hexEncodeBytes(crypto.FromECDSA(privateKey)), nil
	case FormatOnionV3, FormatOnionV3Adr, FormatOnionV3Pub, FormatOnionV3Sec:
		o, err := tor.NewOnion(crypto.FromECDSA(privateKey))
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	}
	return nil, fmt.Errorf("unknown format: %s", format)
}
func hexEncodeBytes(b []byte) []byte {
	buff := make([]byte, len(b)*2)
	hex.Encode(buff, b)
	return buff
}
func b64Encode(b []byte) []byte {
	enc := base64.StdEncoding
	buff := make([]byte, enc.EncodedLen(len(b)))
	enc.Encode(buff, b)
	return buff
}
func seededRandBytesFunc(privateKey *ecdsa.PrivateKey, len int) (func() []byte, error) {
	return rand.SeededRandBytesGen(crypto.FromECDSA(privateKey), len)
}
