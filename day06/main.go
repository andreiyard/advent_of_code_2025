package main

import (
	"github.com/andreiyard/advent_of_code_2025/utils"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Implement Operator enum with String() method
type Operator int

const (
	OperatorSum Operator = iota
	OperatorMul
)

var OperatorLiteral = map[Operator]string{
	OperatorSum: "+",
	OperatorMul: "*",
}

func (o Operator) String() string {
	return OperatorLiteral[o]
}

func (o Operator) GetApplyFunction() func(int, int) int {
	if o == OperatorMul {
		return func(i1, i2 int) int {
			return i1 * i2
		}
	} else if o == OperatorSum {
		return func(i1, i2 int) int {
			return i1 + i2
		}
	} else {
		panic("not implemented!")
	}
}

func parse(data string) (numbers [][]int, operators []Operator) {
	// Returns numbers and operations (for each column)
	data = strings.TrimSpace(data)
	lines := strings.Split(data, "\n")
	nColumns := len(strings.Fields(lines[0]))
	numbers = make([][]int, nColumns)
	operators = make([]Operator, nColumns)

	// find where operators start
	operatorsStartLineNum := 0
	for i, line := range lines {
		if strings.ContainsAny(line, "+*") {
			operatorsStartLineNum = i
		}
	}

	numbersLines := lines[:operatorsStartLineNum]
	operatorsLines := lines[operatorsStartLineNum:]

	// for each column
	for n := range nColumns {
		numbersSlice := make([]int, len(numbersLines))

		for i, line := range numbersLines {
			num, err := strconv.Atoi(strings.Fields(line)[n])
			if err != nil {
				panic("invalid number")
			}
			numbersSlice[i] = num
		}
		numbers[n] = numbersSlice

		oper := strings.Fields(operatorsLines[0])[n]
		switch oper {
		case "+":
			operators[n] = OperatorSum
		case "*":
			operators[n] = OperatorMul
		default:
			panic("invalid operator")
		}

		slog.Debug("parsed column", "n", n, "numbers", numbersSlice, "operator", operators[n])
	}

	return
}

func parsePart2(data string) (numbers [][]int) {
	data = strings.TrimSpace(data)
	lines := strings.Split(data, "\n")
	nColumns := len(strings.Fields(lines[0]))
	numbers = make([][]int, nColumns)

	lines = lines[:len(lines)-1]

	// for each char column
	nColumn := 0
	for i := range len(lines[0]) {

		numString := ""
		for _, line := range lines {
			char := line[i]
			if char == ' ' {
				continue
			}
			numString += string(char)
		}
		// if no digits found -> increase nColumn
		// if found -> convert to int and append
		if numString == "" {
			nColumn += 1
		} else {
			num, err := strconv.Atoi(numString)
			if err != nil {
				panic("invalid num during parsing")
			}
			numbersSlice := numbers[nColumn]
			numbersSlice = append(numbersSlice, num)
			numbers[nColumn] = numbersSlice
		}
		slog.Debug("parsed char column", "i", i, "n", nColumn, "numbers", numbers[nColumn])
	}
	return
}

func reduce(numbers []int, op Operator) (result int) {
	apply := op.GetApplyFunction()
	if op == OperatorMul {
		result = 1
	}
	for _, num := range numbers {
		result = apply(result, num)
	}
	slog.Debug("reduced numbers", "numbers", numbers, "op", op, "res", result)
	return
}

func calculate(numbers [][]int, operator []Operator) int {
	sum := 0
	nColumns := len(numbers)

	for n := range nColumns {
		sum += reduce(numbers[n], operator[n])
	}

	return sum
}

func main() {
	utils.SetupLoggingEnv()
	filename := utils.GetFilenameFromArgs()
	dataBytes, err := os.ReadFile(filename)
	utils.Check(err)
	data := string(dataBytes)
	slog.Debug("Got data", "data", data)

	numbers, operators := parse(data)
	slog.Debug("Parsed input", "numbers", numbers, "operators", operators)
	result1 := calculate(numbers, operators)

	numbers = parsePart2(data)
	slog.Debug("Parsed input for part 2", "numbers", numbers, "operators", operators)
	result2 := calculate(numbers, operators)

	slog.Warn("part 1 result", "sum", result1)
	slog.Warn("part 2 result", "sum", result2)
}
