package common

import (
	"regexp"
	"sort"
	"strings"
)

var pathSortWeight = regexp.MustCompile(`[a-zA_Z0-9]`)

// SortPaths sort the path string slice
func SortPaths(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		a := strings.Split(paths[i], "/")
		b := strings.Split(paths[j], "/")

		return comparePaths(a, b)
	})
}

func comparePaths(a, b []string) bool {
	if len(a) > 0 && len(b) > 0 {
		if a[0] == b[0] {
			return comparePaths(a[1:], b[1:])
		}
		return comparePath(a[0], b[0])
	}

	// a or be are not empty
	return len(a) > 0
}

func comparePath(a, b string) bool {
	if len(a) > 0 && len(b) > 0 {
		if a[0] == b[0] {
			return comparePath(a[1:], b[1:])
		}

		aa := pathSortWeight.MatchString(string(a[0]))
		bb := pathSortWeight.MatchString(string(b[0]))

		if aa && bb {
			return a[0] > b[0]
		}
		return aa
	}

	// a or be are not empty
	return len(a) > 0
}
