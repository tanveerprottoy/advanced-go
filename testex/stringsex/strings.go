package stringsex

import (
	"bytes"
	"strings"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result strings.Builder

	for i := range length {
		result.WriteByte(charset[i%len(charset)])
	}

	return result.String()
}

func ConcatenateBuffer(first string, second string) string {
	var buffer bytes.Buffer

	buffer.WriteString(first)
	buffer.WriteString(second)

	return buffer.String()
}

func ConcatenateJoin(first string, second string) string {
	return strings.Join([]string{first, second}, "")
}
