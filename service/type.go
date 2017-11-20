package service

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Details struct {
	Weight       int `bson:"weight" json:"weight"`
	WeightUnits  string `bson:"weight_units" json:"weight_units"`
	ModelNum     int `bson:"model_num" json:"model_num"`
	Manufacterer string `bson:"manufacturer" json:"manufacturer"`
	Color        string `bson:"color" json:"color"`
}

type Price struct {
	Retail int `bson:"retail" json:"retail"`
	Sale   int `bson:"sale" json:"sale"`
	RangeDate `bson:",inline"`
}

type RangeDate struct {
	Start time.Time `bson:"start,omitempty" json:"start"`
	End   time.Time `bson:"end,omitempty" json:"end"`
}

type Product struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Slug            string  `bson:"slug" json:"slug"`
	Sku             string  `bson:"sku" json:"sku"`
	Name            string  `bson:"name" json:"name"`
	Description     string  `bson:"description" json:"description"`
	Details         Details  `bson:"details" json:"details"`
	TotalReviews    int      `bson:"total_reviews" json:"total_reviews"`
	AverageReview   float64  `bson:"average_review" json:"average_review"`
	Price           Price    `bson:"price" json:"price"`
	HistoricPrices  []Price  `bson:"price_history" json:"price_history"`
	PrimaryCategory bson.ObjectId `bson:"primary_category" json:"primary_category"`
	CategoryIds     []bson.ObjectId `bson:"category_ids" json:"category_ids"`
	MainCatId       bson.ObjectId `bson:"main_cat_id" json:"main_cat_id"`
	Tags            []string `bson:"tags" json:"tags"`
}

type CategoryId struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`
	Slug string  `bson:"slug" json:"slug"`
	Name string `bson:"name" json:"name"`
}
type Category struct {
	CategoryId
	Description string `bson:"description" json:"description"`
	ParentId bson.ObjectId `bson:"parent_id" json:"parent_id"`
	Ancestors []CategoryId `bson:"ancestors" json:"ancestors"`
}


