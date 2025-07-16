package services

import (
	"context"
	"errors"
	"restaurant-app/config"
	"restaurant-app/models"
	"restaurant-app/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ✅ Get the collection safely after DB is connected
func getUserCollection() *mongo.Collection {
	return config.DB.Collection("users")
}

func RegisterUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := getUserCollection()

	// Check for duplicate email
	count, _ := collection.CountDocuments(ctx, bson.M{"email": user.Email})
	if count > 0 {
		return errors.New("email already exists")
	}

	user.ID = primitive.NewObjectID()
	user.Password = utils.HashPassword(user.Password)

	_, err := collection.InsertOne(ctx, user)
	return err
}

func LoginUser(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := getUserCollection()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func GetUserByID(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := getUserCollection()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return &user, nil
}

func UpdateUserProfile(userID string, updates bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := getUserCollection()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	return err
}

func ChangePassword(userID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := getUserCollection()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	newHashed := utils.HashPassword(newPassword)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"password": newHashed}})
	return err
}
