/*
pather is a command-line tool to make working with Unix paths easier.

Usage of pather:
  -d, --detailed-list=false: use a (detailed) long listing format
  -l, --list=false: use a long listing format
*/
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

type Command int

// TODO: figure out why iota does not work here.
const (
	showPath         Command = 0
	showList         Command = 1
	showDetailedList Command = 2
)

func (t Command) String() string {
	var s string
	switch {
	case t&showPath == showPath:
		s = "showPath"
	case t&showList == showList:
		s = "showList"
	case t&showDetailedList == showDetailedList:
		s = "showDetailedList"
	}

	return s
}

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
			if runtime.GOOS != "darwin" {
				if !strings.Contains(scanner.Text(), "PATH=") {
					continue
				}
			}

			if strings.Contains(scanner.Text(), path) {
				p := pathSetBy + source + " (line " + strconv.Itoa(i) + ")"
				pathChan <- p
				return
			}
		}
	}

	// Path wasn't found in any of the known/usual sources.
	pathChan <- pathSetBy + "unknown"
	return
}

// returnPathList returns a slice of path segments that are colon-separated.
// The strings should be printed to stdout.
func returnPathList() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}

// returnDetailedPathList returns a slice of path segments with extra details.
// The strings should be printed to stdout.
func returnDetailedPathList() []string {
	pathList := strings.Split(os.Getenv("PATH"), ":")

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

func printPathList(pathList []string) {
	for _, p := range pathList {
		fmt.Println(p)
	}
}

func executeCommand(cmd Command) {
	switch cmd {
	case showPath:
		// If no special options, we just print the usual PATH env var.
		fmt.Println(os.Getenv("PATH"))
	case showList:
		printPathList(returnPathList())
	case showDetailedList:
		printPathList(returnDetailedPathList())
	}
}

func userInterface() {
	// Don't do anything on unsupported platforms (that we haven't tested yet).
	if !(runtime.GOOS == "darwin" || runtime.GOOS == "linux") {
		fmt.Println("Sorry, pather only supports Linux and OS X for now.")
	}

	const listUsage = "use a long listing format"
	useList := pflag.BoolP("list", "l", false, listUsage)

	const detailedUsage = "use a (detailed) long listing format"
	useDetailedList := pflag.BoolP("detailed-list", "d", false, detailedUsage)

	pflag.Parse()

	cmd := showPath

	if *useList {
		cmd = showList
	}

	if *useDetailedList {
		cmd = showDetailedList
	}

	executeCommand(cmd)
}

func main() {
	userInterface()
}
