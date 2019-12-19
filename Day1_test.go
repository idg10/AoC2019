package main

import "testing"

func TestFuelForMass(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, c := range cases {
		got := fuelForMass(c.in)
		if got != c.want {
			t.Errorf("FuelForMass(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestFuelForMassAndFuel(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{14, 2},
		{1969, 966},
		{100756, 50346},
	}
	for _, c := range cases {
		got := fuelForMassAndFuel(c.in)
		if got != c.want {
			t.Errorf("FuelForMass(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestDay1(t *testing.T) {
	totalFuelIfFuelFree, totalFuel, _ := day1("IdgInput")
	if totalFuelIfFuelFree != 3329926 {
		t.Errorf("Part 1: %v, want %v", totalFuelIfFuelFree, 3329926)
	}
	if totalFuel != 4992008 {
		t.Errorf("Part 2: %v, want %v", totalFuel, 4992008)
	}
}
