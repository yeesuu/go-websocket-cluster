package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go-websocket-cluster/service"
	"log"
	"net/http"
)
func main()  {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	port := viper.GetString("server.port")
	redisHost := viper.GetString("redis.host")
	redisPort := viper.GetString("redis.port")
	redisPassword := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.db")
	redisMessageChannel := viper.GetString("message.channel")
	redisAddr := redisHost+":"+redisPort
	redisService := service.NewRedisService(redisAddr, redisPassword, redisDB)
	hubService := service.NewHubService(redisService)
	go hubService.SubscribeMessage(redisMessageChannel)
	go hubService.Run()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		service.ServeWs(hubService, writer, request, redisMessageChannel)
	})
	fmt.Println("server run on port:"+port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}