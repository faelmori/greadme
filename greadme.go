package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// Regex for capturing titles and badges
var titleRegex = regexp.MustCompile(`^(#+)\s+(.*)`)
var badgeRegex = regexp.MustCompile(`!\[.*\]\(https://img\.shields\.io.*\)`)

type ReadmeData struct {
	ProjectName string
	Org         string
	Repo        string
}

func refReadme(readmeFile string) (string, error) {
	data, err := getDataFromReadme(readmeFile)
	if err != nil {
		fmt.Println("Error getting project details:", err)
		return "", err
	}
	tmplObj := strings.Join([]string{`
# {{.ProjectName}}

![Version](https://img.shields.io/github/v/release/{{.Org}}/{{.Repo}})
![Build Status](https://img.shields.io/github/actions/workflow/status/{{.Org}}/{{.Repo}}/build.yml?branch=main)
![License](https://img.shields.io/github/license/{{.Org}}/{{.Repo}})

A brief description of the project, explaining its purpose and main features.

> **Note:** Important instructions or warnings can be highlighted here.

---

## 📖 Table of Contents

- [✨ Features](#-features)
- [📥 Installation](#-installation)
    - [Supported Platforms](#supported-platforms)
    - [1. Quick Installation](#1-quick-installation)
    - [2. Homebrew](#2-homebrew)
    - [3. Build from Source](#3-build-from-source)
- [☁️ Supported Providers](#️-supported-providers)
- [🚀 Usage](#-usage)
    - [Available Commands](#available-commands)
- [🔑 Provider Credentials](#-provider-credentials)
- [⚙️ Development Guide](#️-development-guide)
- [📌 Contribution](#-contribution)
- [📜 License](#-license)
- [🙌 Acknowledgments](#-acknowledgments)

---

## ✨ Features

{{if .Features}}
{{range .Features}}
- {{.}}
{{else}}
- ✅ Feature 1
- ✅ Feature 2
- ✅ Feature 3
{{end}}

---

## 📥 Installation

### Supported Platforms

{{if .Platforms}}
{{range .Platforms}}
- {{.}}
{{else}}
- **Windows** (if applicable)
- **macOS**
- **Linux**
{{end}}

### 1. Quick Installation

{{if .QuickInstall}}
{{range .QuickInstall}}
{{.}}
{{else}}
`, "```sh", `
 curl -sSL https://example.com/install.sh | sh
`, "```", `
{{end}}

### 2. Homebrew

{{if .Homebrew}}
{{range .Homebrew}}
{{.}}
{{else}}
`, "```sh", `
 brew tap {{.Repo}}
 brew install {{.ProjectName}}
`, "```", `
{{end}}

### 3. Build from Source

{{if .BuildFromSource}}
{{range .BuildFromSource}}
{{.}}
{{else}}
`, "```sh", `
 git clone {{.Repo}}
 cd ./{{.ProjectName}}
 go build
`, "```", `
{{end}}

---

## ☁️ Supported Providers

{{if .Providers}}
{{range .Providers}}
- {{.}}
{{else}}
- AWS
- Google Cloud
- Azure
- Others...
{{end}}

---

## 🚀 Usage

{{if .Usage}}
{{range .Usage}}
{{.}}
{{else}}
`, "```sh", `
 project command1 [flags]
 project command2 [flags]
 project command3 [flags]
`, "```", `
{{end}}

### Available Commands

{{if .Commands}}
{{range .Commands}}
{{.}}
{{else}}
`, "```plaintext", `
 project
 ├── command1  # Description
	 ├── command2  # Description
	 └── command3  # Description
`, "```", `
{{end}}

---

## 🔑 Provider Credentials

To use the project, you need to set up credentials for the providers. You can set them as environment variables or use a configuration file.
Set environment variables:

{{if .EnvVars}}
{{range .EnvVars}}
{{.}}
{{else}}
`, "```sh", `
	 export API_KEY="your-key-here"
`, "```", `
{{end}}

---

## ⚙️ Development Guide

To start developing on the project, follow these steps:

{{if .DevGuide}}
{{range .DevGuide}}
{{.}}
{{else}}
`, "```sh", `
	 echo "dev" > cmd/version
`, "```", `
{{end}}

---

## 📌 Contribution

{{if .Contribution}}
{{range .Contribution}}
{{.}}
{{else}}
We welcome contributions! See the guidelines [here](CONTRIBUTING.md).
{{end}}

---

## 📜 License

{{if .License}}
{{range .License}}
{{.}}
{{else}}
This project is licensed under the [MIT License](LICENSE).
{{end}}

---

## 🙌 Acknowledgments

{{if .Acknowledgments}}
{{range .Acknowledgments}}
{{.}}
{{else}}
Thanks to all contributors and maintainers of the project.
{{end}}

---
`}, "")

	// Fill in the template with the data from the old README
	tmplObj = strings.ReplaceAll(tmplObj, "{{.ProjectName}}", data.ProjectName)
	tmplObj = strings.ReplaceAll(tmplObj, "{{.Org}}", data.Org)
	tmplObj = strings.ReplaceAll(tmplObj, "{{.Repo}}", data.Repo)

	// Parse old README to get the sections and badges, then render the template with the old README data that we get
	var (
		features        []string
		platforms       []string
		quickInstall    []string
		homebrew        []string
		buildFromSource []string
		providers       []string
		usage           []string
		commands        []string
		envVars         []string
		devGuide        []string
		contribution    []string
		license         []string
		acknowledgments []string
	)

	// Fill in the template with the data from the old README
	fillVarsFromReadme(&features, &platforms, &quickInstall, &homebrew, &buildFromSource, &providers, &usage, &commands, &envVars, &devGuide, &contribution, &license, &acknowledgments, readmeFile)

	// Render the template with the old README data
	renderedTmpl, err := renderTemplate(tmplObj, map[string]interface{}{
		"Features":        features,
		"Platforms":       platforms,
		"QuickInstall":    quickInstall,
		"Homebrew":        homebrew,
		"BuildFromSource": buildFromSource,
		"Providers":       providers,
		"Usage":           usage,
		"Commands":        commands,
		"EnvVars":         envVars,
		"DevGuide":        devGuide,
		"Contribution":    contribution,
		"License":         license,
		"Acknowledgments": acknowledgments,
	})
	if err != nil {
		fmt.Println("Error rendering template:", err)
		return "", err
	}

	return renderedTmpl, nil
}

func renderTemplate(tmplObj string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("readme").Parse(tmplObj)
	if err != nil {
		return "", err
	}
	var renderedTmpl strings.Builder
	err = tmpl.Execute(&renderedTmpl, data)
	if err != nil {
		return "", err
	}
	return renderedTmpl.String(), nil
}

func fillVarsFromReadme(features, platforms, quickInstall, homebrew, buildFromSource, providers, usage, commands, envVars, devGuide, contribution, license, acknowledgments *[]string, readmeFile string) {
	// Parse old README to get the sections and badges, then render the template with the old README data that we get
	order, sections, _, err := parseFileOrContent(readmeFile)
	if err != nil {
		fmt.Println("❌ Error reading README:", err)
		return
	}

	// Pause the program to allow the user to check the data
	fmt.Println("Press Enter to continue or Ctrl+C to exit...")
	_, scanlnErr := fmt.Scanln()
	if scanlnErr != nil {
		fmt.Println("Error scanning input:", scanlnErr)
		return
	}

	for _, section := range order {
		switch sectionPattern := strings.ToLower(section); sectionPattern {
		case "features", "feature", "feat":
			*features = sections[section]
		case "platforms", "platform":
			*platforms = sections[section]
		case "quick installation", "quick install", "quick":
			*quickInstall = sections[section]
		case "homebrew":
			*homebrew = sections[section]
		case "build from source", "build":
			*buildFromSource = sections[section]
		case "supported providers", "providers", "provider":
			*providers = sections[section]
		case "usage":
			*usage = sections[section]
		case "available commands", "commands", "command":
			*commands = sections[section]
		case "provider credentials", "credentials", "env vars", "env", "variables":
			*envVars = sections[section]
		case "development guide", "dev guide", "dev", "guide":
			*devGuide = sections[section]
		case "contribution", "contributing", "contributions":
			*contribution = sections[section]
		case "license":
			*license = sections[section]
		case "acknowledgment":
			*acknowledgments = sections[section]
		default:
			fmt.Println("Unknown section:", section)
		}
	}

	return
}

func getDataFromReadme(readmeFile string) (*ReadmeData, error) {
	var readmeData ReadmeData
	var getRepoDetailErr error
	var projectName, org, repo string

	readmeFolder := filepath.Dir(readmeFile)
	gitFolder := filepath.Join(readmeFolder, ".git")
	if _, err := os.Stat(gitFolder); err == nil {
		// Get the project name, org, and repo from the git folder
		projectName, org, repo, getRepoDetailErr = getProjectDetails(gitFolder)
		if getRepoDetailErr != nil {
			return nil, fmt.Errorf("error getting project details from git folder: %v", getRepoDetailErr)
		}
	} else {
		// Ask the user to provide the project name, org, and repo
		projectName, org, repo, getRepoDetailErr = getProjectDetailsFromUser()
		if getRepoDetailErr != nil {
			return nil, fmt.Errorf("error getting project details from user: %v", getRepoDetailErr)
		}
	}

	readmeData = ReadmeData{}
	readmeData.ProjectName = projectName
	readmeData.Org = org
	readmeData.Repo = repo

	return &readmeData, nil
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
	cmdExec := exec.Command("bash", "-c", command)
	output, err := cmdExec.Output()
	if err != nil {
		return "", "", "", fmt.Errorf("error running command: %v", err)
	}
	outputStr := string(output)
	outputStr = strings.TrimSpace(outputStr)
	outputArr := strings.Split(outputStr, " ")
	if len(outputArr) == 3 {
		return outputArr[0], outputArr[1], outputArr[2], nil
	} else {
		fmt.Println("Error getting project details:")
		fmt.Println("Output:", outputStr)
		return "", "", "", fmt.Errorf("error getting project details")
	}
}

func getProjectDetailsFromUser() (string, string, string, error) {
	var projectName, org, repo string

	cobraPromptCmd := exec.Command("cobra", "prompt", "Project Name")
	cobraPromptCmd.Stdin = os.Stdin
	cobraPromptCmd.Stdout = os.Stdout
	cobraPromptCmd.Stderr = os.Stderr
	cobraPromptCmdErr := cobraPromptCmd.Run()

	if cobraPromptCmdErr != nil {
		return "", "", "", cobraPromptCmdErr
	}
	_, scanProjNameErr := fmt.Scanln(&projectName)
	if scanProjNameErr != nil {
		return "", "", "", scanProjNameErr
	}

	cobraPromptCmd = exec.Command("cobra", "prompt", "Organization")
	cobraPromptCmd.Stdin = os.Stdin
	cobraPromptCmd.Stdout = os.Stdout
	cobraPromptCmd.Stderr = os.Stderr
	cobraPromptCmdErr = cobraPromptCmd.Run()
	if cobraPromptCmdErr != nil {
		return "", "", "", cobraPromptCmdErr
	}
	_, scanOrgErr := fmt.Scanln(&org)
	if scanOrgErr != nil {
		return "", "", "", scanOrgErr
	}

	cobraPromptCmd = exec.Command("cobra", "prompt", "Repository")
	cobraPromptCmd.Stdin = os.Stdin
	cobraPromptCmd.Stdout = os.Stdout
	cobraPromptCmd.Stderr = os.Stderr
	cobraPromptCmdErr = cobraPromptCmd.Run()
	if cobraPromptCmdErr != nil {
		return "", "", "", cobraPromptCmdErr
	}
	_, scanRepoErr := fmt.Scanln(&repo)
	if scanRepoErr != nil {
		return "", "", "", scanRepoErr
	}

	return projectName, org, repo, nil
}

func parseFileOrContent(fileOrContent string) ([]string, map[string][]string, []string, error) {
	var reader io.Reader

	if len(fileOrContent) > 255 {
		reader = strings.NewReader(fileOrContent)
	} else {
		file, err := os.Open(fileOrContent)
		if err != nil {
			return nil, nil, nil, err
		}
		defer file.Close()
		reader = file
	}

	order := make([]string, 0)
	sections := make(map[string][]string)
	badges := make([]string, 0)
	var currentSection string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if badgeRegex.MatchString(line) {
			badges = append(badges, line)
		}

		if match := titleRegex.FindStringSubmatch(line); match != nil {
			currentSection = match[2]
			if _, exists := sections[currentSection]; !exists {
				order = append(order, currentSection)
				sections[currentSection] = []string{}
			}
		} else if currentSection != "" {
			sections[currentSection] = append(sections[currentSection], line)
		}
	}

	return order, sections, badges, scanner.Err()
}

func compareAndGenerateImprovedReadme(templateContent, readmeFile, outputFile string) {
	templateOrder, templateSections, templateBadges, err := parseFileOrContent(templateContent)
	if err != nil {
		fmt.Println("❌ Error parsing template content:", err)
		return
	}

	readmeOrder, readmeSections, readmeBadges, err := parseFileOrContent(readmeFile)
	if err != nil {
		fmt.Println("❌ Error reading README:", err)
		return
	}

	fmt.Println("\n🔎 Comparing README with Template...")

	var improvedReadme strings.Builder

	fmt.Println("\n🚀 Checking Badges:")
	for _, badge := range templateBadges {
		if contains(readmeBadges, badge) {
			improvedReadme.WriteString(badge + "\n")
		} else {
			improvedReadme.WriteString("<!-- TODO: Insert missing badge -->\n" + badge + "\n")
			fmt.Printf("  ❌ Missing badge: %s\n", badge)
		}
	}
	improvedReadme.WriteString("\n")

	fmt.Println("\n📌 Checking Sections:")
	for _, section := range templateOrder {
		improvedReadme.WriteString("## " + section + "\n")

		if content, exists := readmeSections[section]; exists {
			if strings.Join(templateSections[section], "\n") != strings.Join(content, "\n") {
				fmt.Printf("  ✏️ Section '%s' needs updates.\n", section)
				improvedReadme.WriteString("<!-- TODO: Review and update this section -->\n")
			}
			improvedReadme.WriteString(strings.Join(content, "\n") + "\n\n")
		} else {
			fmt.Printf("  ❌ Missing section: %s\n", section)
			improvedReadme.WriteString("<!-- TODO: Add missing section -->\n")
			improvedReadme.WriteString(strings.Join(templateSections[section], "\n") + "\n\n")
		}
	}

	for _, section := range readmeOrder {
		if _, exists := templateSections[section]; !exists {
			fmt.Printf("  ❌ Extra section found: %s (not in template)\n", section)
			improvedReadme.WriteString("<!-- TODO: Review if this section should be kept -->\n")
			improvedReadme.WriteString("## " + section + "\n")
			improvedReadme.WriteString(strings.Join(readmeSections[section], "\n") + "\n\n")
		}
	}

	err = writeToFile(outputFile, improvedReadme.String())
	if err != nil {
		fmt.Println("❌ Error writing improved README:", err)
		return
	}

	fmt.Println("\n✅ `IMPROVED_README.md` generated successfully!")
}

func writeToFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "readme-checker",
		Short: "Checks and improves README structure",
		Run: func(cmd *cobra.Command, args []string) {
			templateFile, _ := cmd.Flags().GetString("template")
			readmeFile, _ := cmd.Flags().GetString("readme")
			outputFile, _ := cmd.Flags().GetString("output")

			if templateFile == "" {
				var templateFileErr error
				templateFile, templateFileErr = refReadme(readmeFile)
				if templateFileErr != nil {
					fmt.Println("❌ Error generating template README:", templateFileErr)
					return
				}

				templateFileWriteErr := writeToFile("README_improved.md", templateFile)
				if templateFileWriteErr != nil {
					fmt.Println("❌ Error writing template README:", templateFileWriteErr)
					return
				}

				fmt.Println("📄 Template README generated successfully!")
			} else {
				compareAndGenerateImprovedReadme(templateFile, readmeFile, outputFile)
			}
		},
	}

	rootCmd.Flags().StringP("template", "t", "", "Template README file")
	rootCmd.Flags().StringP("readme", "r", "README_to_check.md", "README file to check")
	rootCmd.Flags().StringP("output", "o", "IMPROVED_README.md", "Output improved README file")

	rootCmd.Execute()
}
