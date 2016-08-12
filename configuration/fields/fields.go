package fields

type ConfigurationFieldType string
const(
    Default ConfigurationFieldType = "default"
    FirstName ConfigurationFieldType = "firstname"
    LastName ConfigurationFieldType = "lastname"
    DateTime ConfigurationFieldType = "datetime"
    String ConfigurationFieldType = "string"
    Number ConfigurationFieldType = "number"
)

type ConfigurationField struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    FieldType   ConfigurationFieldType `json:"-"`
}