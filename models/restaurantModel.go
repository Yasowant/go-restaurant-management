package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Restaurant struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Address   string             `bson:"address" json:"address" validate:"required"`
	Phone     string             `bson:"phone" json:"phone" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	OpenTime  string             `bson:"open_time" json:"open_time" validate:"required"`
	CloseTime string             `bson:"close_time" json:"close_time" validate:"required"`
}
