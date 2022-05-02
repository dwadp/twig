package parser

import (
	"bufio"
	"bytes"
	"strings"
)

type Parser interface {
	Parse([]byte) (string, error)
}

type defaultParser struct{}

func NewDefaultParser() *defaultParser {
	return &defaultParser{}
}

func (d defaultParser) Parse(data []byte) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		content := scanner.Text()
		if strings.HasPrefix(content, "version") && strings.Contains(content, "=") {
			parts := strings.Split(content, "=")
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}
