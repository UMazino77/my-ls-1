package main

import (
	"fmt"
	"os"

	myls "my-ls-1/Functions"
	
)

func main() {
	paths, flags := myls.ParseArgs((os.Args[1:]))

	wd,err := os.Getwd()
	if err!= nil {
		fmt.Println(err)
		return
	}
	if flags["Help"] {
		fmt.Println("Usage: myls [OPTION]... [FILE]...\nList information about the FILEs (the current directory by default).\nSort entries alphabetically if none of -cftuvSUX nor --sort is specified.\n\nMandatory arguments to long options are mandatory for short options too.\n  -R, --recursive     list subdirectories recursively\n  -r, --reverse      reverse order while sorting\n  -a, --all          do not ignore entries starting with .\n  -l                 use a long listing format\n  -t                 sort by time, newest first; see --time")
		return
	} else if len(paths) == 0 {
		paths = append(paths, ".")
	}
	dirSlice, fileSlice := myls.SplitPath(wd,paths)
	myls.FileSlice(fileSlice, flags)
	myls.DirSlice(fileSlice, dirSlice, flags, len(dirSlice)+len(fileSlice))
}
