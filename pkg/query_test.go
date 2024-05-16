package termsql

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadQueryMapDirectory(t *testing.T) {
	validDir, err := os.MkdirTemp("", "valid")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(validDir)

	emptyDir, err := os.MkdirTemp("", "empty")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	invalidDir := "invalidDirectory"

	testCases := []struct {
		name           string
		config         Config
		expectError    bool
		expectedErrMsg string
	}{
		{
			name: "Valid Directory",
			config: Config{
				Directory: &validDir,
			},
			expectError: false,
		},
		{
			name: "Invalid Directory",
			config: Config{
				Directory: &invalidDir,
			},
			expectError:    true,
			expectedErrMsg: "error reading directory: invalidDirectory",
		},
		{
			name: "Empty Directory",
			config: Config{
				Directory: &emptyDir,
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadQueryMapDirectory(tc.config)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for directory '%s', got nil", *tc.config.Directory)
				} else if err.Error() != tc.expectedErrMsg {
					t.Errorf("Expected error message to be '%s', got '%s'", tc.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error for directory '%s', got %v", *tc.config.Directory, err)
				}
			}
		})
	}
}

func TestFindQuery(t *testing.T) {
	qm := QueryMap{
		"testGroup": []*Query{
			{
				Name:          "testQuery",
				Query:         "SELECT * FROM test",
				DatabaseGroup: "testServerGroup",
				DatabasePos:   0,
			},
		},
	}

	testCases := []struct {
		name           string
		groupName      string
		queryName      string
		expectError    bool
		expected       *Query
		expectedErrMsg string
	}{
		{
			name:      "Valid Query",
			groupName: "testGroup",
			queryName: "testQuery",
			expected: &Query{
				Name:          "testQuery",
				Query:         "SELECT * FROM test",
				DatabaseGroup: "testServerGroup",
				DatabasePos:   0,
			},
		},
		{
			name:           "Non-Existent Group",
			groupName:      "nonExistentGroup",
			queryName:      "testQuery",
			expectError:    true,
			expectedErrMsg: fmt.Sprintf("query group %s not found, groups:%v", "nonExistentGroup", []string{"testGroup"}),
		},
		{
			name:           "Non-Existent Query",
			groupName:      "testGroup",
			queryName:      "nonExistentQuery",
			expectError:    true,
			expectedErrMsg: fmt.Sprintf("query %s not found in group:%s, available queries:%v", "nonExistentQuery", "testGroup", []string{"testQuery"}),
		},
	}

	for _, tc := range testCases {
		query, err := qm.FindQuery(tc.groupName, tc.queryName)
		if tc.expectError {
			if err == nil {
				t.Errorf("Expected an error for group '%s' and query '%s', got nil", tc.groupName, tc.queryName)
			} else if err.Error() != tc.expectedErrMsg {
				t.Errorf("Expected error message to be '%s', got '%s'", tc.expectedErrMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Fatalf("Expected no error for group '%s' and query '%s', got %v", tc.groupName, tc.queryName, err)
			}
			if *query != *tc.expected {
				t.Errorf("Expected query to be '%v', got '%v'", tc.expected, query)
			}
		}
	}
}
