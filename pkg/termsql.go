package termsql

import (
	"path/filepath"
)

type Config struct {
	Directory   string
	ServersFile string
}

func RunQuery(c Config, q Query) (map[string]string, error) {
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

func (x Config) BuildServerPath() string {
	return filepath.Join(x.Directory, x.ServersFile)

}
