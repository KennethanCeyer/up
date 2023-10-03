package main

import (
	"fmt"
	"strconv"
	"syscall"
)

type Node interface{}

type IntNode struct {
	Value int
}

type AddNode struct {
	Left  Node
	Right Node
}

type FunctionNode struct {
	Name       string
	Parameters []string
	Body       []Node
}

func Execute(node *FunctionNode) {
	if node.Name != "x" {
		fmt.Println("Unknown function:", node.Name)
		return
	}

	if len(node.Parameters) != 2 {
		fmt.Println("Invalid number of arguments")
		return
	}

	a, err1 := strconv.Atoi(node.Parameters[0])
	b, err2 := strconv.Atoi(node.Parameters[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Error in arguments")
		return
	}

	if addNode, ok := node.Body[0].(AddNode); ok {
		left, lok := addNode.Left.(IntNode)
		right, rok := addNode.Right.(IntNode)
		if !lok || !rok {
			fmt.Println("Error in addition")
			return
		}

		result := left.Value + right.Value

		msg := strconv.Itoa(result) + "\n"
		syscall.Write(syscall.Stdout, []byte(msg))
	}
}

func main() {
	funcNode := &FunctionNode{
		Name:       "x",
		Parameters: []string{"5", "3"},
		Body: []Node{
			AddNode{
				Left:  IntNode{Value: 5},
				Right: IntNode{Value: 3},
			},
		},
	}

	Execute(funcNode)
}
