## **This is just the beginning of a project! Wait a little while for things to work out, or COLLABORATE and help make it real!**

# GReadme


![Version](https://img.shields.io/github/v/release/faelmori/greadme)
![Build Status](https://img.shields.io/github/actions/workflow/status/faelmori/greadme/build.yml?branch=main)
![License](https://img.shields.io/github/license/faelmori/greadme)

A powerful and user-friendly Command Line Interface (CLI) for [GReadme](https://github.com/faelmori/greadme), the simplest way to create and manage READMEs for your projects.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
   - [Supported Platforms](#supported-platforms)
   - [1. Shortcut Installation](#1-shortcut-installation)
   - [2. Homebrew Installation](#2-homebrew-installation)
   - [3. Build from Source](#3-build-from-source)
- [Supported Providers](#supported-providers)
- [Usage](#usage)
   - [Command Overview](#command-overview)
- [Provider Credentials](#provider-credentials)
   - [GitHub](#github)
- [Development Guide](#development-guide)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

---

## Features

- **Simple Installation:** Install GReadme CLI with a single command.
- **Validate your README:** Ensure your README is up-to-date and accurate.
- **Watch for Changes:** Automatically update your README when changes are detected.
- **Improve existing READMEs:** Enhance your README with new features and content.
- **Internationalization:** Generate READMEs in multiple languages.
- **Custom Templates:** Use custom templates to create READMEs for different projects.

---

## Installation

### Supported Platforms

- **macOS**
- **Linux**

### 1. Shortcut Installation

Install GReadme CLI with a single command:

```shell
curl -fsSL get.greadme.io | bash
```

### 2. Homebrew Installation

If Homebrew is not installed, install it first:

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Install GReadme CLI via Homebrew:

```shell
brew tap faelmori/greadme
brew install greadme
```

### 3. Build from Source

For advanced use cases, build and package the binary manually.

#### Requirements

- [Git](https://git-scm.com/downloads)
- [Go](https://go.dev/doc/install)
- [UPX](https://github.com/upx/upx/releases/)

#### Steps

1. **Clone the Repository:**

   ```shell
   git clone https://github.com/faelmori/greadme.git
   ```

2. **Navigate to the Project Directory:**

   ```shell
   cd greadme
   ```

3. **Create a Version Tag (Optional):**

   ```shell
   git tag -a v1.0 -m "Version 1.0"
   ```

4. **Build and Package the Binary:**

   ```shell
   go build -ldflags "-s -w -X main.version=$(git describe --tags --abbrev=0) -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o greadme
   upx greadme --overwrite-output --best -o $(realpath ./)/greadme
   ```

5. **Move the Binary to Your PATH:**

   ```shell
   sudo mv ./greadme /usr/local/bin/
   ```

6. **Reload Shell Configuration:**

   ```shell
   source "$HOME/.$(basename ${SHELL})rc"
   ```

7. **Verify Installation:**

   ```shell
   greadme -h
   ```

---

## Supported Providers

GReadme CLI currently supports the following cloud providers:

- **Local**
- **GitHub**

### Coming Soon

- **AWS**
- **Azure**
- **Google Cloud**
- **Netlify**
- **Vercel**

---

## Usage

### Command Overview

```plaintext
greadme
├── create                 # Create a new README
├── update                 # Update an existing README
├── watch                  # Watch for changes in a README
├── validate               # Validate a README
├── translate              # Translate a README
├── completion             # Generate shell completion scripts
└── help                   # Display help for commands
```

---

## Provider Credentials

Set the appropriate environment variables for your cloud provider before using GReadme CLI.

### GitHub

```shell
export GITHUB_TOKEN=your_access_token
```

---

## Development Guide

### Enable Development Mode

To enable development mode for testing and debugging, set the version to `dev`:

```shell
go build -ldflags "-s -w -X main.version=dev -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o greadme
```

---

## Contributing

We welcome contributions from the community! Please check out our [Contributing Guidelines](https://github.com/faelmori/greadme/blob/main/CONTRIBUTING.md) for more information.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Acknowledgments

- **[GReadme](https://github.com/faelmori/greadme):** The best way to create and manage READMEs for your projects.
- **[Go](https://golang.org/):** The programming language used for development.
- **Community Contributors:** Thank you to all who have contributed to this project.

---

Thank you for using **greadme**! If you have suggestions or encounter issues, please open an issue in the [main repository](https://github.com/faelmori/greadme).

---
