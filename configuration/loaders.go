package configuration

import (
	"fmt"
	"math/rand"
	"time"
)

var loaderFactories = make(map[string]TypeLoaderFactoryFn)

type TypeLoaderFactoryFn func() *TypeLoader

type TypeLoaderLoadFn func(dto *UseTypeDTO)
type TypeLoaderGenerateSingleValueFn func(config *Configuration, set *ResultSet) (interface{}, error)
type TypeLoader struct {
	Load                TypeLoaderLoadFn
	GenerateSingleValue TypeLoaderGenerateSingleValueFn
	LoaderData          interface{} `json:"-"`
}

func addCsvLoaderFactory() {
	fn := func() *TypeLoader {
		loader := &TypeLoader{}

		loader.Load = func(dto *UseTypeDTO) {
			src := dto.LoaderArgs.Args["src"].(string)
			sep := dto.LoaderArgs.Args["separator"].(string)

			content, _ := readContentFromFileAndSplit(src, sep)
			loader.LoaderData = content
		}

		loader.GenerateSingleValue = func(config *Configuration, set *ResultSet) (interface{}, error) {
			content := loader.LoaderData.([]string)

			return GetRandomValue(&content), nil
		}

		return loader
	}

	AddLoaderFactory("csvloader", fn)
}

func addNumberFactory() {
	fn := func() *TypeLoader {
		loader := &TypeLoader{}

		loader.Load = func(dto *UseTypeDTO) {

		}

		loader.GenerateSingleValue = func(config *Configuration, set *ResultSet) (interface{}, error) {
			return rand.Int63(), nil
		}

		return loader
	}

	AddLoaderFactory("number", fn)
}

func addDateTimeFactory() {
	fn := func() *TypeLoader {
		loader := &TypeLoader{}

		loader.Load = func(dto *UseTypeDTO) {

		}

		loader.GenerateSingleValue = func(config *Configuration, set *ResultSet) (interface{}, error) {
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

	fn := func() *TypeLoader {
		loader := &TypeLoader{}

		loader.Load = func(dto *UseTypeDTO) {
			format := dto.LoaderArgs.Args["format"].(string)
			argsRaw := dto.LoaderArgs.Args["args"].([]interface{})

			args := make([]string, len(argsRaw))
			for i, a := range argsRaw {
				args[i] = fmt.Sprintf("%s", a)
			}

			loader.LoaderData = interface{}(strLoaderArgs{Format: format, Args: args})
		}

		loader.GenerateSingleValue = func(config *Configuration, set *ResultSet) (interface{}, error) {
			data := loader.LoaderData.(strLoaderArgs)

			args := make([]interface{}, len(data.Args))
			for i, k := range data.Args {
				value, err := set.Entries.GetEntryWithName(k)
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

func addDefaultLoaderFactories() {
	addCsvLoaderFactory()
	addNumberFactory()
	addStrFormatFactory()
	addDateTimeFactory()
}

func AddLoaderFactory(name string, factory TypeLoaderFactoryFn) {
	loaderFactories[name] = factory
}
