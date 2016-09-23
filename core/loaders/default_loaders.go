package loaders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/elauffenburger/oar/core/common"
	conf "github.com/elauffenburger/oar/core/configuration"
	res "github.com/elauffenburger/oar/core/results"
	"github.com/satori/go.uuid"
)

type FnTypeLoader struct {
	typeLoader

	loadFn                TypeLoaderLoadFn
	generateSingleValueFn TypeLoaderGenerateSingleValueFn
}

func (loader *FnTypeLoader) Load(dto *conf.UseTypeDTO) {
	loader.loadFn(dto)
}

func (loader *FnTypeLoader) GenerateSingleValue(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
	return loader.generateSingleValueFn(config, set)
}

func addCsvLoaderFactory(ctx *TypeLoaderFactoryContext) {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *conf.UseTypeDTO) {
			src := dto.LoaderArgs.Args["src"].(string)
			sep := dto.LoaderArgs.Args["separator"].(string)

			content, _ := common.ReadContentFromFileAndSplit(src, sep)
			loader.LoaderData = content
		}

		loader.generateSingleValueFn = func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
			content := loader.LoaderData.([]string)

			return common.GetRandomValue(&content), nil
		}

		return loader
	}

	ctx.AddLoaderFactory("csvloader", fn)
}

func addNumberFactory(ctx *TypeLoaderFactoryContext) {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *conf.UseTypeDTO) {

		}

		loader.generateSingleValueFn = func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
			return rand.Int63(), nil
		}

		return loader
	}

	ctx.AddLoaderFactory("number", fn)
}

func addDateTimeFactory(ctx *TypeLoaderFactoryContext) {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *conf.UseTypeDTO) {

		}

		loader.generateSingleValueFn = func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
			return time.Unix(rand.Int63(), rand.Int63()), nil
		}

		return loader
	}

	ctx.AddLoaderFactory("datetime", fn)
}

func addStrFormatFactory(ctx *TypeLoaderFactoryContext) {
	type strLoaderArgs struct {
		Format string
		Args   []string
	}

	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *conf.UseTypeDTO) {
			format := dto.LoaderArgs.Args["format"].(string)
			argsRaw := dto.LoaderArgs.Args["args"].([]interface{})

			args := make([]string, len(argsRaw))
			for i, a := range argsRaw {
				args[i] = fmt.Sprintf("%s", a)
			}

			loader.LoaderData = interface{}(strLoaderArgs{Format: format, Args: args})
		}

		loader.generateSingleValueFn = func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
			data := loader.LoaderData.(strLoaderArgs)

			args := make([]interface{}, len(data.Args))
			for i, k := range data.Args {
				value, err := set.Values.GetEntryWithName(k)
				if err != nil {
					// todo do something about this
					continue
				}

				args[i] = value.Value
			}

			return fmt.Sprintf(data.Format, args...), nil
		}

		return loader
	}

	ctx.AddLoaderFactory("strformat", fn)
}

func addAutoIncrementFactory(ctx *TypeLoaderFactoryContext) {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *conf.UseTypeDTO) {
			loader.LoaderData = 0
		}

		loader.generateSingleValueFn = func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
			val := loader.LoaderData.(int)

			val++
			loader.LoaderData = val

			return val, nil
		}

		return loader
	}

	ctx.AddLoaderFactory("autoincrement", fn)
}

func addUUIDFactory(ctx *TypeLoaderFactoryContext) {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *conf.UseTypeDTO) {

		}

		loader.generateSingleValueFn = func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error) {
			return uuid.NewV4().String(), nil
		}

		return loader
	}

	ctx.AddLoaderFactory("uuid", fn)
}

func AddDefaultLoaderFactories(ctx *TypeLoaderFactoryContext) {
	addCsvLoaderFactory(ctx)
	addNumberFactory(ctx)
	addStrFormatFactory(ctx)
	addDateTimeFactory(ctx)
	addAutoIncrementFactory(ctx)
	addUUIDFactory(ctx)
}
