package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CategoryRepository struct {
	C *mgo.Collection
}

func NewCategoryRepository(c *mgo.Collection) *CategoryRepository {
	return &CategoryRepository{c}
}

func (r *CategoryRepository) Create(o *Category) (err error) {
	o.Id = bson.NewObjectId()
	err = r.C.Insert(o)
	return
}

func (r *CategoryRepository) Update(o Category) (err error) {
	err = r.C.Update(bson.M{"_id": o.Id},o)
	return
}


func (r *CategoryRepository) GetById(id string) (Product,error) {
	var product Product
	err := r.C.FindId(bson.ObjectIdHex(id)).One(&product)
	return product,err
}

func (r *CategoryRepository) GetAll() []Product {
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
