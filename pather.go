// pather is a command-line tool to make working with Unix paths easier.
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

// getLinuxSearchSources will return a list of known locations for PATH
// segments in Ubuntu Linux.
func getLinuxSearchSources(home string) []string {
	// TODO: add official support for other distributions.
	bashrc := home + "/.bashrc"
	bashprofile := home + "/.bash_profile"
	profile := home + "/.profile"
	env := "/etc/environment"

	return []string{bashrc, bashprofile, profile, env}
}

// getDarwinSearchSources will return a list of known locations for PATH
// segments in OS X.
func getDarwinSearchSources(home string) []string {
	bashprofile := home + "/.bash_profile"
	paths := "/etc/paths"

	searchSources := []string{bashprofile, paths}

	// Lastly, grab all the files under paths.d.
	walker := fs.Walk("/etc/paths.d")
	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// We want to exclude the top-level directory.
		if walker.Path() != "/etc/paths.d" {
			searchSources = append(searchSources, walker.Path())
		}
	}

	return searchSources
}

// getSearchSources will return a list of search locations that typically set
// path elements. The list depends on the user's OS.
func getSearchSources() []string {
	home := os.Getenv("HOME")
	var searchSources []string

	switch runtime.GOOS {
	case "linux":
		// Ubuntu
		searchSources = getLinuxSearchSources(home)
	case "darwin":
		// OS X
		searchSources = getDarwinSearchSources(home)
	}

	return searchSources
}

// appendSource will search through a list of possible locations, provided by
// getSearchSources(), where the path may have been set and append that data to
// the path string. If it can't be located, "unknown" will be returned.
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

// returnPathList returns a slice of path strings with or without extra details.
// The strings should be printed to stdout.
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

	if !(runtime.GOOS == "darwin" || runtime.GOOS == "linux") {
		fmt.Println("Sorry, detailed list only supports Linux and OS X for now.")
	}

	for _, p := range returnPathList(*detailedList) {
		fmt.Println(p)
	}
}
