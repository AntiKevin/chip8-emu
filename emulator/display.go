package emulator

import "fmt"

type display struct {
	screen [64][32]byte // Tela de 64x32 pixels
}

func newDisplay() *display {
	return &display{
		screen: [64][32]byte{},
	}
}

// drawSprite desenha um sprite na tela e retorna true se houve colisão
func (d *display) drawSprite(x, y, height uint16, I uint16, memory []uint8) bool {
	collision := false

	fmt.Printf("Desenhando sprite em x: %d, y: %d, height: %d, I: %d\n", x, y, height, I)

	for yline := uint16(0); yline < height; yline++ {
		pixel := memory[I+yline]
		for xline := uint16(0); xline < 8; xline++ {
			if (pixel & (0x80 >> xline)) != 0 {
				if d.screen[(x+uint16(xline))%64][(y+yline)%32] == 1 {
					collision = true
				}
				d.screen[(x+uint16(xline))%64][(y+yline)%32] ^= 1
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
