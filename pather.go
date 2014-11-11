package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
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

func appendSource(element string) string {
	element = element + " set by: "
	for _, source := range getSearchSources() {
		b, err := ioutil.ReadFile(source)
		if err != nil {
			// Allow execution to continue as some files are optional.
			//panic(err)
		}

		if strings.Contains(string(b), element) {
			return element + source
		}
	}
	return element + "unknown"
}

func returnPath(shouldList, detailedList bool) {
	if !(shouldList || detailedList) {
		fmt.Println(os.Getenv("PATH"))
		return
	}

	pathList := strings.Split(os.Getenv("PATH"), ":")
	for _, p := range pathList {
		if detailedList {
			p = appendSource(p)
		}

		fmt.Println(p)
	}
}

func main() {
	useList := flag.Bool("l", false, "use a long listing format")
	detailedList := flag.Bool("d", false, "use a (detailed) long listing format")

	flag.Parse()
	returnPath(*useList, *detailedList)
}
