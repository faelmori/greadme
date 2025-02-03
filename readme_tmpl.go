package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type ReadmeData struct {
	Org             string
	Repo            string
	ProjectName     string
	Badges          []string
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

const codeBlock = string("```")

// var defaultTemplate = fmt.Sprintf("\n\n# {{if .ProjectName}}{{.ProjectName}}{{else}}Project Name{{end}}\n\n{{if .Badges}}\n{{range .Badges}}\n{{.}}\n{{else}}\n{{if .Org}}{{if .Repo}}{{if .ProjectName}}\n![Version](https://img.shields.io/github/v/release/{{.Org}}/{{.Repo}})\n![Build Status](https://img.shields.io/github/actions/workflow/status/{{.Org}}/{{.Repo}}/build.yml?branch=main)\n![License](https://img.shields.io/github/license/{{.Org}}/{{.Repo}})\n{{else}}\n![Version](https://img.shields.io/github/v/release/user/repo)\n![Build Status](https://img.shields.io/github/actions/workflow/status/user/repo/build.yml?branch=main)\n![License](https://img.shields.io/github/license/user/repo)\n{{end}}{{end}}{{end}}\n{{end}}\n\nA brief description of the project, explaining its purpose and main features.\n\n> **Note:** Important instructions or warnings can be highlighted here.\n\n---\n\n## 📖 Table of Contents\n\n- [✨ Features](#-features)\n- [📥 Installation](#-installation)\n   - [Supported Platforms](#supported-platforms)\n   - [1. Quick Installation](#1-quick-installation)\n   - [2. Homebrew](#2-homebrew)\n   - [3. Build from Source](#3-build-from-source)\n- [☁️ Supported Providers](#️-supported-providers)\n- [🚀 Usage](#-usage)\n   - [Available Commands](#available-commands)\n- [🔑 Provider Credentials](#-provider-credentials)\n- [⚙️ Development Guide](#️-development-guide)\n- [📌 Contribution](#-contribution)\n- [📜 License](#-license)\n- [🙌 Acknowledgments](#-acknowledgments)\n\n---\n\n## ✨ Features\n\n{{if .Features}}\n{{range .Features}}\n- {{.}}\n{{else}}\n- ✅ Feature 1\n- ✅ Feature 2\n- ✅ Feature 3\n{{end}}\n\n---\n\n## 📥 Installation\n\n### Supported Platforms\n\n{{if .Platforms}}\n{{range .Platforms}}\n- {{.}}\n{{else}}\n- **Windows** (if applicable)\n- **macOS**\n- **Linux**\n{{end}}\n\n### 1. Quick Installation\n\n{{if .QuickInstall}}\n{{range .QuickInstall}}\n{{.}}\n{{else}}\n\n%ssh\ncurl -sSL https://example.com/install.sh | sh\n%s\n{{end}}\n\n### 2. Homebrew\n\n{{if .Homebrew}}\n{{range .Homebrew}}\n{{.}}\n{{else}}\n%ssh\nbrew tap {{.Repo}}\nbrew install {{.ProjectName}}\n%s\n{{end}}\n\n### 3. Build from Source\n\n{{if .BuildFromSource}}\n{{range .BuildFromSource}}\n{{.}}\n{{else}}\n%ssh\ngit clone {{.Repo}}\ncd ./{{.ProjectName}}\ngo build\n%s\n{{end}}\n\n---\n\n## ☁️ Supported Providers\n\n{{if .Providers}}\n{{range .Providers}}\n- {{.}}\n{{else}}\n- AWS\n- Google Cloud\n- Azure\n- Others...\n{{end}}\n\n---\n\n## 🚀 Usage\n\n{{if .Usage}}\n{{range .Usage}}\n{{.}}\n{{else}}\n%ssh\nproject command1 [flags]\nproject command2 [flags]\nproject command3 [flags]\n%s\n{{end}}\n\n### Available Commands\n\n{{if .Commands}}\n{{range .Commands}}\n{{.}}\n{{else}}\n%splaintext\nproject\n├── command1  # Description\n├── command2  # Description\n└── command3  # Description\n%s\n{{end}}\n\n---\n\n## 🔑 Provider Credentials\n\nTo use the project, you need to set up credentials for the providers. You can set them as environment variables or use a configuration file.\n\nSet environment variables:\n\n{{if .EnvVars}}\n{{range .EnvVars}}\n{{.}}\n{{else}}\n%ssh\nexport API_KEY=\"your-key-here\"\n%s\n{{end}}\n\n---\n\n## ⚙️ Development Guide\n\nTo start developing on the project, follow these steps:\n\n{{if .DevGuide}}\n{{range .DevGuide}}\n{{.}}\n{{else}}\n%ssh\necho \"dev\" > cmd/version\n%s\n{{end}}\n\n---\n\n## 📌 Contribution\n\n{{if .Contribution}}\n{{range .Contribution}}\n{{.}}\n{{else}}\nWe welcome contributions! See the guidelines [here](CONTRIBUTING.md).\n{{end}}\n\n---\n\n## 📜 License\n\n{{if .License}}\n{{range .License}}\n{{.}}\n{{else}}\nThis project is licensed under the [MIT License](LICENSE).\n{{end}}\n\n---\n\n## 🙌 Acknowledgments\n\n{{if .Acknowledgments}}\n{{range .Acknowledgments}}\n{{.}}\n{{else}}\nThanks to all contributors and maintainers of the project.\n{{end}}\n\n{{end}}\n",
var defaultTemplate = `

{{if and (ne .Org "") (ne .Repo "") (ne .ProjectName "")}}
# {{.ProjectName}}
{{else}}
# Project Name
{{end}}

{{if gt (len .Badges) 0}}
	{{range .Badges}}
{{.}}
	{{end}}
{{else}}
	{{if and (ne .Org "") (ne .Repo "") (ne .ProjectName "")}}
![Version](https://img.shields.io/github/v/release/{{.Org}}/{{.Repo}})
![Build Status](https://img.shields.io/github/actions/workflow/status/{{.Org}}/{{.Repo}}/build.yml?branch=main)
![License](https://img.shields.io/github/license/{{.Org}}/{{.Repo}})
	{{else}}
![Version](https://img.shields.io/github/v/release/user/repo)
![Build Status](https://img.shields.io/github/actions/workflow/status/user/repo/build.yml?branch=main)
![License](https://img.shields.io/github/license/user/repo)
	{{end}}
{{end}}

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

{{if gt (len .Features) 0}}
	{{range .Features}}
{{.}}
	{{end}}
{{else}}
- ✅ Feature 1
- ✅ Feature 2
- ✅ Feature 3
{{end}}

---

## 📥 Installation

### Supported Platforms

{{if gt (len .Platforms) 0}}
	{{range .Platforms}}
{{.}}
	{{end}}
{{else}}
- **Windows** (if applicable)
- **macOS**
- **Linux**
{{end}}

### 1. Quick Installation

{{if gt (len .QuickInstall) 0}}
	{{range .QuickInstall}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `sh
curl -sSL https://example.com/install.sh | sh
` + codeBlock + `
{{end}}

### 2. Homebrew

{{if gt (len .Homebrew) 0}}
	{{range .Homebrew}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `sh
brew tap {{.Repo}}
brew install {{.ProjectName}}
` + codeBlock + `
{{end}}

### 3. Build from Source

{{if gt (len .BuildFromSource) 0}}
	{{range .BuildFromSource}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `sh
git clone {{.Repo}}
cd ./{{.ProjectName}}
go build
` + codeBlock + `
{{end}}

---

## ☁️ Supported Providers

{{if gt (len .Providers) 0}}
	{{range .Providers}}
{{.}}
	{{end}}
{{else}}
- AWS
- Google Cloud
- Azure
- Others...
{{end}}

---

## 🚀 Usage

{{if gt (len .Usage) 0}}
	{{range .Usage}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `sh
project command1 [flags]
project command2 [flags]
project command3 [flags]
` + codeBlock + `
{{end}}

### Available Commands

{{if gt (len .Commands) 0}}
	{{range .Commands}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `plaintext
project
├── command1  # Description
├── command2  # Description
└── command3  # Description
` + codeBlock + `
{{end}}

---

## 🔑 Provider Credentials

To use the project, you need to set up credentials for the providers. You can set them as environment variables or use a configuration file.

Set environment variables:

{{if gt (len .EnvVars) 0}}
	{{range .EnvVars}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `sh
export API_KEY="your-key-here"
` + codeBlock + `
{{end}}

---

## ⚙️ Development Guide

To start developing on the project, follow these steps:

{{if gt (len .DevGuide) 0}}
	{{range .DevGuide}}
{{.}}
	{{end}}
{{else}}
` + codeBlock + `sh
echo "dev" > cmd/version
` + codeBlock + `
{{end}}

---

## 📌 Contribution

{{if gt (len .Contribution) 0}}
	{{range .Contribution}}
{{.}}
	{{end}}
{{else}}
We welcome contributions! See the guidelines [here](CONTRIBUTING.md).
{{end}}

---

## 📜 License

{{if gt (len .License) 0}}
	{{range .License}}
{{.}}
	{{end}}
{{else}}
This project is licensed under the [MIT License](LICENSE).
{{end}}

---

## 🙌 Acknowledgments

{{if gt (len .Acknowledgments) 0}}
	{{range .Acknowledgments}}
{{.}}
	{{end}}
{{else}}
Thanks to all contributors and maintainers of the project.
{{end}}

---
`

// var defaultTemplate = fmt.Sprintf("\n\n# {{if .ProjectName}}{{.ProjectName}}{{else}}Project Name{{end}}\n\n{{if .Badges}}\n{{range .Badges}}\n{{.}}\n{{else}}\n{{if .Org}}{{if .Repo}}{{if .ProjectName}}\n![Version](https://img.shields.io/github/v/release/{{.Org}}/{{.Repo}})\n![Build Status](https://img.shields.io/github/actions/workflow/status/{{.Org}}/{{.Repo}}/build.yml?branch=main)\n![License](https://img.shields.io/github/license/{{.Org}}/{{.Repo}})\n{{else}}\n![Version](https://img.shields.io/github/v/release/user/repo)\n![Build Status](https://img.shields.io/github/actions/workflow/status/user/repo/build.yml?branch=main)\n![License](https://img.shields.io/github/license/user/repo)\n{{end}}{{end}}{{end}}\n{{end}}\n\nA brief description of the project, explaining its purpose and main features.\n\n> **Note:** Important instructions or warnings can be highlighted here.\n\n---\n\n## 📖 Table of Contents\n\n- [✨ Features](#-features)\n- [📥 Installation](#-installation)\n   - [Supported Platforms](#supported-platforms)\n   - [1. Quick Installation](#1-quick-installation)\n   - [2. Homebrew](#2-homebrew)\n   - [3. Build from Source](#3-build-from-source)\n- [☁️ Supported Providers](#️-supported-providers)\n- [🚀 Usage](#-usage)\n   - [Available Commands](#available-commands)\n- [🔑 Provider Credentials](#-provider-credentials)\n- [⚙️ Development Guide](#️-development-guide)\n- [📌 Contribution](#-contribution)\n- [📜 License](#-license)\n- [🙌 Acknowledgments](#-acknowledgments)\n\n---\n\n## ✨ Features\n\n{{if .Features}}\n{{range .Features}}\n- {{.}}\n{{else}}\n
func ensureTemplateVars(data *ReadmeData) {
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
}

func fillTemplate(data *ReadmeData) (string, error) {
	ensureTemplateVars(data)

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

	// save to file
	outputFileObj, err := os.Create("IMPROVED_README.md")
	if err != nil {
		fmt.Println("❌ Error creating output file:", err)
		return "", err
	}
	defer func(outputFileObj *os.File) {
		_ = outputFileObj.Close()
	}(outputFileObj)

	_, err = outputFileObj.WriteString(tpl.String())
	if err != nil {
		fmt.Println("❌ Error writing to output file:", err)
		return "", err
	}

	return tpl.String(), nil
}
