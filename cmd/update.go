package cmd

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update gmd",
	Run: func(cmd *cobra.Command, args []string) {

		err := update("v0.0.1")
		if err != nil {
			log.Fatal(err)
		}

	},
}

func update(version string) error {

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("error occurred while getting path to executable: %w", err)
	}

	updaterConfig := selfupdate.Config{
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"},
	}

	updater, err := selfupdate.NewUpdater(updaterConfig)
	if err != nil {
		return fmt.Errorf("error occurred while creating updater: %w", err)
	}

	fmt.Println("→ Checking for latest version...")
	latest, found, err := updater.DetectLatest(context.Background(), selfupdate.ParseSlug("kdruelle/gmd"))
	if err != nil {
		return fmt.Errorf("an error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found", runtime.GOOS, runtime.GOARCH)
	}

	if latest.LessOrEqual(version) {
		fmt.Println("✔ Already up to date:", latest.Version())
		return nil
	}

	fmt.Println("→ New version available:", latest.Version())

	if err := updater.UpdateTo(context.Background(), latest, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}
	fmt.Printf("✔ Successfully updated to version %s\n", latest.Version())
	return nil
}
