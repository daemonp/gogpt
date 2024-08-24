# gogpt

## Project Overview

`gogpt` is a command-line tool written in Go that enables you to export the structure and content of a Git repository into a text format. This format is tailored for consumption by AI language models, making it ideal for tasks like automated code reviews, documentation generation, or code analysis. The tool offers a range of customization options to include or exclude specific file types, respect `.gitignore` rules, and handle large files appropriately.

## Project Objectives

- Selective File Inclusion: Use system flags to specify which programming languages' files should be included in the export.
- .gitignore Compliance: Optionally ignore files listed in the project's `.gitignore` files.
- Large File Management: Exclude large files from the output, providing warnings with details about the excluded files.
- Automatic Language Detection: When no specific languages are provided, automatically detect the programming languages used in the repository.
- Human-Readable Logs: Utilize `Zerolog` to provide styled, human-readable logging, with default behavior tailored for both terminal and non-terminal outputs.

## How to Build and Run

To build the `gogpt` CLI tool, navigate to the root of the project and run:

```bash
go build -o gogpt ./cmd/gogpt
```

This will generate an executable named `gogpt` in the root directory.

### Running the Tool

After building the project, you can run the tool with various options:

```bash
./gogpt [options]
```

#### Common Flags

- `-f`: Specify the output file path (default: stdout).
- `-i`: Ignore files listed in `.gitignore`.
- `-l`: Comma-separated list of languages to include (e.g., `go,js,md`).
- `--max-tokens`: Maximum number of tokens per file (default: 1000).
- `-v`: Enable verbose logging.

#### Example Usage

1. Basic Usage
   ```bash
   ./gogpt -l go,js -f output.txt
   ```

   ```bash
   ./gogpt -l go,js | wl-copy
   ```

2. Ignore Files in .gitignore:
   ```bash
   ./gogpt -l go,js -i
   ```

3. Automatic Language Detection:
   ```bash
   ./gogpt
   ```

## Logging

By default, logs are output in a human-readable format to `stderr`. If the output is being piped, logs are adjusted for non-terminal environments.

## Summary

`gogpt` provides a flexible and efficient way to extract and format the contents of a Git repository for use in AI-related projects. It includes intelligent handling of file selection, large files, and repository-specific rules to ensure that the exported data is relevant and manageable.

