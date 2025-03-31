# Commit Helper

A Go tool that automatically generates commit messages for your Git repository changes using AI. This tool analyzes modified and newly created files, generates appropriate commit messages, and commits them to your repository with a single command.

## Features

- Automatically detects modified files in your Git repository
- Discovers newly created files not yet tracked by Git
- Generates meaningful commit messages using AI for each file
- Works on Windows, macOS, and Linux
- Handles directories and binary files gracefully

## Installation

### Prerequisites

- Go 1.16 or higher
- Git installed and available in your PATH
- OpenAI API key (get one at https://platform.openai.com/api-keys)

### Using Go Install

```bash
go install github.com/ygorazambuja/commit-helper@latest
```

### From Source

```bash
git clone https://github.com/ygorazambuja/commit-helper.git
cd commit-helper
go build
```

## Configuration

Before using Commit Helper, you need to set your OpenAI API key as an environment variable:

```bash
# For Linux/macOS
export OPENAI_API_KEY="your-api-key-here"

# For Windows (Command Prompt)
set OPENAI_API_KEY=your-api-key-here

# For Windows (PowerShell)
$env:OPENAI_API_KEY="your-api-key-here"
```

You can add this to your shell profile for permanent configuration.

## Usage

Simply run the command from your Git repository:

```bash
commit-helper
```

The tool will:
1. Detect all modified files in your working directory
2. Find all new files not yet tracked by Git
3. Generate an AI-powered commit message for each file
4. Automatically add and commit each file with its generated message

## How It Works

Commit Helper uses:
- Git commands to detect changes in your repository
- OpenAI's GPT models to analyze code changes and generate commit messages in Portuguese (BR)
- Clean code principles for maintainability and portability

## Limitations

- Each file is committed separately with its own commit message
- Large binary files might not get meaningful commit messages
- Requires an active internet connection for AI processing
- Requires a valid OpenAI API key with sufficient credits

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
