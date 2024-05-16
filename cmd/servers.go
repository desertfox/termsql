package cmd

import (
	"fmt"

	"github.com/desertfox/termsql/cmd/output"
	termsql "github.com/desertfox/termsql/pkg"
	"github.com/spf13/cobra"
)

var (
	tsqlGroup  string
	serversCmd = &cobra.Command{
		Use:     "servers",
		Short:   "servers|s",
		Long:    output.BannerWrap("\nList all registered servers"),
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				output.Error(err.Error())
				return
			}

			for group := range serverList {
				if tsqlGroup != "" && group != tsqlGroup {
					continue
				}

				for i, s := range serverList[group].Servers {
					output.Heading(fmt.Sprintf("Group:%s Position:%d", group, i))
					str, err := termsql.EncodeStringMap(config, s.ToMap())
					if err != nil {
						output.Error(err)
						return
					}
					output.Normal(str)
				}
			}
		},
	}
	serversValidateConfigCmd = &cobra.Command{
		Use:     "validate",
		Short:   "validate|v",
		Long:    output.BannerWrap("\nValidate the server connections"),
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			serverList, err := termsql.LoadServerList(config)
			if err != nil {
				output.Error(err.Error())
				return
			}

			output.Normal("Checking server configuration files")

			for group := range serverList {
				for i, s := range serverList[group].Servers {
					str, err := termsql.EncodeStringMap(config, s.ToMap())
					if err != nil {
						output.Error(err)
						return
					}

					if err := termsql.PingServer(s); err != nil {
						output.Error(fmt.Sprintf("Group:%s,Position:%d", group, i))
						output.Error(str)

						continue
					}
					output.Success(fmt.Sprintf("Group:%s,Position:%d", group, i))
					output.Success(str)
				}
			}

			output.Normal("Finished")
		},
	}
)

func init() {
	serversCmd.Flags().StringVarP(&tsqlGroup, "group", "g", "", "termsql group")
	serversCmd.AddCommand(serversValidateConfigCmd)
}
