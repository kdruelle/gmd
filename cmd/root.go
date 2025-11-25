/*
Copyright © 2024 Saivenkat Ajay D. <ajayds2001@gmail.com>
*/
package cmd

import (
	_ "embed"
	"os"

	"github.com/kdruelle/gmd/tui"
	"github.com/spf13/cobra"
)

var (
	debugfile string
	rootCmd   = &cobra.Command{
		Use:     "gmd",
		Short:   "TUI to manage docker objects",
		Long:    `The Definitive TUI to manage docker objects with ease.`,
		Version: "0.1.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Start(debugfile)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&debugfile, "debug", "d", "", "Create a debug file at the chosen location")
}
