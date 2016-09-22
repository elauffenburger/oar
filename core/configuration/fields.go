package configuration

type ConfigurationField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ConfigurationFields []*ConfigurationField

func NewConfigurationFields() ConfigurationFields {
	return make(ConfigurationFields, 0)
}
