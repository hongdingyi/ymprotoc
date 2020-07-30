package cmd

import (
	"errors"
	"hongdingyi/ymprotoc/internal/flags"
	"github.com/spf13/cobra"
	"log"
)

var (
	lintCmd = &cobra.Command{
		Use:   "check",
		Short: "对proto进行拼写检查",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.HasAvailableFlags() {
				return errors.New("need flag")
			}
			err := App.Lint()
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
)

func init() {
	lintCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "files", "f",
		[]string{}, "相对于导入目录下的源文件")
	lintCmd.PersistentFlags().StringSliceVarP(&flags.SrcDirectories, "directories", "d",
		[]string{}, "相对于导入目录下源文件目录")
}
