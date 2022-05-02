package proxy

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dwadp/twig/parser"
	"github.com/dwadp/twig/utils"
)

func GetVersion() (string, error) {
	if utils.FileExists(".twig") {
		return ParseVersion(".twig", parser.NewDefaultParser())
	} else if utils.FileExists("composer.json") {
		return ParseVersion("composer.json", parser.NewComposerParser())
	}
	return "", nil
}

func ParseVersion(path string, parser parser.Parser) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", fmt.Errorf("error opening configuration file [%s]; %v\n", path, err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading configuration file [%s]; %v\n", path, err)
	}
	return parser.Parse(content)
}
