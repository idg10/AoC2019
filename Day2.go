package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

// ProgramState stores the current state of execution
type ProgramState struct {
	IP int
}

// Opcode is implemented by all opcodes
type Opcode interface {
	Execute(state ProgramState, memory []int) (ProgramState, error)
}

// AddOpcode is operand type 1
type AddOpcode struct {
	Operand1Pos, Operand2Pos, OutputPos int
}

// Execute the AddOpcode operand
func (opcode AddOpcode) Execute(state ProgramState, memory []int) (ProgramState, error) {
	op1 := memory[opcode.Operand1Pos]
	op2 := memory[opcode.Operand2Pos]
	memory[opcode.OutputPos] = op1 + op2
	return ProgramState{state.IP + 4}, nil
}

// MultiplyOpcode is operand type 1
type MultiplyOpcode struct {
	Operand1Pos, Operand2Pos, OutputPos int
}

// Execute the MultiplyOpcode operand
func (opcode MultiplyOpcode) Execute(state ProgramState, memory []int) (ProgramState, error) {
	op1 := memory[opcode.Operand1Pos]
	op2 := memory[opcode.Operand2Pos]
	memory[opcode.OutputPos] = op1 * op2
	return ProgramState{state.IP + 4}, nil
}

// RunStep executes the next instruction
func RunStep(state ProgramState, program []int) (ProgramState, error) {
	var ip int = state.IP
	var opcode Opcode
	switch program[ip] {
	case 1:
		add := AddOpcode{program[ip+1], program[ip+2], program[ip+3]}
		opcode = add
		//fmt.Printf("IP: %v Add @%v (%v) to @%v (%v) -> %v\n", ip, add.Operand1Pos, program[add.Operand1Pos], add.Operand2Pos, program[add.Operand2Pos], add.OutputPos)
	case 2:
		mul := MultiplyOpcode{program[ip+1], program[ip+2], program[ip+3]}
		opcode = mul
		//fmt.Printf("IP: %v Mul @%v (%v) to @%v (%v) -> %v\n", ip, mul.Operand1Pos, program[mul.Operand1Pos], mul.Operand2Pos, program[mul.Operand2Pos], mul.OutputPos)
	case 99:
		return ProgramState{-1}, nil
	}
	return opcode.Execute(state, program)
}

// Run executes the program
func Run(program []int) error {
	state := ProgramState{0}
	var err error
	for state.IP >= 0 {
		state, err = RunStep(state, program)
		if err != nil {
			return err
		}
	}

	return nil
}

func readCommaSeparatedInts(text []byte) ([]int, error) {
	toInt := func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		n := nodes[0]
		t := n.(*parsec.Terminal)
		i, _ := strconv.Atoi(t.Value)
		//s = s[1 : len(s)-1]
		return i
	}

	yint := parsec.And(toInt, parsec.Int())
	ycomma := parsec.Token(`,`, "FIELDSEP")
	ycsv := parsec.Kleene(nil, yint, ycomma)

	s := parsec.NewScanner(text)
	node, _ := ycsv(s)
	nodes := node.([]parsec.ParsecNode)

	//fmt.Println(nodes)

	values := make([]int, len(nodes))
	for i, v := range nodes {
		values[i] = v.(int)
	}

	return values, nil
}

func day2(inputFolder string, arg1 int, arg2 int) (finalValueAtPositionZero int, err error) {
	fileBytes, err := ioutil.ReadFile(path.Join(inputFolder, "Day2.txt"))
	if err != nil {
		return
	}
	program, err := readCommaSeparatedInts(fileBytes)
	if err != nil {
		return
	}

	program[1] = arg1
	program[2] = arg2
	err = Run(program)
	if err != nil {
		return
	}

	return program[0], nil
}

func showDay2(inputFolder string) {
	finalValueAtPositionZero, err := day2(inputFolder, 12, 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Final value (part 1): ", finalValueAtPositionZero)
	}
	var noun, verb int
	noun, verb = 64, 21
	finalValueAtPositionZero, err = day2(inputFolder, noun, verb)
	fmt.Printf("%v, %v: %v\n", noun, verb, finalValueAtPositionZero)

	fmt.Printf("Part 2: %v\n", 100*noun+verb)
}
