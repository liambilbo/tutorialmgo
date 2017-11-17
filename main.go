package main

import (
	"fmt"
	. "github.com/liambilbo/hellomgo/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var productRepository *ProductRepository
var id string

func init() {

	session, err := mgo.Dial("localhost")

	session.SetMode(mgo.Monotonic, false)

	if err != nil {
		log.Fatalf("[MongoDB session] ;: %s \n", err)
	}

	collection := session.DB("tutorial").C("products")
	collection.RemoveAll(nil)
	productRepository = NewProductRepository(collection)

}

func main() {
	createUpdateProduct()
}

func date(d string) time.Time {
	/*
		t, _:=time.Parse("2017-01-20",d)
		return t*/
	return time.Now()
}

func createUpdateProduct() Product {

	product := Product{
		Slug:        "wheelbarrow-9092",
		Name:        "Extra Large Wheelbarrow",
		Sku:         "9092",
		Description: "Heavy duty wheelbarrow...",
		Details: Details{
			Weight:       47,
			WeightUnits:  "lbs",
			ModelNum:     4039283402,
			Manufacterer: "Acme",
			Color:        "Green",
		},
		TotalReviews:  4,
		AverageReview: 4.5,
		Price: Price{
			Retail: 589700,
			Sale:   489700,
		},
		PriceHistory: []Price{
			Price{Sale: 429700, Retail: 529700, RangeDate: &RangeDate{Start: date("2010-4-1"), End: date("2010-4-8")}},
			Price{Sale: 529700, Retail: 529700, RangeDate: &RangeDate{Start: date("2010-4-9"), End: date("2010-4-6")}},
		},
		PrimaryCategory: bson.ObjectIdHex("6a5b1476238d3b4dd5000048"),
		CategoryIds:     []bson.ObjectId{bson.ObjectIdHex("6a5b1476238d3b4dd5000048"), bson.ObjectIdHex("6a5b1476238d3b4dd5000049")},
		MainCatId:       bson.ObjectIdHex("6a5b1476238d3b4dd5000048"),
		Tags:            []string{"tools", "gardening", "soil"},
	}

	if err := productRepository.Create(&product); err != nil {
		log.Fatalf("[Create Product] : %s \n", err)
	}

	id = product.Id.Hex()

	fmt.Printf("[Create Product] id %s \n", id)

	product.Tags = append(product.Tags, "ligth")

	if err := productRepository.Update(product); err != nil {
		log.Fatalf("[Update Product] %s ", err)
	}


	return product
}

func getProductById(id string) Product {
	product,err:=productRepository.GetById(id)
	return product
}
