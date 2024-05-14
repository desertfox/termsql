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
			qm, err := termsql.LoadQueryMapDirectory(config)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render("Available queries:"))
			for group, queries := range qm {
				fmt.Println(ui.BASE_STYLE.Render("group: " + group))
				for _, query := range queries {
					fmt.Println(ui.BASE_STYLE.Render(ui.ToTwoLineString(query)))
				}
			}
		},
	}
	queryCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a query",
		Long:    `Create a query`,
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			serverOptions := make([]huh.Option[string], 0)
			for server := range serverList {
				serverOptions = append(serverOptions, huh.NewOption(server, server))
			}

			q := termsql.Query{}

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

			qm, err := termsql.LoadQueryMapDirectory(config)
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
				config,
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

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTwoLineString(results)))
		},
	}
	loadQueryCmd = &cobra.Command{
		Use:     "load",
		Short:   "Run a saved query",
		Long:    `Run a saved query`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			qm, err := termsql.LoadQueryMapDirectory(config)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			q, err := qm.FindQuery(args[0], args[1])
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			results, err := termsql.RunQuery(config, q)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTwoLineString(results)))
		},
	}
	saveQueryCmd = &cobra.Command{
		Use:     "save",
		Aliases: []string{"s"},
		Short:   "Save a query",
		Long:    `Save a query`,
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			serverOptions := make([]huh.Option[string], 0)
			for server := range serverList {
				serverOptions = append(serverOptions, huh.NewOption(server, server))
			}

			q := termsql.Query{}

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

			qm, err := termsql.LoadQueryMapDirectory(config)
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
			//write data to file
			if _, err := file.Write(data); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

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
	queryCmd.AddCommand(loadQueryCmd)
	queryCmd.AddCommand(saveQueryCmd)
}
