# config-sync

Sync configuration files across machines using git.

## ⚠️ IMPORTANT SECURITY WARNINGS ⚠️

### Be Careful What You Track

**NEVER track files containing:**
- API keys, tokens, or secrets
- Passwords or credentials
- Private SSH keys
- Sensitive personal data
- Anything you wouldn't want public

**If your config files contain secrets, DO NOT use a public GitHub repository.**

Use a **private** repository or consider using secret management tools (like `envchain`, `1password`, `vault`, etc.) for sensitive data.

### No Liability

This software is provided "as is", without warranty of any kind. The authors and contributors are **not liable** for any damages, data loss, security breaches, or issues arising from the use of this software. You are solely responsible for:

- What files you choose to sync
- Where you store your repository (public vs private)
- Protecting sensitive information
- Backing up your data

## Installation

```bash
go install github.com/mohammed-io/config-sync@v0.0.3
```

Or build from source:

```bash
git clone https://github.com/mohammed-io/config-sync.git
cd config-sync
go build -o config-sync .
```

## Usage

### 1. Set the Git Repository

```bash
config-sync set-origin-repo git@github.com:your-username/your-config-repo.git
```

This initializes a git repository in `~/.config-sync` and sets the remote origin.

### 2. Track Files

```bash
# Track a single file
config-sync track ~/.claude/CLAUDE.md

# Track multiple files
config-sync track ~/.vimrc ~/.tmux.conf ~/.gitconfig
```

### 3. Push to Sync

```bash
config-sync push
```

This copies tracked files to `~/.config-sync/synced-files/`, commits, and pushes to git.

### 4. Pull on Other Machines

```bash
config-sync pull
```

Then manually copy files from `~/.config-sync/synced-files/` to their destinations (or create a restore command).

### Untrack Files

```bash
config-sync untrack ~/.vimrc
```

## Example: Syncing Claude Code Config

```bash
# Set up a private repo (IMPORTANT: use private for sensitive configs)
config-sync set-origin-repo git@github.com:your-username/my-config.git

# Track your Claude Code instructions
config-sync track ~/.claude/CLAUDE.md

# Push to sync
config-sync push
```

On another machine:

```bash
# Clone or pull
config-sync pull

# Files are now in ~/.config-sync/synced-files/
# Copy them to the right location manually
```

## Development

This project uses [mise](https://mise.jdx.dev/) for tool management.

### Setup

```bash
# Install mise (if not already installed)
curl https://mise.run | sh

# Install Go and other tools via mise
mise install
```

### Building

```bash
# Build the binary
go build -o config-sync .

# Or use mise to run tasks
mise run build
```

### Project Structure

```
config-sync/
├── main.go              # CLI commands and main entry point
├── json_config.go       # Config management (JsonConfig)
├── git_runner.go        # Git operations (GitRunner interface)
├── shorthand_path.go    # Path utilities (tilde expansion)
├── README.md
└── LICENSE
```

## How It Works

- Tracked files are stored in `~/.config-sync/synced-files/`
- Each file is placed in a subfolder named after the MD5 hash of its path
- `config.json` tracks which files are being synced
- Git operations run in `~/.config-sync/`

## License

MIT License - see [LICENSE](LICENSE) file for details.
