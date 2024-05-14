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

	queryMap, err := LoadQueryMapDirectory(Config{Directory: dir, ServersFile: "serverfile.yaml"})
	if err != nil {
		t.Fatal(err)
	}

	fileParts := strings.Split(file.Name(), ".")
	pathParts := strings.Split(fileParts[0], "/")

	if query != queryMap[pathParts[len(pathParts)-1]][query.DatabasePos] {
		t.Errorf("expected %v, got %v", query, queryMap[file.Name()][query.DatabasePos])
	}
}

func TestFindQuery(t *testing.T) {
	qm := QueryMap{
		"testGroup": []Query{
			{
				Name:          "testQuery",
				Query:         "SELECT * FROM test",
				DatabaseGroup: "testGroup",
				DatabasePos:   0,
			},
		},
	}

	query, err := qm.FindQuery("testGroup", "testQuery")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if query.Name != "testQuery" {
		t.Errorf("Expected query name to be 'testQuery', got '%s'", query.Name)
	}

	_, err = qm.FindQuery("nonExistentGroup", "testQuery")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	_, err = qm.FindQuery("testGroup", "nonExistentQuery")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}
