package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Atom interface {
	IsOperator() bool
}

type Operator string

const (
	ADD Operator = "+"
	DIV Operator = "/"
	MUL Operator = "*"
)

func (o Operator) IsOperator() bool {
	return true
}

type Number float64

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
func Evaluate(stack *Stack) float64 {
	nums := &Stack{}

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
			}
		} else {
			nums.Push(curr)
		}
	}

	return float64(nums.Peek().(Number))
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

func main() {
	expression, _ := Parse("1 2 + 3 4 / * 10 +")
	fmt.Println(expression.Valid())
}
