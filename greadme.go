package main

import (
	"bytes"
	"fmt"
	"github.com/faelmori/greadme/gmdtree"
	"github.com/spf13/cobra"
	"os/exec"
	"path/filepath"

	"os"
	"strings"
	"text/template"
)

// ReadmeData struct to hold README data
type ReadmeData struct {
	Order    []string          // Order of sections in the README
	Sections map[string]string // Map of section titles to their content
	Badges   []string          // List of badges to include in the README

	Org         string // Organization name
	Repo        string // Repository name
	ProjectName string // Project name

	Features        []string // List of features
	Platforms       []string // List of supported platforms
	QuickInstall    []string // Quick installation instructions
	Homebrew        []string // Homebrew installation instructions
	BuildFromSource []string // Instructions to build from source
	Providers       []string // List of providers
	Usage           []string // Usage instructions
	Commands        []string // List of commands
	EnvVars         []string // List of environment variables
	DevGuide        []string // Developer guide
	Contribution    []string // Contribution guidelines
	License         []string // License information
	Acknowledgments []string // Acknowledgments
}

// getProjectDetails retrieves the project name, organization, and repository from the git configuration.
// It changes the current directory to the git folder, runs a git command to get the details, and then changes back to the original directory.
func getProjectDetails(gitFolder string) (string, string, string, error) {
	var projectName, org, repo string
	var runCommandErr error

	// Save the current directory to return to it later
	currentDir, _ := os.Getwd()

	// Change directory to the gitFolder to run git commands
	chDirReadmeErr := os.Chdir(gitFolder)
	if chDirReadmeErr != nil {
		return "", "", "", fmt.Errorf("error changing directory to git folder: %v", chDirReadmeErr)
	}
	gitCommand := "git remote -v | grep origin | head -n 1 | awk '{print $2}' | sed 's/https:\\/\\///' | sed 's/git@//' | sed 's/\\.git$//' | tr '/' ' ' | tr ':' ' '"
	projectName, org, repo, runCommandErr = runCommand(gitCommand)

	// Change back to the original directory
	chDirReturnErr := os.Chdir(currentDir)
	if chDirReturnErr != nil {
		return "", "", "", fmt.Errorf("error changing back to original directory: %v", chDirReturnErr)
	}

	if runCommandErr != nil {
		return "", "", "", runCommandErr
	}
	return projectName, org, repo, nil
}

// runCommand executes a shell command and returns the output split into project name, organization, and repository.
// It captures both stdout and stderr, and returns an error if the command fails or the output cannot be parsed.
func runCommand(command string) (string, string, string, error) {
	cmd := exec.Command("sh", "-c", command)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", "", "", fmt.Errorf("error running command: %v", err)
	}
	if stderr.String() != "" {
		return "", "", "", fmt.Errorf("error running command: %v", stderr.String())
	}
	tmpRespList := strings.Split(strings.ReplaceAll(out.String(), "\n", ""), " ")
	if len(tmpRespList) != 3 {
		return "", "", "", fmt.Errorf("error parsing command output: %v", out.String())
	}
	repo := strings.Join(tmpRespList[1:], "/")
	org := tmpRespList[1]
	projectName := tmpRespList[2]
	return projectName, org, repo, nil
}

// logReadDataSummary logs the summary of the parsed Markdown tree to a file named "log_md_tree.txt".
// It creates the file, writes the tree structure to it, and closes the file.
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

// scanChildren scans the children of a Markdown node and populates the ReadmeData struct with the content.
// It categorizes the content based on the section titles and updates the corresponding fields in the ReadmeData struct.
func scanChildren(node *gmdtree.MarkdownNode, data *ReadmeData) {
	for _, child := range node.Children {
		content := strings.Join(child.Content, "\n")
		if content != "" {
			content = "<!-- TODO: Review and update this section -->\n" + content
		}
		data.Sections[child.Title] = content
		targetTitle := strings.ToLower(child.Title)

		if strings.Contains(targetTitle, "feature") {
			data.Features = child.Content
		}
		if strings.Contains(targetTitle, "platform") {
			data.Platforms = child.Content
		}
		if strings.Contains(targetTitle, "homebrew") || strings.Contains(targetTitle, "source") || strings.Contains(targetTitle, "quick") {
			parentNode := gmdtree.FindParent(child, child.Level+1)
			if strings.Contains(parentNode.Title, "install") {
				if strings.Contains(targetTitle, "homebrew") {
					data.Homebrew = child.Content
				}
				if strings.Contains(targetTitle, "source") {
					data.BuildFromSource = child.Content
				}
				if strings.Contains(targetTitle, "quick") {
					data.QuickInstall = child.Content
				}
			}
		}
		if strings.Contains(targetTitle, "provider") {
			data.Providers = child.Content
		}
		if strings.Contains(targetTitle, "usage") {
			data.Usage = child.Content
		}
		if strings.Contains(targetTitle, "command") || strings.Contains(targetTitle, "available") {
			data.Commands = child.Content
		}
		if strings.Contains(targetTitle, "env") {
			data.EnvVars = child.Content
		}
		if strings.Contains(targetTitle, "dev") {
			data.DevGuide = child.Content
		}
		if strings.Contains(targetTitle, "contrib") {
			data.Contribution = child.Content
		}
		if strings.Contains(targetTitle, "license") {
			data.License = child.Content
		}
		if strings.Contains(targetTitle, "acknowledgment") {
			data.Acknowledgments = child.Content
		}
	}
}

// extractReadmeData parses the README file and extracts the data into a ReadmeData struct.
// It also retrieves project details from the git configuration and logs the parsed Markdown tree.
func extractReadmeData(readmeFile string) (*ReadmeData, error) {
	root, err := gmdtree.ParseMarkdown(readmeFile)
	if err != nil {
		return nil, err
	}

	data := ReadmeData{
		Org:         "",
		Repo:        "",
		ProjectName: "",

		Order:    make([]string, 0),
		Sections: make(map[string]string),
		Badges:   make([]string, 0),

		Features:        make([]string, 0),
		Platforms:       make([]string, 0),
		QuickInstall:    make([]string, 0),
		Homebrew:        make([]string, 0),
		BuildFromSource: make([]string, 0),
		Providers:       make([]string, 0),
		Usage:           make([]string, 0),
		Commands:        make([]string, 0),
		EnvVars:         make([]string, 0),
		DevGuide:        make([]string, 0),
		Contribution:    make([]string, 0),
		License:         make([]string, 0),
		Acknowledgments: make([]string, 0),
	}

	gitFolder := filepath.Dir(readmeFile)
	projectName, org, repo, projectDetailsErr := getProjectDetails(gitFolder)
	if projectDetailsErr != nil {
		fmt.Println("⚠️ Error getting project details:", projectDetailsErr)
	} else {
		data.ProjectName = projectName
		data.Org = org
		data.Repo = repo
	}

	for _, node := range root.Children {
		scanChildren(node, &data)
	}

	_ = logReadDataSummary(root, "")

	return &data, nil
}

// generateImprovedReadme generates an improved README file using a template and the extracted README data.
// It reads the template file, parses it, renders the template with the data, and writes the output to the specified file.
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

// main is the entry point of the application. It sets up the Cobra CLI and executes the root command.
// The root command generates an improved README file based on the provided template and README file.
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
