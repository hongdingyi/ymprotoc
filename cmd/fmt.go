package cmd

import (
	"errors"
	"hongdingyi/ymprotoc/internal/flags"
	"github.com/spf13/cobra"
)

var (
	fmtCmd = &cobra.Command{
		Use:   "fmt",
		Short: "格式化proto",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.HasAvailableFlags() {
				return errors.New("need flag")
			}
			App.Format()
			return nil
		},
	}
)

func init() {
	fmtCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "files", "f",
		[]string{}, "相对于导入目录下的源文件")
	fmtCmd.PersistentFlags().StringSliceVarP(&flags.SrcDirectories, "directories", "d",
		[]string{}, "相对于导入目录下源文件目录")
}
