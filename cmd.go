package main

import (
	"regexp"
	"strconv"
)

var assignRe = regexp.MustCompile(`(?mi)^/sweepstakes\s*([\d]+)\s*$`)

func parseCmd(comment string) int {
	if m := assignRe.FindAllStringSubmatch(comment, -1); len(m) > 0 {
		v, _ := strconv.Atoi(m[0][1])
		return v
	}

	return -1
}
