package termsql

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type History []QueryLog

type QueryLog struct {
	Query *Query    `yaml:"query"`
	RunAt time.Time `yaml:"run_at"`
}

func LoadHistory(c Config) (History, error) {
	_, err := os.Stat(*c.Directory)
	if err != nil && os.IsNotExist(err) {
		return History{}, fmt.Errorf("no directory found: %s", *c.Directory)
	} else if err != nil {
		return History{}, err
	}

	if _, err := os.Stat(c.BuildHistoryPath()); os.IsNotExist(err) {
		_, err := os.Create(c.BuildHistoryPath())
		return History{}, err
	}

	file, err := os.ReadFile(c.BuildHistoryPath())
	if err != nil {
		return History{}, err
	}

	var history History
	err = yaml.Unmarshal(file, &history)
	if err != nil {
		return History{}, err
	}

	return history, nil
}

func (x History) WriteHistory(c Config) error {
	filePath := c.BuildHistoryPath()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(x)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (x QueryLog) String() string {
	q, _ := encodeJSON(x.Query.ToMap())

	return fmt.Sprintf("%s: %s", x.RunAt.Format("2006-01-02 15:04:05"), q)
}

func (x *History) Add(q *Query) {
	ql := append([]QueryLog{
		{
			Query: q,
			RunAt: time.Now(),
		},
	}, *x...)

	*x = ql
}

func UpdateHistory(c Config, q *Query) error {
	h, err := LoadHistory(c)
	if err != nil {
		return err
	}

	h.Add(q)

	err = h.WriteHistory(c)
	if err != nil {
		return err
	}

	return nil
}
