package cobra

import (
	"bytes"
	"crypto/ecdsa"
	rand2 "crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/chronicleprotocol/keeman/eth"
	"github.com/chronicleprotocol/keeman/hdwallet"
	"github.com/chronicleprotocol/keeman/rand"
	"github.com/chronicleprotocol/keeman/ssb"
	"github.com/chronicleprotocol/keeman/tor"
)

var FormatList = []string{
	formatSec,
	formatPub,
	formatAddr,
	formatEth,
	FormatEthStatic,
	formatEthPlain,
	formatSSB,
	formatCaps,
	formatOnionV3,
}

const (
	formatSec       = "sec"
	formatPub       = "pub"
	formatAddr      = "addr"
	formatEth       = "eth"
	FormatEthStatic = "eth-static" // Will generate the same JSON keystore file
	formatEthPlain  = "eth-plain"
	formatSSB       = "ssb"
	formatCaps      = "caps"
	formatOnionV3   = "onion"
)

const (
	iteratorEthereum = "eth"
	iteratorMetaMask = "mm"
	iteratorLedger   = "ll"
)

func NewDerive(opts *Options) *cobra.Command {
	var prefix, suffix, password, format, iterator, encoding string
	var lineNum, num int
	cmd := &cobra.Command{
		Use:     "derive [--prefix path] [--suffix path] [--format " + strings.Join(FormatList, "|") + "] [--password] path...",
		Aliases: []string{"der", "d"},
		Short:   "Derive values from the provided mnemonic phrase",
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 && num > 0 {
				return errors.New("cannot use --num with positional arguments")
			} else if len(args) == 0 && num < 1 {
				num = 1
			}
			// Generate the specified number of arguments to iterate over.
			for i := 0; i < num; i++ {
				args = append(args, fmt.Sprintf("%d", i))
			}

			if lineNum < 1 {
				return errors.New("line number must be greater than 0")
			}
			// Read mnemonic from <lineNum> line of the provided text file. First line number is 1.
			mnemonic, err := lineFromFile(opts.InputFile, lineNum-1)
			if err != nil {
				return err
			}
			// Derive key pair from the mnemonic.
			wallet, err := hdwallet.NewFromMnemonic(mnemonic)
			if err != nil {
				return err
			}

			if iterator != "" && (prefix != "" || suffix != "") {
				return errors.New("--iterator cannot be used with --prefix or --suffix")
			} else if iterator+prefix+suffix == "" {
				iterator = iteratorEthereum
			}
			// Use iterator name to override prefix and suffix.
			prefix, suffix, format, encoding, err = applyIterator(iterator, prefix, suffix, format, encoding)
			if err != nil {
				return err
			}

			// Iterate over the provided paths and derive the key pairs.
			for _, arg := range args {
				dp, err := combineDerivationPath(prefix, arg, suffix)
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
				var e encoder
				switch encoding {
				case "base64", "b64":
					e = base64.StdEncoding
				case "base64url", "b64u":
					e = base64.URLEncoding
				case "base32", "b32":
					e = base32.StdEncoding
				case "base32hex", "b32h":
					e = base32.HexEncoding
				case "hex":
					e = &hexEncoder{}
				case "":
					e = &plainEncoder{}
				}
				fmt.Println(e.EncodeToString(b))
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
		&suffix,
		"suffix",
		"s",
		"",
		"derivation path suffix",
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
		formatEth,
		"output format",
	)
	cmd.Flags().StringVarP(
		&iterator,
		"iterator",
		"i",
		"",
		"which iterator to use",
	)
	cmd.Flags().IntVarP(
		&num,
		"num",
		"n",
		0,
		"how many addresses to generate (in addition to positional arguments)",
	)
	cmd.Flags().StringVarP(
		&encoding,
		"encode",
		"e",
		"",
		"how many addresses to generate (in addition to positional arguments)",
	)
	return cmd
}

func applyIterator(iterator, prefix, suffix, format, encoding string) (string, string, string, string, error) {
	switch iterator {
	case iteratorEthereum, iteratorMetaMask:
		p, ok := hdwallet.PrefixList["Ethereum"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
	case iteratorLedger:
		p, ok := hdwallet.PrefixList["EthereumLedgerLive"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = "'/0/0"
	case "f":
		p, ok := hdwallet.PrefixList["Feed"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
	case "fo":
		p, ok := hdwallet.PrefixList["FeedOnion"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
		format = formatOnionV3
	case "r":
		p, ok := hdwallet.PrefixList["Relay"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
	case "ro":
		p, ok := hdwallet.PrefixList["RelayOnion"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
		format = formatOnionV3
	case "m":
		p, ok := hdwallet.PrefixList["Monitor"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
	case "mo":
		p, ok := hdwallet.PrefixList["MonitorOnion"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
		format = formatOnionV3
	case "bl":
		p, ok := hdwallet.PrefixList["BootstrapLibP2P"]
		if !ok {
			return "", "", "", "", fmt.Errorf("prefix %q not found", prefix)
		}
		prefix = p
		suffix = ""
		format = formatSec
		encoding = "hex"
	default:
		if iterator != "" {
			return "", "", "", "", fmt.Errorf("iterator %q not found", iterator)
		}
	}
	return prefix, suffix, format, encoding, nil
}

func formattedBytes(format string, privateKey *ecdsa.PrivateKey, password string) ([]byte, error) {
	switch format {
	case formatPub:
		return crypto.FromECDSAPub(&privateKey.PublicKey), nil
	case formatSec:
		return crypto.FromECDSA(privateKey), nil
	case formatAddr:
		return []byte(crypto.PubkeyToAddress(privateKey.PublicKey).String()), nil
	case formatEth, FormatEthStatic:
		if format == FormatEthStatic {
			defer func(r io.Reader) { rand2.Reader = r }(rand2.Reader)
			bytesFunc, err := rand.SeededRandBytesGen(crypto.FromECDSA(privateKey), 64)
			if err != nil {
				return nil, err
			}
			rand2.Reader = bytes.NewReader(bytesFunc())
		}
		k, err := eth.NewKeyWithID(privateKey)
		if err != nil {
			return nil, err
		}
		return keystore.EncryptKey(
			k,
			password,
			keystore.StandardScryptN,
			keystore.StandardScryptP,
		)
	case formatEthPlain:
		k, err := eth.NewKeyWithID(privateKey)
		if err != nil {
			return nil, err
		}
		return json.Marshal(k)
	case formatSSB:
		o, err := ssb.NewSecret(crypto.FromECDSA(privateKey))
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	case formatCaps:
		o, err := ssb.NewCaps(crypto.FromECDSA(privateKey))
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	case formatOnionV3:
		o, err := tor.NewOnion(crypto.FromECDSA(privateKey))
		if err != nil {
			return nil, err
		}
		return json.Marshal(o)
	}
	return nil, fmt.Errorf("unknown format: %s", format)
}

type encoder interface {
	EncodeToString(src []byte) string
}

type hexEncoder struct{}

func (hexEncoder) EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

type plainEncoder struct{}

func (plainEncoder) EncodeToString(src []byte) string {
	return string(src)
}

func combineDerivationPath(prefix string, arg string, suffix string) (accounts.DerivationPath, error) {
	prefix = strings.TrimPrefix(prefix, "m/")
	prefix = strings.Trim(prefix, "/")
	if prefix != "" {
		prefix += "/"
	}
	suffix = strings.Trim(suffix, "/")
	if suffix != "" && !strings.HasPrefix(suffix, "'") {
		suffix = "/" + suffix
	}
	p := "m/" + prefix + strings.Trim(arg, "/") + suffix
	path, err := accounts.ParseDerivationPath(p)
	if err != nil {
		return nil, fmt.Errorf(" %s is an invalid derivation path: %w", p, err)
	}
	return path, err
}
