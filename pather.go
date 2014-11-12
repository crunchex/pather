package main

import (
	"bufio"
	"fmt"
	"github.com/ogier/pflag"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func getSearchSources() []string {
	// OS X
	if runtime.GOOS == "darwin" {
		return []string{"/etc/paths"}
	}

	// Ubuntu
	home := os.Getenv("HOME")
	bashrc := home + "/.bashrc"
	bashprofile := home + "/.bash_profile"
	profile := home + "/.profile"
	env := "/etc/environment"

	return []string{bashrc, bashprofile, profile, env}
}

func appendSource(element string, c chan string) {
	elementSetBy := element + " set by: "
	for _, source := range getSearchSources() {
		f, err := os.Open(source)
		if err != nil {
			// Allow execution to continue as some files are optional.
			//panic(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for i := 1; scanner.Scan(); i++ {
			if strings.Contains(scanner.Text(), "PATH=") {
				if strings.Contains(scanner.Text(), element) {
					c <- elementSetBy + source + " (line " + strconv.Itoa(i) + ")"
					return
				}
			}
		}
	}
	c <- elementSetBy + "unknown"
	return
}

func returnPathList(detailedList bool) []string {
	pathList := strings.Split(os.Getenv("PATH"), ":")
	if !detailedList {
		return pathList
	}

	pathChan := make(chan string, len(pathList))
	for _, path := range pathList {
		go appendSource(path, pathChan)
	}

	var appendedPathList []string
	for i := 0; i < len(pathList); i++ {
		appendedPathList = append(appendedPathList, <-pathChan)
	}

	return appendedPathList
}

func main() {
	const listUsage = "use a long listing format"
	useList := pflag.BoolP("list", "l", false, listUsage)

	const detailedUsage = "use a (detailed) long listing format"
	detailedList := pflag.BoolP("detailed-list", "d", false, detailedUsage)

	pflag.Parse()

	if !(*useList || *detailedList) {
		fmt.Println(os.Getenv("PATH"))
		return
	}

	for _, p := range returnPathList(*detailedList) {
		fmt.Println(p)
	}
}
