package forms

import (
	"github.com/charmbracelet/huh"
	termsql "github.com/desertfox/termsql/pkg"
)

func SelectSeverGroup(q *termsql.Query, serverList termsql.ServerList) {
	var (
		serverOptions []huh.Option[string] = make([]huh.Option[string], 0)
		optionsInt    []huh.Option[int]    = make([]huh.Option[int], len(serverList[q.DatabaseGroup].Servers))
	)

	for server := range serverList {
		serverOptions = append(serverOptions, huh.NewOption(server, server))
	}

	huh.NewSelect[string]().
		Title("Select server group").
		Options(serverOptions...).
		Value(&q.DatabaseGroup).Run()

	for pos, server := range serverList[q.DatabaseGroup].Servers {
		optionsInt = append(optionsInt, huh.NewOption(server.Db, pos))
	}

	huh.NewSelect[int]().
		Title("Select database").
		Options(optionsInt...).
		Value(&q.DatabasePos).Run()

}

func SelectQueryGroup(qm termsql.QueryMap) string {
	var (
		groupOptions []huh.Option[string] = make([]huh.Option[string], 0)
		queryGroup   string
	)

	for group := range qm {
		groupOptions = append(groupOptions, huh.NewOption(group, group))
	}

	huh.NewSelect[string]().
		Title("Select query group").
		Options(groupOptions...).
		Value(&queryGroup).Run()

	return queryGroup
}

func SelectOrCreateQueryGroup(qm termsql.QueryMap) string {
	var (
		groupOptions []huh.Option[string] = make([]huh.Option[string], 0)
		queryGroup   string
	)

	for group := range qm {
		groupOptions = append(groupOptions, huh.NewOption(group, group))
	}
	groupOptions = append(groupOptions, huh.NewOption("Create new group", "Create new group"))

	huh.NewSelect[string]().
		Title("Select query group").
		Options(groupOptions...).
		Value(&queryGroup).Run()

	if queryGroup == "Create new group" {
		huh.NewInput().
			Title("Enter new group name").
			Value(&queryGroup).Run()
	}

	return queryGroup
}

func SelectQuery(qm termsql.QueryMap, queryGroup string) string {
	var (
		queryOptions []huh.Option[string] = make([]huh.Option[string], 0)
		queryName    string
	)

	for _, query := range qm[queryGroup] {
		queryOptions = append(queryOptions, huh.NewOption(query.Name, query.Name))
	}

	huh.NewSelect[string]().
		Title("Select query").
		Options(queryOptions...).
		Value(&queryName).Run()

	return queryName
}

func UpdateQueryDetails(q *termsql.Query) {
	huh.NewInput().
		Title("Enter query name").
		Value(&q.Name).Run()

	huh.NewInput().
		Title("Enter query").
		Value(&q.Query).Run()
}

func SelectStringForm(title string, strings []string) string {
	var selectOptions []huh.Option[string] = make([]huh.Option[string], 0, len(strings))
	for _, s := range strings {
		selectOptions = append(selectOptions, huh.NewOption(s, s))
	}

	var selected string
	huh.NewSelect[string]().
		Title(title).
		Options(selectOptions...).
		Value(&selected).Run()

	return selected
}
