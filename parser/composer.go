package parser

import (
	"bytes"
	"encoding/json"
)

type platform struct {
	PHP string `json:"php"`
}

type config struct {
	Platform platform `json:"platform"`
}

type composer struct {
	Config config `json:"config"`
}

type composerParser struct{}

func NewComposerParser() *composerParser {
	return &composerParser{}
}

func (c composerParser) Parse(b []byte) (string, error) {
	conf := composer{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&conf); err != nil {
		return "", err
	}
	return conf.Config.Platform.PHP, nil
}
