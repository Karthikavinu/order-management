package main

import (
	"log"
	"os"

	"github.com/karthikavinu/ordermanagement/controllers"
	"github.com/karthikavinu/ordermanagement/middleware"
	"github.com/karthikavinu/ordermanagement/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	//app := controllers.NewApplication(database.ProductData(database.Client, "Orders"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/getorder", controllers.GetOrder())
	router.GET("/getallorders", controllers.GetAllOrders())
	router.POST("/placeorder", controllers.PlaceOrder())

	log.Fatal(router.Run(":" + port))
}
