/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:       "azure-env",
	Short:     "A tool to inject azure key vault secrets in env variables",
	Long:      `A tool to inject azure key vault secrets in env variables`,
	Version:   "1.0.0",
	ValidArgs: []string{"exec"},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	//fmt.Printf("%v\n", verbosity)
	//	//if envFile == "" {
	//	//	fmt.Printf("nothing\n")
	//	//} else {
	//	//	fmt.Printf("line: %s\n", envFile)
	//	//}
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	rootCmd.DisableSuggestions = false
	if err != nil {
		os.Exit(1)
	}
}

var envFile string

var verbosity bool

func init() {
	azureGroup := cobra.Group{ID: "azure", Title: "Azure Commands"}
	rootCmd.AddGroup(&azureGroup)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.akv-env-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&verbosity, "verbosity", "v", false, "Should verbosity loggin be enabled")
	//fmt.Printf("%v\n", *verbosity)
}
