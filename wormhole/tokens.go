package wormhole

import (
	"encoding/json"
	"reflect"
	"time"

	paseto "github.com/o1egl/paseto/v2"
)

// NewToken returns a paseto.JSONToken, combined with the payload struct it receives
func NewToken(payload interface{}) (jsonToken paseto.JSONToken, err error) {
	jsonToken = paseto.JSONToken{
		IssuedAt:   time.Now(),
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(5 * time.Second),
	}
	v := reflect.ValueOf(payload)
	for i := 0; i < v.NumField(); i++ {
		jsonToken.Set(v.Type().Field(i).Name, v.Field(i).Interface())
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return jsonToken, err
	}
	jsonToken.Set("body", body)
	return jsonToken, nil
}

// SchemaFromJSONToken fills in a schema struct with values from
// additional claims in the PASETO token
// if not all values of the struct are present and strict is true,
// an error will be returned
func SchemaFromJSONToken(jsonToken paseto.JSONToken,
	schema interface{}) error {
	var body []byte
	err := jsonToken.Get("body", &body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &schema)
	return nil
}
