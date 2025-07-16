package services

import (
	"context"
	"restaurant-app/config"
	"restaurant-app/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getRestaurantCollection() *mongo.Collection {
	return config.DB.Collection("restaurants")
}

func CreateRestaurant(restaurant models.Restaurant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	restaurant.ID = primitive.NewObjectID()

	_, err := getRestaurantCollection().InsertOne(ctx, restaurant)
	return err
}
