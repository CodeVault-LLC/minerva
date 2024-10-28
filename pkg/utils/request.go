package utils

import (
	"strings"
)

func ParseQuery(query string) (map[string][]string, error) {
	parsedQuery := make(map[string][]string)
	parts := strings.Fields(query)

	for _, part := range parts {
		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				key := strings.ToLower(kv[0])
				value := kv[1]
				parsedQuery[key] = append(parsedQuery[key], value)
			}
		} else {
			// Handle unstructured terms
			parsedQuery["default"] = append(parsedQuery["default"], part)
		}
	}

	return parsedQuery, nil
}
