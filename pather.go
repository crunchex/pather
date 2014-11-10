package main

import "fmt"
import "os"

func main() {
	fmt.Printf("%s", os.Getenv("PATH"))
}
