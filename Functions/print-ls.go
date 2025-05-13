package functions

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func AddSingleQuotes(name string) (bool, string) {
	// runes := []rune{' ', '*', '?', '(', ')', '$', '\\', '\'', '&', '|', '<', '>', '~', '[', ']'}
	// for _, r := range runes {
	// 	if strings.ContainsRune(name, r) {
	// 		return true, "'" + name + "'"
	// 	}
	// }
	return false, name
}

func isBlockDevice(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeDevice != 0 && info.Mode()&os.ModeCharDevice == 0, nil
}

func maxlen(path string, slice []LongFormatInfo) (bool, []int, bool) {
	a := false
	major, minor, k := "", "", false
	max0, max1, max2, max3, max4, max5, max6 := 0, 0, 0, 0, 0, 0, 0
	for _, v := range slice {
		info, err := os.Stat(path + "/" + v.FileName)
		if err == nil {
			dev := info.Sys().(*syscall.Stat_t)
			devNbr := dev.Rdev
			major, minor = strconv.Itoa(int(Major(devNbr)))+",", strconv.Itoa(int(Minor(devNbr)))
			if minor == "0" && major == "0," {
				major = ""
				minor = ""
			} else if minor == "0" {
				minor = strconv.Itoa(int(v.Size))
			}
			if major == "0," {
				major = ""
			}
		}
		if len(major) > max6 && major != "" {
			max6 = len(major)
		}
		if len(minor) > max4 {
			max4 = len(minor)
			k = true
		}
		if len(strings.ToLower(fmt.Sprintf("%v", v.Permissions))) > max0 {
			max0 = len(strings.ToLower(fmt.Sprintf("%v", v.Permissions)))
		}
		if len(v.NumberLinks) > max1 {
			max1 = len(v.NumberLinks)
		}
		if len(v.User) > max2 {
			max2 = len(v.User)
		}
		if len(v.Group) > max3 {
			max3 = len(v.Group)
		}
		if len(fmt.Sprintf("%d", v.Size)) > max4 {
			max4 = len(fmt.Sprintf("%d", v.Size))
		}
		if len(formattime(v.Time)) > max5 {
			max5 = len(formattime(v.Time))
		}
		if ok, _ := AddSingleQuotes(v.FileName); ok && !a {
			a = true
		}
	}
	return a, []int{max0, max1, max2, max3, max4, max5, max6}, k
}

func isArch(s string) bool {
	a := []string{".zip", ".tar.gz"}
	for _, v := range a {
		if strings.HasSuffix(s, v) {
			return true
		}
	}
	return false
}

func Major(dev uint64) uint64 {
	return (dev >> 8) & 0xfff
}

func Minor(dev uint64) uint64 {
	return (dev & 0xff) | ((dev >> 12) & 0xfff00)
}

func ACL(path string) (bool, error) {
	size, err := syscall.Listxattr(path, nil)
	if err != nil {
		if err == syscall.ENOTSUP {
			return false, nil
		}
		return false, err
	}
	
	if size <= 0 {
		return false, nil
	}
	
	buf := make([]byte, size)
	size, err = syscall.Listxattr(path, buf)
	if err != nil {
		return false, err
	}
	
	var offset int
	for offset < size {
		end := offset
		for end < size && buf[end] != 0 {
			end++
		}
		
		attrName := string(buf[offset:end])
		
		if attrName == "system.posix_acl_access" || attrName == "system.posix_acl_default" {
			return true, nil
		}
		
		offset = end + 1
	}
	
	return false, nil
}

func Color(name string, permission any) string {
	// /* for the string*/ red, cyan, green, blue, reset, yellow, white, black := "\033[1;31m", "\033[1;36m", "\033[1;32m", "\033[1;34m", "\033[1;m", "\033[1;33m", "\033[37m", "\033[30m"
	// /* for the background*/ orangebg, yellowbg, blackbg := "\033[48;5;208m", "\033[48;5;226m", "\x1b[40m"
	// if fmt.Sprintf("%s", permission)[0] == 'c' {
	// 	return blackbg + yellow + name + reset
	// } else if fmt.Sprintf("%s", permission)[0] == 'l' {
	// 	return cyan + name + reset
	// } else if fmt.Sprintf("%s", permission)[0] == 'd' {
	// 	return blue + name + reset
	// } else if fmt.Sprintf("%s", permission)[3] == 's' {
	// 	return orangebg + white + name + reset
	// } else if fmt.Sprintf("%s", permission)[6] == 's' {
	// 	return yellowbg + black + name + reset
	// } else if isArch(name) {
	// 	return red + name + reset
	// } else if fmt.Sprintf("%s", permission)[0] == '-' && fmt.Sprintf("%s", permission)[3] != 'x' {
	// 	return name
	// } else if fmt.Sprintf("%s", permission)[0] == '-' {
	// 	return green + name + reset
	// }
	return name
}

func formattime(z time.Time) string {
	a, b, c, d, res := z.Month(), z.Day(), z.Year(), fmt.Sprintf("%02d:%02d", z.Hour(), z.Minute()), ""
	ok := time.Now().Sub(z)
	ko := ok.Hours()
	if ko < 4380 {
		res = fmt.Sprintf("%s %2d %5s", fmt.Sprintf("%v", a)[:3], b, d)
	} else {
		res = fmt.Sprintf("%s %2d %5d", fmt.Sprintf("%v", a)[:3], b, c)
	}
	return res
}

func LongFormat(slice []LongFormatInfo, path string) {
	c, a, w := maxlen(path, slice)

	d := map[byte]int{'u': 3, 'g': 6, 'o': 9}
	for _, item := range slice {
		z := JoinPaths(path, item.FileName)
		permissions := fmt.Sprintf("%v", item.Permissions)
		_, item.FileName = (AddSingleQuotes(item.FileName))
		block, _ := isBlockDevice(z)
		if block {
			permissions = "b" + permissions[1:]
		}
		k := []byte(strings.ToLower(permissions))
		zz, _ := os.Stat(path)
		if len(slice) == 1 && (!zz.IsDir() || k[0] == 'l') {
			//fmt.Println("ggg")
			z = path
			item.FileName = path
		}
		if k[1] == 't' {
			temp := k[1]
			k = append(k[0:1], k[2:len(k)-1]...)
			k = append(k, temp)
		}
		if k[0] == 'd' && k[1] == 'c' {
			k = k[1:]
		}
		if k[0] == 'u' || k[0] == 'g' || k[0] == 'o' {
			k[d[k[0]]] = 's'
			k[0] = '-'
		}
		minor, major := strconv.Itoa(int(item.Size)), ""
		if k[0] == 'b' || k[0] == 'c' {
			info, err := os.Stat(z)
			if err == nil {
				dev := info.Sys().(*syscall.Stat_t)
				devNbr := dev.Rdev
				major, minor = strconv.Itoa(int(Major(devNbr)))+",", strconv.Itoa(int(Minor(devNbr)))
				if minor == "0" && major == "0," {
					major = ""
					minor = ""
				} else if minor == "0" {
					minor = strconv.Itoa(int(item.Size))
				}
				if major == "0," {
					major = ""
				}
			}
		}

		acl ,err2 := ACL(path+"/"+item.FileName)

		if err2 != nil {
			fmt.Println(err2)
			os.Exit(1)
		}

		if acl {
			k = append(k, '+')
		}

		l := formattime(item.Time)
		fmt.Printf("%-*s %*s %-*s %-*s ",
			a[0], string(k),
			a[1], item.NumberLinks,
			a[2], item.User,
			a[3], item.Group)
		if w {
			fmt.Printf("%*s ",
				a[6], major)
		}
		fmt.Printf("%*s %-*s ",
			a[4], minor,
			a[5], l)
		target, err := os.Readlink(z)
		if !c || (c && item.FileName[0] == '\'') {
			if err != nil {
				fmt.Printf("%s\n", Color(item.FileName, string(k)))
			} else {
				fmt.Printf("%s -> %s\n", Color(item.FileName, "l"), Color(target, "---x---"))
			}
		} else {
			if err != nil {
				fmt.Printf(" %s\n", Color(item.FileName, string(k)))
			} else {
				fmt.Printf(" %s -> %s\n", Color(item.FileName, "l"), Color(target, "---x---"))
			}
		}
	}
}

func ShortFormat(masterSlice []LongFormatInfo, file int) {
	for _, item := range masterSlice {
		_, item.FileName = AddSingleQuotes(item.FileName)
		fmt.Printf("%v  ", Color(item.FileName, item.Permissions))
	}
	fmt.Println()
}
