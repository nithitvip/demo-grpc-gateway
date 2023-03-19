package service

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}
