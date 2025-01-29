package pear

import (
	"os"
	"strings"
)

func truncateRecords(recs []ErrorWithStackFrame) []ErrorWithStackFrame {
	bigPaths := make([]string, len(recs))
	for i, rec := range recs {
		bigPaths[i] = rec.File
	}
	smallPaths := truncatePaths(bigPaths)
	for i := range recs {
		recs[i].File = smallPaths[i]
	}
	return recs
}

// chop off the insignificant leading folders from file paths
func truncatePaths(paths []string) []string {
	rows := make([][]string, len(paths))
	for i, p := range paths {
		rows[i] = strings.Split(p, string(os.PathSeparator))
	}
	j := 0

outer:
	for i := range len(rows[0]) {
		word := rows[0][i]
		for _, row := range rows {
			if len(row) < i {
				break outer
			}
			if row[i] != word {
				break outer
			}
		}
		j = i
	}
	shortPaths := make([]string, len(rows))
	for i, row := range rows {
		shortPaths[i] = strings.Join(row[j:], string(os.PathSeparator))
	}
	return shortPaths
}
