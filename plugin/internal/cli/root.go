/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"kubectl-login/internal/config"

	"github.com/spf13/cobra"
)

func ExecuteCommand() error {
	root := NewRootCmd()
	return root.Execute()
}

func NewRootCmd() *cobra.Command {
	cfg := config.ConfigInit()

	var rootCmd = &cobra.Command{
		Use:   "kubectl-login",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return LoginExecute(cfg)
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	rootCmd.PersistentFlags().StringVarP(&cfg.Username, "username", "u", "", "Username")
	rootCmd.PersistentFlags().StringVarP(&cfg.Password, "password", "p", "", "Password")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Interactive, "interactive", "i", false, "Intercactive input")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Verbose output")

	return rootCmd
}
