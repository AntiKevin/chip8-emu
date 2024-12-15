package emulator

import "fmt"

type cpu struct {
	V          [16]uint8  // 16 registradores de 8 bits (V0 a VF)
	I          uint16     // Registrador de endereço (16 bits)
	PC         uint16     // Program Counter (PC), aponta para a próxima instrução
	Stack      [16]uint16 // Pilha de 16 níveis para sub-rotinas
	SP         uint8      // Stack Pointer (SP), aponta para o topo da pilha
	DelayTimer uint8      // Timer de atraso, decrementa a uma taxa fixa
	SoundTimer uint8      // Timer de som, emite som enquanto for maior que 0
}

func newCPU() *cpu {
	return &cpu{
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
func (cpu *cpu) fetchOpcode(memory []uint8) uint16 {
	highByte := uint16(memory[cpu.PC])
	lowByte := uint16(memory[cpu.PC+1])
	opcode := (highByte << 8) | lowByte

	cpu.PC += 2

	return opcode
}

// ExecuteOpcode executa a instrução baseada no opcode
func (cpu *cpu) executeOpcode(opcode uint16, memory []uint8) {
	switch opcode & 0xF000 {
	case 0x0000:
		if opcode == 0x00E0 {
			// Limpar a tela
		} else if opcode == 0x00EE {
			// Retornar de uma sub-rotina
		}
	case 0x1000:
		// Salto para o endereço NNN
		cpu.PC = opcode & 0x0FFF
	case 0x2000:
		// Chamar sub-rotina no endereço NNN
		cpu.Stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = opcode & 0x0FFF
	default:
		fmt.Printf("Opcode desconhecido: 0x%X\n", opcode)
	}
}
