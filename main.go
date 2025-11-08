package main

import (
	"fmt"
	"os"

	"github.com/prodXCE/gobox/cmd"
	"github.com/prodXCE/gobox/isolation"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		fmt.Println("Child: Detected, running isolation...")

		if len(os.Args) < 4 {
			fmt.Println("Child: Not enough args for child process")
			os.Exit(1)
		}

		isolation.Child(os.Args[2], os.Args[3:])

		return
	}

	cmd.Execute()
}
