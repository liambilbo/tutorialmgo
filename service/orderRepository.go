package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
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

func (r *OrderRepository) GetReportByMonthAfter(since time.Time) map[string]struct{Number int;Subtotal int}{
	result := make(map[string]struct{Number int;Subtotal int})

	pipe:=[]bson.M{
		bson.M{"$match":bson.M{"purchase_data":bson.M{"$gt":since}}},
		bson.M{"$group":bson.M{"_id":bson.M{"year":bson.M{"$year":"$purchase_data"},
		                                    "month":bson.M{"$month":"$purchase_data"}},
								"count":bson.M{"$sum":1},
								"subtotal":bson.M{"$sum":"$sub_total"}}},
		bson.M{"$sort":bson.M{"_id":-1}},

	}

	var resp []bson.M
	r.C.Pipe(pipe).All(&resp)

	for _,v:= range resp {
		id:=v["_id"].(bson.M)
		year:=strconv.Itoa(id["year"].(int))
		month:=strconv.Itoa(id["month"].(int))
		result[year+"-"+month]= struct {
			Number   int
			Subtotal int
		}{Number:v["count"].(int) , Subtotal:v["subtotal"].(int) }
	}

	return result

}



