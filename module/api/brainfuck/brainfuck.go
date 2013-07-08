// brainfuck.go

// brainfuck interpreter
// All Hail Urban Müller
package brainfuck

import (
	"errors"
	"strings"
)

// https://en.wikipedia.org/wiki/Brainfuck
//
// >	increment the data pointer (to point to the next cell to the right).
// <	decrement the data pointer (to point to the next cell to the left).
// +	increment (increase by one) the byte at the data pointer.
// -	decrement (decrease by one) the byte at the data pointer.
// .	output the byte at the data pointer.
// ,	accept one byte of input, storing its value in the byte at the data pointer.
// [	if the byte at the data pointer is zero, then instead of moving the instruction pointer forward to the next command, jump it forward to the command after the matching ] command.
// ]	if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward to the next command, jump it back to the command after the matching [ command.
type BrainfuckInterpreter struct {
	memory []byte
	pos    int
	src    []string
	input  *strings.Reader
	loop   int
}

// Create a new BrainfuckInterpreter instance.
// Takes the the brainfuck source and an input string as parameters
func NewBrainfuckInterpreter(src string, input string) *BrainfuckInterpreter {
	return &BrainfuckInterpreter{
		memory: make([]byte, 256),
		pos:    0, src: strings.Split(src, ""),
		input: strings.NewReader(input)}
}

func (self *BrainfuckInterpreter) parseSource(srcPos int) ([]byte, error) {
	self.loop++
	if self.loop > 350000 {
		return nil, errors.New("Too many loops")
	}

	output := make([]byte, 0)

	for i := srcPos; i < len(self.src); i++ {
		switch self.src[i] {
		case ">":
			self.pos++
			if self.pos == len(self.memory) {
				self.memory = append(self.memory, 0)
			}
		case "<":
			self.pos--
			if self.pos < 0 {
				return nil, errors.New("negative postion")
			}
		case "+":
			self.memory[self.pos]++
		case "-":
			self.memory[self.pos]--
		case ".":
			output = append(output, self.memory[self.pos])
		case ",":
			// Leave cell untouched when input is EOF, like the original implementation of Urban Müller
			if b, err := self.input.ReadByte(); err == nil {
				self.memory[self.pos] = b
			}
		case "[":
			if self.memory[self.pos] != 0 {
				ret, err := self.parseSource(i + 1)
				if err != nil {
					return nil, err
				}
				output = append(output, ret...)
				i--
			} else {
				next, err := self.parseLoop(i)
				if err != nil {
					return nil, err
				}
				i = next
			}
		case "]":
			return output, nil
		}
	}

	return output, nil
}

func (self *BrainfuckInterpreter) parseLoop(srcPos int) (int, error) {
	if self.src[srcPos] != "[" {
		return 0, nil
	}

	depth := 0
	for i := srcPos; i < len(self.src); i++ {
		switch self.src[i] {
		case "[":
			depth++
		case "]":
			depth--
		}
		if depth == 0 {
			return i, nil
		}
	}

	return 0, errors.New("unclosed loop")
}

// Run the interpreter
// Returns the output of the program
func (self *BrainfuckInterpreter) Start() ([]byte, error) {
	return self.parseSource(0)
}
