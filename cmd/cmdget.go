/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	FileList         string
	IsPurgeErrorFile bool
	MinSize          int64
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get --filelist=  --workers=",
	Long:  `default value:  --filelist="filelist.txt"  --workers=3`,
	Run: func(cmd *cobra.Command, args []string) {
		kvHeaders = ParseHeader()
		DebugInfo("User-Agent", UserAgent)
		DebugInfo("Headers", kvHeaders)

		GetRun()

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVar(&FileList, "filelist", "filelist.txt", "file name")
	getCmd.PersistentFlags().BoolVar(&IsPurgeErrorFile, "purge", false, "if purge unsuccessful files automatically")
	getCmd.PersistentFlags().Int64Var(&MinSize, "minsize", 0, "if download file size <= minsize, remove it, 0 means unlimited")
}
