package cmd

import (
	"github.com/desertfox/termsql/cmd/output"
	termsql "github.com/desertfox/termsql/pkg"
	"github.com/spf13/cobra"
)

var (
	historyCmd = &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "history|h",
		Long:    output.BannerWrap("\nList run queries"),
		Run: func(cmd *cobra.Command, args []string) {
			history, err := termsql.LoadHistory(config)
			if err != nil {
				output.Error(err.Error())
				return
			}

			for _, h := range history {
				output.Success(h)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(historyCmd)
}
