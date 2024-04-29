package ui

func ToTable(rows map[string]string) string {
	var output string
	for k, v := range rows {
		output += k + ": " + v + "\n"
	}
	return output
}
