pub struct Memory {
    pub data: [u8; 4096],
}

impl Memory {
    pub fn new() -> Self {
        Self {
            data: [0; 4096],
        }
    }

    pub fn load_program(&mut self, program: &[u8]) {
        for (i, &byte) in program.iter().enumerate() {
            self.data[0x200 + i] = byte;
        }
    }
}
