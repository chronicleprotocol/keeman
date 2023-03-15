package cobra

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

const bitBlockSize = 32

func GenerateSeed(opts *Options) *cobra.Command {
	var bitSizeMultiplier, bits int
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "g"},
		Args:    cobra.NoArgs,
		Short:   "Generate HD seed phrase with a specific bit size.",
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
			mnemonic, err := bip39.NewMnemonic(entropy)
			if err != nil {
				return err
			}
			fmt.Println(mnemonic)
			if opts.Verbose {
				log.Println(mnemonic)
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
	return cmd
}
