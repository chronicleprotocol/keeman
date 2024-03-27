# Keeman

## Installation

```shell
go install github.com/chronicleprotocol/keeman@latest
```

## Examples

### Mnemonic Phrase Generation

Generate a mnemonic phrase with 256 bits of entropy:

```shell
keeman generate -b 256
```

### Key Derivation

Derive a key pair from the provided mnemonic phrase:
```shell
echo "your mnemonic phrase" | keeman derive
```

Derive an Ethereum Key and encrypt with a static salt (creates reproducible keystore JSONs):
```shell
echo "your mnemonic phrase" | keeman derive --format eth
```

Create a plain text Ethereum key:
```shell
echo "your mnemonic phrase" | keeman derive --format eth-plain
```

Generate and save a service keys to files usable by tor
```shell
MNEMONIC=$(keeman generate)
keeman derive --format onion <<<"$MNEMONIC" | jq -r .hostname > hostname
keeman derive --format onion <<<"$MNEMONIC" | jq -r .public_key | base64 -d > hs_ed25519_public_key
keeman derive --format onion <<<"$MNEMONIC" | jq -r .secret_key | base64 -d > hs_ed25519_secret_key
```

Generate the mnemonic phrase and derive four keys in one command (and show the mnemonic phrase on stderr):
```shell
keeman gen | tee >(cat >&2) | keeman der 0 1 2 3
```

## Help Page

```text
Usage:
  keeman [command]

Available Commands:
  derive      Derive values from the provided mnemonic phrase
  generate    Generate HD seed phrase with a specific bit size
  help        Help about any command
  list        List word count and first word from the input, omitting the comments

Flags:
  -h, --help           help for keeman
  -i, --input string   input file path
  -v, --verbose        verbose logging

Use "keeman [command] --help" for more information about a command.
```

## Help Page for `generate` Command

```text
Generate HD seed phrase with a specific bit size

Usage:
  keeman generate [flags]

Aliases:
  generate, gen, g

Flags:
  -b, --bits int         number of bits of entropy <128;256> (has priority over --multiplier)
  -f, --format string    output format (default "mnemonic")
  -h, --help             help for generate
  -l, --lang string      word list language (default "en")
  -k, --multiplier int   number of 32 bit size blocks for entropy <4;8> (ignored when --bits is used) (default 4)

Global Flags:
  -i, --input string   input file path
  -v, --verbose        verbose logging
```

## Help Page for `derive` Command

```text
Derive values from the provided mnemonic phrase

Usage:
  keeman derive [--prefix path] [--suffix path] [--format sec|pub|addr|eth|eth-static|eth-plain|ssb|caps|onion] [--password] path... [flags]

Aliases:
  derive, der, d

Flags:
  -e, --encode string     encoding to use
  -f, --format string     output format (default "eth")
  -h, --help              help for derive
  -t, --iterator string   which iterator to use
  -l, --line int          which seed line to take from the input file (default 1)
  -n, --num int           how many addresses to generate (in addition to positional arguments)
  -w, --password string   encryption password
  -p, --prefix string     derivation path prefix
  -s, --suffix string     derivation path suffix

Global Flags:
  -i, --input string   input file path
  -v, --verbose        verbose logging
```

## Help Page for `list` Command

```text
List word count and first word from the input, omitting the comments

Usage:
  keeman list [--all] [flags]

Aliases:
  list, l

Flags:
  -a, --all         all data
  -h, --help        help for list
      --index int   data index

Global Flags:
  -i, --input string   input file path
  -v, --verbose        verbose logging
```
