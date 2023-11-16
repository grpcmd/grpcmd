package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Outputs the version",
	Run: func(cmd *cobra.Command, args []string) {
		root := cmd.Root()
		root.SetArgs([]string{"--version"})
		root.Execute()
	},
}