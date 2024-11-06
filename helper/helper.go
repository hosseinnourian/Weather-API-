package helper

import (
	"fmt"
	"os"
)

var (
	BaseUrl = "http://api.weatherstack.com/current"
)

func CreateRequestUrl(cityName string) string {

	requestUrl := fmt.Sprintf("%s?access_key=%s&query=%s", BaseUrl, os.Getenv("KEY"), cityName)
	return requestUrl
}
