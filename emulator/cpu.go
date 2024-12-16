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
	display    *display   // ponteeiro para tela
	keypad     *keypad    // ponteiro para o teclado
	delayTimer uint8      // Timer de atraso, decrementa a uma taxa fixa
	soundTimer uint8      // Timer de som, emite som enquanto for maior que 0
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
		display:    &display{},
		keypad:     &keypad{},
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
// ExecuteOpcode executa a instrução baseada no opcode
func (cpu *cpu) executeOpcode(opcode uint16, memory []uint8) {
	switch opcode & 0xF000 {
	case 0x0000:
		if opcode == 0x00E0 {
			// Limpar a tela
			// Limpar a tela (normalmente define a memória de tela como 0)
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
		cpu.Stack[cpu.SP] = cpu.PC // Armazena o endereço de retorno
		cpu.SP++                   // Incrementa o SP após armazenar o endereço
		cpu.PC = opcode & 0x0FFF   // Atualiza o PC para o endereço da sub-rotina
	case 0x3000:
		// Saltar se Vx == NN
		// Se o valor no registrador Vx for igual ao valor NN, o PC é incrementado para pular a instrução
		x := (opcode & 0x0F00) >> 8
		kk := opcode & 0x00FF
		if cpu.V[x] == uint8(kk) {
			cpu.PC += 2
		}
	case 0x4000:
		// Saltar se Vx != NN
		// Se o valor no registrador Vx for diferente do valor NN, o PC é incrementado para pular a instrução
		x := (opcode & 0x0F00) >> 8
		kk := opcode & 0x00FF
		if cpu.V[x] != uint8(kk) {
			cpu.PC += 2
		}
	case 0x5000:
		// Saltar se Vx == Vy
		// Se os valores nos registradores Vx e Vy forem iguais, o PC é incrementado para pular a instrução
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if cpu.V[x] == cpu.V[y] {
			cpu.PC += 2
		}
	case 0x6000:
		// Definir Vx com NN
		// O registrador Vx é definido como o valor imediato NN
		x := (opcode & 0x0F00) >> 8
		kk := opcode & 0x00FF
		cpu.V[x] = uint8(kk)
	case 0x7000:
		// Adicionar NN a Vx
		// O valor NN é somado ao valor armazenado em Vx
		x := (opcode & 0x0F00) >> 8
		kk := opcode & 0x00FF
		cpu.V[x] += uint8(kk)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			// Definir Vx com Vy
			// O registrador Vx é copiado do valor armazenado em Vy
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			cpu.V[x] = cpu.V[y]
		case 0x0001:
			// Definir Vx com Vx OR Vy
			// A operação OR bit-a-bit é realizada entre Vx e Vy e o resultado é armazenado em Vx
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			cpu.V[x] |= cpu.V[y]
		case 0x0002:
			// Definir Vx com Vx AND Vy
			// A operação AND bit-a-bit é realizada entre Vx e Vy e o resultado é armazenado em Vx
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			cpu.V[x] &= cpu.V[y]
		case 0x0003:
			// Definir Vx com Vx XOR Vy
			// A operação XOR bit-a-bit é realizada entre Vx e Vy e o resultado é armazenado em Vx
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			cpu.V[x] ^= cpu.V[y]
		case 0x0004:
			// Adicionar Vy a Vx com carry
			// O valor de Vy é somado a Vx e, se houver overflow, o carry é armazenado em VF
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			sum := uint16(cpu.V[x]) + uint16(cpu.V[y])
			if sum > 255 {
				cpu.V[0xF] = 1 // Carry
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = uint8(sum & 0xFF)
		case 0x0005:
			// Subtrair Vy de Vx com borrow
			// O valor de Vy é subtraído de Vx, e o borrow é armazenado em VF
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			if cpu.V[x] > cpu.V[y] {
				cpu.V[0xF] = 1 // No borrow
			} else {
				cpu.V[0xF] = 0 // Borrow
			}
			cpu.V[x] -= cpu.V[y]
		case 0x0006:
			// Deslocar Vx para a direita com carry
			// O valor de Vx é deslocado para a direita e o bit de carry é armazenado em VF
			x := (opcode & 0x0F00) >> 8
			cpu.V[0xF] = cpu.V[x] & 0x01
			cpu.V[x] >>= 1
		case 0x0007:
			// Subtrair Vx de Vy com borrow
			// O valor de Vx é subtraído de Vy, e o borrow é armazenado em VF
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			if cpu.V[y] > cpu.V[x] {
				cpu.V[0xF] = 1 // No borrow
			} else {
				cpu.V[0xF] = 0 // Borrow
			}
			cpu.V[y] -= cpu.V[x]
		case 0x000E:
			// Deslocar Vx para a esquerda com carry
			// O valor de Vx é deslocado para a esquerda e o bit de carry é armazenado em VF
			x := (opcode & 0x0F00) >> 8
			cpu.V[0xF] = cpu.V[x] >> 7
			cpu.V[x] <<= 1
		}
	case 0x9000:
		// Saltar se Vx != Vy
		// Se os valores nos registradores Vx e Vy forem diferentes, o PC é incrementado para pular a instrução
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if cpu.V[x] != cpu.V[y] {
			cpu.PC += 2
		}
	case 0xA000:
		// Definir I para o endereço NNN
		cpu.I = opcode & 0x0FFF
	case 0xB000:
		// Saltar para o endereço NNN + V0
		// O PC é atualizado para o endereço NNN somado ao valor de V0
		cpu.PC = (opcode & 0x0FFF) + uint16(cpu.V[0])
	case 0xC000:
		// Gerar um número aleatório e fazer um AND com NN
		x := (opcode & 0x0F00) >> 8
		kk := opcode & 0x00FF
		cpu.V[x] = uint8(rand.Intn(256)) & uint8(kk)
	case 0xD000:
		// Desenhar sprite
		// Desenha um sprite na tela baseado nos dados armazenados na memória a partir de I
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		height := opcode & 0x000F
		cpu.display.drawSprite(x, y, height, cpu.I, memory)
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			// Pular se a tecla Vx estiver pressionada
			// Se a tecla correspondente a Vx estiver pressionada, o PC é incrementado para pular a instrução
			x := (opcode & 0x0F00) >> 8
			if cpu.keypad.keys[cpu.V[x]] {
				cpu.PC += 2
			}
		case 0x00A1:
			// Pular se a tecla Vx não estiver pressionada
			// Se a tecla correspondente a Vx não estiver pressionada, o PC é incrementado para pular a instrução
			x := (opcode & 0x0F00) >> 8
			if !cpu.keypad.keys[cpu.V[x]] {
				cpu.PC += 2
			}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			// Definir Vx com o valor do timer de atraso
			// O valor do timer de atraso é copiado para o registrador Vx
			x := (opcode & 0x0F00) >> 8
			cpu.V[x] = cpu.delayTimer
		case 0x000A:
			// Esperar até que uma tecla seja pressionada
			// O emulador aguarda até que uma tecla seja pressionada e então o valor é armazenado em Vx
			x := (opcode & 0x0F00) >> 8
			// Aqui, você precisará implementar um método para aguardar a tecla ser pressionada
			cpu.V[x] = cpu.keypad.waitForKeypress()
		case 0x0015:
			// Definir o timer de atraso com o valor de Vx
			// O valor de Vx é copiado para o timer de atraso
			x := (opcode & 0x0F00) >> 8
			cpu.delayTimer = cpu.V[x]
		case 0x0018:
			// Definir o timer de som com o valor de Vx
			// O valor de Vx é copiado para o timer de som
			x := (opcode & 0x0F00) >> 8
			cpu.soundTimer = cpu.V[x]
		case 0x001E:
			// Adicionar Vx ao registrador I
			// O valor de Vx é somado ao registrador I
			x := (opcode & 0x0F00) >> 8
			cpu.I += uint16(cpu.V[x])
		case 0x0029:
			// Definir I para o endereço do sprite do caractere
			// O endereço do sprite de Vx é armazenado no registrador I
			x := (opcode & 0x0F00) >> 8
			cpu.I = uint16(cpu.V[x]) * 5
		case 0x0033:
			// Armazenar BCD de Vx na memória
			// O valor de Vx é armazenado em formato BCD (Binary Coded Decimal) na memória
			x := (opcode & 0x0F00) >> 8
			value := cpu.V[x]
			memory[cpu.I] = value / 100
			memory[cpu.I+1] = (value / 10) % 10
			memory[cpu.I+2] = value % 10
		case 0x0055:
			// Armazenar V0 a Vx na memória
			// Os valores de V0 a Vx são armazenados na memória a partir de I
			x := (opcode & 0x0F00) >> 8
			for i := 0; i <= int(x); i++ {
				memory[cpu.I+uint16(i)] = cpu.V[i]
			}
		case 0x0065:
			// Carregar V0 a Vx da memória
			// Os valores de V0 a Vx são carregados da memória a partir de I
			x := (opcode & 0x0F00) >> 8
			for i := 0; i <= int(x); i++ {
				cpu.V[i] = memory[cpu.I+uint16(i)]
			}
		}
	}
}
