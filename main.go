package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	// resolving to absolute path
	absRootfs, err := filepath.Abs(rootfsPath)
	if err != nil {
		fmt.Println("Parent: Error resolving rootfs path: %v\n", err)
		os.Exit(1)
	}

	// first fix:
	// getting the path to the current running exec
	// implement what /proc/self/exe does on linux
	// by using os.Executable()
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Parent: Error finding executable: %v\n", err)
		os.Exit(1)
	}

	childArgs := append([]string{"child", absRootfs}, args...)
	cmd := exec.Command(exePath, childArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUTS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Parent: Error running child: %v\n", err)
		os.Exit(1)
	}
}

func runChild(rootfsPath string, args []string) {
	fmt.Printf("Child: Setting up jail in %s and running %v\n", rootfsPath, args)

	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		fmt.Printf("Child:Mount private error: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(rootfsPath); err != nil {
		fmt.Printf("Child: os.Chdir to rootfs error: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Chroot("."); err != nil {
		fmt.Printf("Child: Chroot error: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Printf("Child: Mount proc error: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("tmpfs", "tmp", "tmpfs", 0, ""); err != nil {
		fmt.Printf("Child: Mount tmpfs error: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Sethostname([]byte("gobox")); err != nil {
		fmt.Printf("Child: Sethostname error: %v\n", err)
		os.Exit(1)
	}

	cmdPath := args[0]
	cmdArgs := args

	if err := syscall.Exec(cmdPath, cmdArgs, os.Environ()); err != nil {
		fmt.Printf("Child: Exec error: %v\n", err)
		os.Exit(1)
	}
}
