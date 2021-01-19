package main

import (
	"fmt"
	"github.com/MaxFuhrich/containerNetworkExample/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	//Endpoints
	router.GET("/", func(context *gin.Context) {
		message := "Hello there!"
		fmt.Println(message)
		context.JSON(http.StatusOK, message)
	})
	router.GET("/hello", func(context *gin.Context) {
		message := "Hello there!"
		fmt.Println(message)
		context.JSON(http.StatusOK, message)
	})
	router.GET("/add", func(context *gin.Context) {
		fmt.Println("Endpoint /add called!")
		controller.AddTime(context)
	})
	router.GET("/history", controller.History)
	err := router.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
