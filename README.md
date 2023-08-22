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
  generate    Generate HD seed phrase with a specific bit size.
  help        Help about any command
  list        List word count and first word from the input, omitting the comments.

Flags:
  -h, --help           help for keeman
  -c, --input string   input file path
  -v, --verbose        verbose logging

Use "keeman [command] --help" for more information about a command.
```
