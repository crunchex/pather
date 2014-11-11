package main

import (
	"fmt"
	"strings"
	"os"
	"flag"
)

func returnPath(shouldList bool) {
	if !shouldList {
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

	flag.Parse()
	returnPath(*useList)
}
