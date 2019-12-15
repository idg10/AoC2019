package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

func fuelForMass(mass int) int {
	massOverThree := mass / 3
	return massOverThree - 2
}

func fuelForMassAndFuel(mass int) int {
	total := 0
	for mass > 0 {
		mass = fuelForMass(mass)
		if mass > 0 {
			total += mass
		}
	}

	return total
}

func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}

func day1(inputFolder string) (totalFuelIfFuelFree int, totalFuel int, err error) {
	var r io.Reader
	r, err = os.Open(path.Join(inputFolder, "Day1.txt"))
	if err != nil {
		return
	}
	masses, err := readInts(r)
	if err != nil {
		return
	}

	totalFuelIfFuelFree = 0
	totalFuel = 0
	for _, mass := range masses {
		totalFuelIfFuelFree += fuelForMass(mass)
		totalFuel += fuelForMassAndFuel(mass)
	}

	return
}

func showDay1(inputFolder string) {
	totalFuelIfFuelFree, totalFuel, err := day1(inputFolder)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Total fuel (part 1): ", totalFuelIfFuelFree)
		fmt.Println("Total fuel (part 2): ", totalFuel)
	}
}
