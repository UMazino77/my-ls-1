package functions

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strings"
	"syscall"
	"time"
)

type LongFormatInfo struct {
	Permissions fs.FileMode
	NumberLinks string
	User        string
	Group       string
	Size        int64
	Time        time.Time
	FileName    string
}

func MasterSlice(list []fs.FileInfo, flags map[string]bool, total *int) []LongFormatInfo {
	masterSlice := []LongFormatInfo{}
	var User, Group, NumberLinks string
	for _, item := range list {
		if !flags["All"] && strings.HasPrefix(item.Name(), ".") {
			continue
		}

		if stat, ok := item.Sys().(*syscall.Stat_t); ok {
			*total += int(stat.Blocks)
			User = fmt.Sprintf("%d", stat.Uid)
			Group = fmt.Sprintf("%d", stat.Gid)
			NumberLinks = fmt.Sprintf("%d", stat.Nlink)
		}
		if user, err := user.LookupId(User); err == nil {
			User = user.Username
		}
		if group, err := user.LookupGroupId(Group); err == nil {
			Group = group.Name
		}
		masterSlice = append(masterSlice, LongFormatInfo{item.Mode(), NumberLinks, User, Group, item.Size(), item.ModTime(), item.Name()})
	}
	return masterSlice
}

func MyLs(path string, flags map[string]bool, totalPath int, rec bool) {
	list, total := CheckPath(path, flags)
	masterSlice := MasterSlice(list, flags, &total)
	SortLs(masterSlice)
	if flags["Time"] {
		SortByTime(masterSlice)
	}
	if flags["Reverse"] {
		ReverseSorting(masterSlice)
	}
	_, err := os.ReadDir(path)
	if err == nil && (flags["Recursive"] || totalPath > 1) {
		_, l := AddSingleQuotes(path)
		fmt.Printf("%v:\n", l)
	}
	if flags["LongFormat"] {
		_, err := os.ReadDir(path)
		if total != -1 && err == nil {
			fmt.Println("total", total/2)
		}
		LongFormat(masterSlice, path)
	} else {
		ShortFormat(masterSlice, totalPath)
	}
	if flags["Recursive"] {
		for _, item := range masterSlice {
			if fmt.Sprint(item.Permissions)[0] == 'd' {
				if !rec {
					fmt.Println("")
				}
				rec = true
				Recursive(item, path, flags, totalPath, rec)
			} else {
				rec = false
			}
		}
	}
}

func Recursive(item LongFormatInfo, path string, flags map[string]bool, totalPath int, rec bool) {
	if !flags["All"] && (strings.HasPrefix(item.FileName, ".")) || item.FileName == "." || item.FileName == ".." {
		return
	}
	MyLs(JoinPaths(path, item.FileName), flags, totalPath, true)
}

func JoinPaths(s, t string) string {
	for s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	for t[0] == '/' {
		t = t[1:]
	}
	return s + "/" + t
}
