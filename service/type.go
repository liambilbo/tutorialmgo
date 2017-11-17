package service

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Details struct {
	Weight       int
	WeightUnits  string
	ModelNum     int
	Manufacterer string
	Color        string
}

type Price struct {
	Retail int
	Sale   int
	*RangeDate
}

type RangeDate struct {
	Start time.Time
	End   time.Time
}

type Product struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Slug            string
	Sku             string
	Name            string
	Description     string
	Details         Details
	TotalReviews    int
	AverageReview   float64
	Price           Price
	PriceHistory    []Price
	PrimaryCategory bson.ObjectId
	CategoryIds     []bson.ObjectId
	MainCatId       bson.ObjectId
	Tags            []string
}
