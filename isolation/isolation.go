package isolation

import (
	"fmt"
	"os"
	"syscall"
)

func Child(rootfsPath string, args []string) {
	fmt.Printf("Child: Setting up jail in %s and running %v\n", rootfsPath, args)

	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		fmt.Printf("Child: Mount private error: %v\n", err)
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
