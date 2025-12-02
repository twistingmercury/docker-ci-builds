/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = ""
	buildDate = ""
	gitCommit = ""
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Run,
}

func VersionCmd() *cobra.Command {
	return versionCmd
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Run(cmd *cobra.Command, args []string) {
	fmt.Printf("version: %s\ndate: %s\ncommit: %s\n", version, buildDate, gitCommit)
}
