type Node interface{}

type FunctionNode struct {
	Name       string
	Parameters []string
	Body       []Node
}

func Parser(tokens []Token) *FunctionNode {
}
