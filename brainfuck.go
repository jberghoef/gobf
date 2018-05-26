package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

// Brainfuck ...
type Brainfuck struct {
	program []string
	cells   []uint64
	cell    *uint64
	pos     int
	step    int
	output  string
}

// BrainfuckConstructor will return a Brainfuck struct allowing
// execution of Brainfuck code
func BrainfuckConstructor(program string) (b Brainfuck) {
	minified := syntax.ReplaceAllString(string(program), "")

	b.program = strings.Split(minified, "")
	b.cells = make([]uint64, 1)
	b.cell = &b.cells[b.pos]

	return
}

// Execute the parsed Brainfuck code.
func (b *Brainfuck) Execute() string {
	writer := uilive.New()
	writer.Start()

	for b.step < len(b.program) {
		b.handle(b.program[b.step])
		b.cell = &b.cells[b.pos]

		fmt.Fprintf(writer, "%s\n%s\n%s\n",
			b.PrettyCellState(), b.PrettyProgramState(), b.output)
		time.Sleep(time.Millisecond * time.Duration(delay))
	}

	writer.Stop()
	return b.output
}

func (b *Brainfuck) handle(action string) {
	switch action {
	case ">":
		b.incrementPosition()
	case "<":
		b.decrementPosition()
	case "+":
		b.incrementCell()
	case "-":
		b.decrementCell()
	case ".":
		b.outputValue()
	case ",":
		b.inputValue()
	case "[":
		if *b.cell == 0 {
			for i := b.step; i < len(b.program); i++ {
				if b.program[i] == "]" {
					b.step = i + 1
					break
				}
			}
		}
	case "]":
		if *b.cell > 0 {
			for i := b.step; i > 0; i-- {
				if b.program[i] == "[" {
					b.step = i
					break
				}
			}
		}
	default:
		fmt.Printf("'%s' is invalid syntax.", action)
	}

	b.step = b.step + 1
}

func (b *Brainfuck) incrementPosition() {
	b.pos++
	if b.pos >= len(b.cells) {
		b.cells = append(b.cells, 0)
	}
}

func (b *Brainfuck) decrementPosition() {
	b.pos--
	if b.pos < 0 {
		panic("Invalid syntax, negative pointer.")
	}
}

func (b *Brainfuck) incrementCell() {
	*b.cell++
}

func (b *Brainfuck) decrementCell() {
	*b.cell--
}

func (b *Brainfuck) outputValue() {
	var buf bytes.Buffer
	buf.WriteString(b.output)
	buf.WriteString(string(*b.cell))
	b.output = buf.String()
}

func (b *Brainfuck) inputValue() {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	input, err := buf.ReadBytes('\n')

	if err != nil {
		panic(err)
	} else {
		*b.cell = binary.BigEndian.Uint64(input)
	}
}

// PrettyCellState will return a pretty formatted table containing
// all the cell states at the time of calling
func (b Brainfuck) PrettyCellState() string {
	var buf bytes.Buffer
	red := color.New(color.FgRed).SprintFunc()

	for i, cell := range b.cells {
		str := strings.TrimSpace(string(cell))

		if i == b.pos {
			buf.WriteString(red(str))
		} else {
			buf.WriteString(str)
		}

		buf.WriteString(strings.Repeat(" ", 8-len(str)))
	}

	buf.WriteString("\n")

	for i, cell := range b.cells {
		str := fmt.Sprint(cell)

		if i == b.pos {
			buf.WriteString(red(str))
		} else {
			buf.WriteString(str)
		}

		buf.WriteString(strings.Repeat(" ", 8-len(str)))
	}

	return buf.String()
}

// PrettyProgramState will return a pretty formatted line of text
// indicating the position in the program
func (b Brainfuck) PrettyProgramState() string {
	var buf bytes.Buffer
	red := color.New(color.FgRed).SprintFunc()

	for i, step := range b.program {
		if i == b.step {
			buf.WriteString(red(step))
		} else {
			buf.WriteString(step)
		}
	}

	return buf.String()
}
