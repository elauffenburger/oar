package core

import (
	"fmt"
	"math/rand"
	"time"
)

type FnTypeLoader struct {
	typeLoader

	loadFn                TypeLoaderLoadFn
	generateSingleValueFn TypeLoaderGenerateSingleValueFn
}

func (loader *FnTypeLoader) Load(dto *UseTypeDTO) {
	loader.loadFn(dto)
}

func (loader *FnTypeLoader) GenerateSingleValue(config *Configuration, set *ResultsRow) (interface{}, error) {
	return loader.generateSingleValueFn(config, set)
}

func addCsvLoaderFactory() {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *UseTypeDTO) {
			src := dto.LoaderArgs.Args["src"].(string)
			sep := dto.LoaderArgs.Args["separator"].(string)

			content, _ := readContentFromFileAndSplit(src, sep)
			loader.LoaderData = content
		}

		loader.generateSingleValueFn = func(config *Configuration, set *ResultsRow) (interface{}, error) {
			content := loader.LoaderData.([]string)

			return GetRandomValue(&content), nil
		}

		return loader
	}

	AddLoaderFactory("csvloader", fn)
}

func addNumberFactory() {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *UseTypeDTO) {

		}

		loader.generateSingleValueFn = func(config *Configuration, set *ResultsRow) (interface{}, error) {
			return rand.Int63(), nil
		}

		return loader
	}

	AddLoaderFactory("number", fn)
}

func addDateTimeFactory() {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *UseTypeDTO) {

		}

		loader.generateSingleValueFn = func(config *Configuration, set *ResultsRow) (interface{}, error) {
			return time.Unix(rand.Int63(), rand.Int63()), nil
		}

		return loader
	}

	AddLoaderFactory("datetime", fn)
}

func addStrFormatFactory() {
	type strLoaderArgs struct {
		Format string
		Args   []string
	}

	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *UseTypeDTO) {
			format := dto.LoaderArgs.Args["format"].(string)
			argsRaw := dto.LoaderArgs.Args["args"].([]interface{})

			args := make([]string, len(argsRaw))
			for i, a := range argsRaw {
				args[i] = fmt.Sprintf("%s", a)
			}

			loader.LoaderData = interface{}(strLoaderArgs{Format: format, Args: args})
		}

		loader.generateSingleValueFn = func(config *Configuration, set *ResultsRow) (interface{}, error) {
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

	AddLoaderFactory("strformat", fn)
}

func addAutoIncrementFactory() {
	fn := func() TypeLoader {
		loader := &FnTypeLoader{}

		loader.loadFn = func(dto *UseTypeDTO) {
			loader.LoaderData = 0
		}

		loader.generateSingleValueFn = func(config *Configuration, set *ResultsRow) (interface{}, error) {
			val := loader.LoaderData.(int)

			val++
			loader.LoaderData = val

			return val, nil
		}

		return loader
	}

	AddLoaderFactory("autoincrement", fn)
}

func addDefaultLoaderFactories() {
	addCsvLoaderFactory()
	addNumberFactory()
	addStrFormatFactory()
	addDateTimeFactory()
	addAutoIncrementFactory()
}
