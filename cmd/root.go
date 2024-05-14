package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	termsql "github.com/desertfox/termsql/pkg"
	"github.com/spf13/cobra"
)

var (
	termSQLDirectory   string
	termSQLServersFile string
	config             termsql.Config
	rootCmd            = &cobra.Command{
		Use:   "termsql",
		Short: "termsql is a command line tool for interacting with SQL Server",
		Long:  `termsql is a command line tool for interacting with SQL Server`,
	}
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory")
		return
	}
	defaultDirectory := filepath.Join(home, ".termsql")

	rootCmd.PersistentFlags().StringVarP(&termSQLDirectory, "directory-config", "d", defaultDirectory, "Directory where termsql files are stored")
	rootCmd.PersistentFlags().StringVarP(&termSQLServersFile, "server-config", "s", "servers.yaml", "termsql servers")

	rootCmd.AddCommand(serversCmd)
	rootCmd.AddCommand(queryCmd)

	config = termsql.Config{
		Directory:   termSQLDirectory,
		ServersFile: termSQLServersFile,
	}
}

func Execute() {
	rootCmd.Execute()
}
