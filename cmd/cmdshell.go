/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	CmdList string
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "--cmdlist=",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ShellRun()
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)

	shellCmd.PersistentFlags().StringVar(&CmdList, "cmdlist", "cmdlist.txt", "shell cmd list")
}
