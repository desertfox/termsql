package termsql

import (
	"fmt"

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

func (x ServerList) FindServer(q *Query) (Server, error) {
	if _, ok := x[q.DatabaseGroup]; !ok {
		return Server{}, fmt.Errorf("server group \"%s\" not found, groups:%v", q.DatabaseGroup, x.Keys())
	}

	return x[q.DatabaseGroup].Servers[0], nil
}

func (x ServerList) FindServerGroup(group string) (Servers, error) {
	if _, ok := x[group]; !ok {
		return Servers{}, fmt.Errorf("server group %s not found, groups:%v", group, x.Keys())
	}

	return x[group], nil
}

func (x ServerList) Keys() []string {
	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	return keys
}

func (x Server) ToMap() map[string]string {
	return map[string]string{
		"db":          x.Db,
		"host":        x.Host,
		"port":        fmt.Sprintf("%d", x.Port),
		"pass":        "***Redacted***",
		"user":        x.User,
		"client_key":  x.ClientKey,
		"client_cert": x.ClientCert,
		"ca_file":     x.CaFile,
	}
}
