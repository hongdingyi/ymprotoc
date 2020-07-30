package app

import (
	"fmt"
	"github.com/googleapis/api-linter/lint"
	"github.com/googleapis/api-linter/rules"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"hongdingyi/ymprotoc/internal/compile"
	"hongdingyi/ymprotoc/internal/conf"
	"hongdingyi/ymprotoc/internal/flags"
	"hongdingyi/ymprotoc/internal/format"
	"hongdingyi/ymprotoc/internal/proto"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type App struct {
	config    *conf.Config
	compiler  *compile.Compiler
	formatter *format.Formatter
}

func NewApp(config *conf.Config, compiler *compile.Compiler, formatter *format.Formatter) *App {
	return &App{
		config:    config,
		compiler:  compiler,
		formatter: formatter,
	}
}

func (a *App) Format() {
	var (
		err error

		absFiles []string
	)
	err = a.config.Load()
	if err != nil {
		log.Fatal(err)
	}

	absFiles = a.specialFile()
	if len(absFiles) > 0 {
		a.formatter.Format(absFiles)
	}
}

func (a *App) Gen() {
	var (
		err error
	)
	err = a.config.Load()
	if err != nil {
		log.Fatal(err)
	}
	protos := a.specialFile()
	if len(protos) == 0 {
		protos = a.config.Protos
	}
	descSource, err := proto.DescriptorSourceFromProtoFiles(a.config.Includes, protos...)
	if err != nil {
		log.Fatalf("Failed to process proto source files. %v", err)
	}

	err = a.compiler.Compile(descSource)
	if err != nil {
		log.Fatalf("compile error %v", err)
	}
	log.Println("build  proto success")
	return
}

func (a *App) Lint() (err error) {
	var (
		absFiles []string
	)
	err = a.config.Load()
	if err != nil {
		log.Fatal(err)
	}

	absFiles = a.specialFile()

	if len(absFiles) == 0 {
		return
	}

	for _, itr := range absFiles {
		log.Println("lint file:", itr)
	}

	lintCfgs := lint.Configs{}
	// Add configs for the enabled rules.
	lintCfgs = append(lintCfgs, lint.Config{
		EnabledRules: a.config.Lint.Rules.Enable,
	})
	// Add configs for the disabled rules.
	lintCfgs = append(lintCfgs, lint.Config{
		DisabledRules: a.config.Lint.Rules.Disable,
	})
	// Prepare proto import lookup.
	fs, err := loadFileDescriptors()
	if err != nil {
		return err
	}
	lookupImport := func(name string) (*desc.FileDescriptor, error) {
		if f, found := fs[name]; found {
			return f, nil
		}
		return nil, fmt.Errorf("%q is not found", name)
	}

	var errorsWithPos []protoparse.ErrorWithPos
	var lock sync.Mutex
	// Parse proto files into `protoreflect` file descriptors.
	p := protoparse.Parser{
		ImportPaths:           a.config.Includes,
		IncludeSourceCodeInfo: true,
		LookupImport:          lookupImport,
		ErrorReporter: func(errorWithPos protoparse.ErrorWithPos) error {
			// Protoparse isn't concurrent right now but just to be safe for the future.
			lock.Lock()
			errorsWithPos = append(errorsWithPos, errorWithPos)
			lock.Unlock()
			// Continue parsing. The error returned will be protoparse.ErrInvalidSource.
			return nil
		},
	}
	// Resolve file absolute paths to relative ones.
	protoFiles, err := protoparse.ResolveFilenames(a.config.Includes, absFiles...)
	if err != nil {
		return err
	}
	fd, err := p.ParseFiles(protoFiles...)
	if err != nil {
		if err == protoparse.ErrInvalidSource {
			if len(errorsWithPos) == 0 {
				return errors.New("got protoparse.ErrInvalidSource but no ErrorWithPos errors")
			}
			// TODO: There's multiple ways to deal with this but this prints all the errors at least
			errStrings := make([]string, len(errorsWithPos))
			for i, errorWithPos := range errorsWithPos {
				errStrings[i] = errorWithPos.Error()
			}
			return errors.New(strings.Join(errStrings, "\n"))
		}
		return err
	}

	// Create a linter to lint the file descriptors.
	globalRule := lint.NewRuleRegistry()
	rules.Add(globalRule)
	l := lint.New(globalRule, lintCfgs)
	results, err := l.LintProtos(fd...)
	if err != nil {
		return err
	}

	// Determine the output for writing the results.
	// Stdout is the default output.
	w := os.Stdout

	// Determine the format for printing the results.
	// YAML format is the default.
	marshal := getOutputFormatFunc("yaml")

	// Print the results.
	b, err := marshal(results)
	if err != nil {
		return err
	}
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}

func (a *App) Config() {
	err := a.config.Output()
	log.Fatal(err)
	return
}

func (a *App) specialFile() []string {
	var (
		absFiles []string
	)
	//文件参数
	sourceFiles := map[string]struct{}{}
	for _, itr := range flags.SrcFiles {
		sourceFiles[itr] = struct{}{}
	}
	for _, itr := range a.config.Protos {
		if _, ok := sourceFiles[itr]; ok {
			absPath := filepath.Join(a.config.ImportPath, itr)
			absPath = filepath.ToSlash(absPath)
			absFiles = append(absFiles, absPath)
		}
	}
	//目录参数
	for _, itr := range a.config.Protos {
		for _, dir := range flags.SrcDirectories {
			if strings.Contains(itr, dir) {
				absPath := filepath.Join(a.config.ImportPath, itr)
				absPath = filepath.ToSlash(absPath)
				absFiles = append(absFiles, absPath)
			}
		}
	}

	return absFiles
}
