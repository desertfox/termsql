package termsql

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadQueryMapDirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "test*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	query := Query{
		Name:          "test",
		Query:         "SELECT * FROM test",
		DatabaseGroup: "test",
		DatabasePos:   0,
	}

	data, err := yaml.Marshal([]Query{query})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Write(data); err != nil {
		t.Fatal(err)
	}
	file.Close()

	queryMap, err := LoadQueryMapDirectory(dir, "serverList")
	if err != nil {
		t.Fatal(err)
	}

	fileParts := strings.Split(file.Name(), ".")
	pathParts := strings.Split(fileParts[0], "/")

	if query != queryMap[pathParts[len(pathParts)-1]][query.DatabasePos] {
		t.Errorf("expected %v, got %v", query, queryMap[file.Name()][query.DatabasePos])
	}
}
