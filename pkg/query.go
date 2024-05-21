package termsql

import (
	"fmt"
)

type QueryMap map[string][]*Query

type Query struct {
	Name          string `yaml:"name"`
	Query         string `yaml:"query"`
	DatabaseGroup string `yaml:"server_group"`
	DatabasePos   int    `yaml:"server_position"`
}

func (x QueryMap) FindQuery(group, queryName string) (*Query, error) {
	queries, err := x.FindQueryGroup(group)
	if err != nil {
		return &Query{}, err
	}

	for _, query := range queries {
		if query.Name == queryName {
			return query, nil
		}
	}

	queryNames := make([]string, 0, len(queries))
	for _, query := range queries {
		queryNames = append(queryNames, query.Name)
	}

	return nil, fmt.Errorf("query %s not found in group:%s, available queries:%v", queryName, group, queryNames)
}

func (x QueryMap) FindQueryGroup(group string) ([]*Query, error) {
	if _, ok := x[group]; !ok {
		return nil, fmt.Errorf("query group %s not found, groups:%v", group, x.Keys())
	}

	return x[group], nil
}

func (x QueryMap) Keys() []string {
	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	return keys
}

func (x QueryMap) AddQuery(group string, query *Query) {
	if _, ok := x[group]; !ok {
		x[group] = []*Query{}
	}

	x[group] = append(x[group], query)
}

func (x Query) ToMap() map[string]string {
	return map[string]string{
		"name":           x.Name,
		"database_group": x.DatabaseGroup,
		"database_pos":   fmt.Sprintf("%d", x.DatabasePos),
		"query":          x.Query,
	}
}
