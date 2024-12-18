package emulator

//keypad.go

type keypad struct {
	keys [16]bool
}

func newKeypad() *keypad {
	return &keypad{
		keys: [16]bool{},
	}
}

func (k *keypad) isKeyPressed(key byte) bool {
	return k.keys[key]
}

func (k *keypad) waitForKeyPress() byte {
	for {
		for i, pressed := range k.keys {
			if pressed {
				return byte(i)
			}
		}
	}
}
