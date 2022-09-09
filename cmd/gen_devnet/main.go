package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/node/app"
	"github.com/bnb-chain/node/app/config"
	bnbInit "github.com/bnb-chain/node/cmd/bnbchaind/init"
	"github.com/bnb-chain/node/common/utils"
)

var (
	chainID = "Binance-Chain-Kongo"
	nodeNum = 1
)

func main() {
	fmt.Println("start generate devNet configs")
	cwd, _ := os.Getwd()
	devNetHomeDir := path.Join(cwd, "build", "devNet")
	fmt.Println("devNet home dir:", devNetHomeDir)
	// clear devNetHomeDir
	err := os.RemoveAll(devNetHomeDir)
	if err != nil {
		panic(err)
	}
	// init nodes
	cdc := app.Codec
	ctx := app.ServerContext
	appInit := app.BinanceAppInit()
	ctxConfig := ctx.Config
	sdkConfig := sdk.GetConfig()
	ctx.Bech32PrefixAccAddr = "tbnb"
	ctx.Bech32PrefixAccPub = "tbnbp"
	sdkConfig.SetBech32PrefixForAccount(ctx.Bech32PrefixAccAddr, ctx.Bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(ctx.Bech32PrefixValAddr, ctx.Bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(ctx.Bech32PrefixConsAddr, ctx.Bech32PrefixConsPub)
	sdkConfig.Seal()
	var appState json.RawMessage
	var seeds string
	genesisTime := utils.Now()
	ServerContext := config.NewDefaultContext()
	ServerContext.Bech32PrefixAccAddr = ctx.Bech32PrefixAccAddr
	ServerContext.Bech32PrefixAccPub = ctx.Bech32PrefixAccPub

	for i := 0; i < nodeNum; i++ {
		nodeName := fmt.Sprintf("node%d", i)
		nodeDir := path.Join(devNetHomeDir, nodeName, "testnoded")
		cliDir := path.Join(devNetHomeDir, nodeName, "testnodecli")
		ctxConfig.SetRoot(nodeDir)
		for _, subdir := range []string{"data", "config"} {
			err = os.MkdirAll(path.Join(nodeDir, subdir), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		// app.toml
		binanceChainConfig := ServerContext.BinanceChainConfig
		binanceChainConfig.UpgradeConfig.BEP3Height = 1
		binanceChainConfig.UpgradeConfig.BEP8Height = 1
		binanceChainConfig.UpgradeConfig.BEP12Height = 1
		binanceChainConfig.UpgradeConfig.BEP67Height = 1
		binanceChainConfig.UpgradeConfig.BEP70Height = 1
		binanceChainConfig.UpgradeConfig.BEP82Height = 1
		binanceChainConfig.UpgradeConfig.BEP84Height = 1
		binanceChainConfig.UpgradeConfig.BEP87Height = 1
		binanceChainConfig.UpgradeConfig.FixFailAckPackageHeight = 1
		binanceChainConfig.UpgradeConfig.EnableAccountScriptsForCrossChainTransferHeight = 1
		binanceChainConfig.UpgradeConfig.BEP128Height = 1
		binanceChainConfig.UpgradeConfig.BEP151Height = 1
		binanceChainConfig.UpgradeConfig.BEP153Height = 50
		binanceChainConfig.BreatheBlockInterval = 100
		binanceChainConfig.LogToConsole = false
		binanceChainConfig.CrossChainConfig.BscIbcChainId = 714
		appConfigFilePath := filepath.Join(ctxConfig.RootDir, "config", "app.toml")
		config.WriteConfigFile(appConfigFilePath, binanceChainConfig)
		// pk
		nodeID, pubKey := bnbInit.InitializeNodeValidatorFiles(ctxConfig)
		ctxConfig.Moniker = nodeName
		valOperAddr, secret := bnbInit.CreateValOperAccount(cliDir, ctxConfig.Moniker)
		fmt.Printf("%v secret: %v\n", nodeName, secret)
		if i == 0 {
			memo := fmt.Sprintf("%s@%s:26656", nodeID, "127.0.0.1")
			genTx := bnbInit.PrepareCreateValidatorTx(cdc, chainID, ctxConfig.Moniker, memo, valOperAddr, pubKey)
			appState, err = appInit.AppGenState(cdc, []json.RawMessage{genTx})
			if err != nil {
				panic(err)
			}
			seeds = fmt.Sprintf("%s@127.0.0.1:26656", nodeID)
		} else {
			ctxConfig.P2P.Seeds = seeds
		}
		genFile := ctxConfig.GenesisFile()
		// genesis.json
		err = bnbInit.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genesisTime)
		if err != nil {
			panic(err)
		}
		// edit ctxConfig
		ctxConfig.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", 26656+4*i)
		ctxConfig.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", 26657+4*i)
		ctxConfig.ProxyApp = fmt.Sprintf("tcp://0.0.0.0:%d", 26658+4*i)
		ctxConfig.LogLevel = "main:info,oracle:info,stake:info,*:error"
		// config.toml
		bnbInit.WriteConfigFile(ctxConfig)
	}
}
