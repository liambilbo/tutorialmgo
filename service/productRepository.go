package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProductRepository struct {
	C *mgo.Collection
}

func NewProductRepository(c *mgo.Collection) *ProductRepository {
	return &ProductRepository{c}
}

func (r *ProductRepository) Create(o *Product) (err error) {
	if o.Id ==*new(bson.ObjectId) {
		o.Id = bson.NewObjectId()
	}

	err = r.C.Insert(o)
	return
}

func (r *ProductRepository) Update(o Product) (err error) {
	err = r.C.Update(bson.M{"_id": o.Id},o)
	return

}


func (r *ProductRepository) GetById(id string) (Product,error) {
	var product Product
	err := r.C.FindId(bson.ObjectIdHex(id)).One(&product)
	return product,err
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

func (r *ProductRepository) GetByCategoryId(categoryId string) []Product {
	query := r.C.Find(bson.M{"category_ids":bson.ObjectIdHex(categoryId)})
	iter := query.Iter()
	var result Product
	var products []Product

	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Id)
		products=append(products,result)
	}
	return products
}

func (r *ProductRepository) GetBySlug(slug string) (Product,error) {
	var product Product
	err := r.C.Find(bson.M{"slug":slug}).One(&product)
	return product,err
}

func (r *ProductRepository) GetByColorAndManufacturer(color string,manufacturer string,page Page) []Product {
	var products []Product
	r.C.Find(bson.M{"$or":[]bson.M{bson.M{"details.color":color},bson.M{"details.manufacturer":manufacturer}}}).Skip(page.skip()).Limit(page.limit()).All(&products)
	return products
}


func (r *ProductRepository) GetByColorExists(exist bool,page Page) []Product {
	var products []Product
	r.C.Find(bson.M{"details.color":bson.M{"$exists":exist}}).Skip(page.skip()).Limit(page.limit()).All(&products)
	return products
}

func (r *ProductRepository) GetByFirstTag(tag string,page Page) []Product {
	var products []Product
	r.C.Find(bson.M{"tags.0":tag}).Skip(page.skip()).Limit(page.limit()).All(&products)
	return products
}


func  (r *ProductRepository) LoadProductsByCategory() {
	resp := []bson.M{}


	pipe :=r.C.Pipe([]bson.M{
		bson.M{"$project":bson.M{"category_ids":1}},
		bson.M{"$unwind":"$category_ids"},
		bson.M{"$group":bson.M{"_id":"$category_ids",
			"count":bson.M{"$sum":1}}},
		bson.M{"$out":"countsByCategory"},
	})

	pipe.All(&resp)

	return
}

