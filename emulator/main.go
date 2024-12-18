package emulator

import (
	"time"
)

func RunEmulator(romPath string) {
	println("Iniciando emulador...")

	display := newDisplay()
	keypad := newKeypad()
	cpu := newCPU(display, keypad)
	memory := newMemory()
	opcode := uint16(0)

	println("Carregando ROM...")
	memory.LoadROM(cpu, romPath, memory.ram[:])
	println("ROM carregada com sucesso!")

	win := newWindow()
	defer win.destroy()

	print("Iniciando loop principal...")
	for {
		opcode = cpu.fetchOpcode(memory.ram[:])
		cpu.executeOpcode(opcode, memory.ram[:])
		win.update(cpu.display)

		time.Sleep(16 * time.Millisecond) // Aproximadamente 60Hz
	}
}
