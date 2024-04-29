package termsql

import (
	"database/sql"
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
	queries, err := x.FindQueryGroup(group)
	if err != nil {
		return Query{}, err
	}

	for _, query := range queries {
		if query.Name == queryName {
			return query, nil
		}
	}

	return Query{}, fmt.Errorf("query %s not found in group %s", queryName, group)
}

func (x QueryMap) FindQueryGroup(group string) ([]Query, error) {
	if _, ok := x[group]; !ok {
		return []Query{}, fmt.Errorf("group %s not found", group)
	}

	return x[group], nil
}

func (x Query) Run(db *sql.DB, params ...string) (map[string]string, error) {
	q := x.Query
	for _, p := range params {
		q = strings.Replace(q, "?", p, 1)
	}

	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("error running query: %s", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %s", err)
	}

	var (
		results    = make(map[string]string, len(columns))
		result     = make([]string, len(columns))
		resultPtrs = make([]interface{}, len(columns))
	)

	for i := 0; i < len(columns); i++ {
		resultPtrs[i] = &result[i]
	}

	for rows.Next() {
		if err := rows.Scan(resultPtrs...); err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}

		for i, col := range result {
			results[columns[i]] = col
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	return results, nil
}
