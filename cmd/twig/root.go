package twig

import (
	"fmt"
	"os"

	"github.com/dwadp/twig/pkg/config"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
	err error

	Version string = "2.0.0-dev"

	isPrintingVersion bool = false
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd)

	rootCmd.Flags().BoolVarP(&isPrintingVersion, "version", "v", false, "Print the version of twig")
}

var rootCmd = &cobra.Command{
	Use:   "twig",
	Short: "Twig - A multi PHP & Composer Command Line executable helper",
	Long: `Twig will help you to run any version of PHP on any project that you have without having to type the php version.
   
for every time you need to run the PHP CLI command.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if isPrintingVersion {
			cmd.Printf("twig %s\n", Version)
			os.Exit(0)
		}

		cmd.Println("To use twig. Run `twig --help`")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	cfg, err = config.NewConfig("twig", config.WithStore(config.NewInMemStore()))
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to initialize configuration: %q\n", err)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if err := cfg.Read(); err != nil {
		cmd.Printf("failed to read configuration file: %v\n, make sure you have the configuration file exists and if you didn't just run the `twig init` command and try again.\n", err)
		os.Exit(1)
	}
}
