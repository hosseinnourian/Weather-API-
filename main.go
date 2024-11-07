package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"time"
	"weather-wrapper/helper"
	"weather-wrapper/model"
	"weather-wrapper/pkg"
)

const cacheExpiry = 24 * time.Hour

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
}

func main() {
	city := flag.String("city", "", "City")
	flag.Parse()
	if *city == "" {
		log.Error(errors.New("no city provided"))
		return
	}
	redisClient := pkg.NewRedisClient()

	//read city name from cache
	res, err := redisClient.Get(context.Background(), *city).Result()
	if err == nil {

		log.Info("Cache hit!")
		fmt.Println(res)
		return
	} else if !errors.Is(err, redis.Nil) {

		log.Error("Error fetching from cache:", err)
	}

	//if not exist, make request and cache request
	log.Info("Cache miss, making HTTP request")
	requestUrl := helper.CreateRequestUrl(*city)

	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Error("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Non-OK HTTP response status:", resp.StatusCode)
		return
	}

	var apiResponse model.ApiResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Error("Error parsing JSON response:", err)
		return
	}

	if apiResponse.Location.Name == "Wrong" {
		log.Error("Invalid city name or data returned from API")
		return
	}

	if err := redisClient.Set(context.Background(), *city, string(body), cacheExpiry).Err(); err != nil {
		log.Error("Error caching response:", err)
	}

	// Print response body
	log.Info("Weather data:", string(body))
}
