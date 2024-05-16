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
	ServersFile string = "servers.yaml"
	HistoryFile string = "history.yaml"
)

type Config struct {
	Directory      *string
	OutputEncoding *int
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

	db, err := Connect(s)
	if err != nil {
		return nil, err
	}

	return q.Run(db)
}

func Run(c Config, q *Query) (string, error) {
	results, err := RunQuery(c, q)
	if err != nil {
		return "", fmt.Errorf("error running query: %w", err)
	}
	return EncodeStringMap(c, results)
}

func (x Config) BuildServerPath() string {
	return filepath.Join(*x.Directory, ServersFile)

}

func PingServer(s Server) error {
	db, err := Connect(s)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

func (x Config) BuildHistoryPath() string {
	return filepath.Join(*x.Directory, HistoryFile)
}
