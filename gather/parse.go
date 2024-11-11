package gather

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
)

var grammars = map[string]*sitter.Language{
	"go": golang.GetLanguage(),
	"rust": rust.GetLanguage(),
	"python": python.GetLanguage(),
}

var patterns = map[string]string{
	"go": "(function_declaration) @capture",
	"rust": "(function_item) @capture",
	"python": "(function_definition) @capture",
}

func Parse(source []byte, language string) (sections [][]byte, err error) {
	grammar, ok := grammars[language]
	if !ok {
		return nil, fmt.Errorf("grammar does not exist for language: %s", language)
	}

	pattern, ok := patterns[language]
	if !ok {
		return nil, fmt.Errorf("pattern does not exist for language: %s", language)
	}

	query, err := sitter.NewQuery([]byte(pattern), grammar)
	if err != nil {
		return nil, fmt.Errorf("query parsing: %w", err)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(grammar)
	tree, err := parser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, fmt.Errorf("tree parsing: %w", err)
	}

	root := tree.RootNode()

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		match = cursor.FilterPredicates(match, source)
		for _, capture := range match.Captures {
			start := capture.Node.StartByte()
			end := capture.Node.EndByte()
			sections = append(sections, source[start:end])
		}
	}

	return
}

