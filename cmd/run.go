package cmd

import (
	"chip8-emu/emulator"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Roda o emulador",
	Run: func(cmd *cobra.Command, args []string) {
		emulator.RunEmulator()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
