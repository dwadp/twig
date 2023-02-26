package twig

import (
	"fmt"
	"github.com/dwadp/twig/pkg/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	cfg *config.Config
	err error
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd)
}

var rootCmd = &cobra.Command{
	Use:   "twig",
	Short: "Twig - A multi PHP Command Line executable helper",
	Long: `Twig will help you to run any version of PHP on any project that you have without having to type the php version.
   
for every time you need to run the PHP CLI command.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("To use twig. Run `twig --help`")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	cfg, err = config.NewConfig("twig", config.WithStore(config.NewInMemStore()))
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v\n", err)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if err := cfg.Read(); err != nil {
		log.Fatalf("Failed to read configuration file: %v\n. Make sure you have the configuration file exists and if you didn't just run the `twig init` command and try again.\n", err)
	}
}
