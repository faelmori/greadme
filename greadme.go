package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

var titleRegex = regexp.MustCompile(`^(#+)\s+(.\*)`)
var badgeRegex = regexp.MustCompile(`!\\[.\*\\]\\(https\://img\\.shields\\.io.\*\\)`)

type ReadmeData struct {
	Org             string
	Repo            string
	ProjectName     string
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

func generateImprovedReadme(templateFile, readmeFile, outputFile string) {
	readmeData, err := extractReadmeData(readmeFile)

	if err != nil {
		fmt.Println("❌ Error extracting README data:", err)
		return
	}

	var tmplStr string
	var templateObj *template.Template
	var willUseDefaultTemplate bool

	if templateFile != "" {
		fileContent, fileContentErr := os.ReadFile(templateFile)
		if fileContentErr == nil {
			willUseDefaultTemplate = false
			tmplStr = string(fileContent)
		} else {
			willUseDefaultTemplate = true
		}
	} else {
		willUseDefaultTemplate = true
	}

	if willUseDefaultTemplate {
		tmplStr = defaultTemplate
	}

	tmplObj, tmplObjErr := template.New("readme_template").Parse(tmplStr)
	if tmplObjErr != nil {
		fmt.Println("############################################")
		preLineNumberIndex := strings.Index(tmplObjErr.Error(), "readme_template:")
		lineNumberStr := strings.Split(tmplObjErr.Error()[preLineNumberIndex:], ":")[1]
		lineNumber, lineNumberErr := strconv.Atoi(strings.TrimSpace(lineNumberStr))
		if lineNumberErr != nil {
			fmt.Println("❌ Error parsing template error line number:", lineNumberErr)
			fmt.Println("############################################")
			return
		}
		if lineNumber-1 < 0 {
			fmt.Println("❌ Error parsing template at line 1")
			fmt.Println("############################################")
			return
		}

		fmt.Println("❌ Error parsing template at line", lineNumber)

		tmplStrLinesList := strings.Split(tmplStr, "\n")
		tmplStrLinesListInverse := make([]string, 0)
		maxLines := 5

		for i := lineNumber - 1; i >= 0; i-- {
			if i == lineNumber-1 || i == lineNumber || i <= 0 || maxLines == 0 {
				break
			}
			i--

			if strings.TrimSpace(tmplStrLinesList[i]) != "" {
				tmplStrLinesListInverse = append(tmplStrLinesListInverse, strings.Join([]string{strconv.Itoa(i+1) + ":", tmplStrLinesList[i]}, " "))
				maxLines--
			} else {
				continue
			}
		}

		if len(tmplStrLinesListInverse) > 0 {
			fmt.Println("❌ Template part with error:")
			for i := len(tmplStrLinesListInverse) - 1; i >= 0; i-- {
				fmt.Println(tmplStrLinesListInverse[i])
			}
		} else {
			fmt.Println("❌ Template part with error:")
			fmt.Println(tmplStrLinesList[lineNumber-1])
		}

		fmt.Println("❌ Error:", tmplObjErr)
		fmt.Println("############################################")
		return
	}
	if tmplObj != nil {
		templateObj = tmplObj
	} else {
		fmt.Println("❌ Error parsing template: tmplObj is nil without error")
		return
	}

	outputFileObj, outputFileObjErr := os.Create(outputFile)
	if outputFileObjErr != nil {
		fmt.Println("❌ Error creating output file:", outputFileObjErr)
		return
	}
	defer func(outputFileObj *os.File) {
		_ = outputFileObj.Close()
	}(outputFileObj)

	templateExecuteErr := templateObj.Execute(outputFileObj, readmeData)
	if templateExecuteErr != nil {
		fmt.Println("❌ Error executing template:", templateExecuteErr)
		return
	}

	fmt.Println("✅ `IMPROVED_README.md` generated successfully!")
}

func extractReadmeData(readmeFile string) (*ReadmeData, error) {
	order, sections, _, err := parseFileOrContent(readmeFile)

	if err != nil {
		return nil, err
	}

	data := &ReadmeData{}
	for _, section := range order {
		content := strings.Join(sections[section], "\n")
		if content == "" {
			content = "<!-- TODO: Add content for " + section + " -->"
		}

		switch strings.ToLower(section) {
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

	return data, nil
}

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
	cmdExecErr := rootCmd.Execute()
	if cmdExecErr != nil {
		fmt.Println("❌ Error executing command:", cmdExecErr)
		return
	}
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
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
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
