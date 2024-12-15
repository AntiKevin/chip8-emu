package cmd

import (
	"chip8-emu/emulator"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:       "run",
	Short:     "Roda o emulador",
	Args:      cobra.MatchAll(cobra.MinimumNArgs(1)),
	ValidArgs: []string{"rom"},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			println("Por favor, informe o caminho do arquivo ROM")
			return
		}
		emulator.RunEmulator(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
