/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/henokv/azure-env/internal"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Args:  cobra.MinimumNArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runCmdFunc,
}

func runCmdFunc(cmd *cobra.Command, args []string) {
	var runner *exec.Cmd
	if len(args) == 1 {
		runner = exec.Command(args[0])
	} else {
		runner = exec.Command(args[0], args[1:]...)
	}
	//runner.Env
	secrets, otherEnv, err := internal.GetEnvAsSecret()
	if err != nil {

	}
	env := internal.GetFullRenderedEnv(secrets, otherEnv)
	runner.Env = env
	//runner.Env = os.Environ()
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	runner.Stdin = os.Stdin
	runner.Run()
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
