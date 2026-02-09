package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// Global config instance
var appConfig JsonConfig

// configFolder returns the ShorthandPath for ~/.config-sync
func configFolder() ShorthandPath {
	return ShorthandPath{}.New("~/.config-sync")
}

// Cobra commands
var trackCmd = &cobra.Command{
	Use:   "track [files...]",
	Short: "Add files to sync config",
	Long: "Add files to be tracked and synced across machines.\n\n" +
		"WARNING: Be careful not to track files containing secrets, API keys, passwords,\n" +
		"or sensitive data. These files will be stored in a git repository and potentially\n" +
		"shared with others. Only track configuration files that are safe to be public or\n" +
		"shared within your trusted team.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := appConfig.Track(args); err != nil {
			log.Fatalf("Track failed: %v", err)
		}
	},
}

var untrackCmd = &cobra.Command{
	Use:   "untrack [files...]",
	Short: "Remove files from sync config",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := appConfig.Untrack(args); err != nil {
			log.Fatalf("Untrack failed: %v", err)
		}
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Git pull and restore files to their locations",
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()
		if err := git.Pull(); err != nil {
			log.Fatalf("Pull failed: %v", err)
		}
		if err := appConfig.RestoreFiles(); err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
		log.Println("Pull and restore completed successfully")
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Sync tracked files and push to git",
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()

		// Auto-init git repo if needed
		if err := git.Init(); err != nil {
			log.Fatalf("Git init failed: %v", err)
		}

		// Sync files to synced-folder
		if err := appConfig.SyncFiles(); err != nil {
			log.Fatalf("Sync failed: %v", err)
		}

		// Git add, commit, push
		commitMsg := fmt.Sprintf("config-sync: update files [%s]", time.Now().UTC().Format(time.RFC3339))
		if err := git.Add(); err != nil {
			log.Fatalf("Git add failed: %v", err)
		}
		if err := git.Commit(commitMsg); err != nil {
			log.Fatalf("Git commit failed: %v", err)
		}
		if err := git.Push(); err != nil {
			log.Fatalf("Push failed: %v", err)
		}

		log.Println("Push completed successfully")
	},
}

var setOriginCmd = &cobra.Command{
	Use:   "set-origin-repo <url>",
	Short: "Set git remote origin",
	Args:  cobra.ExactArgs(1),
	Long: "Set the git remote origin for the synced-files repository.\n\n" +
		"This initializes a git repository in ~/.config-sync/synced-files if it doesn't exist,\n" +
		"then sets the remote origin to the provided URL.\n\n" +
		"Example:\n  config-sync set-origin-repo git@github.com:user/config-repo.git",
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()
		force, _ := cmd.Flags().GetBool("force")
		if err := git.SetOrigin(args[0], force); err != nil {
			log.Fatalf("Set origin failed: %v", err)
		}
		log.Printf("Origin set to: %s\n", args[0])
		log.Println("You can now use 'config-sync push' to sync your files.")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config-sync for the first time",
	Long: "Initialize config-sync by creating the necessary directory structure and git repository.\n\n" +
		"This creates:\n" +
		"  - ~/.config-sync/ directory\n" +
		"  - config.json for tracking files\n" +
		"  - synced-files/ directory for your configs\n" +
		"  - Local git repository\n\n" +
		"After init, use 'set-origin-repo' to connect to a remote repository.",
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()

		// Initialize git repo
		if err := git.Init(); err != nil {
			log.Fatalf("Git init failed: %v", err)
		}

		// Create config
		if err := appConfig.Create(configFolder()); err != nil {
			log.Fatalf("Config creation failed: %v", err)
		}

		log.Printf("\n✓ config-sync initialized at %s\n", configFolder().TildePath)
		log.Printf("Next steps:")
		log.Printf("  config-sync set-origin-repo <url>  # Connect to a remote repo")
		log.Printf("  config-sync track ~/.vimrc          # Start tracking files")
		log.Printf("  config-sync push                    # Push to remote")
	},
}

var initFromCmd = &cobra.Command{
	Use:   "init-from <url>",
	Short: "Clone an existing config-sync repository",
	Long: "Clone an existing config-sync repository from a git URL.\n\n" +
		"This will clone the repository into ~/.config-sync, making it ready to use.\n" +
		"Use this on a new machine to quickly set up config-sync.\n\n" +
		"Example:\n  config-sync init-from git@github.com:user/config-repo.git",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()
		force, _ := cmd.Flags().GetBool("force")

		// Security check: warn if repo appears to be public
		if !force {
			if httpsURL := sshToHTTPS(args[0]); httpsURL != "" {
				if isPublicRepo(httpsURL) {
					log.Printf("⚠️  WARNING: Repository appears to be publicly accessible!\n")
					log.Printf("If this repo contains sensitive configs, use --force only if you understand the risks.")
					return
				}
			}
		}

		if err := git.Clone(args[0]); err != nil {
			log.Fatalf("Clone failed: %v", err)
		}

		// Load the config after cloning
		if err := appConfig.Initialize(configFolder()); err != nil {
			if os.IsNotExist(err) {
				// Config doesn't exist in cloned repo, create it
				log.Printf("Config not found in cloned repository, creating new config...")
				if err := appConfig.Create(configFolder()); err != nil {
					log.Fatalf("Config creation failed: %v", err)
				}
			} else {
				log.Fatalf("Config initialization failed: %v", err)
			}
		}

		log.Printf("\n✓ Repository cloned successfully!")
		log.Printf("You can now use:")
		log.Printf("  config-sync pull    # To sync files from the repo")
		log.Printf("  config-sync track  # To add new files to track")
	},
}

var checkUpdatesCmd = &cobra.Command{
	Use:   "check-updates",
	Short: "Check if config is out of sync",
	Long: "Check if there are unpushed local changes or unpulled remote changes.\n\n" +
		"This is a lightweight check that doesn't modify any files.\n" +
		"Useful for running in shell prompts or startup scripts.\n\n" +
		"Exits silently if not initialized.",
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()

		// Check if git repo exists
		if _, err := os.Stat(filepath.Join(configFolder().FullPath, ".git")); os.IsNotExist(err) {
			return // Silent exit if not initialized
		}

		hasUnpushed, err := git.HasUnpushedChanges()
		if err != nil {
			return // Silent on error
		}

		hasUnpulled, err := git.HasUnpulledChanges()
		if err != nil {
			return // Silent on error
		}

		if !hasUnpushed && !hasUnpulled {
			return // Silent exit - everything up to date
		}

		// Build message
		var msgs []string
		if hasUnpushed && hasUnpulled {
			msgs = append(msgs, "Your config is out of sync:")
			msgs = append(msgs, "  • You have local changes not pushed")
			msgs = append(msgs, "  • There are remote changes not pulled")
			msgs = append(msgs, "")
			msgs = append(msgs, "Run: config-sync pull && config-sync push")
		} else if hasUnpushed {
			msgs = append(msgs, "You have local changes not pushed")
			msgs = append(msgs, "")
			msgs = append(msgs, "Run: config-sync push")
		} else if hasUnpulled {
			msgs = append(msgs, "There are remote changes not pulled")
			msgs = append(msgs, "")
			msgs = append(msgs, "Run: config-sync pull")
		}

		for _, msg := range msgs {
			log.Println(msg)
		}
	},
}

func init() {
	setOriginCmd.Flags().Bool("force", false, "Bypass public repository warning")
	initFromCmd.Flags().Bool("force", false, "Bypass public repository warning")
}

var rootCmd = &cobra.Command{
	Use:   "config-sync",
	Short: "Sync config files across machines",
	Long: `config-sync helps you track and sync configuration files across machines.
Files are stored in ~/.config-sync/synced-files and can be managed with git.

Version: ` + Version,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization check for init, init-from, check-updates, help, completion, and version commands
		skipInitCheck := map[string]bool{
			"init":          true,
			"init-from":     true,
			"check-updates": true,
			"help":          true,
			"completion":    true,
			"version":       true,
		}
		if skipInitCheck[cmd.Name()] {
			return nil
		}

		// Try to load existing config
		err := appConfig.Initialize(configFolder())
		if err != nil {
			if os.IsNotExist(err) || !appConfig.IsInitialized() {
				return fmt.Errorf("config not initialized. Run one of:\n  config-sync init              # Start fresh\n  config-sync init-from <url>   # Clone existing repo")
			}
			return err
		}
		return nil
	},
}

func main() {
	rootCmd.AddCommand(initCmd, initFromCmd, checkUpdatesCmd, trackCmd, untrackCmd, pullCmd, pushCmd, setOriginCmd)
	rootCmd.Execute()
}
