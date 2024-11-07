package up

import (
	"fmt"
	"strings"
)

func VisualizeTokens(tokens []Token) {
	var tokenStrings []string

	for _, token := range tokens {
		tokenStrings = append(tokenStrings, fmt.Sprintf("{Type: %s, Value: %q}", token.Type, token.Value))
	}

	visualizedText := strings.Join(tokenStrings, "\n")
	fmt.Println(visualizedText)
}
