package emulator

func RunEmulator() {
	cpu := newCPU()
	memory := newMemory()
	opcode := cpu.fetchOpcode(memory.ram[:])

	println("Opcode:", opcode)
}
