package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "rmd",
	Short: "Remote MD",
	Long:  "Remote MD is a remote support tool developed by key.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rmd login            to get a session.")
		fmt.Println("rmd ls               to get a list of all machines")
		fmt.Println("rmd ssh machineID    to ssh to a target machine")
		fmt.Println("rmd vnc machineID    to vnc to a target machine")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func init() {

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(
		&cfgFile, "config", "c", "",
		"For development purposes, use the YAML configuration file.",
	)

	rootCmd.AddCommand(Login)
	rootCmd.AddCommand(LS)
	rootCmd.AddCommand(SSH)
	rootCmd.AddCommand(VNC)
}
