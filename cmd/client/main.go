/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "client"
	app.Usage = "interact with node process"
	app.HideVersion = true

	minerCommands := cli.Command{
		Name:  "miner",
		Usage: "miner commands",
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "start miner",
				Flags:  rpcFlags(threadsFlag),
				Action: rpcAction("miner", "start"),
			},
			{
				Name:   "stop",
				Usage:  "stop miner",
				Flags:  rpcFlags(),
				Action: rpcAction("miner", "stop"),
			},
			{
				Name:   "hashrate",
				Usage:  "get miner hashrate",
				Flags:  rpcFlags(),
				Action: rpcAction("miner", "hashrate"),
			},
			{
				Name:   "getthreads",
				Usage:  "get miner thread number",
				Flags:  rpcFlags(),
				Action: rpcAction("miner", "getThreads"),
			},
			{
				Name:   "setthreads",
				Usage:  "set miner thread number",
				Flags:  rpcFlags(threadsFlag),
				Action: rpcAction("miner", "setThreads"),
			},
			{
				Name:   "setcoinbase",
				Usage:  "set miner coinbase",
				Flags:  rpcFlags(coinbaseFlag),
				Action: rpcAction("miner", "setCoinbase"),
			},
			{
				Name:   "getcoinbase",
				Usage:  "get miner coinbase",
				Flags:  rpcFlags(),
				Action: rpcAction("miner", "getCoinbase"),
			},
			{
				Name:   "status",
				Usage:  "get miner status",
				Flags:  rpcFlags(),
				Action: rpcAction("miner", "status"),
			},
		},
	}

	htlcCommands := cli.Command{
		Name:  "htlc",
		Usage: "Hash time lock contract commands",
		Subcommands: []cli.Command{
			{
				Name:   "create",
				Usage:  "create HTLC",
				Flags:  rpcFlags(fromFlag, toFlag, amountFlag, feeFlag, nonceFlag, hashFlag, timeLockFlag),
				Action: rpcActionSystemContract("htlc", "create", handleCallResult),
			},
			{
				Name:   "withdraw",
				Usage:  "withdraw from HTLC",
				Flags:  rpcFlags(fromFlag, feeFlag, nonceFlag, hashFlag, preimageFlag),
				Action: rpcActionSystemContract("htlc", "withdraw", handleCallResult),
			},
			{
				Name:   "refund",
				Usage:  "refund from HTLC",
				Flags:  rpcFlags(fromFlag, feeFlag, nonceFlag, hashFlag),
				Action: rpcActionSystemContract("htlc", "refund", handleCallResult),
			},
			{
				Name:   "get",
				Usage:  "get HTLC information",
				Flags:  rpcFlags(fromFlag, feeFlag, nonceFlag, hashFlag),
				Action: rpcActionSystemContract("htlc", "get", handleCallResult),
			},
			{
				Name:  "decode",
				Usage: "decode HTLC contract information",
				Flags: []cli.Flag{
					payloadFlag,
				},
				Action: decodeHTLC,
			},
			{
				Name:   "key",
				Usage:  "generate preimage key and key hash",
				Action: generateHTLCKey,
			},
			{
				Name:  "time",
				Usage: "generate unix timestamp",
				Flags: []cli.Flag{
					timeLockFlag,
				},
				Action: generateHTLCTime,
			},
		},
	}

	domainCommands := cli.Command{
		Name:  "domain",
		Usage: "system domain name commands",
		Subcommands: []cli.Command{
			{
				Name:   "register",
				Usage:  "register a domain name",
				Flags:  rpcFlags(fromFlag, feeFlag, nameFlag, nonceFlag),
				Action: rpcActionSystemContract("domain", "create", handleCallResult),
			},
			{
				Name:   "owner",
				Usage:  "get the domain name owner",
				Flags:  rpcFlags(fromFlag, feeFlag, nameFlag, nonceFlag),
				Action: rpcActionSystemContract("domain", "getOwner", handleCallResult),
			},
		},
	}

	subChainCommands := cli.Command{
		Name:  "subchain",
		Usage: "system sub chain commands",
		Subcommands: []cli.Command{
			{
				Name:   "register",
				Usage:  "register a sub chain",
				Flags:  rpcFlags(fromFlag, feeFlag, nonceFlag, subChainJSONFileFlag),
				Action: rpcActionSystemContract("subchain", "register", handleCallResult),
			},
			{
				Name:   "query",
				Usage:  "query sub chain",
				Flags:  rpcFlags(fromFlag, feeFlag, nonceFlag, nameFlag),
				Action: rpcActionSystemContract("subchain", "query", handleCallResult),
			},
			{
				Name:   "config",
				Usage:  "generate sub chain config file",
				Flags:  rpcFlags(coinbaseFlag, privateKeyFlag, nameFlag, outPutFlag, shardFlag, staticNodesFlag),
				Action: createSubChainConfigFile,
			},
		},
	}

	p2pCommands := cli.Command{
		Name:  "p2p",
		Usage: "p2p commands",
		Subcommands: []cli.Command{
			{
				Name:   "peers",
				Usage:  "get p2p peer connections",
				Flags:  rpcFlags(),
				Action: rpcAction("network", "getPeerCount"),
			},
			{
				Name:   "peersinfo",
				Usage:  "get p2p peers information",
				Flags:  rpcFlags(),
				Action: rpcAction("network", "getPeersInfo"),
			},
			{
				Name:   "networkversion",
				Usage:  "get current network version",
				Flags:  rpcFlags(),
				Action: rpcAction("network", "getNetworkVersion"),
			},
			{
				Name:   "protocolversion",
				Usage:  "get seele protocol version",
				Flags:  rpcFlags(),
				Action: rpcAction("network", "getProtocolVersion"),
			},
		},
	}

	app.Commands = []cli.Command{
		htlcCommands,
		minerCommands,
		p2pCommands,
		domainCommands,
		subChainCommands,
		{
			Name:   "getinfo",
			Usage:  "get node info",
			Flags:  rpcFlags(),
			Action: rpcAction("seele", "getInfo"),
		},
		{
			Name:   "getbalance",
			Usage:  "get balance info",
			Flags:  rpcFlags(accountFlag),
			Action: rpcAction("seele", "getBalance"),
		},
		{
			Name:   "sendtx",
			Usage:  "send transaction to node",
			Flags:  rpcFlags(fromFlag, toFlag, amountFlag, feeFlag, payloadFlag, nonceFlag),
			Action: rpcActionEx("seele", "addTx", makeTransaction, onTxAdded),
		},
		{
			Name:   "getnonce",
			Usage:  "get account nonce",
			Flags:  rpcFlags(accountFlag),
			Action: rpcAction("seele", "getAccountNonce"),
		},
		{
			Name:   "call",
			Usage:  "call contract",
			Flags:  rpcFlags(toFlag, payloadFlag, heightFlag),
			Action: rpcAction("seele", "call"),
		},
		{
			Name:   "getblockheight",
			Usage:  "get block height",
			Flags:  rpcFlags(),
			Action: rpcAction("seele", "getBlockHeight"),
		},
		{
			Name:   "getblock",
			Usage:  "get block by height or hash",
			Flags:  rpcFlags(hashFlag, heightFlag, fulltxFlag),
			Action: rpcAction("seele", "getBlock"),
		},
		{
			Name:   "getlogs",
			Usage:  "get logs",
			Flags:  rpcFlags(heightFlag, contractFlag, topicFlag),
			Action: rpcAction("seele", "getLogs"),
		},
		{
			Name:   "gettxpoolcontent",
			Usage:  "get transaction pool contents",
			Flags:  rpcFlags(),
			Action: rpcAction("debug", "getTxPoolContent"),
		},
		{
			Name:   "gettxpoolcount",
			Usage:  "get transaction pool transaction count",
			Flags:  rpcFlags(),
			Action: rpcAction("debug", "getTxPoolTxCount"),
		},
		{
			Name:   "getblocktxcount",
			Usage:  "get block transaction count by block height or block hash",
			Flags:  rpcFlags(hashFlag, heightFlag),
			Action: rpcAction("txpool", "getBlockTransactionCount"),
		},
		{
			Name:   "gettxinblock",
			Usage:  "get transaction by block height or block hash with index of the transaction in the block",
			Flags:  rpcFlags(hashFlag, heightFlag, indexFlag),
			Action: rpcAction("txpool", "getTransactionByBlockIndex"),
		},
		{
			Name:   "gettxbyhash",
			Usage:  "get transaction by transaction hash",
			Flags:  rpcFlags(hashFlag),
			Action: rpcAction("txpool", "getTransactionByHash"),
		},
		{
			Name:   "getdebtbyhash",
			Usage:  "get debt by debt hash",
			Flags:  rpcFlags(hashFlag),
			Action: rpcAction("txpool", "getDebtByHash"),
		},
		{
			Name:   "getdebts",
			Usage:  "get pending debts",
			Flags:  rpcFlags(),
			Action: rpcAction("debug", "getPendingDebts"),
		},
		{
			Name:   "getreceipt",
			Usage:  "get receipt by transaction hash",
			Flags:  rpcFlags(hashFlag),
			Action: rpcAction("txpool", "getReceiptByTxHash"),
		},
		{
			Name:   "getpendingtxs",
			Usage:  "get pending transactions",
			Flags:  rpcFlags(),
			Action: rpcAction("debug", "getPendingTransactions"),
		},
		{
			Name:  "getshardnum",
			Usage: "get account shard number",
			Flags: []cli.Flag{
				accountFlag,
				privateKeyFlag,
			},
			Action: GetAccountShardNumAction,
		},
		{
			Name:  "savekey",
			Usage: "save private key to a keystore file",
			Flags: []cli.Flag{
				privateKeyFlag,
				fileNameFlag,
			},
			Action: SaveKeyAction,
		},
		{
			Name:  "sign",
			Usage: "generate a signed transaction and print it out",
			Flags: []cli.Flag{
				addressFlag,
				privateKeyFlag,
				toFlag,
				amountFlag,
				feeFlag,
				payloadFlag,
				nonceFlag,
			},
			Action: SignTxAction,
		},
		{
			Name:  "key",
			Usage: "generate key with or without shard number",
			Flags: []cli.Flag{
				shardFlag,
			},
			Action: GenerateKeyAction,
		},
		{
			Name:   "dumpheap",
			Usage:  "dump heap for profiling, return the file path",
			Flags:  rpcFlags(dumpFileFlag, gcBeforeDumpFlag),
			Action: rpcAction("debug", "dumpHeap"),
		},
	}

	// sort commands and flags by name
	sortCommands(app.Commands)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func sortCommands(commands []cli.Command) {
	sort.Sort(cli.CommandsByName(commands))

	for _, command := range commands {
		if len(command.Subcommands) > 0 {
			sortCommands(command.Subcommands)
		}

		sort.Sort(cli.FlagsByName(command.Flags))
	}
}
