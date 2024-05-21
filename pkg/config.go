package termsql

import "path/filepath"

type Config struct {
	Directory      *string
	OutputEncoding *int
}

func (x Config) BuildServerPath() string {
	return filepath.Join(*x.Directory, ServersFile)

}

func (x Config) BuildHistoryPath() string {
	return filepath.Join(*x.Directory, HistoryFile)
}
