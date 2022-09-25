package models

import (
	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type		string		`bson:"type" json:"type"`
	Coordinates	[]float64	`bson:"coordinates" json:"coordinates"`
}

type Shop struct {
	ShopId		primitive.ObjectID	`bson:"_id" json"_id"`
	CatalogueId	primitive.ObjectID	`bson:"catalogueId" json"catalogueId"`
	Location	Location			`bson:"location" json:"location"`
	Name		string				`bson:"name" json:"name"`
	Owner		string				`bson:"owner" json:"owner"`
}

func (l *Location) Distance(longitude, latitude float64) float64 {	
	radlat1 := float64(math.Pi * l.Coordinates[1] / 180)
	radlat2 := float64(math.Pi * latitude / 180)

	theta := float64(l.Coordinates[0] - longitude)
	radtheta := float64(math.Pi * theta / 180)
	
	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta);
	dist = math.Min(dist, 1)

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515 * 1609.34

	return dist
}
