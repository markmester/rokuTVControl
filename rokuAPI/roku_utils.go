package rokuAPI

import (
	"bufio"
	"strings"
	"regexp"
)

func ParseApps(apps_raw string) (apps map[string]string) {
	apps = make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(apps_raw))
	r := regexp.MustCompile(`<app id="(\d*)".*>(.*)</app>`)
	for scanner.Scan() {
		match := r.FindStringSubmatch(scanner.Text())
		if len(match) == 3 {
			apps[match[2]] = match[1]
		}
	}

	return apps
}
