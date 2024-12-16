package emulator

type display struct {
	screen [64][32]byte // Tela de 64x32 pixels
}

// drawSprite desenha um sprite na tela e retorna true se houve colisão
func (d *display) drawSprite(x, y, height uint16, I uint16, memory []uint8) bool {
	collision := false

	for row := uint16(0); row < height; row++ {
		spriteRow := memory[I+row] // Acessando diretamente a memória via o índice I
		for col := uint16(0); col < 8; col++ {
			if spriteRow&(0x80>>col) != 0 {
				px := (x + col) % 64
				py := (y + row) % 32

				// Se o pixel já estiver aceso, houve uma colisão
				if d.screen[px][py] == 1 {
					collision = true
				}

				// Inverte o pixel
				d.screen[px][py] ^= 1
			}
		}
	}

	return collision
}

// clearScreen limpa a tela, definindo todos os pixels como 0
func (d *display) clearScreen() {
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			d.screen[x][y] = 0
		}
	}
}
