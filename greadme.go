package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"regexp"
	"strings"
)

var titleRegex = regexp.MustCompile(`^(#+)\s+(.\*)`)
var badgeRegex = regexp.MustCompile(`!\\[.\*\\]\\(https\://img\\.shields\\.io.\*\\)`)

func generateImprovedReadme(templateFile, readmeFile, outputFile string) {
	readmeData, err := extractReadmeData(readmeFile)

	if err != nil {
		fmt.Println("❌ Error extracting README data:", err)
		return
	}

	var tmplStr string
	//var templateObj *template.Template
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
	} else {
		readmeFile = templateFile

		if tmplStr == "" {
			fmt.Println("❌ Error reading template file:", templateFile)
			return
		}
	}

	readmeData, readmeDataErr := extractReadmeData(readmeFile)
	if readmeDataErr != nil {
		fmt.Println("❌ Error extracting README data:", readmeDataErr)
		return
	}
	templateStrFilled, templateStrFilledErr := fillTemplate(readmeData)
	if templateStrFilledErr != nil {
		fmt.Println("❌ Error filling template:", templateStrFilledErr)
		return
	}

	if templateStrFilled != tmplStr {
		outputFileObj, outputFileObjErr := os.Create(outputFile)
		if outputFileObjErr != nil {
			fmt.Println("❌ Error creating output file:", outputFileObjErr)
			return
		}
		defer func(outputFileObj *os.File) {
			_ = outputFileObj.Close()
		}(outputFileObj)

		_, _ = outputFileObj.WriteString(templateStrFilled)
	} else {
		fmt.Println("✅ README structure is already improved!")
		return
	}

	fmt.Println("✅ `IMPROVED_README.md` generated successfully!")

	//tmplObj, tmplObjErr := template.New("readme_template").Parse(tmplStr)
	//if tmplObjErr != nil {
	//	fmt.Println("############################################")
	//	preLineNumberIndex := strings.Index(tmplObjErr.Error(), "readme_template:")
	//	lineNumberStr := strings.Split(tmplObjErr.Error()[preLineNumberIndex:], ":")[1]
	//	lineNumber, lineNumberErr := strconv.Atoi(strings.TrimSpace(lineNumberStr))
	//	if lineNumberErr != nil {
	//		fmt.Println("❌ Error parsing template error line number:", lineNumberErr)
	//		fmt.Println("############################################")
	//		return
	//	}
	//	if lineNumber-1 < 0 {
	//		fmt.Println("❌ Error parsing template at line 1")
	//		fmt.Println("############################################")
	//		return
	//	}
	//
	//	fmt.Println("❌ Error parsing template at line", lineNumber)
	//
	//	tmplStrLinesList := strings.Split(tmplStr, "\n")
	//	tmplStrLinesListInverse := make([]string, 0)
	//	maxLines := 5
	//
	//	for i := lineNumber - 1; i >= 0; i-- {
	//		if i == lineNumber-1 || i == lineNumber || i <= 0 || maxLines == 0 {
	//			break
	//		}
	//		i--
	//
	//		if strings.TrimSpace(tmplStrLinesList[i]) != "" {
	//			tmplStrLinesListInverse = append(tmplStrLinesListInverse, strings.Join([]string{strconv.Itoa(i+1) + ":", tmplStrLinesList[i]}, " "))
	//			maxLines--
	//		} else {
	//			continue
	//		}
	//	}
	//
	//	if len(tmplStrLinesListInverse) > 0 {
	//		fmt.Println("❌ Template part with error:")
	//		for i := len(tmplStrLinesListInverse) - 1; i >= 0; i-- {
	//			fmt.Println(tmplStrLinesListInverse[i])
	//		}
	//	} else {
	//		fmt.Println("❌ Template part with error:")
	//
	//		fmt.Println(tmplStrLinesList[lineNumber-6])
	//		fmt.Println(tmplStrLinesList[lineNumber-5])
	//		fmt.Println(tmplStrLinesList[lineNumber-4])
	//		fmt.Println(tmplStrLinesList[lineNumber-3])
	//		fmt.Println(tmplStrLinesList[lineNumber-2])
	//		fmt.Println(tmplStrLinesList[lineNumber-1])
	//
	//		outputFileObj, outputFileObjErr := os.Create("IMPROVED_README_ERROR.md")
	//		if outputFileObjErr != nil {
	//			fmt.Println("❌ Error creating output file:", outputFileObjErr)
	//			return
	//		}
	//		defer func(outputFileObj *os.File) {
	//			_ = outputFileObj.Close()
	//		}(outputFileObj)
	//
	//		_, _ = outputFileObj.WriteString(tmplStr)
	//
	//		fmt.Println("❌ Error template content written to `IMPROVED_README_ERROR.md`")
	//	}
	//
	//	fmt.Println("❌ Error:", tmplObjErr)
	//	fmt.Println("############################################")
	//	return
	//}
	//if tmplObj != nil {
	//	templateObj = tmplObj
	//} else {
	//	fmt.Println("❌ Error parsing template: tmplObj is nil without error")
	//	return
	//}
	//
	//outputFileObj, outputFileObjErr := os.Create(outputFile)
	//if outputFileObjErr != nil {
	//	fmt.Println("❌ Error creating output file:", outputFileObjErr)
	//	return
	//}
	//defer func(outputFileObj *os.File) {
	//	_ = outputFileObj.Close()
	//}(outputFileObj)
	//
	//templateExecuteErr := templateObj.Execute(outputFileObj, readmeData)
	//if templateExecuteErr != nil {
	//	fmt.Println("❌ Error executing template:", templateExecuteErr)
	//	return
	//}
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
			data.Features = append(data.Features, content)
		case "platforms":
			data.Platforms = append(data.Platforms, content)
		case "quick installation":
			data.QuickInstall = append(data.QuickInstall, content)
		case "homebrew":
			data.Homebrew = append(data.Homebrew, content)
		case "build from source":
			data.BuildFromSource = append(data.BuildFromSource, content)
		case "supported providers":
			data.Providers = append(data.Providers, content)
		case "usage":
			data.Usage = append(data.Usage, content)
		case "available commands":
			data.Commands = append(data.Commands, content)
		case "provider credentials":
			data.EnvVars = append(data.EnvVars, content)
		case "development guide":
			data.DevGuide = append(data.DevGuide, content)
		case "contribution":
			data.Contribution = append(data.Contribution, content)
		case "license":
			data.License = append(data.License, content)
		case "acknowledgments":
			data.Acknowledgments = append(data.Acknowledgments, content)
		}
	}

	return data, nil
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
