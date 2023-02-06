package cmd

import (
	"os"

	"github.com/dyng/ramen/internal/common"
	conf "github.com/dyng/ramen/internal/config"
	"github.com/dyng/ramen/internal/view"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
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

func init() {
	rootCmd.AddCommand(versionCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		common.Exit("Root command failed: %v", err)
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
		"Should ramen run in debug mode",
	)
	flags.StringVarP(
		&config.ConfigFile,
		"config-file",
		"c",
		conf.DefaultConfigFile,
		"Path to ramen's config file",
	)
	flags.StringVarP(
		&config.Network,
		"network",
		"n",
		conf.DefaultNetwork,
		"Specify the chain that ramen will connect to",
	)
	flags.StringVarP(
		&config.Provider,
		"provider",
		"p",
		conf.DefaultProvider,
		"Specify a blockchain provider",
	)
	flags.StringVar(
		&config.ApiKey,
		"apikey",
		"",
		"ApiKey for specified Ethereum JSON-RPC provider",
	)
	flags.StringVar(
		&config.EtherscanApiKey,
		"etherscan-apikey",
		"",
		"ApiKey for Etherscan API",
	)

	return &cmd
}

func run(cmd *cobra.Command, args []string) {
	// recovery
	defer logPanicAndExit()

	// setup logger
	initLogger()

	// read and parse configurations from config file 
	err := conf.ParseConfig(config)
	if err != nil {
		common.Exit("Cannot parse config file: %v", err)
	}

	// validate config
	err = config.Validate()
	if err != nil {
		common.Exit("Invalid config: %v", err)
	}

	// start application
	view.NewApp(config).Start()
}

func initLogger() {
	// FIXME: use log file in config
	path := "/tmp/ramen.log"
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		common.Exit("Cannot create log file at path %s: %v", path, err)
	}

	handler := common.ErrorStackHandler(log.StreamHandler(file, log.TerminalFormat(false)))
	if config.DebugMode {
		handler = log.LvlFilterHandler(log.LvlDebug, handler)
	} else {
		handler = log.LvlFilterHandler(log.LvlInfo, handler)
	}
	log.Root().SetHandler(handler)
}

func logPanicAndExit() {
	if r := recover(); r != nil {
		log.Error("Unexpected error occurs", "error", errors.Errorf("%v", r))
		common.Exit("Exit due to unexpected error: %v", r)
	}
}
