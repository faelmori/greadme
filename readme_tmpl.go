package main

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

const defaultTemplate = `

# {{if .ProjectName}}{{.ProjectName}}{{else}}Project Name{{end}}

{{if .Badges}}
{{range .Badges}}
{{.}}
{{else}}
{{if .Org}}{{if .Repo}}{{if .ProjectName}}
![Version](https://img.shields.io/github/v/release/{{.Org}}/{{.Repo}})
![Build Status](https://img.shields.io/github/actions/workflow/status/{{.Org}}/{{.Repo}}/build.yml?branch=main)
![License](https://img.shields.io/github/license/{{.Org}}/{{.Repo}})
{{else}}
![Version](https://img.shields.io/github/v/release/user/repo)
![Build Status](https://img.shields.io/github/actions/workflow/status/user/repo/build.yml?branch=main)
![License](https://img.shields.io/github/license/user/repo)
{{end}}{{end}}{{end}}
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
` + "```" + `sh
curl -sSL https://example.com/install.sh | sh
` + "```" + `
{{end}}

### 2. Homebrew

{{if .Homebrew}}
{{range .Homebrew}}
{{.}}
{{else}}
` + "```" + `sh
brew tap {{.Repo}}
brew install {{.ProjectName}}
` + "```" + `
{{end}}

### 3. Build from Source

{{if .BuildFromSource}}
{{range .BuildFromSource}}
{{.}}
{{else}}
` + "```" + `sh
git clone {{.Repo}}
cd ./{{.ProjectName}}
go build
` + "```" + `
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
` + "```" + `sh
project command1 [flags]
project command2 [flags]
project command3 [flags]
` + "```" + `
{{end}}

### Available Commands

{{if .Commands}}
{{range .Commands}}
{{.}}
{{else}}
` + "```" + `plaintext
project
├── command1  # Description
├── command2  # Description
└── command3  # Description
` + "```" + `
{{end}}

---

## 🔑 Provider Credentials

To use the project, you need to set up credentials for the providers. You can set them as environment variables or use a configuration file.

Set environment variables:

{{if .EnvVars}}
{{range .EnvVars}}
{{.}}
{{else}}
` + "```" + `sh
export API_KEY="your-key-here"
` + "```" + `
{{end}}

---

## ⚙️ Development Guide

To start developing on the project, follow these steps:

{{if .DevGuide}}
{{range .DevGuide}}
{{.}}
{{else}}
` + "```" + `sh
echo "dev" > cmd/version
` + "```" + `
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

{{end}}
`
