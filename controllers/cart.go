package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/karthikavinu/ordermanagement/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		}
		validate := validator.New()
		if err := validate.Struct(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		ordert_id := order.Order_ID

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		c.IndentedJSON(500, ordert_id)

		var filledcart models.Order
		err := ProductCollection.FindOne(ctx, bson.M{"order_id": order.Order_ID}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Order ID not found")
			c.IndentedJSON(500, err.Error())

			return
		}
		c.IndentedJSON(200, filledcart)

		ctx.Done()
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {

		var allorders models.Order
		if err := c.BindJSON(&allorders); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		}

		validate := validator.New()
		if err := validate.Struct(&allorders); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		user_id := allorders.User_ID

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "user_id", Value: user_id}}}}
		pointcursor, err := ProductCollection.Aggregate(ctx, mongo.Pipeline{filter_match})
		if err != nil {
			log.Println(err)
		}
		var listing []models.Order
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.IndentedJSON(200, err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, orderhistory := range listing {

			c.IndentedJSON(200, orderhistory)
		}
		ctx.Done()
	}
}

func PlaceOrder() gin.HandlerFunc {

	return func(c *gin.Context) {
		UserQueryID := c.Query("user_id")
		if UserQueryID == "" {
			log.Println("UserID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}
		var product models.Order
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		}
		validate := validator.New()
		if err := validate.Struct(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		product.ID = primitive.NewObjectID()
		product.Order_ID = product.ID.Hex()
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		_, insertErr := ProductCollection.InsertOne(ctx, product)
		defer cancel()

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order could not be placed"})
			return
		}

		c.IndentedJSON(200, "Successully placed the order")
		c.JSON(http.StatusFound, product)

	}
}
