package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "tsql-cli",
	Short: "tsql-cli is a command line tool for interacting with SQL Server",
	Long:  `tsql-cli is a command line tool for interacting with SQL Server`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	rootCmd.Execute()
}
