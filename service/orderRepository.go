package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OrderRepository struct {
	C *mgo.Collection
}

func NewOrderRepository(c *mgo.Collection) *OrderRepository {
	return &OrderRepository{c}
}

func (r *OrderRepository) Create(o *Order) (err error) {
	if o.Id ==*new(bson.ObjectId) {
		o.Id = bson.NewObjectId()
	}
	err = r.C.Insert(o)
	return
}

func (r *OrderRepository) Update(o Order) (err error) {
	err = r.C.Update(bson.M{"_id": o.Id},o)
	return
}


func (r *OrderRepository) GetById(id string) (Order,error) {
	var order Order
	err := r.C.FindId(bson.ObjectIdHex(id)).One(&order)
	return order,err
}


func (r *OrderRepository) getByQuery(query *mgo.Query) []Order {
	iter := query.Iter()
	var result Order
	var orders []Order
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Id)
		orders = append(orders, result)
	}
	return orders
}



