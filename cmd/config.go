package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "init",
		Short: "生成yaml文件",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			App.Config()
		},
	}
)
