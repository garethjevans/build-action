package pkg

import "strings"

func ParseEnvVars(in string) map[string]string {
	m := make(map[string]string)
	in = strings.TrimSpace(in)

	for _, field := range strings.Fields(in) {
		split := strings.SplitN(field, "=", 2)
		m[split[0]] = split[1]
	}

	return m
}
