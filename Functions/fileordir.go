package functions

import (
	"fmt"
	"os"
)

func SplitPath(paths []string) ([]string, []string) {
	dirSlice, fileSlice := []string{}, []string{}
	for _, path := range paths {
		if _, err := os.ReadDir(path); err == nil {
			dirSlice = append(dirSlice, path)
		} else {
			fileSlice = append(fileSlice, path)
		}
	}
	return dirSlice, fileSlice
}

func FileSlice(fileSlice []string, flags map[string]bool) {
	SortPath(fileSlice)
	for _, path := range fileSlice {
		MyLs(path, flags, -1,false)
	}
	if len(fileSlice) != 0 && !flags["LongFormat"] {
		fmt.Println()
	}
	
}

func DirSlice(fileSlice, dirSlice []string, flags map[string]bool, totalPath int) {
	SortPath(dirSlice)
	if len(fileSlice) != 0 && len(dirSlice) != 0 {
		fmt.Println()
	}
	for i, path := range dirSlice {
		MyLs(path, flags, totalPath, false)
		if i != len(dirSlice)-1 {
			fmt.Println()
		}
	}
}

