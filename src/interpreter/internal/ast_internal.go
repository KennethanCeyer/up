package interpreter

import "fmt"

func VisualizeNode(node Node) {
	fmt.Println("+------------+------------------------------------------+")
	fmt.Println("| Node Type  | Value                                    |")
	fmt.Println("+------------+------------------------------------------+")
	printNode(node, "")
	fmt.Println("+------------+------------------------------------------+")
}

func printNode(node Node, prefix string) {
	switch n := node.(type) {
	case *ProgramNode:
		for _, fn := range n.Functions {
			printNode(fn, prefix)
		}
	case *FuncDeclarationNode:
		printTableRow("Function", n.Name)
		for _, param := range n.Parameters {
			printNode(param, "  "+prefix)
		}
		for _, stmt := range n.Body {
			printNode(stmt, "  "+prefix)
		}
	case *ParameterNode:
		printTableRow("Parameter", n.String())
	case *ForLoopNode:
		printTableRow("ForLoop", "for "+n.Variable+" in range(...)")
		for _, stmt := range n.Body {
			printNode(stmt, "  "+prefix)
		}
	case *AssignmentNode:
		printTableRow("Assignment", n.VarName+": "+n.Type)
	case *ReturnNode:
		printTableRow("Return", n.Value.String())
	case *BinOpNode:
		printTableRow("BinOp", n.Op)
	case *FunctionCallNode:
		printTableRow("FunctionCall", n.FunctionName+"(...)")
	case *IdentifierNode:
		printTableRow("Identifier", n.Name)
	case *IntNode:
		printTableRow("Int", n.String())
	case *StringNode:
		printTableRow("String", n.String())
	default:
		printTableRow("Unknown", "")
	}
}

func printTableRow(nodeType string, value string) {
	truncatedValue := value
	if len(value) > 40 {
		truncatedValue = value[:37] + "..."
	}
	fmt.Printf("| %-10s | %-40s |\n", nodeType, truncatedValue)
}
