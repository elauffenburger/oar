package loaders

import (
	conf "github.com/elauffenburger/oar/core/configuration"
	res "github.com/elauffenburger/oar/core/results"
)

type TypeLoaderFactoryFn func() TypeLoader

type TypeLoaderLoadFn func(dto *conf.UseTypeDTO)
type TypeLoaderGenerateSingleValueFn func(config *conf.Configuration, set *res.ResultsRow) (interface{}, error)
type TypeLoader interface {
	Load(dto *conf.UseTypeDTO)
	GenerateSingleValue(config *conf.Configuration, set *res.ResultsRow) (interface{}, error)
}

type typeLoader struct {
	LoaderData interface{} `json:"-"`
}

type TypeLoaderFactoryContext map[string]TypeLoaderFactoryFn

func (ctx *TypeLoaderFactoryContext) AddLoaderFactory(name string, factory TypeLoaderFactoryFn) {
	(*ctx)[name] = factory
}
