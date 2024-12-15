package emulator

func RunEmulator() {
	cpu := NewCPU()
	memory := make([]uint8, 4096) // Memória do CHIP-8 (4 KB)
	opcode := cpu.FetchOpcode(memory)

	println("Opcode:", opcode)
}
