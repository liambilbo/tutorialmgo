package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type ProductRepository struct {
	C *mgo.Collection
}

func NewProductRepository(c *mgo.Collection) *ProductRepository {
	return &ProductRepository{c}
}

func (r *ProductRepository) Create(o *Product) (err error) {
	o.Id = bson.NewObjectId()
	err = r.C.Insert(o)
	return
}

func (r *ProductRepository) Update(o Product) (err error) {
	err = r.C.Update(bson.M{"_id": o.Id},
		bson.M{"slug": o.Slug,
			"sku":         o.Sku,
			"name":        o.Name,
			"description": o.Description,
			"details": bson.M{
				"weight":       o.Details.Weight,
				"weight_units": o.Details.WeightUnits,
				"model_num":    o.Details.ModelNum,
				"maufacterer":  o.Details.Manufacterer,
				"color":        o.Details.Color,
			}})
	return
}

func (r *ProductRepository) GetById(id string) Product {
	var product Product
	if err := r.C.FindId(bson.ObjectIdHex(id)).One(&product); err != nil {
		log.Fatalf("[Get Product By ID] %s ", err)
	}
	return product
}

func (r *ProductRepository) GetAll() []Product {
	query := r.C.Find(bson.M{})
	iter := query.Iter()
	var result Product
	var products []Product
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Id)
		products = append(products, result)
	}
	return products
}
