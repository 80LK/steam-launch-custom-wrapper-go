package main

import (
	"os"
)

func main() {
	var cmd = os.Args[1]

	switch cmd {
	case "build":
		build()
		return
	case "patch":
		patch()
		return
	}
}
