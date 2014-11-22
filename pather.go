package main

import (
	"bufio"
	"fmt"
	"github.com/kr/fs"
	"github.com/ogier/pflag"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func getSearchSources() []string {
	home := os.Getenv("HOME")

	// OS X
	if runtime.GOOS == "darwin" {
		var sources []string

		sources = append(sources, "/etc/paths")

		walker := fs.Walk("/etc/paths.d")
		for walker.Step() {
			if err := walker.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			if walker.Path() != "/etc/paths.d" {
				sources = append(sources, walker.Path())
			}
		}

		sources = append(sources, home+"/.bash_profile")

		return sources
	}

	// Ubuntu
	bashrc := home + "/.bashrc"
	bashprofile := home + "/.bash_profile"
	profile := home + "/.profile"
	env := "/etc/environment"

	return []string{bashrc, bashprofile, profile, env}
}

func appendSource(path string, pathChan chan string) {
	pathSetBy := path + " set by: "

	for _, source := range getSearchSources() {
		f, err := os.Open(source)
		if err != nil {
			// Allow execution to continue as some files are optional.
			//panic(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for i := 1; scanner.Scan(); i++ {
			if runtime.GOOS == "darwin" {
				if strings.Contains(scanner.Text(), path) {
					p := pathSetBy + source + " (line " + strconv.Itoa(i) + ")"
					pathChan <- p
					return
				}
			} else {
				if strings.Contains(scanner.Text(), "PATH=") {
					if strings.Contains(scanner.Text(), path) {
						p := pathSetBy + source + " (line " + strconv.Itoa(i) + ")"
						pathChan <- p
						return
					}
				}
			}
		}
	}

	pathChan <- pathSetBy + "unknown"
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
