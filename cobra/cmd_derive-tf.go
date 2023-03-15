package cobra

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/spf13/cobra"
	refs "github.com/ssbc/go-ssb-refs"
)

func NewDeriveTf() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "derive-tf",
		Aliases: []string{"dtf"},
		Short:   "Derive keys from HD mnemonic (terraform external data style).",
		RunE: func(_ *cobra.Command, _ []string) error {
			var q query
			err := json.NewDecoder(os.Stdin).Decode(&q)
			if err != nil {
				return fmt.Errorf("input decoding failed: %w", err)
			}
			wallet, err := hdwallet.NewFromMnemonic(q.Mnemonic)
			if err != nil {
				return err
			}
			dp, err := accounts.ParseDerivationPath(q.Path)
			if err != nil {
				return err
			}
			acc, err := wallet.Derive(dp, false)
			if err != nil {
				return err
			}
			privateKey, err := wallet.PrivateKey(acc)
			if err != nil {
				return err
			}
			b, err := formattedBytes(q.Format, privateKey, q.Password)
			if err != nil {
				return err
			}
			addr := acc.Address.String()
			if q.Format == FormatSSB {
				var ssb ssbSecret
				if err := json.Unmarshal(b, &ssb); err != nil {
					return err
				}
				addr = ssb.ID.String()
			} else if q.Format == FormatLibP2P {
				seed := make([]byte, hex.DecodedLen(len(b)))
				_, err := hex.Decode(seed, b)
				if err != nil {
					return err
				}
				if len(seed) != ed25519.SeedSize {
					return fmt.Errorf("invalid privKeySeed value, 32 bytes expected")
				}
				seedReader := bytes.NewReader(seed)
				_, pub, err := crypto.GenerateEd25519Key(seedReader)
				if err != nil {
					return err
				}
				id, err := peer.IDFromPublicKey(pub)
				if err != nil {
					return err
				}
				addr = id.String()
			} else if strings.HasPrefix(q.Format, FormatOnionV3+"-") {
				var o onion
				if err := json.Unmarshal(b, &o); err != nil {
					return err
				}
				addr = o.Hostname
				switch q.Format {
				case FormatOnionV3Adr:
					b = []byte(addr)
				case FormatOnionV3Pub:
					b = o.PublicKey
				case FormatOnionV3Sec:
					b = o.SecretKey
				}
			}
			fmt.Printf(
				`{"output":"%s","path":"%s","addr":"%s"}`,
				base64.StdEncoding.EncodeToString(b),
				dp.String(),
				addr,
			)
			return nil
		},
	}
	return cmd
}

type query struct {
	Mnemonic string `json:"mnemonic"`
	Path     string `json:"path"`
	Password string `json:"password"`
	Format   string `json:"format"`
}

type ssbSecret struct {
	Curve   string       `json:"curve"`
	ID      refs.FeedRef `json:"id"`
	Private string       `json:"private"`
	Public  string       `json:"public"`
}

type onion struct {
	Prefix    string `json:"prefix"`
	Hostname  string `json:"hostname"`
	PublicKey []byte `json:"public_key"`
	SecretKey []byte `json:"secret_key"`
}
