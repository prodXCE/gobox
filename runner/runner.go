package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func Run(rootfsPath string, args []string) {
	fmt.Printf("Parent: Running command %v in %s\n", args, rootfsPath)

	absRootfs, err := filepath.Abs(rootfsPath)
	if err != nil {
		fmt.Printf("Parent: Error resolving rootfs path: %v\n", err)
		os.Exit(1)
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Parent: Error finding executable: %v\n", err)
		os.Exit(1)
	}

	// Note: We are just passing "child", "rootfs", and the command
	childArgs := append([]string{"child", absRootfs}, args...)
	cmd := exec.Command(exePath, childArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Our stable, sudo-based namespace set
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
