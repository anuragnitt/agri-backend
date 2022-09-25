package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ProductId	primitive.ObjectID	`bson:"_id" json:"_id"`
	Name		string				`bson:"name" json:"name"`
	Price		uint32				`bson:"price" json:"price"`
	Category	string				`bson:"category" json:"category"`
}
