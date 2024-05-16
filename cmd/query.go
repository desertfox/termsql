package cmd

import (
	"github.com/desertfox/termsql/cmd/forms"
	"github.com/desertfox/termsql/cmd/output"
	termsql "github.com/desertfox/termsql/pkg"

	"github.com/spf13/cobra"
)

var (
	params      []string
	serverGroup string
	serverPos   int
	queryCmd    = &cobra.Command{
		Use:   "query",
		Short: "query|q",
		Long:  output.BannerWrap("\nQuery Interface for executing saved and raw queries"),
		Example: `	List saved queries
		termsql query
	Saved query
		termsql query query_group query_name
	Raw query
		termsql query raw server_group server_pos "select * from table"`,
		Aliases: []string{"q"},
		Run: func(cmd *cobra.Command, args []string) {
			qm, err := termsql.LoadQueryMapDirectory(config)
			if err != nil {
				output.Error(err)
				return
			}

			for group, queries := range qm {
				output.Normal("Group: " + group)
				for _, query := range queries {
					s, err := termsql.EncodeStringMap(config, query.ToMap())
					if err != nil {
						output.Error(err)
						return
					}
					output.Normal(s)
				}
			}
		},
	}
	queryCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create|c",
		Long:    output.BannerWrap("\nCreate and save a new query"),
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				output.Error(err)
				return
			}

			q := &termsql.Query{}

			forms.SelectSeverGroup(q, serverList)

			qm, err := termsql.LoadQueryMapDirectory(config)
			if err != nil {
				output.Error(err)
				return
			}

			forms.UpdateQueryDetails(q)

			qm.AddQuery(forms.SelectOrCreateQueryGroup(qm), q)

			results, err := termsql.Run(config, q)
			if err != nil {
				output.Error(err)
				return
			}

			output.Success(results)

			if err := termsql.WriteQueryMapToFile(config, qm); err != nil {
				output.Error(err)
				return
			}

			if err := termsql.UpdateHistory(config, q); err != nil {
				output.Error(err)
				return
			}

			output.Success("Query saved")
		},
	}
	rawQueryCmd = &cobra.Command{
		Use:     "raw",
		Short:   "raw|r",
		Long:    output.BannerWrap("\nRun a raw query"),
		Aliases: []string{"r"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			q := &termsql.Query{
				Query:         args[0],
				DatabaseGroup: serverGroup,
				DatabasePos:   serverPos,
			}

			results, err := termsql.Run(config, q)
			if err != nil {
				output.Error(err)
				return
			}

			if err := termsql.UpdateHistory(config, q); err != nil {
				output.Error(err)
				return
			}

			output.Success(results)
		},
	}
	loadQueryCmd = &cobra.Command{
		Use:     "load",
		Short:   "load|l",
		Long:    output.BannerWrap("\nLoad and run a saved query"),
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			qm, err := termsql.LoadQueryMapDirectory(config)
			if err != nil {
				output.Error(err)
				return
			}

			if len(args) != 2 {
				args = append(args, forms.SelectQueryGroup(qm))
				args = append(args, forms.SelectQuery(qm, args[0]))
			}

			q, err := qm.FindQuery(args[0], args[1])
			if err != nil {
				output.Error(err)
				return
			}

			results, err := termsql.Run(config, q)
			if err != nil {
				output.Error(err)
				return
			}

			if err := termsql.UpdateHistory(config, q); err != nil {
				output.Error(err)
				return
			}

			output.Success(results)
		},
	}
	saveQueryCmd = &cobra.Command{
		Use:     "save",
		Aliases: []string{"s"},
		Short:   "save|s",
		Long:    output.BannerWrap("\nSave a query"),
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				output.Error(err)
				return
			}

			q := &termsql.Query{}

			forms.SelectSeverGroup(q, serverList)

			qm, err := termsql.LoadQueryMapDirectory(config)
			if err != nil {
				output.Error(err)
				return
			}

			forms.UpdateQueryDetails(q)

			qm.AddQuery(forms.SelectQueryGroup(qm), q)

			termsql.WriteQueryMapToFile(config, qm)
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
