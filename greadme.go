package main

import (
	"fmt"
	"github.com/faelmori/greadme/gmdtree"
	"github.com/spf13/cobra"

	"os"
	"strings"
	"text/template"
)

// ReadmeData struct to hold README data
type ReadmeData struct {
	Org             string
	Repo            string
	ProjectName     string
	Badges          []string
	Features        string
	Platforms       string
	QuickInstall    string
	Homebrew        string
	BuildFromSource string
	Providers       string
	Usage           string
	Commands        string
	EnvVars         string
	DevGuide        string
	Contribution    string
	License         string
	Acknowledgments string
}

// logReadDataSummary 📌 **COMPLEMENTO: Log de sumário**
func logReadDataSummary(node *gmdtree.MarkdownNode, indent string) error {
	logMdTree := gmdtree.GetMarkdownTree(node, indent)

	logMdTreeFile, logMdTreeFileErr := os.Create("log_md_tree.txt")
	if logMdTreeFileErr != nil {
		fmt.Println("❌ Error creating log file:", logMdTreeFileErr)
		return logMdTreeFileErr
	}

	_, writeStringErr := logMdTreeFile.WriteString(logMdTree)
	if writeStringErr != nil {
		fmt.Println("❌ Error writing log file:", writeStringErr)
		return writeStringErr
	}
	_ = logMdTreeFile.Close()

	fmt.Println("✅ `log_md_tree.txt` generated successfully!")
	return nil
}

// extractReadmeData 📌 **PASSO 2: Integrar extração ao preenchimento do README**
func extractReadmeData(readmeFile string) (*ReadmeData, error) {
	root, err := gmdtree.ParseMarkdown(readmeFile)
	if err != nil {
		return nil, err
	}

	data := &ReadmeData{}

	if root.Level == 0 || root.Level == 1 {
		for _, section := range root.Children {
			content := strings.Join(section.Content, "\n")
			if content == "" {
				content = "<!-- TODO: Add content for " + section.Title + " -->"
			}

			if section.Type == "title" {
				switch strings.ToLower(strings.TrimSpace(section.Title)) {
				case "features":
					data.Features = content
				case "platforms":
					data.Platforms = content
				case "quick installation":
					data.QuickInstall = content
				case "homebrew":
					data.Homebrew = content
				case "build from source":
					data.BuildFromSource = content
				case "supported providers":
					data.Providers = content
				case "usage":
					data.Usage = content
				case "available commands":
					data.Commands = content
				case "provider credentials":
					data.EnvVars = content
				case "development guide":
					data.DevGuide = content
				case "contribution":
					data.Contribution = content
				case "license":
					data.License = content
				case "acknowledgments":
					data.Acknowledgments = content
				}
			} else if section.Type == "list" {
				if strings.Contains(strings.ToLower(content), "github.com") {
					data.Repo = content
				} else if strings.Contains(strings.ToLower(content), "github.io") {
					data.Org = content
				} else {
					data.ProjectName = content
				}
			}
		}
	} else {
		content := strings.Join(root.Content, "\n")
		switch strings.ToLower(strings.TrimSpace(root.Title)) {
		case "features":
			data.Features = content
		case "platforms":
			data.Platforms = content
		case "quick installation":
			data.QuickInstall = content
		case "homebrew":
			data.Homebrew = content
		case "build from source":
			data.BuildFromSource = content
		case "supported providers":
			data.Providers = content
		case "usage":
			data.Usage = content
		case "available commands":
			data.Commands = content
		case "provider credentials":
			data.EnvVars = content
		case "development guide":
			data.DevGuide = content
		case "contribution":
			data.Contribution = content
		case "license":
			data.License = content
		case "acknowledgments":
			data.Acknowledgments = content
		}
	}

	_ = logReadDataSummary(root, "")

	return data, nil
}

// generateImprovedReadme 📌 **PASSO 3: Gerar o IMPROVED_README.md corretamente**
func generateImprovedReadme(templateFile, readmeFile, outputFile string) {
	readmeData, err := extractReadmeData(readmeFile)
	if err != nil {
		fmt.Println("❌ Error extracting README data:", err)
		return
	}

	tmplStr, err := os.ReadFile(templateFile)
	if err != nil {
		fmt.Println("⚠️ Using default embedded template as fallback.")
		tmplStr = []byte(defaultTemplate)
	}

	tmplObj, err := template.New("readme").Parse(string(tmplStr))
	if err != nil {
		fmt.Println("❌ Error parsing template:", err)
		return
	}

	var rendered strings.Builder
	err = tmplObj.Execute(&rendered, readmeData)
	if err != nil {
		fmt.Println("❌ Error rendering template:", err)
		return
	}

	err = os.WriteFile(outputFile, []byte(rendered.String()), 0644)
	if err != nil {
		fmt.Println("❌ Error writing improved README:", err)
		return
	}

	fmt.Println("✅ `IMPROVED_README.md` generated successfully!")
}

// main 📌 **PASSO 4: CLI com Cobra**
func main() {
	var rootCmd = &cobra.Command{
		Use:   "readme-checker",
		Short: "Checks and improves README structure",
		Run: func(cmd *cobra.Command, args []string) {
			templateFile, _ := cmd.Flags().GetString("template")
			readmeFile, _ := cmd.Flags().GetString("readme")
			outputFile, _ := cmd.Flags().GetString("output")

			generateImprovedReadme(templateFile, readmeFile, outputFile)
		},
	}

	rootCmd.Flags().StringP("template", "t", "", "Template README file (leave empty to use built-in template)")
	rootCmd.Flags().StringP("readme", "r", "README_to_check.md", "README file to check")
	rootCmd.Flags().StringP("output", "o", "IMPROVED_README.md", "Output improved README file")

	executeCmdErr := rootCmd.Execute()
	if executeCmdErr != nil {
		fmt.Println("❌ Error executing command:", executeCmdErr)
	}
}
