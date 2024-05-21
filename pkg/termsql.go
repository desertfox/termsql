package termsql

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Encoding int

const (
	JSON Encoding = iota
	YAML
	CSV
	ServersFile string = "servers.yaml"
	HistoryFile string = "history.yaml"
)

var (
	DEBUG bool = os.Getenv("DEBUG_TERMSQL") != ""
)

func Run(c Config, q *Query) (string, error) {
	results, err := RunQuery(c, q)
	if err != nil {
		return "", fmt.Errorf("error running query: %w", err)
	}
	return EncodeStringMap(c, results)
}

func RunQuery(c Config, q *Query) (map[string]string, error) {
	serverList, err := LoadServerList(c)
	if err != nil {
		return nil, err
	}

	s, err := serverList.FindServer(q)
	if err != nil {
		return nil, err
	}

	db, err := MySQLConnect(s)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return RunQueryDynamic(db, q)
}

func LoadQueryMapDirectory(c Config) (QueryMap, error) {
	files, err := os.ReadDir(*c.Directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %s", *c.Directory)
	}

	var QueryMaps QueryMap = make(QueryMap, 0)
	for _, entry := range files {
		if entry.IsDir() && filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		if entry.Name() != ServersFile && entry.Name() != HistoryFile {
			filePath := filepath.Join(*c.Directory, entry.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", filePath, err)
				continue
			}

			var q []*Query
			if err := yaml.Unmarshal(data, &q); err != nil {
				fmt.Println("Error unmarshalling YAML file:", filePath, err)
				continue
			}

			parts := strings.Split(entry.Name(), ".")

			QueryMaps[parts[0]] = q
		}
	}

	return QueryMaps, nil
}

func WriteQueryMapToFile(c Config, qm QueryMap) error {
	for group, queries := range qm {
		if group == "" {
			return fmt.Errorf("group cannot be empty")
		}
		filePath := filepath.Join(*c.Directory, group+".yaml")
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
		}
		defer file.Close()

		data, err := yaml.Marshal(&queries)
		if err != nil {
			return fmt.Errorf("error marshaling to YAML: %s", err)
		}

		if _, err := file.Write(data); err != nil {
			return fmt.Errorf("error writing to file: %s", err)
		}
	}

	return nil
}

func LoadServerList(c Config) (ServerList, error) {
	_, err := os.Stat(*c.Directory)
	if err != nil && os.IsNotExist(err) {
		return ServerList{}, fmt.Errorf("no directory found: %s", *c.Directory)
	} else if err != nil {
		return ServerList{}, err
	}

	file, err := os.ReadFile(c.BuildServerPath())
	if err != nil {
		return ServerList{}, err
	}

	var serverList ServerList
	err = yaml.Unmarshal(file, &serverList)
	if err != nil {
		return ServerList{}, err
	}

	if len(serverList) == 0 {
		return ServerList{}, NoServersFoundError{}
	}

	return serverList, nil
}
