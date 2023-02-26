/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/henokv/azure-env/internal"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

// kvCmd represents the kv command
var readCmd = &cobra.Command{
	Use:   "read <reference>",
	Short: "Read references from akv or acs",
	Long:  `Read references from akv or acs`,
	Args:  readCmdValidator,
	Example: fmt.Sprintf(`Get a key vault secret, stored in azure key vault 'knox' with secret name 'gold':

        %s read azure://knox.vault.azure.net/gold`, rootCmd.Name()),
	Run: readCmdFunc,
}

func readCmdValidator(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected %d arguments but got %d instead", 1, len(args))
	} else if !strings.HasPrefix(args[0], "azure://") {
		return fmt.Errorf("arg should begin with azure:// prefix but got '%s' instead", args[0])
	}
	return nil
}

//var readRef string

func readCmdFunc(cmd *cobra.Command, args []string) {
	ref := args[0]
	vaultUrl, secretName, err := internal.DecodeRef(ref)
	if err != nil {
		log.Fatalf("%s", err)
	}
	secret, err := internal.GetSecret(vaultUrl, secretName)
	if err != nil {
		log.Fatalf("%s", err)
	}
	//log.Printf("%v", *secret.Value)
	fmt.Println(*secret.Value)
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kvCmd.PersistentFlags().String("foo", "", "A help for foo")

	//readCmd.Flags().StringVar(&readRef, "out-file", "", "readCmd")
	//readCmd.MarkFlagRequired("ref")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
