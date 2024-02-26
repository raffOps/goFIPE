package utils

import (
	"fmt"
	"net/url"
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
