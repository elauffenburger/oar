package core

var loaderFactories = make(map[string]TypeLoaderFactoryFn)

type TypeLoaderFactoryFn func() TypeLoader

type TypeLoaderLoadFn func(dto *UseTypeDTO)
type TypeLoaderGenerateSingleValueFn func(config *Configuration, set *ResultsRow) (interface{}, error)
type TypeLoader interface {
	Load(dto *UseTypeDTO)
	GenerateSingleValue(config *Configuration, set *ResultsRow) (interface{}, error)
}

type typeLoader struct {
	LoaderData interface{} `json:"-"`
}

func AddLoaderFactory(name string, factory TypeLoaderFactoryFn) {
	loaderFactories[name] = factory
}

func generateTypeLoaders(config *Configuration) map[string]TypeLoader {
	types := make(map[string]TypeLoader)

	// load custom types
	for typename, t := range config.Types {
		if _, exists := types[typename]; exists {
			continue
		}

		loadername := t.LoaderArgs.Name
		factory, ok := loaderFactories[loadername]

		if !ok {
			// todo handle unknown loader
			continue
		}

		// create a loader for this type
		loader := factory()
		loader.Load(&t)

		types[typename] = loader
	}

	return types
}
