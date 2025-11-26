package cmd

import (
	"fmt"
	"runtime"

	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update gmd",
	Run: func(cmd *cobra.Command, args []string) {
		u := &updater.Updater{
			Provider: &provider.Github{
				RepositoryURL: "github.com/kdruelle/gmd",
				ArchiveName:   fmt.Sprintf("gmd_%s_%s_%s.tar.gz", version, runtime.GOOS, runtime.GOARCH),
			},
			ExecutableName: "gmd",
			Version:        version,
		}

		if _, err := u.Update(); err != nil {
			fmt.Println(err)
		}

	},
}
