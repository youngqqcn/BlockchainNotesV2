package main

import (
	"fmt"
	"github.com/spf13/viper"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
	clix "github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"mytokenapp/mytokenapp"
	"os"
	"path/filepath"
)

var configFile string // 配置文件路径
var accountDbDirPath string

//// 获取命令行参数，初始化 configFile
//func init() {
//	//flag.StringVar(&configFile, "config", "config/config.toml", "Path to config.toml")
//	flag.StringVar(&accountDbDirPath, "accdb", ".", "Path to save accountdb")
//}

func main() {

	//flag.Parse()
	root := commands.RootCmd
	root.AddCommand(commands.GenNodeKeyCmd)
	root.AddCommand(commands.GenValidatorCmd)
	root.AddCommand(commands.InitFilesCmd)
	root.AddCommand(commands.ResetAllCmd)
	root.AddCommand(commands.ShowNodeIDCmd)
	root.AddCommand(commands.TestnetFilesCmd)

	app := mytokenapp.NewMyTokenApp(".")
	provider := makeNodeProvider(app)
	root.AddCommand(commands.NewRunNodeCmd(provider))

	fmt.Println("starting node ")

	exec := clix.PrepareBaseCmd(root, "yqq", ".")
	exec.Execute()

	//
	//node, err := newTendermint(app, configFile)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%v", err)
	//	os.Exit(2)
	//}
	//
	//err = node.Start()
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%v", err)
	//	os.Exit(3)
	//}
	//defer func() {
	//	node.Stop()
	//	node.Wait()
	//}()
	//
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//<-c
	//os.Exit(0)
}

func makeNodeProvider(app abcitypes.Application) nm.Provider {
	return func(config *cfg.Config, logger log.Logger) (*nm.Node, error) {

		nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
		if err != nil {
			return nil, err
		}

		return nm.NewNode(config,
			privval.LoadOrGenFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile()),
			nodeKey,
			proxy.NewLocalClientCreator(app),
			nm.DefaultGenesisDocProviderFunc(config),
			nm.DefaultDBProvider,
			nm.DefaultMetricsProvider(config.Instrumentation),
			logger,
		)
	}
}

func newTendermint(app abcitypes.Application, configFile string) (*nm.Node, error) {
	// read config
	config := cfg.DefaultConfig()
	config.RootDir = filepath.Dir(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read config file: %w", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("viper failed to unmarshal config: %w", err)
	}
	if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var err error
	logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	// read private validator
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(), // 用于签名共识消息
		config.PrivValidatorStateFile(),
	)

	// read node key 获取节点key 用于 p2p 网络身份识别
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load node's key: %w", err)
	}

	// create node  创建节点实例
	node, err := nm.NewNode(
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app), //直接创建Client实例， 而不是通过Socket或 gRPC ?
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Tendermint node: %w", err)
	}

	return node, nil
}
