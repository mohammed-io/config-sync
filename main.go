package main

import (
	"fmt"
	"log"
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
	Short: "Git pull in synced-files folder",
	Run: func(cmd *cobra.Command, args []string) {
		git := NewGitRunner()
		if err := git.Pull(); err != nil {
			log.Fatalf("Pull failed: %v", err)
		}
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

func init() {
	setOriginCmd.Flags().Bool("force", false, "Bypass public repository warning")
}

var rootCmd = &cobra.Command{
	Use:   "config-sync",
	Short: "Sync config files across machines",
	Long: `config-sync helps you track and sync configuration files across machines.
Files are stored in ~/.config-sync/synced-files and can be managed with git.

Version: ` + Version,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for help and completion commands
		if cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "version" {
			return nil
		}
		return appConfig.Initialize(configFolder())
	},
}

func main() {
	rootCmd.AddCommand(trackCmd, untrackCmd, pullCmd, pushCmd, setOriginCmd)
	rootCmd.Execute()
}
