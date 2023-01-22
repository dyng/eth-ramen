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
	// setup bootstrap logger
	file, _ := os.OpenFile("ramen.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug,
		log.StreamHandler(file, log.TerminalFormat(false))))

	if err := rootCmd.Execute(); err != nil {
		log.Crit("failed to run roo command", "error", err)
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

	flags.StringVar(
		&config.Network,
		"network",
		conf.DefaultNetwork,
		"Specify a chain to connect",
	)
	flags.StringVar(
		&config.Provider,
		"provider",
		conf.DefaultProvider,
		"Specify a blockchain provider",
	)

	return &cmd
}

func run(cmd *cobra.Command, args []string) {
	view.NewApp(config).Start()
}
