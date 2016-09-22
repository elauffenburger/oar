package configuration

type UseTypeDTO struct {
	LoaderArgs UseTypeLoaderArgsDTO `json:"loader"`
}

type UseTypeLoaderArgsDTO struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type Configuration struct {
	Fields  ConfigurationFields   `json:"fields"`
	NumRows int                   `json:"numRows"`
	Options map[string]string     `json:"options"`
	Types   map[string]UseTypeDTO `json:"types"`
}

type OutputType string

const (
	JSON OutputType = "json"
	SQL  OutputType = "sql"
)

func NewConfiguration() *Configuration {
	return &Configuration{Options: make(map[string]string), Fields: NewConfigurationFields()}
}
