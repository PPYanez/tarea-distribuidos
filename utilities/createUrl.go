package utilities

import (
	"log"
	"net/url"
)

func CreateUrl(section string, params map[string]string) string {
	url, err := url.Parse("http://localhost:5000/api/" + section)
	if err != nil {
		log.Fatal("URL no válida")
	}

	if params == nil {
		return url.String()
	}

	values := url.Query()
	for key, value := range params {
		values.Add(key, value)
	}

	url.RawQuery = values.Encode()

	return url.String()
}
