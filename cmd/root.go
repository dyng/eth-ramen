package cmd

import (
	"os"

	conf "github.com/dyng/ramen/internal/config"
	"github.com/dyng/ramen/internal/view"
	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/cobra"
)

const (
	appName = "ramen"
	appDesc = "A graphic CLI for interaction with Ethereum easily and happily, by builders, for builders.🍜"
)

var (
	config = conf.NewConfig()
	rootCmd = NewRootCmd()
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Crit("Failed to run roo command", "error", err)
	}
}

func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   appName,
		Short: appDesc,
		Long:  appDesc,
		Run:   run,
	}

	flags := cmd.Flags()

	flags.BoolVar(
		&config.DebugMode,
		"debug",
		false,
		"Should ramen run in debug mode (default: false)",
	)
	flags.StringVar(
		&config.ConfigFile,
		"config-file",
		conf.DefaultConfigFile,
		"Path to ramen's config file (default: ~/.ramen.json)",
	)
	flags.StringVar(
		&config.Network,
		"network",
		conf.DefaultNetwork,
		"Specify the chain that ramen will connect to (default: mainnet)",
	)
	flags.StringVar(
		&config.Provider,
		"provider",
		conf.DefaultProvider,
		"Specify a blockchain provider (default: alchemy)",
	)
	flags.StringVar(
		&config.ApiKey,
		"apikey",
		"",
		"ApiKey for specified Ethereum JSON-RPC provider (required, default: empty)",
	)
	flags.StringVar(
		&config.EtherscanApiKey,
		"etherscan-apikey",
		"",
		"ApiKey for Etherscan API (required, default: empty)",
	)

	return &cmd
}

func run(cmd *cobra.Command, args []string) {
	// setup logger
	initLogger()

	// read and parse configurations from config file 
	err := conf.ParseConfig(config)
	if err != nil {
		log.Crit("Error occurs during parsing config file", "error", err)
	}

	// start application
	view.NewApp(config).Start()
}

func initLogger() {
	file, err := os.OpenFile("/tmp/ramen.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Crit("Failed to create log file", "error", err)
	}

	handler := log.StreamHandler(file, log.TerminalFormat(false))
	if config.DebugMode {
		handler = log.LvlFilterHandler(log.LvlDebug, handler)
	} else {
		handler = log.LvlFilterHandler(log.LvlInfo, handler)
	}
	log.Root().SetHandler(handler)
}
