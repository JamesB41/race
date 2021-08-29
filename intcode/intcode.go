package intcode

import (
	"fmt"
	"strconv"
	"time"
)


type IntCode struct {
	Program [] int
	RetVal int
	InChannel chan int
	OutChannel chan int
	AnswerChannel chan int
	DoneChannel chan bool
	Id string
	done bool
}

func (i *IntCode) LoadProgram (p []string) {
	for j := 0; j < len(p); j++ {
		i.Program = append(i.Program, atoi(p[j]))
	}
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (i *IntCode) getParam(inst string, programPos int, paramPos int) int {
	var mode string

	switch paramPos {
	case 1:
		mode = string(inst[2])
	case 2:
		mode = string(inst[1])
	case 3:
		mode = string(inst[0])
	}

	if mode == "1" {
		// Immediate Mode
		return i.Program[programPos+paramPos]
	} else {
		// Position Mode
		return i.Program[i.Program[programPos+paramPos]]
	}
}

func getOpcode(inst string) int {
	opcode := inst[len(inst) - 2:]
	return atoi(opcode)
}

func (i *IntCode) Run() {
	pos := 0

	if i.Id == "E" {
		go func() {
			<- i.DoneChannel
			fmt.Println("E knows we're done!")
			i.done = true
		}()
	}

	for {
		instruction := fmt.Sprintf("%05v", strconv.Itoa(i.Program[pos]))
		opcode := getOpcode(instruction)

		switch opcode {
		case 1:
			x := i.getParam(instruction, pos, 1)
			y := i.getParam(instruction, pos, 2)
			z := i.Program[pos + 3]

			i.add(x, y, z)
			pos += 4
		case 2:
			x := i.getParam(instruction, pos, 1)
			y := i.getParam(instruction, pos, 2)
			z := i.Program[pos + 3]

			i.mult(x, y, z)
			pos += 4
		case 3:
			x := i.Program[pos + 1]
			i.inp(x)
			pos += 2
		case 4:
			x := i.getParam(instruction, pos, 1)
			i.out(x)
			pos += 2
		case 5:
			x := i.getParam(instruction, pos, 1)
			y := i.getParam(instruction, pos, 2)

			if x != 0 {
				pos = y
			} else {
				pos += 3
			}
		case 6:
			x := i.getParam(instruction, pos, 1)
			y := i.getParam(instruction, pos, 2)

			if x == 0 {
				pos = y
			} else {
				pos += 3
			}
		case 7:
			x := i.getParam(instruction, pos, 1)
			y := i.getParam(instruction, pos, 2)
			z := i.Program[pos + 3]

			if x < y {
				i.Program[z] = 1
			} else {
				i.Program[z] = 0
			}
			pos += 4
		case 8:
			x := i.getParam(instruction, pos, 1)
			y := i.getParam(instruction, pos, 2)
			z := i.Program[pos + 3]

			if x == y {
				i.Program[z] = 1
			} else {
				i.Program[z] = 0
			}
			pos += 4

		case 99:
			fmt.Println("halting", i.Id, time.Now().String())

			if i.Id == "A" {
				i.DoneChannel <- true
			}

			return
		}
	}
}

func (i *IntCode) add(x int, y int, z int) {
	i.Program[z] = x + y
}

func (i *IntCode) mult(x int, y int, z int) {
	i.Program[z] = x * y
}

func (i *IntCode) inp(x int) {
	input := <- i.InChannel
	i.Program[x] = input
}

func (i *IntCode) out(x int) {
	if i.done {
		fmt.Println("Outputting when done from", i.Id)
	}

	if i.done && i.Id == "E" {
		i.AnswerChannel <- x
	} else {
		i.OutChannel <- x
	}
	i.RetVal = x
}