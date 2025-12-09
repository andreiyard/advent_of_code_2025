package main

import (
	"log/slog"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/andreiyard/advent_of_code_2025/utils"
)

type Point3D struct {
	x, y, z int
}

type Pair struct {
	p1, p2 Point3D
}

type Circuit map[Point3D]struct{}

func (p Pair) Dist() float64 {
	dx := float64(p.p2.x - p.p1.x)
	dy := float64(p.p2.y - p.p1.y)
	dz := float64(p.p2.z - p.p1.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func GetAllPairs(points []Point3D) (pairs []Pair) {
	for i, p1 := range points {
		for _, p2 := range points[i+1:] {
			pairs = append(pairs, Pair{p1, p2})
		}
	}
	return
}

func parse(data string) []Point3D {
	data = strings.TrimSpace(data)
	lines := strings.Split(data, "\n")
	points := make([]Point3D, len(lines))

	for i, line := range lines {
		parts := strings.Split(line, ",")
		x, err := strconv.Atoi(parts[0])
		utils.Check(err)
		y, err := strconv.Atoi(parts[1])
		utils.Check(err)
		z, err := strconv.Atoi(parts[2])
		utils.Check(err)
		points[i] = Point3D{x, y, z}
	}
	return points
}

func checkContains(circuit Circuit, point Point3D) (exists bool) {
	_, exists = circuit[point]
	return
}

func findCircuitContains(circuits []Circuit, point Point3D) (Circuit, bool) {
	for _, c := range circuits {
		if checkContains(c, point) {
			return c, true
		}
	}
	return Circuit{}, false
}

func part1(points []Point3D) int {
	sum := 1
	nShortest := 1000

	pairs := GetAllPairs(points)
	slog.Debug("Generate pairs", "pairs", pairs)

	// Sort all pairs based on len
	slices.SortFunc(pairs, func(p1, p2 Pair) int {
		return int(p1.Dist() - p2.Dist())
	})

	// Take only N shortest distance pairs
	slog.Debug("sorted pairs", "pairs", pairs)

	// Take slice of points that belong to nShortes pairs and create slice of slices [[p1] [p2] [p3] ...]
	// Iterate over pairs and merge slices of points that should be connected
	// At the end we will have the list of circuits which we can sort
	circuits := make([]Circuit, 0)
	n := 0
outer:
	for _, pair := range pairs {
		if n > nShortest {
			break
		}
		//slog.Debug("Current circuits", "circuits", circuits)
		// Check if one of pair points already in some circuit
		// if exists -> append other point
		// if doesnt -> create new circuit with both points
		for _, c := range circuits {
			// if already connected -> skip counter
			if checkContains(c, pair.p1) && checkContains(c, pair.p2) {
				continue outer
			}
			if checkContains(c, pair.p1) {
				// also check if any other circuit has other point
				otherCircuit, otherContains := findCircuitContains(circuits, pair.p2)
				if otherContains {
					// union circuits
					for k, _ := range otherCircuit {
						c[k] = struct{}{}
						delete(otherCircuit, k)
					}
				} else {
					c[pair.p2] = struct{}{}
				}
				n += 1
				continue outer
			}
			if checkContains(c, pair.p2) {
				otherCircuit, otherContains := findCircuitContains(circuits, pair.p1)
				if otherContains {
					// union circuits
					for k, _ := range otherCircuit {
						c[k] = struct{}{}
						delete(otherCircuit, k)
					}
				} else {
					c[pair.p1] = struct{}{}
				}
				n += 1
				continue outer
			}
		}

		// if wasn't found anywhere
		n += 1
		circuits = append(circuits, map[Point3D]struct{}{pair.p1: struct{}{}, pair.p2: struct{}{}})
	}

	// Sort circuits based on number of elementss (descending)
	slices.SortFunc(circuits, func(a, b Circuit) int {
		return len(b) - len(a)
	})
	for _, c := range circuits[:3] {
		sum *= len(c)
		slog.Debug("One of biggest circuits", "len", len(c), "circuit", c)
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

	points := parse(data)
	slog.Debug("Parsed input", "points", points)
	result1 := part1(points)

	//numbers = parsePart2(data)
	//slog.Debug("Parsed input for part 2", "numbers", numbers, "operators", operators)
	//result2 := calculate(numbers, operators)

	slog.Warn("part 1 result", "sum", result1)
	//slog.Warn("part 2 result", "sum", result2)
}
