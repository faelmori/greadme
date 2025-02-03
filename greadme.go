package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var titleRegex = regexp.MustCompile(`^(#+)\s+(.*)`)
var badgeRegex = regexp.MustCompile(`!\[.*\]\(https://img\.shields\.io.*\)`)

func parseFile(filename string) (map[string][]string, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	sections := make(map[string][]string)
	var badges []string
	badges = make([]string, 0)
	var currentSection string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Capture badges
		if badgeRegex.MatchString(line) {
			badges = append(badges, line)
		}

		// Capture section titles
		if match := titleRegex.FindStringSubmatch(line); match != nil {
			currentSection = match[2] // Section name
			sections[currentSection] = []string{}
		} else if currentSection != "" {
			sections[currentSection] = append(sections[currentSection], line)
		}
	}

	return sections, badges, scanner.Err()
}

func parseContent(content string) (map[string][]string, []string, error) {
	sections := make(map[string][]string)
	var badges []string
	badges = make([]string, 0)
	var currentSection string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Capture badges
		if badgeRegex.MatchString(line) {
			badges = append(badges, line)
		}

		// Capture section titles
		if match := titleRegex.FindStringSubmatch(line); match != nil {
			currentSection = match[2] // Section name
			sections[currentSection] = []string{}
		} else if currentSection != "" {
			sections[currentSection] = append(sections[currentSection], line)
		}
	}

	return sections, badges, scanner.Err()
}

func compareAndGenerateImprovedReadme(readmeFile, outputFile string) {
	templateSections, templateBadges, err := parseContent(refReadme)
	if err != nil {
		fmt.Println("❌ Error reading template:", err)
		return
	}

	readmeSections, readmeBadges, err := parseFile(readmeFile)
	if err != nil {
		fmt.Println("❌ Error reading README:", err)
		return
	}

	fmt.Println("\n🔎 Comparing README with Template...")

	var improvedReadme strings.Builder

	// Add badges to improved README in the correct order
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

	// Check sections and maintain order
	fmt.Println("\n📌 Checking Sections:")
	for section, templateContent := range templateSections {
		improvedReadme.WriteString("## " + section + "\n")

		if content, exists := readmeSections[section]; exists {
			if strings.Join(templateContent, "\n") != strings.Join(content, "\n") {
				fmt.Printf("  ✏️ Section '%s' needs updates.\n", section)
				improvedReadme.WriteString("<!-- TODO: Review and update this section -->\n")
			}
			improvedReadme.WriteString(strings.Join(content, "\n") + "\n\n")
		} else {
			fmt.Printf("  ❌ Missing section: %s\n", section)
			improvedReadme.WriteString("<!-- TODO: Add missing section -->\n")
			improvedReadme.WriteString(strings.Join(templateContent, "\n") + "\n\n")
		}
	}

	// Save the improved README
	err = os.WriteFile(outputFile, []byte(improvedReadme.String()), 0644)
	if err != nil {
		fmt.Println("❌ Error writing improved README:", err)
		return
	}

	fmt.Println("\n✅ `IMPROVED_README.md` generated successfully!")
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
			readmeFile, _ := cmd.Flags().GetString("readme")
			outputFile, _ := cmd.Flags().GetString("output")

			compareAndGenerateImprovedReadme(readmeFile, outputFile)
		},
	}

	rootCmd.Flags().StringP("readme", "r", "README_to_check.md", "README file to check")
	rootCmd.Flags().StringP("output", "o", "IMPROVED_README.md", "Output improved README file")

	rootCmd.Execute()
}
