/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	timeBoot  time.Time
	IsDebug   bool
	Workers   int
	UserAgent string
	Header    []string
	kvHeaders map[string]string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godown [subcommand --...] --debug --user-agent=  --header=  --header= --header=",
	Short: "subcommands: get, shell",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		//fmt.Println("\n *** elapse:", time.Since(timeBoot), "***")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	timeBoot = GetTimeNow()
	Bootstrap()

	rootCmd.PersistentFlags().BoolVar(&IsDebug, "debug", false, "print debug info")
	rootCmd.PersistentFlags().StringVar(&UserAgent, "user-agent", defaultUserAgent, "http client's user agent")
	rootCmd.PersistentFlags().StringArrayVar(&Header, "header", []string{}, "key:val , i.e.: \"Content-Type: image/png\" ")
	rootCmd.PersistentFlags().IntVar(&Workers, "workers", 3, "max worker count")
}
