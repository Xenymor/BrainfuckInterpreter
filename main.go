package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/amarsinghrathour/go-stack/stack"
	"github.com/hashicorp/go-set/v3"
)

func main() {
	var codeFile string

	fmt.Print("Enter Brainfuck code file path (default: \"test.bf\"): ")
	fmt.Scanln(&codeFile)

	code, err := os.ReadFile(ternary(codeFile == "", "test.bf", codeFile))
	if err != nil {
		panic(err)
	}

	chars := []rune{'>', '<', '+', '-', '.', ',', '[', ']'}

	validChars = set.From(chars)

	interpret(code)
}

var validChars *set.Set[rune]

func interpret(file []byte) {
	code := make([]rune, 0, len(file))

	firstIsNum := false
	size := 1 << 10

	for i, curr := range file {
		newRune := rune(curr)
		if validChars.Contains(newRune) {
			code = append(code, newRune)
		} else if i == 0 && curr >= '0' && curr <= '9' {
			firstIsNum = true
		} else if firstIsNum && i == 1 && curr >= '0' && curr <= '9' {
			d0 := int(file[0] - '0')
			d1 := int(file[1] - '0')
			exp := d0*10 + d1
			size = 1 << exp
		}
	}

	fmt.Printf("Running code with pointer range of 0-%v\n", size)

	values := make([]int, size)

	execute(code, &values)
}

var reader = bufio.NewReader(os.Stdin)

func execute(code []rune, values *[]int) {
	pointer := 0

	bracketPairs := map[int]int{}

	bracketStack := stack.NewLinkedStack()

	for i := 0; i < len(code); i++ {
		curr := code[i]
		if curr == '[' {
			bracketStack.Push(i)
		}
		if curr == ']' {
			data, ok := bracketStack.Pop()
			if !ok {
				fmt.Println("\nInvalid Code: Too many closing brackets")
				return
			}
			first := data.(int)
			bracketPairs[first] = i
			bracketPairs[i] = first
		}
	}

	if !bracketStack.IsEmpty() {
		fmt.Println("\nInvalid Code: Not all brackets are closed")
		return
	}

	for i := 0; i < len(code); i++ {
		curr := code[i]
		switch curr {
		case '>':
			pointer++
			if pointer >= len(*values) {
				fmt.Printf("\nOutOfBoundsException: Pointer %v out of bounds for length %v at symbol %v\n", pointer, len(*values), i)
				return
			}
		case '<':
			pointer--
			if pointer < 0 {
				fmt.Printf("\nOutOfBoundsException: Pointer %v out of bounds for length %v at symbol %v\n", pointer, len(*values), i)
				return
			}
		case '+':
			(*values)[pointer]++
		case '-':
			(*values)[pointer]--
		case '.':
			fmt.Print(string(rune((*values)[pointer])))
		case ',':
			b, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}
			(*values)[pointer] = int(b)
		case '[':
			if (*values)[pointer] == 0 {
				i = bracketPairs[i]
			}
		case ']':
			if (*values)[pointer] != 0 {
				i = bracketPairs[i]
			}
		}
	}
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
