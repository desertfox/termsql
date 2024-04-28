package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type QueryMap map[string][]Query

type Query struct {
	Name          string `yaml:"name"`
	Query         string `yaml:"query"`
	DatabaseGroup string `yaml:"server_group"`
	DatabasePos   int    `yaml:"server_position"`
}

func LoadQueryMapDirectory(p, serverList string) (QueryMap, error) {
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %s", p)
	}

	var QueryMaps QueryMap = make(QueryMap, 0)
	for _, entry := range files {
		if !entry.IsDir() && entry.Name() != serverList && filepath.Ext(entry.Name()) == ".yaml" {
			filePath := filepath.Join(p, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", filePath, err)
				continue
			}

			var q []Query
			err = yaml.Unmarshal(data, &q)
			if err != nil {
				fmt.Println("Error unmarshalling YAML file:", filePath, err)
				continue
			}

			parts := strings.Split(entry.Name(), ".")

			QueryMaps[parts[0]] = q
		}
	}

	return QueryMaps, nil
}

func (x QueryMap) FindQuery(group, queryName string) (Query, error) {
	if _, ok := x[group]; !ok {
		return Query{}, fmt.Errorf("group %s not found", group)
	}

	for _, query := range x[group] {
		if query.Name == queryName {
			return query, nil
		}
	}

	return Query{}, fmt.Errorf("query %s not found in group %s", queryName, group)
}
