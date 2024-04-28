package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerList map[string]map[string][]Server

type Server struct {
	Db         string `yaml:"db"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	HumanName  string `yaml:"humanname"`
	Pass       string `yaml:"pass"`
	User       string `yaml:"user"`
	ClientKey  string `yaml:"client_key"`
	ClientCert string `yaml:"client_cert"`
	CaFile     string `yaml:"ca_file"`
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

	return serverList, nil
}

func (x ServerList) FindServer(searchGroup string, position int) (Server, error) {
	if _, ok := x[searchGroup]; !ok {
		return Server{}, fmt.Errorf("group %s not found %v", searchGroup, x)
	}
	if len(x[searchGroup]["servers"]) > 0 && position-1 >= len(x[searchGroup]["servers"]) {
		return Server{}, fmt.Errorf("server %d not found in group %s, %v", position, searchGroup, x)
	}

	return x[searchGroup]["servers"][position], nil
}
