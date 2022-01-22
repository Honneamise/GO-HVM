package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"strconv"
	"../hvmlib"
)

const (
	TITLE = "(\\_/)\n(-.-)\n(\")(\")\nHVM : THE HONNY VIRTUAL MACHINE COMPILER\nVer. 1.0"
	USAGE = "USAGE : hvmc.bin src_file dst_file"
)

var LabelsAddress = map[string]uint16{}

/**********/
func compile(lines []string) error {

	var ADC uint16 = hvmlib.PC_BASE

	for count, line := range lines {

		tokenizer := bufio.NewScanner(strings.NewReader(line))

		tokenizer.Split(bufio.ScanWords)

		for tokenizer.Scan() {

			token := tokenizer.Text()

			if token[0] == '#' { //comment
				for tokenizer.Scan() {
				}

			} else if token[len(token)-1] == ':' { //label
				_, exist := LabelsAddress[token[:len(token)-1]]
				if !exist {
					LabelsAddress[token[:len(token)-1]] = ADC
				} else {
					str := fmt.Sprintf("line:%d token:\"%s\" duplicate label", count+1, token)
					return errors.New(str)
				}

			} else if val, exist := hvmlib.InstructionSize[token]; exist { //instruction

				ADC += uint16(val)

				for tokenizer.Scan() {
				}

			} else { //invalid token

				str := fmt.Sprintf("line:%d token:\"%s\" unahandled token", count+1, token)
				return errors.New(str)
			}
		}

		err := tokenizer.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

/**********/
func assembleInstruction(token string, tokenizer *bufio.Scanner, writer *bufio.Writer) error {

	switch token {

	case "LDI":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x00+reg)
		} else {
			return errors.New("invalid register")
		}

		if !tokenizer.Scan() {
			return errors.New("missing value")
		}

		val,err := strconv.ParseUint(tokenizer.Text(), 16, 8)
		if err!=nil {
			return err
		}
		writer.WriteByte(byte(val))

	case "LDM":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x10+reg)
		} else {
			return errors.New("invalid register")
		}

		if !tokenizer.Scan() {
			return errors.New("missing value")
		}

		val,err := strconv.ParseUint(tokenizer.Text(), 16, 16)
		if err!=nil {
			return err
		}
		writer.WriteByte(byte(val>>8))
		writer.WriteByte(byte(val&0x00FF))
	
	case "STR":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x20+reg)
		} else {
			return errors.New("invalid register")
		}

		if !tokenizer.Scan() {
			return errors.New("missing value")
		}

		val,err := strconv.ParseUint(tokenizer.Text(), 16, 16)
		if err!=nil {
			return err
		}
		writer.WriteByte(byte(val>>8))
		writer.WriteByte(byte(val&0x00FF))

	case "MOV":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x30+reg)
		} else {
			return errors.New("invalid register")
		}

	case "SAV":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x40+reg)
		} else {
			return errors.New("invalid register")
		}

	case "ADD":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x50+reg)
		} else {
			return errors.New("invalid register")
		}

	case "SUB":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x60+reg)
		} else {
			return errors.New("invalid register")
		}

	case "MUL":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x70+reg)
		} else {
			return errors.New("invalid register")
		}

	case "DIV":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x80+reg)
		} else {
			return errors.New("invalid register")
		}

	case "AND":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0x90+reg)
		} else {
			return errors.New("invalid register")
		}

	case "OR":
		if !tokenizer.Scan() {
			return errors.New("missing register")
		}

		regname := tokenizer.Text()

		if reg,exist := hvmlib.RegistersMap[regname]; exist {
			writer.WriteByte(0xA0+reg)
		} else {
			return errors.New("invalid register")
		}

	case "INC":
		writer.WriteByte(0xB0)

	case "DEC":
		writer.WriteByte(0xB1)

	case "LSH":
		writer.WriteByte(0xB2)

	case "RSH":
		writer.WriteByte(0xB3)

	case "NEG":
		writer.WriteByte(0xB4)

	case "BNEG":
		writer.WriteByte(0xB5)

	case "BZ":
		writer.WriteByte(0xC0)

	case "BP":
		writer.WriteByte(0xC1)

	case "BN":
		writer.WriteByte(0xC2)

	case "CALL":
		if !tokenizer.Scan() {
			return errors.New("missing label parameter")
		}

		if add,exist := LabelsAddress[tokenizer.Text()]; exist {
			writer.WriteByte(0xC3)
			writer.WriteByte(byte(add>>8))
			writer.WriteByte(byte(add&0x00FF))
		} else {
			return errors.New("label not exist")
		}

	case "RET":
		writer.WriteByte(0xC4)

	case "JMP":
		if !tokenizer.Scan() {
			return errors.New("missing label parameter")
		}

		if add,exist := LabelsAddress[tokenizer.Text()]; exist {
			writer.WriteByte(0xC5)
			writer.WriteByte(byte(add>>8))
			writer.WriteByte(byte(add&0x00FF))
		} else {
			return errors.New("label not exist")
		}

	case "BUSR":
		writer.WriteByte(0xD0)

	case "BUSW":
		writer.WriteByte(0xD1)

	case "BUSCLR":
		writer.WriteByte(0xD2)

	case "NOP":
		writer.WriteByte(0xE0)

	case "HALT":
		writer.WriteByte(0xE1)

	default:
		return fmt.Errorf("unknow instruction %s", token)
	}

	writer.Flush()

	return nil
}

/**********/
func assemble(lines []string, file string) error {

	f, err := os.Create(file)

	if err != nil {
		return err
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	
	for count, line := range lines {

		tokenizer := bufio.NewScanner(strings.NewReader(line))

		tokenizer.Split(bufio.ScanWords)

		for tokenizer.Scan() {

			token := tokenizer.Text()

			if token[0] == '#' { //comment
				for tokenizer.Scan() {
				}

			} else if token[len(token)-1] == ':' { //label
				
				continue

			} else { //instruction

				err := assembleInstruction(token, tokenizer, writer)

				if err != nil {
					_err := fmt.Errorf("line:%d %w", count+1, err)
					return _err
				}

			} 

		}

		err := tokenizer.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

/**********/
func main() {

	fmt.Println(TITLE)

	if len(os.Args) != 3 {
		fmt.Println(USAGE)
		return
	}

	var lines []string
	
	f, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = compile(lines)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = assemble(lines, os.Args[2])

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("COMPILED SUCCESSFULL !!!")
}
