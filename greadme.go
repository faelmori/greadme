package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var titleRegex = regexp.MustCompile(`^(#+)\s+(.\*)`)
var badgeRegex = regexp.MustCompile(`!\\[.\*\\]\\(https\://img\\.shields\\.io.\*\\)`)

func ensureTemplateVars(data *ReadmeData) *ReadmeData {
	if data == nil {
		data = &ReadmeData{}
	}

	if data.Features == nil {
		data.Features = []string{"✅ Feature 1", "✅ Feature 2", "✅ Feature 3"}
	}
	if data.Platforms == nil {
		data.Platforms = []string{"- **Windows** (if applicable)", "- **macOS**", "- **Linux**"}
	}
	if data.QuickInstall == nil {
		data.QuickInstall = []string{"```sh\ncurl -sSL https://example.com/install.sh | sh\n```"}
	}
	if data.Homebrew == nil {
		data.Homebrew = []string{"```sh\nbrew tap {{.Repo}}\nbrew install {{.ProjectName}}\n```"}
	}
	if data.BuildFromSource == nil {
		data.BuildFromSource = []string{"```sh\ngit clone {{.Repo}}\ncd ./{{.ProjectName}}\ngo build\n```"}
	}
	if data.Providers == nil {
		data.Providers = []string{"- AWS", "- Google Cloud", "- Azure", "- Others..."}
	}
	if data.Usage == nil {
		data.Usage = []string{"```sh\nproject command1 [flags]\nproject command2 [flags]\nproject command3 [flags]\n```"}
	}
	if data.Commands == nil {
		data.Commands = []string{"```plaintext\nproject\n├── command1  # Description\n├── command2  # Description\n└── command3  # Description\n```"}
	}
	if data.EnvVars == nil {
		data.EnvVars = []string{"```sh\nexport API_KEY=\"your-key-here\"\n```"}
	}
	if data.DevGuide == nil {
		data.DevGuide = []string{"```sh\necho \"dev\" > cmd/version\n```"}
	}
	if data.Contribution == nil {
		data.Contribution = []string{"We welcome contributions! See the guidelines [here](CONTRIBUTING.md)."}
	}
	if data.License == nil {
		data.License = []string{"This project is licensed under the [MIT License](LICENSE)."}
	}
	if data.Acknowledgments == nil {
		data.Acknowledgments = []string{"Thanks to all contributors and maintainers of the project."}
	}

	return data
}

func fillTemplate(data *ReadmeData) (string, error) {
	tmpl := template.New("readme")
	if tmpl == nil {
		fmt.Println("❌ Error creating template on first method")
		return "", fmt.Errorf("error creating template")
	}

	tmplParsed, tmplParsedErr := tmpl.Parse(defaultTemplate)
	if tmplParsedErr != nil {
		fmt.Println("❌ Error parsing template on first method:", tmplParsedErr)
		return "", tmplParsedErr
	}

	var tpl bytes.Buffer
	executeErr := tmplParsed.Execute(&tpl, data)
	if executeErr != nil {
		fmt.Println("❌ Error executing template on first method:", executeErr)
		return "", executeErr
	}

	return tpl.String(), nil
}

func generateImprovedReadme(templateFile string, readmeFile string, outputFile string) {
	readmeData, err := extractReadmeData(readmeFile)

	readmeData = ensureTemplateVars(readmeData)

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
}

func logReadDataSummary(order []string, sections map[string][]string, badges []string) error {
	// log the summary of the sections and badges in a README_LOG.txt file
	logFile, logFileErr := os.Create("README_LOG.txt")
	if logFileErr != nil {
		fmt.Println("❌ Error creating log file:", logFileErr)
		return logFileErr
	}
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)

	_, _ = logFile.WriteString("📝 Order (" + fmt.Sprint(len(order)) + "):\n")
	for _, section := range order {
		_, _ = logFile.WriteString(" - " + section + "\n")
	}

	_, _ = logFile.WriteString("📄 Sections (" + fmt.Sprint(len(sections)) + "):\n")
	for key, value := range sections {
		_, _ = logFile.WriteString(" - " + key + " (" + fmt.Sprint(len(value)) + " lines)\n")
		for _, line := range value {
			_, _ = logFile.WriteString("   " + line + "\n")
		}
	}

	_, _ = logFile.WriteString("🏷️ Badges (" + fmt.Sprint(len(badges)) + "):\n")
	for _, badge := range badges {
		_, _ = logFile.WriteString(" - " + badge + "\n")
	}

	return nil
}

func extractReadmeData(readmeFile string) (*ReadmeData, error) {
	order, sections, badges, err := parseFileOrContent(readmeFile)
	if err != nil {
		return nil, err
	}

	_ = logReadDataSummary(order, sections, badges)

	data := &ReadmeData{}
	for _, section := range order {
		for _, badge := range badges {
			data.Badges = append(data.Badges, badge)
		}

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
