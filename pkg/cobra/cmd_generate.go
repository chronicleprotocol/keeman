package cobra

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"

	"github.com/chronicleprotocol/keeman/pkg/hdwallet/wl"
)

const (
	bitBlockSize = 32

	genFormatMnemonic = "mnemonic"
)

func GenerateSeed(opts *Options) *cobra.Command {
	var format, lang string
	var bitSizeMultiplier, bits int
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "g"},
		Args:    cobra.NoArgs,
		Short:   "Generate HD seed phrase with a specific bit size",
		RunE: func(_ *cobra.Command, args []string) error {
			if bits > 0 {
				bitSizeMultiplier = bits / bitBlockSize
				if bitBlockSize*bitSizeMultiplier != bits {
					return fmt.Errorf("entropy size must be a multiple of %d", bitBlockSize)
				}
			}
			if bitSizeMultiplier < 4 || bitSizeMultiplier > 8 {
				return fmt.Errorf("entropy size multiplier must be between 4 and 8")
			}

			bitSize := bitBlockSize * bitSizeMultiplier
			log.Printf("entropy bit size: %d * %d = %d\n", bitSizeMultiplier, bitBlockSize, bitSize)

			entropy, err := bip39.NewEntropy(bitSize)
			if err != nil {
				return err
			}
			switch lang {
			case "en":
				bip39.SetWordList(wordlists.English)
			case "es":
				bip39.SetWordList(wordlists.Spanish)
			case "fr":
				bip39.SetWordList(wordlists.French)
			case "it":
				bip39.SetWordList(wordlists.Italian)
			case "ja":
				bip39.SetWordList(wordlists.Japanese)
			case "ko":
				bip39.SetWordList(wordlists.Korean)
			case "zh":
				bip39.SetWordList(wordlists.ChineseSimplified)
			case "zh-tw":
				bip39.SetWordList(wordlists.ChineseTraditional)
			case "cs":
				bip39.SetWordList(wordlists.Czech)
			case "pl":
				bip39.SetWordList(wl.Polish)
			default:
				return fmt.Errorf("unsupported language: %s", lang)
			}
			switch format {
			case formatEth:
			case FormatEthStatic:
			case formatEthPlain:
			case genFormatMnemonic:
				mnemonic, err := bip39.NewMnemonic(entropy)
				if err != nil {
					return err
				}
				fmt.Println(mnemonic)
				if opts.Verbose {
					log.Println(mnemonic)
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(
		&bitSizeMultiplier,
		"multiplier",
		"k",
		4,
		"number of 32 bit size blocks for entropy <4;8> (ignored when --bits is used)",
	)
	cmd.Flags().IntVarP(
		&bits,
		"bits",
		"b",
		0,
		"number of bits of entropy <128;256> (has priority over --multiplier)",
	)
	cmd.Flags().StringVarP(
		&format,
		"format",
		"f",
		genFormatMnemonic,
		"output format",
	)
	cmd.Flags().StringVarP(
		&lang,
		"lang",
		"l",
		"en",
		"word list language",
	)
	return cmd
}
