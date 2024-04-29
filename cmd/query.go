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
	params   []string
	queryCmd = &cobra.Command{
		Use:     "query",
		Short:   "query|q",
		Long:    `Run a query`,
		Aliases: []string{"q"},
		Run: func(cmd *cobra.Command, args []string) {
			qm, err := termsql.LoadQueryMapDirectory(termSQLDirectory, termSQLServersFile)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			options := make([]huh.Option[string], 0)
			for item := range qm {
				options = append(options, huh.NewOption(item, item))
			}

			var group string
			if len(args) == 0 {
				huh.NewSelect[string]().
					Title("Pick sql group.").
					Options(
						options...,
					).
					Value(&group).Run()
			} else {
				group = args[0]
			}

			options = make([]huh.Option[string], 0)
			for _, q := range qm[group] {
				options = append(options, huh.NewOption(q.Name, q.Name))
			}

			var query string
			if len(args) < 1 {
				huh.NewSelect[string]().
					Title("Pick sql query.").
					Options(
						options...,
					).
					Value(&query).Run()
			} else {
				query = args[0]
			}

			q, err := qm.FindQuery(group, query)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			serverList, err := termsql.LoadServerList(filepath.Join(termSQLDirectory, termSQLServersFile))
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			s, err := serverList.FindServer(q.DatabaseGroup, q.DatabasePos)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			db, err := termsql.Connect(s)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}
			result, err := q.Run(db, params...)
			if err != nil {
				fmt.Println(ui.ERROR_STYLE.Render(err.Error()))
				return
			}

			fmt.Println(ui.BASE_STYLE.Render(ui.ToTable(result)))
		},
	}
	queryCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a query",
		Long:    `Create a query`,
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(filepath.Join(termSQLDirectory, termSQLServersFile))
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

			qm, err := termsql.LoadQueryMapDirectory(termSQLDirectory, "servers.yaml")
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
)

func init() {
	queryCmd.Flags().StringArrayVarP(&params, "params", "p", nil, "Query parameters")
	queryCmd.AddCommand(queryCreateCmd)
}
