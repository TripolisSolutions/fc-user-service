package main

import "strconv"

// ParseIntWithFallback get int from string, if nil, return default value
func ParseIntWithFallback(v string, def int) int {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return def
	}
	return int(i)
}
