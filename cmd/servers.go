package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	termsql "github.com/desertfox/termsql/pkg"
	"github.com/desertfox/termsql/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	tsqlGroup  string
	serversCmd = &cobra.Command{
		Use:     "servers",
		Short:   "servers|s",
		Long:    `List all registered servers`,
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := loadServerConfig()
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			if tsqlGroup == "" {
				for group := range serverList {
					for i, s := range serverList[group].Servers {
						fmt.Println(ui.BASE_STYLE.Render(fmt.Sprintf("Group   : %s\nPosition: %d\n", group, i) + s.ToTable()))
					}
				}
			} else {
				for i, s := range serverList[tsqlGroup].Servers {
					fmt.Println(ui.BASE_STYLE.Render(fmt.Sprintf("Group: %s\nPosition: %d\n", tsqlGroup, i) + s.ToTable()))
				}
			}
		},
	}
	serverExplorerCmd = &cobra.Command{
		Use:     "explore",
		Short:   "explore|e",
		Aliases: []string{"e"},
		Long:    `Explore a server`,
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := loadServerConfig()
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			var (
				serverOptions []huh.Option[string] = make([]huh.Option[string], 0)
				serverGroup   string
			)

			for server := range serverList {
				serverOptions = append(serverOptions, huh.NewOption(server, server))
			}

			huh.NewSelect[string]().
				Title("Select server group").
				Options(serverOptions...).
				Value(&serverGroup).Run()

			server := serverList[serverGroup]
			db, err := termsql.Connect(server.Servers[0])
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			q := termsql.Query{
				Query: "show tables;",
			}
			result, err := q.Run(db)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			tableOptions := make([]huh.Option[string], 0)
			for table := range result {
				tableOptions = append(tableOptions, huh.NewOption(table, table))
			}

			table := ""
			huh.NewSelect[string]().
				Title("Select table").
				Options(tableOptions...).
				Value(&table).Run()

			q.Query = fmt.Sprintf("show create table %s;", result[table])
			result, err = q.Run(db)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTable(result)))
		},
	}
)

func init() {
	serversCmd.Flags().StringVarP(&tsqlGroup, "group", "g", "", "termsql group")

	serversCmd.AddCommand(serverExplorerCmd)
}

func loadServerConfig() (termsql.ServerList, error) {
	_, err := os.Stat(termSQLDirectory)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("no directory found: %s", termSQLDirectory)
	} else if err != nil {
		return nil, err
	}

	return termsql.LoadServerList(filepath.Join(termSQLDirectory, termSQLServersFile))
}
