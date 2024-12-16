package emulator

func RunEmulator(romPath string) {
	println("Iniciando emulador...")

	cpu := newCPU()
	memory := newMemory()
	opcode := uint16(0)

	println("Carregando ROM...")
	memory.LoadROM(cpu, romPath, memory.ram[:])
	println("ROM carregada com sucesso!")

	print("Iniciando loop principal...")
	for {
		opcode = cpu.fetchOpcode(memory.ram[:])
		cpu.executeOpcode(opcode, memory.ram[:])
	}
}
