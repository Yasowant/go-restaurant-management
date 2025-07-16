package controllers

import (
	"net/http"
	"restaurant-app/models"
	"restaurant-app/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var restaurantValidator = validator.New()

func CreateRestaurant(c *gin.Context) {
	var r models.Restaurant
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := restaurantValidator.Struct(r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateRestaurant(r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
}
