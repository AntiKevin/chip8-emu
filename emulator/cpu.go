package emulator

type CPU struct {
	V          [16]uint8  // 16 registradores de 8 bits (V0 a VF)
	I          uint16     // Registrador de endereço (16 bits)
	PC         uint16     // Program Counter (PC), aponta para a próxima instrução
	Stack      [16]uint16 // Pilha de 16 níveis para sub-rotinas
	SP         uint8      // Stack Pointer (SP), aponta para o topo da pilha
	DelayTimer uint8      // Timer de atraso, decrementa a uma taxa fixa
	SoundTimer uint8      // Timer de som, emite som enquanto for maior que 0
}

func NewCPU() *CPU {
	return &CPU{
		V:          [16]uint8{},
		I:          0,
		PC:         0x200,
		Stack:      [16]uint16{},
		SP:         0,
		DelayTimer: 0,
		SoundTimer: 0,
	}
}

// FetchOpcode busca e retorna o opcode da memória e incrementa o PC
func (cpu *CPU) FetchOpcode(memory []uint8) uint16 {
	highByte := uint16(memory[cpu.PC])
	lowByte := uint16(memory[cpu.PC+1])
	opcode := (highByte << 8) | lowByte

	cpu.PC += 2

	return opcode
}
