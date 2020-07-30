package cmd

import (
	"fmt"
	"hongdingyi/ymprotoc/internal/app"
	"github.com/spf13/cobra"
)

var (
	App   *app.App
	major = 8
	minor = 8
	patch = 8
)

var (
	RootCmd = &cobra.Command{
		Use:  "ymprotoc",
		Args: cobra.NoArgs,
	}
)

func init() {
	App, _ = app.InitApp()
}

func Run() {
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(generateCmd)
	RootCmd.Version = fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	RootCmd.Execute()
}
