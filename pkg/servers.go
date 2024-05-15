package termsql

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerList map[string]Servers

type Servers struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Db         string `yaml:"db"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Pass       string `yaml:"pass"`
	User       string `yaml:"user"`
	ClientKey  string `yaml:"client_key"`
	ClientCert string `yaml:"client_cert"`
	CaFile     string `yaml:"ca_file"`
}

type NoServersFoundError struct{}

func (x NoServersFoundError) Error() string {
	d, _ := yaml.Marshal(DummyServer())
	return "no servers found, example: \n" + string(d)
}

func DummyServer() Server {
	return Server{
		Db:         "db",
		Host:       "host",
		Port:       3306,
		Pass:       "pass",
		User:       "user",
		ClientKey:  "client_key",
		ClientCert: "client_cert",
		CaFile:     "ca_file",
	}
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

func (x ServerList) FindServer(q Query) (Server, error) {
	if _, ok := x[q.DatabaseGroup]; !ok {
		return Server{}, fmt.Errorf("server group \"%s\" not found, groups:%v", q.DatabaseGroup, x.Keys())
	}

	return x[q.DatabaseGroup].Servers[0], nil
}

func (x ServerList) Keys() []string {
	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	return keys
}

func (x Server) String() string {
	return fmt.Sprintf("Db:%s,Host:%s,Port:%d,User:%s,Pass:%s,ClientKey:%s,ClientCert:%s,CaFile:%s",
		x.Db, x.Host, x.Port, x.User, "***", x.ClientKey, x.ClientCert, x.CaFile)
}
