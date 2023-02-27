/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/henokv/azure-env/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <command> (args...)",
	Args:  cobra.MinimumNArgs(1),
	Short: "Runs commands with all env vars converted if azure reference found ",
	Long:  `Runs commands with all env vars converted if azure reference found.`,
	Example: fmt.Sprintf(`Runs 'terraform plan' but will add all env vars following pattern 'azure://' with key vault ref:

        %s run terraform plan`, rootCmd.Name()),
	//Run: runCmdFunc,
	RunE:    runCmdFunc,
	Aliases: []string{"exec"},
	GroupID: "azure",
}

func runCmdFunc(cmd *cobra.Command, args []string) (error error) {
	internal.SetVerbosity(verbosity)
	var runner *exec.Cmd
	if len(args) == 1 {
		runner = exec.Command(args[0])
	} else {
		runner = exec.Command(args[0], args[1:]...)
	}
	//runner.Env
	secrets, otherEnv, err := internal.GetEnvAsSecret()
	if err != nil {
		return err
		log.Fatalf("hi")
	}
	env := internal.GetFullRenderedEnv(secrets, otherEnv)
	runner.Env = env
	//runner.Env = os.Environ()
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	runner.Stdin = os.Stdin
	runner.Run()
	return nil
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
	runCmd.Flags().StringVarP(&envFile, "env-file", "f", "", "the path to an env file which needs to be used")
}
