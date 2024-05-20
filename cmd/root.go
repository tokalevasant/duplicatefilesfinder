/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var paths []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dupfinder",
	Short: "Find duplicate files in given paths",
	Long: `	
		Find duplicate files in given paths. Files that have the same name or same hash will be considered as duplicates. 
		The result is presented as a table with following details: 
		1. File name
		2. Size
		3. Full path
		4. Last modified date
		`,
	Run: func(cmd *cobra.Command, args []string) {
		
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dupfinder.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringArrayVarP(&paths,"paths", "p", paths, "list of paths to search")
	
	if err := rootCmd.MarkFlagRequired("paths"); err != nil  {
		fmt.Println(err)
	}
}
