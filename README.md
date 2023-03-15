# Keeman

```text
Usage:
  keeman [command]

Available Commands:
  derive      Derive a key pair from the provided mnemonic phrase
  derive-tf   Derive keys from HD mnemonic (terraform external data style).
  generate    Generate HD seed phrase with a specific bit size
  help        Help about any command
  list        List word count and first word from the input, omitting the comments

Flags:
  -h, --help            help
  -i, --input string    input file path
  -o, --output string   output file path
  -v, --verbose         verbose logging
```

### Examples

Generate a mnemonic phrase with 256 bits of entropy:
```shell
keeman generate -b 256
```

Derive a key pair from the provided mnemonic phrase:
```shell
echo "your mnemonic phrase" | keeman derive
```

Generate the mnemonic phrase and derive four keys in one command (and show the mnemonic phrase on stderr):
```shell
keeman gen | tee >(cat >&2) | keeman der 0 1 2 3
```

## Mnemonic Phrase Generation
Generate HD seed phrase with a specific bit size.

```text
Usage:
  keeman generate [flags]

Aliases:
  generate, gen, g

Flags:
  -b, --bits int         number of bits of entropy <128;256>
  -k, --multiplier int   number of 32 bit size blocks for entropy <4;8> (default 4)

Global Flags:
  -h, --help           help
  -i, --input string   input file path
  -v, --verbose        verbose logging
```

## Key Derivation
Derive a key pair from the provided mnemonic phrase.

```text
Usage:
  keeman derive [--prefix path] [--format eth|ssb|caps|shs|b32|privhex|libp2p|onion|onion-adr|onion-pub|onion-sec] [--password] path... [flags]

Aliases:
  derive, der, d

Flags:
  -f, --format string     output format (default "eth")
  -l, --line int          which seed line to take from the input file (default 1)
  -w, --password string   encryption password
  -p, --prefix string     derivation path prefix

Global Flags:
  -h, --help           help
  -i, --input string   input file path
  -v, --verbose        verbose logging
```