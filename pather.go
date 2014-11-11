package main

import (
	"fmt"
	"strings"
	"os"
	"flag"
)

func returnPath(shouldList, detailedList bool) {
	if !(shouldList || detailedList) {
		fmt.Println(os.Getenv("PATH"))
		return
	}

	pathList := strings.Split(os.Getenv("PATH"), ":")
	for _, p := range pathList {
		fmt.Println(p)
	}
}

func main() {
	useList := flag.Bool("l", false, "use a long listing format")
	detailedList := flag.Bool("d", false, "use a (detailed) long listing format")

	flag.Parse()
	returnPath(*useList, *detailedList)
}
