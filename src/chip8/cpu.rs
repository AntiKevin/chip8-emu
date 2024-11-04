pub struct CPU {
    pub v: [u8; 16],         // 16 registradores de 8 bits (V0 a VF)
    pub i: u16,               // Registrador de endereço (16 bits)
    pub pc: u16,              // Program Counter (PC), aponta para a próxima instrução
    pub stack: [u16; 16],     // Pilha de 16 níveis para sub-rotinas
    pub sp: u8,               // Stack Pointer (SP), aponta para o topo da pilha
    pub delay_timer: u8,      // Timer de atraso, decrementa a uma taxa fixa
    pub sound_timer: u8,      // Timer de som, emite som enquanto for maior que 0
}

impl CPU {
    pub fn new() -> Self {
        Self {
            v: [0; 16],
            i: 0,
            pc: 0x200,
            stack: [0; 16],
            sp: 0,
            delay_timer: 0,
            sound_timer: 0,
        }
    }

    pub fn fetch_opcode(&mut self, memory: &[u8]) -> u16 {
        let high_byte: u16 = memory[self.pc as usize] as u16;
        let low_byte: u16 = memory[(self.pc + 1) as usize] as u16;
        let opcode: u16 = (high_byte << 8) | low_byte;
        
        self.pc += 2;

        opcode
    }
}