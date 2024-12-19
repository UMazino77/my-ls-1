package functions

import (
	"strings"
)

func SortPath(slice []string) []string {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			if strings.ToLower(getKey(slice[i])) > strings.ToLower(getKey(slice[j])) {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
	return slice
}

func SortLs(slice []LongFormatInfo) {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if strings.ToLower(getKey(slice[i].FileName)) == "" && strings.ToLower(getKey(slice[j].FileName)) == "" {
				if slice[i].FileName > slice[j].FileName {
					slice[i], slice[j] = slice[j], slice[i]
				}
			} else if strings.ToLower(getKey(slice[i].FileName)) > strings.ToLower(getKey(slice[j].FileName)) {
				slice[i], slice[j] = slice[j], slice[i]
			} else if strings.ToLower(getKey(slice[i].FileName)) == strings.ToLower(getKey(slice[j].FileName)) {
				if slice[i].Time.Before(slice[j].Time) {
					slice[i], slice[j] = slice[j], slice[i]
				}
			}
		}
	}
}

func getKey(filename string) string {
	for i := 0; i < len(filename); i++ {
		if !IsLetter(rune(filename[i])) && !IsDigit(rune(filename[i])) {
			filename = filename[:i] + filename[i+1:]
			i--
		}
	}
	return filename
}

func SortByTime(slice []LongFormatInfo) {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[j].Time.After(slice[i].Time) {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
	for i := 0; i < len(slice)-1; i++ {
		if slice[i].Time.Compare(slice[i+1].Time) == 0 {
			if getKey(slice[i].FileName) > getKey(slice[i+1].FileName) {
				slice[i], slice[i+1] = slice[i+1], slice[i]
				i = -1
			}
		}
	}
}

func ReverseSorting(slice []LongFormatInfo) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func IsLetter(r rune) bool {
	return (r >='a' && r<='z') || (r >='A' && r<='Z')
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9' 
}
