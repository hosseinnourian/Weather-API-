package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"io/ioutil"
	"net/http"
	"time"
	"weather-wrapper/helper"
	"weather-wrapper/pkg"
)

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
	if errors.Is(err, redis.Nil) {
		log.Error(errors.New("no cache found"))
		log.Info("make http request")
	} else {
		//if existed, return to user
		log.Infof("cache hit!")
		fmt.Println(res)
		return
	}

	//if not exist, make request and cache request

	requestUrl := helper.CreateRequestUrl(*city)

	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		err := redisClient.Set(context.Background(), *city, string(body), 86400*time.Second).Err()
		if err != nil {
			log.Error(err)
		}
	}

	log.Infof(string(body))
}
