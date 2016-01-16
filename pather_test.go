package main

import (
	"os"
	"runtime"
	"testing"
)

/*\/\/ Helpers /\/\*/

const nullSourcesArrayMessage = "returned null instead of sources array."
const emptySourcesArrayMessage = "returned an empty sources array."

func wrongOsMessage() string {
	return "skipping test not applicable for current OS (" + runtime.GOOS + ")."
}

/*\/\/ Tests /\/\*/

func TestGetLinuxSearchSources(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip(wrongOsMessage())
	}

	home := os.Getenv("HOME")
	searchSources := getLinuxSearchSources(home)

	if searchSources == nil {
		t.Log(nullSourcesArrayMessage)
		t.Fail()
	}

	if len(searchSources) <= 0 {
		t.Log(emptySourcesArrayMessage)
		t.Fail()
	}
}

func TestGetDarwinSearchSources(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip(wrongOsMessage())
	}

	home := os.Getenv("HOME")
	searchSources := getDarwinSearchSources(home)

	if searchSources == nil {
		t.Log(nullSourcesArrayMessage)
		t.Fail()
	}

	if len(searchSources) <= 0 {
		t.Log(emptySourcesArrayMessage)
		t.Fail()
	}
}
