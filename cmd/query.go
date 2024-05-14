package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	termsql "github.com/desertfox/termsql/pkg"
	"github.com/desertfox/termsql/pkg/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	params      []string
	serverGroup string
	serverPos   int
	queryCmd    = &cobra.Command{
		Use:   "query",
		Short: "query|q",
		Long: `query|q

Interactive mode: termsql query
Saved query     : termsql query query_group query_name
Raw query       : termsql q raw server_group server_pos "select * from table"`,
		Aliases: []string{"q"},
		Run: func(cmd *cobra.Command, args []string) {
			qm, err := termsql.LoadQueryMapDirectory(termSQLDirectory, termSQLServersFile)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			var (
				group string
				query string
			)

			options := make([]huh.Option[string], 0)
			for item := range qm {
				options = append(options, huh.NewOption(item, item))
			}

			huh.NewSelect[string]().
				Title("Pick sql group.").
				Options(
					options...,
				).
				Value(&group).Run()

			options = make([]huh.Option[string], 0)
			for _, q := range qm[group] {
				options = append(options, huh.NewOption(q.Name, q.Name))
			}

			huh.NewSelect[string]().
				Title("Pick sql query.").
				Options(
					options...,
				).
				Value(&query).Run()

			q, err := qm.FindQuery(group, query)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			results, err := termsql.RunQuery(termsql.Config{
				Directory:   termSQLDirectory,
				ServersFile: termSQLServersFile,
			}, q)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTable(results)))
		},
	}
	queryCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a query",
		Long:    `Create a query`,
		Run: func(cmd *cobra.Command, args []string) {
			config := termsql.Config{
				Directory:   termSQLDirectory,
				ServersFile: termSQLServersFile,
			}
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			serverOptions := make([]huh.Option[string], 0)
			for server := range serverList {
				serverOptions = append(serverOptions, huh.NewOption(server, server))
			}

			var q termsql.Query
			huh.NewSelect[string]().
				Title("Select server group").
				Options(serverOptions...).
				Value(&q.DatabaseGroup).Run()

			optionsInt := make([]huh.Option[int], 0)
			for pos, server := range serverList[q.DatabaseGroup].Servers {
				optionsInt = append(optionsInt, huh.NewOption(server.Db, pos))
			}

			huh.NewSelect[int]().
				Title("Select database").
				Options(optionsInt...).
				Value(&q.DatabasePos).Run()

			qm, err := termsql.LoadQueryMapDirectory(config.Directory, config.ServersFile)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			groupOptions := make([]huh.Option[string], 0)
			queryGroup := ""
			for group := range qm {
				groupOptions = append(groupOptions, huh.NewOption(group, group))
			}
			groupOptions = append(groupOptions, huh.NewOption("New group", "New group"))
			huh.NewSelect[string]().
				Title("Select query group").
				Options(groupOptions...).
				Value(&queryGroup).Run()

			if queryGroup == "New group" {
				queryGroup = ""
				huh.NewInput().
					Title("Enter query group").
					Value(&queryGroup).Run()
			}

			huh.NewInput().
				Title("Enter query alias").
				Value(&q.Name).Run()

			huh.NewInput().
				Title("Enter query").
				Value(&q.Query).Run()

			qs, _ := qm.FindQueryGroup(queryGroup)
			qs = append(qs, q)

			filePath := filepath.Join(termSQLDirectory, queryGroup+".yaml")
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			data, err := yaml.Marshal(&qs)
			if err != nil {
				fmt.Println("Error marshaling to YAML:", err)
				return
			}

			if _, err := file.Write(data); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

			fmt.Println("Query saved to", filePath)
		},
	}
	rawQueryCmd = &cobra.Command{
		Use:   "raw",
		Short: "Run a raw query",
		Long:  `Run a raw query`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			results, err := termsql.RunQuery(
				termsql.Config{
					Directory:   termSQLDirectory,
					ServersFile: termSQLServersFile,
				},
				termsql.Query{
					Query:         args[0],
					DatabaseGroup: serverGroup,
					DatabasePos:   serverPos,
				},
			)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTable(results)))
		},
	}
	savedQueryCmd = &cobra.Command{
		Use:     "saved",
		Short:   "Run a saved query",
		Long:    `Run a saved query`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			qm := GetQueryMap()

			q, err := qm.FindQuery(args[0], args[1])
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			results, err := termsql.RunQuery(
				termsql.Config{
					Directory:   termSQLDirectory,
					ServersFile: termSQLServersFile,
				},
				q,
			)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTable(results)))

		},
	}
)

func init() {
	queryCmd.Flags().StringArrayVarP(&params, "params", "p", nil, "Query parameters")

	rawQueryCmd.Flags().StringVar(&serverGroup, "server", "", "Server group")
	rawQueryCmd.MarkFlagRequired("server")
	rawQueryCmd.Flags().IntVar(&serverPos, "pos", 0, "Server position, default 0")

	queryCmd.AddCommand(queryCreateCmd)
	queryCmd.AddCommand(rawQueryCmd)
	queryCmd.AddCommand(savedQueryCmd)
}
