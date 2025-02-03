package main

const codeBlock = string("```")

var defaultTemplate = `

{{if or (ne .Org "") (ne .Repo "") (ne .ProjectName "")}}
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
![Version](https://img.shields.io/github/v/release/{{.Repo}})
![Build Status](https://img.shields.io/github/actions/workflow/status/{{.Repo}}/build.yml?branch=main)
![License](https://img.shields.io/github/license/{{.Repo}})
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
