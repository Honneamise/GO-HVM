package hvmlib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const (
	MEM_SIZE = 0xFFFF
	PC_BASE  = 0x800
	PC_MAX   = 0xFFFF
	SP_MAX   = 0xFF

	FL_ZERO = 0x00
	FL_POS  = 0x01
	FL_NEG  = 0x02
	FL_ERR  = 0xFF

	ST_RUN       = 0x01
	ST_SUSPENDED = 0x02
	ST_HALTED    = 0x03
)

type Hvm struct {
	ACC byte
	R   []byte
	MEM []byte

	PC uint16
	SP byte
	FL byte

	INPUT  bufio.Reader
	OUTPUT bufio.Writer

	BUS bufio.ReadWriter

	STATUS byte
}

type Fn func(h *Hvm, op uint8) error

var FunctionsMap = [...]Fn{
	0x00: Fn0,
	0x01: Fn1,
	0x02: Fn2,
	0x03: Fn3,
	0x04: Fn4,
	0x05: Fn5,
	0x06: Fn6,
	0x07: Fn7,
	0x08: Fn8,
	0x09: Fn9,
	0x0A: FnA,
	0x0B: FnB,
	0x0C: FnC,
	0x0D: FnD,
	0x0E: FnE,
	0x0F: FnF,
}

var InstructionSize = map[string]byte{
	"LDI": 2,
	"LDM": 3,
	"STR": 3,

	"MOV": 1,
	"SAV": 1,
	"ADD": 1,
	"SUB": 1,
	"MUL": 1,
	"DIV": 1,
	"AND": 1,
	"OR":  1,

	"INC":  1,
	"DEC":  1,
	"LSH":  1,
	"RSH":  1,
	"NEG":  1,
	"BNEG": 1,

	"BZ":   1,
	"BP":   1,
	"BN":   1,
	"CALL": 3,
	"RET":  1,
	"JMP":  3,

	"BUSR":   1,
	"BUSW":   1,
	"BUSCLR": 1,

	"NOP":  1,
	"HALT": 1,
}

var RegistersMap = map[string]byte{
	"R0":0x00,
	"R1":0x01,
	"R2":0x02,
	"R3":0x03,
	"R4":0x04,
	"R5":0x05,
	"R6":0x06,
	"R7":0x07,
	"R8":0x08,
	"R9":0x09,
	"RA":0x0A,
	"RB":0x0B,
	"RC":0x0C,
	"RD":0x0D,
	"RE":0x0E,
	"RF":0x0F,
}


/***************************/
/* VIRTUAL MACHINE SECTION */
/***************************/
func Create(file string) (*Hvm, error) {

	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	fi, err := f.Stat()

	if err != nil {
		return nil, err
	}

	if fi.Size() >= MEM_SIZE-PC_BASE {
		return nil, errors.New("file exceed max file size")
	}

	var h Hvm

	h.STATUS = ST_RUN

	h.PC = PC_BASE

	h.R = make([]byte, 16)

	h.MEM = make([]byte, MEM_SIZE)

	h.BUS = *bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))

	n, err := f.Read(h.MEM[PC_BASE:])

	if err != nil {
		return nil, err
	}

	if int64(n) != fi.Size() {
		return nil, errors.New("incomplete file read")
	}

	return &h, nil
}

/**********/
func (h *Hvm) Dump() {

	fmt.Println("===Hvm dump===")

	fmt.Printf("STATUS : %02X\n", h.STATUS)

	fmt.Printf("ACC : %02X\t", h.ACC)
	fmt.Printf("FL : %02X", h.FL)

	for i, v := range h.R {
		if i%4 == 0 {
			fmt.Println()
		}
		fmt.Printf("R[%X] : %02X\t", i, v)
	}

	fmt.Printf("\nPC : %04X\t", h.PC)
	fmt.Printf("SP : %02X\n", h.SP)
}

/**********/
func (h *Hvm) UpdateFlag() {

	switch {
	case h.ACC == 0:
		h.FL = FL_ZERO
	case (h.ACC >> 7) == 0:
		h.FL = FL_POS
	case (h.ACC >> 7) == 1:
		h.FL = FL_NEG
	default:
		h.FL = FL_ERR
	}
}

/**********/
func (h *Hvm) Execute() error {

	for {

		//run
		if h.STATUS == ST_RUN {

			op := h.MEM[h.PC]

			h.PC++

			err := h.InsExecute(op)

			if err != nil {
				return err
			}

			h.UpdateFlag()

			continue
		}

		//suspended
		if h.STATUS == ST_SUSPENDED {

			continue

		}

		//halted
		if h.STATUS == ST_HALTED {

			break

		}
	}

	return nil
}

/***********************/
/* INSTRUCTIONS SECTION*/
/***********************/
func (h *Hvm) InsExecute(op uint8) error {

	fn := FunctionsMap[op>>4]

	if fn == nil {
		return errors.New("invalid opcode")
	}

	return fn(h, op)
}

/**********/
func Fn0(h *Hvm, op uint8) error {
	//R[l] <- PC PC++

	h.R[op&0x0F] = h.MEM[h.PC]
	h.PC++

	return nil
}

/**********/
func Fn1(h *Hvm, op uint8) error {
	//R[l] <- MEM[PC,PC+1]; PC+=2

	add := uint16(h.MEM[h.PC])<<8 + uint16(h.MEM[h.PC+1])

	h.R[op&0x0F] = h.MEM[add]

	h.PC += 2

	return nil
}

/**********/
func Fn2(h *Hvm, op uint8) error {
	//MEM[PC,PC+1] <- R[l]; PC+=2

	add := uint16(h.MEM[h.PC])<<8 + uint16(h.MEM[h.PC+1])

	h.MEM[add] = h.R[op&0x0F]

	h.PC += 2

	return nil
}

/**********/
func Fn3(h *Hvm, op uint8) error {
	//ACC <- R[l]

	h.ACC = h.R[op&0x0F]

	return nil
}

/**********/
func Fn4(h *Hvm, op uint8) error {
	//R[l] <- ACC

	h.R[op&0x0F] = h.ACC

	return nil
}

/**********/
func Fn5(h *Hvm, op uint8) error {
	//ACC <- ACC + R[l]

	h.ACC = h.ACC + h.R[op&0x0F]

	return nil
}

/**********/
func Fn6(h *Hvm, op uint8) error {
	//ACC <- ACC - R[l]

	h.ACC = h.ACC - h.R[op&0x0F]

	return nil
}

/**********/
func Fn7(h *Hvm, op uint8) error {
	//ACC <- ACC * R[l]

	h.ACC = h.ACC * h.R[op&0x0F]

	return nil
}

/**********/
func Fn8(h *Hvm, op uint8) error {
	//ACC <- ACC / R[l]

	h.ACC = h.ACC / h.R[op&0x0F]

	return nil
}

/**********/
func Fn9(h *Hvm, op uint8) error {
	//ACC <- ACC & R[l]

	h.ACC = h.ACC & h.R[op&0x0F]

	return nil
}

/**********/
func FnA(h *Hvm, op uint8) error {
	//ACC <- ACC | R[l]

	h.ACC = h.ACC | h.R[op&0x0F]

	return nil
}

/**********/
func FnB(h *Hvm, op uint8) error {

	switch op {

	case 0xB0:
		//ACC++

		h.ACC++

	case 0xB1:
		//ACC--

		h.ACC--

	case 0xB2:
		//ACC <- ACC<<1

		h.ACC = h.ACC << 1

	case 0xB3:
		//ACC <- ACC>>1

		h.ACC = h.ACC >> 1

	case 0xB4:
		//ACC <- -ACC

		h.ACC = -h.ACC

	case 0xB5:
		//ACC <- ~ACC

		h.ACC = ^h.ACC

	default:
		return errors.New("invalid opcode")

	}

	return nil
}

/**********/
func FnC(h *Hvm, op uint8) error {

	switch op {

	case 0xC0:
		//if FL==FL_ZERO then PC++

		if h.FL == FL_ZERO {
			h.PC++
		}

	case 0xC1:
		//if FL==FL_POS then PC++

		if h.FL == FL_POS {
			h.PC++
		}

	case 0xC2:
		//if FL==FL_NEG then PC++

		if h.FL == FL_NEG {
			h.PC++
		}

	case 0xC3:
		//MEM[SP] <- PC[high byte]; MEM[SP+1] <- PC[low byte]; SP+=2
		//PC[high byte] <- MEM[PC]; PC[low byte] <- MEM[PC+1]
		if h.SP == 0xFF {
			return errors.New("sp out of range")
		}

		h.MEM[h.SP] = byte(h.PC>>8)
		h.MEM[h.SP+1] = byte(h.PC&0x00FF)
		h.SP += 2

		h.PC = uint16(h.MEM[h.PC])<<8  + uint16(h.MEM[h.PC+1])

	case 0xC4:
		//PC[high byte] <- MEM[SP-2]; PC[low byte] <- MEM[SP-1]; SP-=2

		if h.SP == 0x00 {
			return errors.New("sp out of range")
		}

		h.PC = uint16(h.MEM[h.SP-2])<<8 +uint16(h.MEM[h.SP-1])
		h.SP -= 2

		h.PC +=2
		
	case 0xC5:
		//PC[high byte] <- MEM[PC]; PC[low byte] <- MEM[PC+1]
		if h.SP == 0xFF {
			return errors.New("sp out of range")
		}

		h.PC = uint16(h.MEM[h.PC])<<8  + uint16(h.MEM[h.PC+1])

	default:
		return errors.New("invalid opcode")

	}

	return nil
}

/**********/
func FnD(h *Hvm, op uint8) error {

	switch op {

	case 0xD0:
		//ACC <- BUS
		c, err := h.BUS.ReadByte()

		if err != nil {
			return err
		}

		h.ACC = c

	case 0xD1:
		//BUS <- ACC
		h.BUS.WriteByte(h.ACC)
		h.BUS.Flush()

	case 0xD2:
		//BUS <- ACC
		h.BUS.Reader.Reset(bufio.NewReader(os.Stdin))
		h.BUS.Writer.Reset(bufio.NewWriter(os.Stdout))

	default:
		return errors.New("invalid opcode")

	}

	return nil
}

/**********/
func FnE(h *Hvm, op uint8) error {

	switch op {

	case 0xE0:
		//No op

	case 0xE1:
		//STATUS = ST_HALTED
		h.STATUS = ST_HALTED

	default:
		return errors.New("invalid opcode")

	}

	return nil
}

/**********/
func FnF(h *Hvm, op uint8) error {

	return errors.New("reserved opcode")

}
