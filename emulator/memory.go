package emulator

//memory.go

type memory struct {
	ram [4096]uint8
}

func newMemory() *memory {
	return &memory{
		ram: [4096]uint8{},
	}
}

func (m *memory) loadROM(rom []uint8) {
	for i, b := range rom {
		m.ram[0x200+i] = b
	}
}
