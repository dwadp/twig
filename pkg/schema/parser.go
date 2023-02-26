package schema

import (
	"errors"
	"github.com/valyala/fastjson"
	"strings"
)

var (
	ErrNoVersionDefined = errors.New("could not get any version")
)

func GetVersion(d []byte) (versions string, err error) {
	version := fastjson.GetString(d, "require", "php")

	if version == "" {
		version = fastjson.GetString(d, "config", "platform", "php")
	}

	if version == "" {
		return versions, ErrNoVersionDefined
	}

	// If we encounter a "|" replace it with "||" to match the "semver" package constraints parser
	return strings.ReplaceAll(version, "|", "||"), nil
}
