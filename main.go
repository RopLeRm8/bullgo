package main

import (
	"bullgo/examples"
	"os"
)

// Main to test examples

func main() {
	switch os.Args[1] {
	case "example-error":
		examples.BasicError()
	case "example-success":
		examples.BasicSuccess()
	}
}
