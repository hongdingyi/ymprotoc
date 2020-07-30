//+build wireinject

package app

import (
	"github.com/google/wire"
	"hongdingyi/ymprotoc/internal/compile"
	"hongdingyi/ymprotoc/internal/conf"
	"hongdingyi/ymprotoc/internal/format"
)

func InitApp() (*App, error) {
	panic(wire.Build(conf.NewConfig, compile.NewCompiler, format.NewFormatter, NewApp))
}
