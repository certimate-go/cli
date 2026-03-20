package models

// Access represents provider access credentials
type Access struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Provider string       `json:"provider"`
	Config   AccessConfig `json:"config"`
	Created  CustomTime   `json:"created"`
	Updated  CustomTime   `json:"updated"`
}

// AccessConfig is a generic access configuration
type AccessConfig map[string]interface{}

// SensitiveFields returns a list of field names that should be masked
var SensitiveFields = map[string]bool{
	"secretAccessKey":    true,
	"secretKey":          true,
	"accessKeySecret":    true,
	"apiSecret":          true,
	"apiKey":             true,
	"password":           true,
	"token":              true,
	"privateKey":         true,
	"secret":             true,
	"authorizationToken": true,
	"secretId":           true,
}

// MaskSensitive creates a copy with sensitive fields masked
func (a *Access) MaskSensitive() *Access {
	masked := &Access{
		ID:       a.ID,
		Name:     a.Name,
		Provider: a.Provider,
		Config:   make(AccessConfig),
		Created:  a.Created,
		Updated:  a.Updated,
	}

	for k, v := range a.Config {
		if SensitiveFields[k] {
			masked.Config[k] = "********"
		} else {
			masked.Config[k] = v
		}
	}

	return masked
}
