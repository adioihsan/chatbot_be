package helper

import "strings"

func Prefixify(q string) string {
	fields := strings.Fields(q)
	for i, t := range fields {
		// ignore operators from websearch syntax
		switch t {
		case "AND", "OR", "NOT":
			continue
		default:
			// remove punctuation that breaks tsquery (basic sanitize)
			t = strings.Trim(t, `"'()[]{}!&|:*`)
			if t == "" {
				continue
			}
			fields[i] = t + ":*"
		}
	}
	return strings.Join(fields, " ")
}
