package api

import (
	"fmt"

	"github.com/mssola/user_agent"
)

func ExtractDeviceName(uaString string) string {
	ua := user_agent.New(uaString)

	name, version := ua.Browser()
	os := ua.OS()

	return fmt.Sprintf("%s %s (%s)", name, version, os)
}
