package cmd

import (
	"github.com/gandarfh/httui/internal/tui"
	"github.com/spf13/cobra"
)

func loginCmd() *cobra.Command {
	init := &cobra.Command{
		Use:     "login",
		Short:   "Login in httui platform.",
		Long:    "Login to sync your httui account in with your terminal.",
		Example: "httui login",
		Aliases: []string{"login"},
		RunE: func(cmd *cobra.Command, args []string) error {
			tui.Login()
			return nil
		},
	}
	return init
}
