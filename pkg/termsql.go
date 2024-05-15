package termsql

import (
	"fmt"
	"path/filepath"
)

type Encoding int

const (
	JSON Encoding = iota
	YAML
	CSV
)

type Config struct {
	Directory      *string
	ServersFile    *string
	OutputEncoding *int
}

func runQuery(c Config, q Query) (map[string]string, error) {
	serverList, err := LoadServerList(c)
	if err != nil {
		return nil, err
	}

	s, err := serverList.FindServer(q)
	if err != nil {
		return nil, err
	}

	db, err := Connect(s)
	if err != nil {
		return nil, err
	}

	return q.Run(db)
}

func Run(c Config, q Query) (string, error) {
	results, err := runQuery(c, q)
	if err != nil {
		return "", fmt.Errorf("error running query: %w", err)
	}
	return EncodeStringMap(c, results)
}

func (x Config) BuildServerPath() string {
	return filepath.Join(*x.Directory, *x.ServersFile)

}
