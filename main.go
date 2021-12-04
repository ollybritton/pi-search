package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Atom interface {
	IsOperator() bool
}

type Operator string

const (
	ADD  Operator = "+"
	DIV  Operator = "/"
	MUL  Operator = "*"
	SQRT Operator = "√"
)

func RandomOperator() Operator {
	return []Operator{
		ADD,
		DIV,
		MUL,
	}[rand.Intn(3)]
}

func (o Operator) IsOperator() bool {
	return true
}

type Number float64

func RandomWholeNumber(min, max int) Number {
	return Number(float64(rand.Intn(max-min) + min))
}

func (n Number) IsOperator() bool {
	return false
}

type Stack struct {
	items []Atom
}

func NewStack(items ...Atom) *Stack {
	return &Stack{items: items}
}

func (s *Stack) Pop() Atom {
	atom := s.items[0]
	s.items = s.items[1:]

	return atom
}

func (s *Stack) Len() int {
	return len(s.items)
}

func (s *Stack) Peek() Atom {
	return s.items[0]
}

func (s *Stack) Push(atom Atom) {
	s.items = append([]Atom{atom}, s.items...)
}

func (s *Stack) Copy() *Stack {
	items := make([]Atom, s.Len())
	copy(items, s.items)

	return &Stack{items: items}
}

func (s *Stack) String() string {
	var out []string

	for _, atom := range s.items {
		if atom.IsOperator() {
			out = append(out, string(atom.(Operator)))
		} else {
			out = append(out, fmt.Sprint(atom.(Number)))
		}
	}

	return strings.Join(out, " ")
}

// Parse parses a string of space-separated operators and numbers in postfix notation to a stack.
func Parse(expression string) (*Stack, error) {
	unparsedAtoms := strings.Split(expression, " ")
	parsedAtoms := []Atom{}

	for _, unparsedAtom := range unparsedAtoms {
		var parsedAtom Atom

		switch unparsedAtom {
		case "+":
			parsedAtom = ADD
		case "*":
			parsedAtom = MUL
		case "/":
			parsedAtom = DIV
		case "√":
			parsedAtom = SQRT
		default:
			num, err := strconv.ParseFloat(unparsedAtom, 64)
			if err != nil {
				return nil, fmt.Errorf("couldn't parse %q: %w", unparsedAtom, err)
			}

			parsedAtom = Number(num)
		}

		parsedAtoms = append(parsedAtoms, parsedAtom)
	}

	stack := &Stack{}
	length := len(parsedAtoms)

	for i := range parsedAtoms {
		stack.Push(parsedAtoms[length-i-1])
	}

	return stack, nil
}

// Valid returns true if the stack represents valid a RPN/infix expression.
// Algorithm from: https://stackoverflow.com/questions/14506831/whats-the-fastest-way-to-check-if-input-string-is-a-correct-rpn-expression
func (s *Stack) Valid() bool {
	size := 0

	for _, atom := range s.items {
		valence := 0

		switch atom {
		case ADD:
			valence = 2
		case MUL:
			valence = 2
		case DIV:
			valence = 2
		case SQRT:
			valence = 1
		default:
			valence = 0
		}

		size += 1 - valence

		if size <= 0 {
			return false
		}
	}

	return size == 1
}

// Evaluate evaluates a stack of atoms in postfix notation.
func Evaluate(s *Stack) float64 {
	nums := &Stack{}
	stack := s.Copy()

	for stack.Len() > 0 {
		curr := stack.Pop()

		if curr.IsOperator() {
			switch curr {
			case ADD:
				x := nums.Pop().(Number)
				y := nums.Pop().(Number)
				nums.Push(y + x)
			case MUL:
				x := nums.Pop().(Number)
				y := nums.Pop().(Number)
				nums.Push(y * x)
			case DIV:
				x := nums.Pop().(Number)
				y := nums.Pop().(Number)
				nums.Push(y / x)
			case SQRT:
				x := nums.Pop().(Number)
				nums.Push(Number(math.Sqrt(float64(x))))
			}
		} else {
			nums.Push(curr)
		}
	}

	return float64(nums.Peek().(Number))
}

// Generate generates a random, valid RPN string of length n.
func Generate(length int) *Stack {
	stack := NewStack(generateRecursive(1, 10, length)...)

	return stack
}

// generateRecursive
func generateRecursive(min, max, length int) []Atom {
	switch {
	case length < 1:
		return []Atom{}
	case length == 1:
		return []Atom{RandomWholeNumber(min, max)}
	case length == 2:
		return []Atom{RandomWholeNumber(min, max), SQRT}
	case length == 3:
		return []Atom{RandomWholeNumber(min, max), RandomWholeNumber(min, max), RandomOperator()}
	default:
		if rand.Intn(4) == 0 {
			return append(
				generateRecursive(min, max, length-1),
				SQRT,
			)
		} else {
			return append(
				generateRecursive(min, max, length/2),
				append(
					generateRecursive(min, max, length/2),
					RandomOperator(),
				)...,
			)
		}
	}
}

func Improve(expression *Stack, target, val, diff float64) (bool, float64, float64, *Stack) {
	for i, atom := range expression.items {
		if atom.IsOperator() {
			continue
		}

		num := atom.(Number)
		expression.items[i] = num + 1

		newVal := Evaluate(expression)
		newDiff := math.Abs(target - newVal)

		if newDiff < diff {
			return true, newVal, newDiff, expression
		}

		expression.items[i] = num
	}

	return false, val, diff, expression
}

// Search searches for approximations to the input number using basic math operations.
// Precision is the number of decimal places.
func Search(approximate float64, precision int, minLength, maxLength, minNum, maxNum int) {
	epsilon := math.Pow10(-precision)

	for i := 0; i < 10; i++ {
		go func() {
			for {
				expression := Generate(rand.Intn(maxLength-minLength) + minLength)
				val := Evaluate(expression)
				diff := math.Abs(approximate - val)

				if diff < epsilon {
					fmt.Printf("%f,%f,%s\n", diff/epsilon, val, expression.String())
				}
			}
		}()
	}

	select {}
}

func generateDistribtuion() {
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"num", "expression"})

	for i := 0; i < 1_000_000; i++ {
		expression := Generate(5)
		val := Evaluate(expression)
		writer.Write([]string{fmt.Sprint(val), expression.String()})
	}

	writer.Flush()
}

func main() {
	Search(math.Pi, 5, 10, 20, 1, 100)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
