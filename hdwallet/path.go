package hdwallet

// Paths:
//   m/<env=[0,1,...]>'/<purpose>/<role>[/<idx>]
//
// Key purpose (prefixes):
//   eth:       m/0'/0
//   libp2p:    m/0'/1
//   caps:      m/0'/2 - no idx
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

var PrefixList = map[string]string{
	"Feed":            "m/0'/0/2",
	"FeedOnion":       "m/0'/3/2",
	"Relay":           "m/0'/0/5",
	"RelayOnion":      "m/0'/3/5",
	"BootstrapLibP2P": "m/0'/1/1",
	"Monitor":         "m/0'/0/8",
	"MonitorOnion":    "m/0'/3/8",

	"Akroma":                    "m/44'/200625'/0'/0",
	"Atheios":                   "m/44'/1620'/0'/0",
	"Callisto":                  "m/44'/820'/0'/0",
	"Ellaism":                   "m/44'/163'/0'/0",
	"EOSClassic":                "m/44'/2018'/0'/0",
	"Ether":                     "m/44'/1313114'/0'/0",
	"Ethereum":                  "m/44'/60'/0'/0",
	"EthereumClassic":           "m/44'/61'/0'/0",
	"EthereumClassicLedger":     "m/44'/60'/160720'/0",
	"EthereumClassicLedgerLive": "m/44'/61'",
	"EthereumLedger":            "m/44'/60'/0'",
	"EthereumLedgerLive":        "m/44'/60'",
	"EthereumSocial":            "m/44'/1128'/0'/0",
	"EthereumTestnetRopsten":    "m/44'/1'/0'/0",
	"EtherGem":                  "m/44'/1987'/0'/0",
	"EtherSocialNetwork":        "m/44'/31102'/0'/0",
	"Expanse":                   "m/44'/40'/0'/0",
	"GoChain":                   "m/44'/6060'/0'/0",
	"Iolite":                    "m/44'/1171337'/0'/0",
	"MetaMask":                  "m/44'/60'/0'/0",
	"MixBlockchain":             "m/44'/76'/0'/0",
	"Musicoin":                  "m/44'/184'/0'/0",
	"PIRL":                      "m/44'/164'/0'/0",
	"RSKMainnet":                "m/44'/137'/0'/0",
	"ThunderCore":               "m/44'/1001'/0'/0",
	"TomoChain":                 "m/44'/889'/0'/0",
	"Ubiq":                      "m/44'/108'/0'/0",
}
