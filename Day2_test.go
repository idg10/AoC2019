package main

import (
	"reflect"
	"testing"
)

func TestOpcodes(t *testing.T) {
	cases := []struct {
		op         Opcode
		programIn  []int
		stateIn    ProgramState
		programOut []int
		stateOut   ProgramState
	}{
		// Adds

		// Simple example
		{
			AddOpcode{4, 5, 6},
			[]int{1, 4, 5, 6, 111, 222, 0},
			ProgramState{0},
			[]int{1, 4, 5, 6, 111, 222, 333},
			ProgramState{4},
		},
		// Simple example at non-zero offset
		{
			AddOpcode{5, 6, 7},
			[]int{0, 1, 5, 6, 7, 111, 222, 0},
			ProgramState{1},
			[]int{0, 1, 5, 6, 7, 111, 222, 333},
			ProgramState{5},
		},
		// First example in AoC
		{
			AddOpcode{10, 20, 30},
			append(append(append(append(make([]int, 10), 42), make([]int, 9)...), 99), make([]int, 10)...),
			ProgramState{0},
			append(append(append(append(append(make([]int, 10), 42), make([]int, 9)...), 99), make([]int, 9)...), 141),
			ProgramState{4},
		},
		// Add in longer example in AoC
		{
			AddOpcode{9, 10, 3},
			[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			ProgramState{0},
			[]int{1, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
			ProgramState{4},
		},
		// Add in 2nd longer example in AoC
		{
			AddOpcode{0, 0, 0},
			[]int{1, 0, 0, 0, 99},
			ProgramState{0},
			[]int{2, 0, 0, 0, 99},
			ProgramState{4},
		},
		// Add in final longer example in AoC
		{
			AddOpcode{1, 1, 4},
			[]int{1, 1, 1, 4, 99, 5, 6, 0, 99},
			ProgramState{0},
			[]int{1, 1, 1, 4, 2, 5, 6, 0, 99},
			ProgramState{4},
		},

		// Multiply

		// Simple example
		{
			MultiplyOpcode{4, 5, 6},
			[]int{1, 4, 5, 6, 111, 2, 0},
			ProgramState{0},
			[]int{1, 4, 5, 6, 111, 2, 222},
			ProgramState{4},
		},
		// Simple example at non-zero offset
		{
			MultiplyOpcode{5, 6, 7},
			[]int{0, 1, 5, 6, 7, 111, 2, 0},
			ProgramState{1},
			[]int{0, 1, 5, 6, 7, 111, 2, 222},
			ProgramState{5},
		},
		// Multiply in longer example in AoC
		{
			MultiplyOpcode{3, 11, 0},
			[]int{1, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
			ProgramState{4},
			[]int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50},
			ProgramState{8},
		},
		// Multiply in 2nd longer example in AoC
		{
			MultiplyOpcode{3, 0, 3},
			[]int{2, 3, 0, 3, 99},
			ProgramState{0},
			[]int{2, 3, 0, 6, 99},
			ProgramState{4},
		},
		// Multiply in 3rd longer example in AoC
		{
			MultiplyOpcode{4, 4, 5},
			[]int{2, 4, 4, 5, 99, 0},
			ProgramState{0},
			[]int{2, 4, 4, 5, 99, 9801},
			ProgramState{4},
		},
		// Multiply in final longer example in AoC
		{
			MultiplyOpcode{5, 6, 0},
			[]int{1, 1, 1, 4, 2, 5, 6, 0, 99},
			ProgramState{4},
			[]int{30, 1, 1, 4, 2, 5, 6, 0, 99},
			ProgramState{8},
		},
	}

	for _, c := range cases {
		programCopy := make([]int, len(c.programIn))
		copy(programCopy, c.programIn)
		stateOut, _ := c.op.Execute(c.stateIn, programCopy)
		if stateOut != c.stateOut {
			t.Errorf("Execute(%v) state out %v, want %v", c.op, stateOut, c.stateOut)
		}
		if !reflect.DeepEqual(c.programOut, programCopy) {
			t.Errorf("Execute(%v) program after %v, want %v", c.op, programCopy, c.programOut)
		}
	}
}

func TestRun(t *testing.T) {
	cases := []struct {
		programIn  []int
		programOut []int
	}{
		{
			[]int{1, 0, 0, 0, 99},
			[]int{2, 0, 0, 0, 99},
		},
	}

	for _, c := range cases {
		programCopy := make([]int, len(c.programIn))
		copy(programCopy, c.programIn)
		_ = Run(programCopy)
		if !reflect.DeepEqual(c.programOut, programCopy) {
			t.Errorf("Execute program %v, after %v, want %v", c.programOut, programCopy, c.programOut)
		}
	}
}

// Part 1: 4090701
// Part 2: 6421

func TestDay2(t *testing.T) {
	finalValueAtPositionZero, _ := day2("IdgInput", 12, 1)
	if finalValueAtPositionZero != 4090700 {
		t.Errorf("Part 1: %v, want %v", finalValueAtPositionZero, 4090700)
	}
	finalValueAtPositionZero, _ = day2("IdgInput", 64, 21)
	if finalValueAtPositionZero != 19690720 {
		t.Errorf("Part 1: %v, want %v", finalValueAtPositionZero, 4090700)
	}
}
