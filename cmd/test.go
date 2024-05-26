package cmd

import (
	"github.com/spf13/cobra"
)

func testCmd() *cobra.Command {
	init := &cobra.Command{
		Use:     "test",
		Short:   "test",
		Long:    "test",
		Example: "httui test",
		Aliases: []string{"test"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return init
}
