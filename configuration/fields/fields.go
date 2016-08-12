package fields

import "math/rand"
import "time"

type ConfigurationFieldType string

const (
	Default   ConfigurationFieldType = "default"
	FirstName ConfigurationFieldType = "firstname"
	LastName  ConfigurationFieldType = "lastname"
	FullName  ConfigurationFieldType = "fullname"
	DateTime  ConfigurationFieldType = "datetime"
	String    ConfigurationFieldType = "string"
	Number    ConfigurationFieldType = "number"
)

type ConfigurationField struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	FieldType ConfigurationFieldType `json:"-"`
}

type ConfigurationFieldData struct {
	FirstNames NamesList
	LastNames  NamesList
}

type ConfigurationFields []*ConfigurationField

func NewConfigurationFields() ConfigurationFields {
	return make(ConfigurationFields, 0)
}

type NamesList []string

func (list *NamesList) GetRandomValue() string {
	src := rand.New(rand.NewSource(time.Now().Unix()))

	n := len(*list)
	randomint := src.Int()

	return (*list)[randomint%n]
}
