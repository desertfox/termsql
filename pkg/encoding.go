package termsql

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

func EncodeStringMap(c Config, m map[string]string) (string, error) {
	switch Encoding(*c.OutputEncoding) {
	case JSON:
		return encodeJSON(m)
	case YAML:
		return encodeYAML(m)
	case CSV:
		return encodeCSV(m)
	default:
		return "", fmt.Errorf("unknown encoding %d", c.OutputEncoding)
	}
}

func encodeJSON(m map[string]string) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("error encoding JSON: %w", err)
	}
	return string(b), nil
}

func encodeYAML(m map[string]string) (string, error) {
	b, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("error encoding YAML: %w", err)
	}
	return string(b), nil
}

func encodeCSV(m map[string]string) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)
	for k, v := range m {
		if err := w.Write([]string{k, v}); err != nil {
			return "", fmt.Errorf("error encoding CSV: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("error encoding CSV: %w", err)
	}
	return b.String(), nil
}
