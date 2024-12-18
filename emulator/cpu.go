package emulator

import (
	"golang.org/x/exp/rand"
)

type cpu struct {
	V          [16]uint8  // 16 registradores de 8 bits (V0 a VF)
	I          uint16     // Registrador de endereço (16 bits)
	PC         uint16     // Program Counter (PC), aponta para a próxima instrução
	Stack      [16]uint16 // Pilha de 16 níveis para sub-rotinas
	SP         uint8      // Stack Pointer (SP), aponta para o topo da pilha
	DelayTimer uint8      // Timer de atraso, decrementa a uma taxa fixa
	SoundTimer uint8      // Timer de som, emite som enquanto for maior que 0
	display    *display   // ponteiro para tela
	keypad     *keypad    // ponteiro para o teclado
}

func newCPU(display *display, keypad *keypad) *cpu {
	return &cpu{
		V:          [16]uint8{},
		I:          0,
		PC:         0x200,
		Stack:      [16]uint16{},
		SP:         0,
		DelayTimer: 0,
		SoundTimer: 0,
		display:    display,
		keypad:     keypad,
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
			cpu.display.clearScreen()
		} else if opcode == 0x00EE {
			// Retornar de uma sub-rotina
			cpu.SP--                   // Decrementa o SP antes de acessar o Stack
			cpu.PC = cpu.Stack[cpu.SP] // Atualiza o PC com o endereço de retorno
		}
	case 0x1000:
		// Salto para o endereço NNN
		cpu.PC = opcode & 0x0FFF
	case 0x2000:
		// Chamar sub-rotina no endereço NNN
		cpu.Stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = opcode & 0x0FFF
	case 0x3000:
		// Saltar se Vx == NN
		x := (opcode & 0x0F00) >> 8
		if cpu.V[x] == uint8(opcode&0x00FF) {
			cpu.PC += 2
		}
	case 0x4000:
		// Saltar se Vx != NN
		x := (opcode & 0x0F00) >> 8
		if cpu.V[x] != uint8(opcode&0x00FF) {
			cpu.PC += 2
		}
	case 0x5000:
		// Saltar se Vx == Vy
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if cpu.V[x] == cpu.V[y] {
			cpu.PC += 2
		}
	case 0x6000:
		// Definir Vx = NN
		x := (opcode & 0x0F00) >> 8
		cpu.V[x] = uint8(opcode & 0x00FF)
	case 0x7000:
		// Adicionar NN a Vx
		x := (opcode & 0x0F00) >> 8
		cpu.V[x] += uint8(opcode & 0x00FF)
	case 0x8000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		switch opcode & 0x000F {
		case 0x0000:
			// Definir Vx = Vy
			cpu.V[x] = cpu.V[y]
		case 0x0001:
			// Definir Vx = Vx OR Vy
			cpu.V[x] |= cpu.V[y]
		case 0x0002:
			// Definir Vx = Vx AND Vy
			cpu.V[x] &= cpu.V[y]
		case 0x0003:
			// Definir Vx = Vx XOR Vy
			cpu.V[x] ^= cpu.V[y]
		case 0x0004:
			// Adicionar Vy a Vx, definir VF = carry
			sum := uint16(cpu.V[x]) + uint16(cpu.V[y])
			if sum > 0xFF {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = uint8(sum & 0xFF)
		case 0x0005:
			// Subtrair Vy de Vx, definir VF = NOT borrow
			if cpu.V[x] > cpu.V[y] {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] -= cpu.V[y]
		case 0x0006:
			// Deslocar Vx para a direita por 1, definir VF = bit menos significativo de Vx antes da mudança
			cpu.V[0xF] = cpu.V[x] & 0x1
			cpu.V[x] >>= 1
		case 0x0007:
			// Definir Vx = Vy - Vx, definir VF = NOT borrow
			if cpu.V[y] > cpu.V[x] {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = cpu.V[y] - cpu.V[x]
		case 0x000E:
			// Deslocar Vx para a esquerda por 1, definir VF = bit mais significativo de Vx antes da mudança
			cpu.V[0xF] = cpu.V[x] >> 7
			cpu.V[x] <<= 1
		}
	case 0x9000:
		// Saltar se Vx != Vy
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if cpu.V[x] != cpu.V[y] {
			cpu.PC += 2
		}
	case 0xA000:
		// Definir I = NNN
		cpu.I = opcode & 0x0FFF
	case 0xB000:
		// Saltar para o endereço NNN + V0
		cpu.PC = (opcode & 0x0FFF) + uint16(cpu.V[0])
	case 0xC000:
		// Definir Vx = rand() AND NN
		x := (opcode & 0x0F00) >> 8
		cpu.V[x] = uint8(rand.Intn(256)) & uint8(opcode&0x00FF)
	case 0xD000:
		// Desenhar sprite na tela
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		height := opcode & 0x000F
		cpu.V[0xF] = 0
		if cpu.display.drawSprite(uint16(cpu.V[x]), uint16(cpu.V[y]), height, cpu.I, memory) {
			cpu.V[0xF] = 1
		}
	case 0xE000:
		x := (opcode & 0x0F00) >> 8
		switch opcode & 0x00FF {
		case 0x009E:
			// Saltar se a tecla com o valor de Vx estiver pressionada
			if cpu.keypad.isKeyPressed(cpu.V[x]) {
				cpu.PC += 2
			}
		case 0x00A1:
			// Saltar se a tecla com o valor de Vx não estiver pressionada
			if !cpu.keypad.isKeyPressed(cpu.V[x]) {
				cpu.PC += 2
			}
		}
	case 0xF000:
		x := (opcode & 0x0F00) >> 8
		switch opcode & 0x00FF {
		case 0x0007:
			// Definir Vx = valor do timer de atraso
			cpu.V[x] = cpu.DelayTimer
		case 0x000A:
			// Esperar até que uma tecla seja pressionada, armazenar o valor da tecla em Vx
			cpu.V[x] = cpu.keypad.waitForKeyPress()
		case 0x0015:
			// Definir timer de atraso = Vx
			cpu.DelayTimer = cpu.V[x]
		case 0x0018:
			// Definir timer de som = Vx
			cpu.SoundTimer = cpu.V[x]
		case 0x001E:
			// Adicionar Vx a I
			cpu.I += uint16(cpu.V[x])
		case 0x0029:
			// Definir I = localização do sprite para o dígito Vx
			cpu.I = uint16(cpu.V[x]) * 5
		case 0x0033:
			// Armazenar BCD de Vx em I, I+1, I+2
			memory[cpu.I] = cpu.V[x] / 100
			memory[cpu.I+1] = (cpu.V[x] / 10) % 10
			memory[cpu.I+2] = (cpu.V[x] % 10)
		case 0x0055:
			// Armazenar V0 até Vx na memória começando em I
			for i := uint16(0); i <= x; i++ {
				memory[cpu.I+i] = cpu.V[i]
			}
		case 0x0065:
			// Ler V0 até Vx da memória começando em I
			for i := uint16(0); i <= x; i++ {
				cpu.V[i] = memory[cpu.I+i]
			}
		}
	}
}
