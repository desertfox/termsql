package termsql

import (
	"fmt"
	"os"
	"reflect"
	"strings"

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

func LoadServerList(p string) (ServerList, error) {
	file, err := os.ReadFile(p)
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

func (x ServerList) FindServer(searchGroup string, position int) (Server, error) {
	if _, ok := x[searchGroup]; !ok {
		return Server{}, fmt.Errorf("server group \"%s\" not found, groups:%v", searchGroup, x.Keys())
	}

	return x[searchGroup].Servers[0], nil
}

func (s Server) ToTable() string {
	var result strings.Builder

	v := reflect.ValueOf(s)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Name

		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				result.WriteString(fmt.Sprintf("%s: %s\n", name, field.String()))
			}
		case reflect.Int:
			if field.Int() != 0 {
				result.WriteString(fmt.Sprintf("%s: %d\n", name, field.Int()))
			}
		}
	}

	return result.String()
}

func (x ServerList) Keys() []string {
	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	return keys
}
