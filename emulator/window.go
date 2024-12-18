package emulator

import (
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type window struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func newWindow() *window {
	if err := sdl.Init(uint32(sdl.INIT_EVERYTHING)); err != nil {
		log.Fatalf("Falha ao inicializar SDL: %s", err)
	}

	win, err := sdl.CreateWindow("CHIP-8 Emulator", int32(sdl.WINDOWPOS_UNDEFINED), int32(sdl.WINDOWPOS_UNDEFINED), 640, 320, uint32(sdl.WINDOW_SHOWN))
	if err != nil {
		log.Fatalf("Falha ao criar janela: %s", err)
	}

	renderer, err := sdl.CreateRenderer(win, -1, uint32(sdl.RENDERER_ACCELERATED))
	if err != nil {
		log.Fatalf("Falha ao criar renderizador: %s", err)
	}

	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA8888), int(sdl.TEXTUREACCESS_STREAMING), 64, 32)
	if err != nil {
		log.Fatalf("Falha ao criar textura: %s", err)
	}

	return &window{
		window:   win,
		renderer: renderer,
		texture:  texture,
	}
}

func (w *window) update(display *display) {
	pixels := make([]byte, 64*32*4)
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			index := (y*64 + x) * 4
			if display.screen[x][y] == 1 {
				pixels[index] = 255   // R
				pixels[index+1] = 255 // G
				pixels[index+2] = 255 // B
				pixels[index+3] = 255 // A
			} else {
				pixels[index] = 0     // R
				pixels[index+1] = 0   // G
				pixels[index+2] = 0   // B
				pixels[index+3] = 255 // A
			}
		}
	}

	w.texture.Update(nil, unsafe.Pointer(&pixels[0]), 64*4)
	w.renderer.Clear()
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

func (w *window) destroy() {
	w.texture.Destroy()
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}
