package emulator

func RunEmulator(romTest string) {
	cpu := newCPU()
	memory := newMemory()
	opcode := cpu.fetchOpcode(memory.ram[:])

	println("opcode:", opcode)
	println("romTest:", romTest)
}
