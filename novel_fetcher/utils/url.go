package utils

import "bytes"

// GenerateURL ...
func GenerateURL(parts ...string) string {
	var buffer bytes.Buffer
	for _, part := range parts {
		buffer.WriteString(part)
	}
	return buffer.String()
}
