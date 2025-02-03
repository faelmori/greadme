package gmdtree

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// MarkdownNode represents a node in the Markdown tree.
type MarkdownNode struct {
	Level    int             // Level of the node (e.g., heading level)
	Title    string          // Title of the node
	Type     string          // Type of the node (e.g., title, list, codeBlock)
	Content  []string        // Content of the node
	Children []*MarkdownNode // Child nodes
}

var (
	// Regular expressions to match different Markdown elements
	titleRegex       = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	listRegex        = regexp.MustCompile(`^\s*[-*+] (.+)$`)
	orderedListRegex = regexp.MustCompile(`^\s*\d+\.\s+(.+)$`)
	codeBlockRegex   = regexp.MustCompile("^```(.*?)$")
	inlineCodeRegex  = regexp.MustCompile("`([^`]+)`")
)

// ParseMarkdown parses a Markdown file and returns the root node of the Markdown tree.
// It takes the file path as input and returns a pointer to the root MarkdownNode and an error if any.
func ParseMarkdown(filePath string) (*MarkdownNode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	root := &MarkdownNode{Level: 0, Title: "Root"}
	var currentNode = root
	codeBlockActive := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Detect code blocks
		if codeBlockRegex.MatchString(line) {
			codeBlockActive = !codeBlockActive
			continue
		}

		// If we are inside a code block, just store the line
		if codeBlockActive {
			currentNode.Content = append(currentNode.Content, line)
			continue
		}

		// Detect titles
		if match := titleRegex.FindStringSubmatch(line); match != nil {
			level := len(match[1])
			title := match[2]

			isTitle := titleRegex.MatchString(line)
			isList := listRegex.MatchString(line)
			isOrderedList := orderedListRegex.MatchString(line)
			isCodeBlock := codeBlockRegex.MatchString(line)
			isInlineCode := inlineCodeRegex.MatchString(line)

			var tp string
			if isTitle {
				tp = "title"
			}
			if isList {
				tp = "list"
			}
			if isOrderedList {
				tp = "orderedList"
			}
			if isCodeBlock {
				tp = "codeBlock"
			}
			if isInlineCode {
				tp = "inlineCode"
			}
			if !isTitle && !isList && !isOrderedList && !isCodeBlock && !isInlineCode {
				tp = "content"
			}

			newNode := &MarkdownNode{Level: level, Title: title, Type: tp}
			parent := FindParent(root, level)
			parent.Children = append(parent.Children, newNode)
			currentNode = newNode
			continue
		}

		// Detect lists
		if listRegex.MatchString(line) || orderedListRegex.MatchString(line) {
			currentNode.Content = append(currentNode.Content, line)
			continue
		}

		// Add common content
		if currentNode != nil {
			currentNode.Content = append(currentNode.Content, line)
		}
	}

	return root, scanner.Err()
}

// FindParent finds the parent node for a given level in the Markdown tree.
// It takes the root node and the level as input and returns the parent MarkdownNode.
func FindParent(root *MarkdownNode, level int) *MarkdownNode {
	var lastNode *MarkdownNode = root
	for len(lastNode.Children) > 0 {
		lastChild := lastNode.Children[len(lastNode.Children)-1]
		if lastChild.Level < level {
			return lastChild
		}
		lastNode = lastChild
	}
	return root
}

// PrintMarkdownTree prints the Markdown tree to the console.
// It takes the root node and an indent string as input.
func PrintMarkdownTree(node *MarkdownNode, indent string) {
	fmt.Printf("%s- [%s] (Level %d)\n", indent, node.Title, node.Level)
	for _, line := range node.Content {
		fmt.Printf("%s  * %s\n", indent, line)
	}
	for _, child := range node.Children {
		PrintMarkdownTree(child, indent+"  ")
	}
}

// GetMarkdownTree returns the Markdown tree as a string.
// It takes the root node and an indent string as input and returns the tree structure as a string.
func GetMarkdownTree(node *MarkdownNode, indent string) string {
	logOutput := fmt.Sprintf("%s- [%s] (Level %d)\n", indent, node.Title, node.Level)
	for _, line := range node.Content {
		logOutput += fmt.Sprintf("%s  %s\n", indent, line)
	}
	for _, child := range node.Children {
		logOutput += GetMarkdownTree(child, indent+"  ")
	}
	return logOutput
}
