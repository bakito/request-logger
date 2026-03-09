package common

import (
	"regexp"
	"slices"
	"strings"
)

var pathSortWeight = regexp.MustCompile(`[a-zA_Z0-9]`)

// SortPaths sort the path string slice.
func SortPaths(paths []string) {
	slices.SortFunc(paths, func(a, b string) int {
		aa := strings.Split(a, "/")
		bb := strings.Split(b, "/")

		return comparePaths(aa, bb)
	})
}

func comparePaths(a, b []string) int {
	if len(a) > 0 && len(b) > 0 {
		if a[0] == b[0] {
			return comparePaths(a[1:], b[1:])
		}
		return comparePath(a[0], b[0])
	}

	// a or be are not empty
	if len(a) == 0 {
		return -1
	}
	return 1
}

func comparePath(a, b string) int {
	if a != "" && b != "" {
		if a[0] == b[0] {
			return comparePath(a[1:], b[1:])
		}

		aa := pathSortWeight.MatchString(string(a[0]))
		bb := pathSortWeight.MatchString(string(b[0]))

		if aa && bb {
			return strings.Compare(string(a[0]), string(b[0]))
		}
		return 1
	}

	// a or be are not empty
	if a != "" {
		return -1
	}
	return 1
}
