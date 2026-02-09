# config-sync

Do you have multiple machines and felt the fatigue of syncing your config across all of them? The config no longer in one `~/.config` folder, now they could be your `~/.claude/CLAUDE.md` and its skills, or across multiple files and folders. This tool solves your problem.

## Limitations

- You cannot change the config destination, it always reads the config from `~/.config`

## ⚠️ IMPORTANT SECURITY WARNINGS ⚠️

### Be Careful What You Track

**NEVER track files containing:**
- API keys, tokens, or secrets
- Passwords or credentials
- Private SSH keys
- Sensitive personal data
- Anything you wouldn't want to be public or shared by others

**If your config files contain secrets, DO NOT use a public GitHub repository.**

Use a **private** repository OR consider using secret management tools (like `envchain`, `1password`, `vault`, etc.) for sensitive data.

### Optional: Encrypt Your Repository

If your configs contain sensitive data and you want extra protection, consider encrypting your repository with [git-crypt](https://github.com/AGWA/git-crypt).

Once enabled, files are encrypted on `push` and decrypted on `pull` automatically. Your `config-sync` commands work exactly the same - no changes needed.

**First time setup:**

```bash
# Install git-crypt
brew install git-crypt           # macOS
sudo apt install git-crypt       # Ubuntu/Debian

# Go to your sync folder and enable encryption
cd ~/.config-sync
git-crypt init

# Mark synced files for encryption
echo "synced-files/** filter=git-crypt diff=git-crypt" >> .gitattributes

# Save the key somewhere safe (you'll need it on other machines)
git-crypt export-key ~/config-sync-key
```

**On each of your other machines:**

```bash
# After cloning, unlock the repository with your key
cd ~/.config-sync
git-crypt unlock ~/config-sync-key
```

> **What this protects:** Your files are encrypted in GitHub and in git history. They're still decrypted on your machine, so keep your `config-sync-key` file safe and don't share it.

### No Liability

This software is provided "as is", without warranty of any kind. The authors and contributors are **not liable** for any damages, data loss, security breaches, or issues arising from the use of this software. You are solely responsible for:

- What files you choose to sync
- Where you store your repository (public vs private)
- Protecting sensitive information
- Backing up your data

## Installation

```bash
go install github.com/mohammed-io/config-sync@v0.0.15
```

Or build from source:

```bash
git clone https://github.com/mohammed-io/config-sync.git
cd config-sync
go build -o config-sync .
```

## Usage

### First Time Setup

**Option A: Start fresh on this machine**

```bash
config-sync init
```

This creates the directory structure and local git repository at `~/.config-sync`.

**Option B: Clone existing repository**

```bash
config-sync init-from git@github.com:your-username/your-config-repo.git
```

This clones an existing config-sync repository to `~/.config-sync`.

### Set Remote Repository (for fresh installs)

```bash
config-sync set-origin-repo git@github.com:your-username/your-config-repo.git
```

### Track Files

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

### Pull on Other Machines

```bash
config-sync pull
```

This restores files from `~/.config-sync/synced-files/` to their original locations.

### Untrack Files

```bash
config-sync untrack ~/.vimrc
```

### Check for Updates

```bash
config-sync check-updates
```

Lightweight check for sync status. Exits silently if up to date or not initialized. Shows a message if you need to pull or push.

**Example output when out of sync:**
```
You have local changes not pushed

Run: config-sync push
```

## Example: Syncing Claude Code Config

**First machine:**

```bash
# Initialize (fresh start)
config-sync init

# Set up a private repo (IMPORTANT: use private for sensitive configs)
config-sync set-origin-repo git@github.com:your-username/my-config.git

# Track your Claude Code instructions
config-sync track ~/.claude/CLAUDE.md

# Push to sync
config-sync push
```

**On another machine:**

```bash
# Clone the existing repository
config-sync init-from git@github.com:your-username/my-config.git

# Pull and restore files
config-sync pull
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
