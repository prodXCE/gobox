package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: gobox run <rootfs-path> <command> [args]...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runParent(os.Args[2], os.Args[3:])
	case "child":
		runChild(os.Args[2], os.Args[3:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Usage: gobox run <rootfs-path> <command> [args...]")
		os.Exit(1)
	}
}

func runParent(rootfsPath string, args []string) {
	fmt.Printf("Parent: Running command %v in %s\n", args, rootfsPath)

	// first fix:
	// getting the path to the current running exec
	// implement what /proc/self/exe does on linux
	// by using os.Executable()
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Parent: Error finding executable: %v\n", err)
		os.Exit(1)
	}
	childArgs := append([]string{"child", rootfsPath}, args...)

	cmd := exec.Command(exePath, childArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Parent: Error running child: %v\n", err)
		os.Exit(1)
	}
}

func runChild(rootfsPath string, args []string) {
	fmt.Printf("Child: Setting up jail in %s and running %v\n", rootfsPath, args)

	if err := os.Chdir("/"); err != nil {
		fmt.Printf("Child: Chdir error: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir("/"); err != nil {
		fmt.Printf("Child: Chdir error: %v\n", err)
		os.Exit(1)
	}

	cmdPath := args[0]
	cmdArgs := args

	if err := syscall.Exec(cmdPath, cmdArgs, os.Environ()); err != nil {
		fmt.Printf("Child: Exec error: %v\n", err)
		os.Exit(1)
	}
}
