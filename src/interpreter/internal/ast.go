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
	Body       Node
}
