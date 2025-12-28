package main

import (
	"slices"
	"strings"
)

func chkTLD(ext string) (bool, int) {
	ext = strings.ToUpper(ext)
	invExt := []string{"\n", " ", "-", "_"}
	if slices.Contains(invExt, ext) || ext == "" { return false, 0 }

	for i, tld := range strings.Split(ianaTldList, "\n")[1:] {
		if ext == tld {
			return true, i
		}
	}
	
	return false, 0
}
