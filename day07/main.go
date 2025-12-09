package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/andreiyard/advent_of_code_2025/utils"
)

type Point struct {
	x, y int
}

type Cells map[Point]string

func (c Cells) getActiveBeams() []Point {
	var points []Point
	for k, v := range c {
		if v == "S" {
			points = append(points, k)
		}
	}
	return points
}

func (c Cells) printGrid() {
	var maxX, maxY int
	for k, _ := range c {
		if k.x > maxX {
			maxX = k.x
		}
		if k.y > maxY {
			maxY = k.y
		}
	}

	for i := range maxX + 1 {
		for j := range maxY + 1 {
			val := c[Point{i, j}]
			switch val {
			case "":
				fmt.Print(".")
			default:
				fmt.Print(val)
			}
		}
		fmt.Print("\n")
	}
}

func (c Cells) extendBeam(start Point, limitX int) (splitted bool) {
	// Start for point and check next x point until end of the grid or split
	for x := start.x; x < limitX; x++ {
		point := Point{x, start.y}
		val := c[point]
		if val == "^" {
			// add new beam points to the grid because split
			left := Point{x, start.y - 1}
			right := Point{x, start.y + 1}
			c[left] = "S"
			c[right] = "S"
			// mark splitted
			splitted = true
			c[point] = "V"
			break
		} else if val == "V" {
			// if already was used to split beam do not count
			// only create splitted beams
			left := Point{x, start.y - 1}
			right := Point{x, start.y + 1}
			c[left] = "S"
			c[right] = "S"
		} else {
			c[point] = "|"
		}
	}
	// remove starting point from grid
	c[start] = "|"
	return
}

func parse(data string) (Cells, int) {
	data = strings.TrimSpace(data)
	lines := strings.Split(data, "\n")
	cells := make(Cells)

	for i, line := range lines {
		for j, char := range line {
			point := Point{i, j}
			if char == '.' {
				continue
			}
			cells[point] = string(char)
		}
	}

	return cells, len(lines)
}

func part1(cells Cells, limitX int) int {
	sum := 0

	// while grid contains beams starts
	for {
		cells.printGrid()
		nextBeams := cells.getActiveBeams()
		slog.Debug("Next beams to check", "beamsStarts", nextBeams)
		if len(nextBeams) == 0 {
			break
		}
		for _, beamStart := range nextBeams {
			if cells.extendBeam(beamStart, limitX) {
				sum++
				slog.Debug("Beam splitted", "beam", beamStart, "sum", sum)
			} else {
				slog.Debug("Beam NOT splitted", "beam", beamStart, "sum", sum)
			}
		}
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

	cells, limitX := parse(data)
	slog.Debug("Parsed input", "cells", cells, "limitX", limitX)
	result1 := part1(cells, limitX)

	//numbers = parsePart2(data)
	//slog.Debug("Parsed input for part 2", "numbers", numbers, "operators", operators)
	//result2 := calculate(numbers, operators)

	slog.Warn("part 1 result", "sum", result1)
	//slog.Warn("part 2 result", "sum", result2)
}
