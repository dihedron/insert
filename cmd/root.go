package cmd

import (
	log "github.com/dihedron/go-log"
	"github.com/dihedron/zed/text"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zed",
	Short: "A simple stream editor with readability in mind",
	Long: `
	zed is a simple stream editor with a limied set of functionalities
	to manipulate text files one line at a time performing the most basic
	pattern matching and operations such as insertions and replacements.`,
	Args: cobra.NoArgs,
	Run:  text.Copy,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zed.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Error opening configuration directory: %v\n", err)
		}

		// Search config in home directory with name ".zed" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".zed")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file:", viper.ConfigFileUsed())
	}
}
