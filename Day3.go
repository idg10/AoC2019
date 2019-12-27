package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path"
	"reflect"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

// Direction indicates a direction
type Direction int

func toInt(nodes []parsec.ParsecNode) parsec.ParsecNode {
	n := nodes[0]
	t := n.(*parsec.Terminal)
	i, _ := strconv.Atoi(t.Value)
	return i
}

// YInt is a Parser that parses a decimal integer and produces an int
var YInt = parsec.And(toInt, parsec.Int())

// Ycomma is a Parser that matches and consumes a comma.
var Ycomma = parsec.Token(`,`, "FIELDSEP")

// YEol is a Parser that matches and consumes an end of line
//var YEol = parsec.TokenExact(`(\r\n)|\r|\n`, "EOL")
var YEol = parsec.OrdChoice(
	nil,
	parsec.AtomExact("\r\n", "CRLF"),
	parsec.AtomExact("\n", "LF"),
	parsec.AtomExact("\r", "CR"))

// YCsvInt is a Parser that matches a comma-separated list of decimal integers and produces an []int
var YCsvInt = parsec.Kleene(nil, YInt, Ycomma)

func toDirection(nodes []parsec.ParsecNode) parsec.ParsecNode {
	n := nodes[0]
	t := n.(*parsec.Terminal)
	switch t.Value {
	case "U":
		return DirectionUp
	case "D":
		return DirectionDown
	case "L":
		return DirectionLeft
	case "R":
		return DirectionRight
	}

	panic("Bad argument to toDirection")
}

// Could also do this with a single Token using a regex matching any of the letters.
// I guess it depends on how you feel about expressing your rules in a mixture of the
// parser world and the world of regex.
var ydirection = parsec.OrdChoice(
	toDirection,
	parsec.Atom("U", "Up"),
	parsec.Atom("D", "Down"),
	parsec.Atom("L", "Left"),
	parsec.Atom("R", "Right"))

func toSegment(nodes []parsec.ParsecNode) parsec.ParsecNode {
	d := nodes[0]
	n := nodes[1]
	direction := d.(Direction)
	i := n.(int)
	return WireSegment{direction, i}
}

// YSegment parses segments of the form U5 or R15 producing a WireSegment
var YSegment = parsec.And(toSegment, ydirection, YInt)

// ParsecNodesToArray is nodifying function that converts a slice of ParsecNodes to
// a slice of a specific type.
func ParsecNodesToArray(t reflect.Type, nodes []parsec.ParsecNode) parsec.ParsecNode {
	count := len(nodes)
	values := reflect.MakeSlice(reflect.SliceOf(t), 0, count)
	for _, v := range nodes {
		values = reflect.Append(values, reflect.ValueOf(v))
	}

	itf := values.Interface()

	return itf
}

// MakeArrayNodifier creates a nodifying function that converts a slice of ParsecNodes
// to a slice of a specific type.
func MakeArrayNodifier(elementType reflect.Type) func(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		return ParsecNodesToArray(elementType, nodes)
	}
}

var toSegments = MakeArrayNodifier(reflect.TypeOf((*WireSegment)(nil)).Elem())

// YWirePath parses a comma-separated list of wire segments, e.g.:
// U10,R42,D4,U1
var YWirePath = parsec.Kleene(toSegments, YSegment, Ycomma)

var toPaths = MakeArrayNodifier(reflect.TypeOf((*[]WireSegment)(nil)).Elem())

// YWirePathLines parses multiple lines each of which contains a comma-separated list
// of wire segments (in the form parsed by YWirePath).
var YWirePathLines = parsec.Kleene(toPaths, YWirePath, YEol)

const (
	// DirectionUp indicates upward movement.
	DirectionUp Direction = 1

	// DirectionDown indicates downward movement.
	DirectionDown Direction = 2

	// DirectionLeft indicates movement to the left.
	DirectionLeft Direction = 3

	// DirectionRight indicates movement to the right.
	DirectionRight Direction = 4
)

// WireSegment represents one segment of a wire.
type WireSegment struct {
	direction Direction
	distance  int
}

// GridPosition represents an x,y position on a grid.
type GridPosition struct {
	x, y int
}

// WireGrid represents the wiring grid.
type WireGrid struct {
	points              map[GridPosition]map[int]int
	collisions          map[GridPosition]struct{}
	atLeastOneCollision bool
}

// OMG I am starting to get fed up of "There is no built-in X, but it’s simple to write your own"
func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(start GridPosition, end GridPosition) int {
	dx := intAbs(end.x - start.x)
	dy := intAbs(end.y - start.y)
	return dx + dy
}

// AddWireAt records the presence of a particular wire at a particular location
func (grid *WireGrid) AddWireAt(wireID int, position GridPosition, runDistance int) {
	ids, occupied := grid.points[position]
	if occupied {
		ids[wireID] = runDistance
		// Do we need to do this? Does append mutate the slice or create a new one?
		grid.points[position] = ids

		if len(ids) > 1 {
			grid.collisions[position] = struct{}{}
			// It wasn't just us in there already, so this is a collision.
			// 	origin := GridPosition{0, 0}
			// 	distanceToThisCollision := manhattan(origin, position)

			// 	if !grid.atLeastOneCollision || distanceToThisCollision < manhattan(origin, grid.closestCollision) {
			// 		// This is either the first collision we've seen, or it's closer than any we've
			// 		// seen before
			// 		grid.atLeastOneCollision = true
			// 		grid.closestCollision = position
			// 	}
			// 	// sort.Search(len(ids), func(i int) bool { return manhattan(origin, grid.collisions[i]) <= distance })
		}
	} else {
		grid.points[position] = map[int]int{wireID: runDistance}
	}
}

func day3Core(text []byte, distance func(grid WireGrid, position GridPosition) int) int {

	s := parsec.NewScanner(text)
	node, _ := YWirePathLines(s)
	items := node.([][]WireSegment)

	grid := WireGrid{map[GridPosition]map[int]int{}, map[GridPosition]struct{}{}, false}
	for wireID, wire := range items {
		x := 0
		y := 0
		runDistance := 0
		for _, segment := range wire {
			var dx, dy int
			switch segment.direction {
			case DirectionDown:
				dx, dy = 0, -1
			case DirectionUp:
				dx, dy = 0, 1
			case DirectionLeft:
				dx, dy = -1, 0
			case DirectionRight:
				dx, dy = 1, 0
			}
			for i := 0; i < segment.distance; i++ {
				x += dx
				y += dy
				runDistance++
				grid.AddWireAt(wireID, GridPosition{x, y}, runDistance)
			}
		}
	}

	var closest = math.MaxInt32
	for position := range grid.collisions {
		itemDistance := distance(grid, position)
		if itemDistance < closest {
			closest = itemDistance
		}
	}

	return closest
}

// Day3Part1 calculates the result for day 3 part 1.
func Day3Part1(text []byte) int {
	return day3Core(text, func(grid WireGrid, position GridPosition) int {
		return manhattan(GridPosition{0, 0}, position)
	})
}

// Day3Part2 calculates the result for day 3 part 2.
func Day3Part2(text []byte) int {
	return day3Core(text, func(grid WireGrid, position GridPosition) int {
		runDistanceSum := 0
		for _, runDistance := range grid.points[position] {
			runDistanceSum += runDistance
		}
		return runDistanceSum
	})
}

// Day3 runs day 3's tests.
func showDay3(inputFolder string) {
	fileBytes, err := ioutil.ReadFile(path.Join(inputFolder, "Day3.txt"))
	if err != nil {
		return
	}

	part1 := Day3Part1(fileBytes)
	part2 := Day3Part2(fileBytes)

	fmt.Printf("Part 1: %v\n", part1)
	fmt.Printf("Part 2: %v\n", part2)
}
