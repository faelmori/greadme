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
	Order    []string
	Sections map[string]string
	Badges   []string

	Org         string
	Repo        string
	ProjectName string

	Features        []string
	Platforms       []string
	QuickInstall    []string
	Homebrew        []string
	BuildFromSource []string
	Providers       []string
	Usage           []string
	Commands        []string
	EnvVars         []string
	DevGuide        []string
	Contribution    []string
	License         []string
	Acknowledgments []string
}

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
	children := root.Children

	gitFolder := filepath.Dir(readmeFile)
	projectName, org, repo, projectDetailsErr := getProjectDetails(gitFolder)
	if projectDetailsErr != nil {
		fmt.Println("⚠️ Error getting project details:", projectDetailsErr)
	} else {
		data.ProjectName = projectName
		data.Org = org
		data.Repo = repo
	}

	for _, section := range children {
		content := strings.Join(section.Content, "\n")
		if content == "" {
			content = "<!-- TODO: Add content for " + section.Title + " -->"
		}

		if section.Level > 1 {
			if section.Type == "title" {
				matchTarget := strings.ToLower(section.Title)
				if strings.Contains(matchTarget, "badge") {
					data.Badges = append(data.Badges, content)
				} else if strings.Contains(matchTarget, "feature") {
					data.Features = append(data.Features, content)
				} else if strings.Contains(matchTarget, "platform") {
					data.Platforms = append(data.Platforms, content)
				} else if strings.Contains(matchTarget, "install") {
					data.QuickInstall = append(data.QuickInstall, content)
				} else if strings.Contains(matchTarget, "homebrew") {
					data.Homebrew = append(data.Homebrew, content)
				} else if strings.Contains(matchTarget, "build") {
					data.BuildFromSource = append(data.BuildFromSource, content)
				} else if strings.Contains(matchTarget, "provider") {
					data.Providers = append(data.Providers, content)
				} else if strings.Contains(matchTarget, "usage") {
					data.Usage = append(data.Usage, content)
				} else if strings.Contains(matchTarget, "command") {
					data.Commands = append(data.Commands, content)
				} else if strings.Contains(matchTarget, "env") {
					data.EnvVars = append(data.EnvVars, content)
				} else if strings.Contains(matchTarget, "dev") {
					data.DevGuide = append(data.DevGuide, content)
				} else if strings.Contains(matchTarget, "contrib") {
					data.Contribution = append(data.Contribution, content)
				} else if strings.Contains(matchTarget, "license") {
					data.License = append(data.License, content)
				} else if strings.Contains(matchTarget, "acknowledgment") {
					data.Acknowledgments = append(data.Acknowledgments, content)
				}
			}
		} else {
			if section.Type == "title" && strings.ToLower(section.Title) != "root" {
				data.ProjectName = section.Title
			}
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
