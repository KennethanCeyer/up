package interpreter

import "fmt"

func VisualizeNode(node Node) {
	visualizedText := node.String()
	fmt.Println(visualizedText)
}
