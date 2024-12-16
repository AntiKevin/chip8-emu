package emulator

//keypad.go

type keypad struct {
	keys [16]bool
}

func (k *keypad) waitForKeypress() byte {
	for {
		for i, pressed := range k.keys {
			if pressed {
				return byte(i)
			}
		}
	}
}
