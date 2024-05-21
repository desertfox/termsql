package cmd

import (
	"fmt"
	"strconv"

	"github.com/desertfox/termsql/cmd/output"
	termsql "github.com/desertfox/termsql/pkg"
	"github.com/spf13/cobra"
)

var (
	historyCmd = &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "history|h",
		Long:    output.BannerWrap("List run queries"),
		Run: func(cmd *cobra.Command, args []string) {
			history, err := termsql.LoadHistory(config)
			if err != nil {
				output.Error(err.Error())
				return
			}

			if len(args) > 0 {
				index, _ := strconv.Atoi(args[0])
				results, err := termsql.Run(config, history[index].Query)
				if err != nil {
					output.Error(err)
					return
				}

				if err := termsql.UpdateHistory(config, history[index].Query); err != nil {
					output.Error(err)
					return
				}

				output.Success(results)
				return
			}

			for i, h := range history {
				output.Normal(fmt.Sprintf("%d: %s", i, h))
			}
		},
	}
	clearHistoryCmd = &cobra.Command{
		Use:     "clear",
		Aliases: []string{"c"},
		Short:   "clear|c",
		Long:    output.BannerWrap("Clear run queries"),
		Run: func(cmd *cobra.Command, args []string) {
			h := termsql.History{}
			if err := h.WriteHistory(config); err != nil {
				output.Error(err.Error())
				return
			}

			output.Success("History cleared")
		},
	}
)

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.AddCommand(clearHistoryCmd)
}
