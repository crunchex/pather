package main

import "fmt"
import "strings"
import "os"
import "flag"

func main() {
	useList := flag.Bool("l", false, "use a long listing format")

	flag.Parse()

	if *useList {
		pathList := strings.Split(os.Getenv("PATH"), ":")
		for _, p := range pathList {
			fmt.Println(p)
		}
	} else {
		fmt.Println(os.Getenv("PATH"))
	}
}
