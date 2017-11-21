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
	if o.Id ==*new(bson.ObjectId) {
		o.Id = bson.NewObjectId()
	}
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


func (r *CategoryRepository) getByQuery(query *mgo.Query) []Category {
	iter := query.Iter()
	var result Category
	var categories []Category
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Id)
		categories = append(categories, result)
	}
	return categories
}
func (r *CategoryRepository) GetAll() []Category {
	query := r.C.Find(bson.M{})
	return r.getByQuery(query)

}

func (r *CategoryRepository) GetCategoByIds(ids []string) []Category {
	var objectIds []bson.ObjectId
	for _,v:=range ids {
		objectIds=append(objectIds,bson.ObjectIdHex(v))
	}
	query:=r.C.Find(bson.M{"_id":bson.M{"$in":objectIds}})
	return r.getByQuery(query)
}
