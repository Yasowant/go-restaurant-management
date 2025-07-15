package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName string             `bson:"full_name" json:"full_name" validate:"required,min=2,max=100"`
	Email    string             `bson:"email" json:"email" validate:"required,email"`
	Phone    string             `bson:"phone" json:"phone" validate:"required"`
	Password string             `bson:"password" json:"password,omitempty"`
	Role     string             `bson:"role" json:"role" validate:"required,oneof=admin user"`
}
