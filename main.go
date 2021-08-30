package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sandbox/intcode"
	"strconv"
	"strings"
)

func main() {
	programA := loadProgram()
	programB := loadProgram()
	programC := loadProgram()
	programD := loadProgram()
	programE := loadProgram()

	arr := []int{5,6,7,8,9}
	phaseSettings := permutations(arr)

	maxValue := -10000000000

	ch0 := make(chan int)
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)

	answerChannel := make(chan int)

	// Part 1
	for _, phaseSetting := range phaseSettings {
		fmt.Println("Attempting phase setting", phaseSetting)
		ampA := intcode.IntCode{Program: programA, InChannel: ch0, OutChannel: ch1, Id: "A", AnswerChannel: answerChannel}
		ampB := intcode.IntCode{Program: programB, InChannel: ch1, OutChannel: ch2, Id: "B"}
		ampC := intcode.IntCode{Program: programC, InChannel: ch2, OutChannel: ch3, Id: "C"}
		ampD := intcode.IntCode{Program: programD, InChannel: ch3, OutChannel: ch4, Id: "D"}
		ampE := intcode.IntCode{Program: programE, InChannel: ch4, OutChannel: ch0, Id: "E"}

		go ampA.Run()
		go ampB.Run()
		go ampC.Run()
		go ampD.Run()
		go ampE.Run()

		ch0 <- phaseSetting[0]
		ch1 <- phaseSetting[1]
		ch2 <- phaseSetting[2]
		ch3 <- phaseSetting[3]
		ch4 <- phaseSetting[4]

		ch0 <- 0

		answer := <- answerChannel

		fmt.Println("Got", answer)

		if answer > maxValue {
			 maxValue = answer
		}
	}

	fmt.Println("Part 2:", maxValue)
}

func loadProgram() []int {
	result := loadFile(7)
	input := strings.Split(result[0], ",")

	var program [] int

	for _, x := range input {
		intValue, _ := strconv.Atoi(x)
		program = append(program, intValue)
	}

	return program
}

//func getOutput(signals []int, inVal int, program []int) int {
//	for i := 0; i < 5; i++ {
//		a := intcode.IntCode{Program: program}
//
//		a.Input = make([]int, 2)
//		a.Input[0] = signals[i]
//		a.Input[1] = inVal
//		a.Run()
//
//		// Set input value for next run
//		inVal = a.RetVal
//	}
//
//	return inVal
//}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func loadFile(day int) []string {
	b, err := ioutil.ReadFile(fmt.Sprintf("day%d.txt", day))

	if err != nil {
		log.Fatal(err)
	}

	str := string(b)
	return strings.Split(str, "\n")
}
