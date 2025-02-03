package gmdtree

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Estrutura de um nó da árvore Markdown
type MarkdownNode struct {
	Level    int
	Title    string
	Content  []string
	Children []*MarkdownNode
}

// Expressões regulares para capturar elementos do Markdown
var (
	titleRegex       = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	listRegex        = regexp.MustCompile(`^\s*[-*+] (.+)$`)
	orderedListRegex = regexp.MustCompile(`^\s*\d+\.\s+(.+)$`)
	codeBlockRegex   = regexp.MustCompile("^```(.*?)$")
	inlineCodeRegex  = regexp.MustCompile("`([^`]+)`")
)

// Função para processar o arquivo Markdown e construir a hierarquia
func ParseMarkdown(filePath string) (*MarkdownNode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	root := &MarkdownNode{Level: 0, Title: "Root"}
	var currentNode *MarkdownNode = root
	codeBlockActive := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Detecta blocos de código
		if codeBlockRegex.MatchString(line) {
			codeBlockActive = !codeBlockActive
			continue
		}

		// Se estamos dentro de um bloco de código, apenas armazena a linha
		if codeBlockActive {
			currentNode.Content = append(currentNode.Content, line)
			continue
		}

		// Detecta títulos
		if match := titleRegex.FindStringSubmatch(line); match != nil {
			level := len(match[1]) // Número de `#` define o nível
			title := match[2]

			newNode := &MarkdownNode{Level: level, Title: title}

			// Encontra o nó pai adequado
			parent := FindParent(root, level)
			parent.Children = append(parent.Children, newNode)
			currentNode = newNode
			continue
		}

		// Detecta listas
		if match := listRegex.FindStringSubmatch(line); match != nil || orderedListRegex.MatchString(line) {
			currentNode.Content = append(currentNode.Content, line)
			continue
		}

		// Detecta inline code
		if match := inlineCodeRegex.FindStringSubmatch(line); match != nil {
			currentNode.Content = append(currentNode.Content, "Code: "+match[1])
			continue
		}

		// Adiciona conteúdo comum
		if currentNode != nil {
			currentNode.Content = append(currentNode.Content, line)
		}
	}

	return root, scanner.Err()
}

// Encontra o nó pai correto na hierarquia
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

// Função para imprimir a árvore Markdown
func PrintMarkdownTree(node *MarkdownNode, indent string) {
	fmt.Printf("%s- [%s] (Level %d)\n", indent, node.Title, node.Level)
	for _, line := range node.Content {
		fmt.Printf("%s  * %s\n", indent, line)
	}
	for _, child := range node.Children {
		PrintMarkdownTree(child, indent+"  ")
	}
}
