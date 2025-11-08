package cmd

import (
	"fmt"
	"os"

	"github.com/prodXCE/gobox/runner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [rootfs path] [command]",
	Short: "Run a command inside a new container",
	Long:  `Runs a command in a new, isolated container.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: gobox run <rootfs-path> <command> [args...]")
			os.Exit(1)
		}

		rootfsPath := args[0]
		command := args[1:]

		runner.Run(rootfsPath, command)
	},
	DisableFlagParsing: true,
}
