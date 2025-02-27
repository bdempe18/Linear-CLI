// consider using init to get the apikey

package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/machinebox/graphql"
)

const apiBase = "https://api.linear.app/graphql"
const fieldFormat = "%s {\n    %s\n}"
const nodeFormat = "%s {\n    nodes {\n        %s\n    }\n}"
const qlFormat = "query {\n    %s\n}"

// ------------
// types
// ------------
type GraphQLLayer interface {
	String() string
}

type Field struct {
	Name     string
	Children []GraphQLLayer
}

type Node struct {
	Name     string
	Children []GraphQLLayer
	Filter   map[string]string
}

type GraphQLResponse interface {
	Show()
}

// ------------
// Stringifying
// ------------

func (f Field) String() string {
	if len(f.Children) == 0 {
		return f.Name
	}

	// node and field have some very similar logic at this point but will later diverge a bit.
	// I will refactor later
	var childStrings []string
	for _, child := range f.Children {
		childStrings = append(childStrings, child.String())
	}

	joinedChildren := strings.Join(childStrings, "\n        ")

	return fmt.Sprintf(fieldFormat, f.Name, joinedChildren)
}

func (n Node) String() string {
	if len(n.Children) == 0 {
		return n.Name
	}

	var childStrings []string
	for _, child := range n.Children {
		childStrings = append(childStrings, child.String())
	}

	joinedChildren := strings.Join(childStrings, "\n        ")

	return fmt.Sprintf(nodeFormat, n.Name, joinedChildren)
}

// ------------
// Factories
// ------------
func GenField(name string, children ...GraphQLLayer) Field {
	return Field{Name: name, Children: children}
}

func GenNode(name string, children ...GraphQLLayer) Node {
	return Node{Name: name, Children: children}
}

// -------------------
// Making the request
// -------------------
func Build(root GraphQLLayer) string {
	return fmt.Sprintf(qlFormat, root.String())
}

func Request(response interface{}, root GraphQLLayer, apiKey string) error {
	query := Build(root)
	request := graphql.NewRequest(query)

	request.Header.Set("Authorization", apiKey)

	client := graphql.NewClient(apiBase)
	if err := client.Run(context.Background(), request, response); err != nil {
		return fmt.Errorf("woops: %v", err)
	}

	return nil
}
