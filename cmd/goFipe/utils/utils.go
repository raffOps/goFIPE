package utils

import (
	"fmt"
	"net/url"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/raffops/gofipe/cmd/goFipe/errs"
)

func EncodeUrl(path string, parameters map[string]interface{}) string {
	base, err := url.Parse(path)
	if err != nil {
		panic(fmt.Sprintf("Error parsing url: %v", err))
	}
	params := url.Values{}
	for key, value := range parameters {
		if value != "" {
			params.Add(key, value.(string))
		}
	}
	base.RawQuery = params.Encode()
	return base.String()
}

func StructToMap(obj interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	for i := 0; i < v.NumField(); i++ {
		out[t.Field(i).Name] = v.Field(i).Interface()
	}
	return out
}

func LoadEnvVariables() *errs.AppError {
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "dev"
	}
	configPath := fmt.Sprintf("./configs/%s.env", env)
	errEnv := godotenv.Load(configPath)
	if errEnv != nil {
		return errs.NewNotFoundError(
			fmt.Sprintf("Config file not found: %s", configPath),
		)
	}
	return nil
}
