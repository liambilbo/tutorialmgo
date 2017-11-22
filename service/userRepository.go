package service

import (
"fmt"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	C *mgo.Collection
}

func NewUserRepository(c *mgo.Collection) *UserRepository {
	return &UserRepository{c}
}

func (r *UserRepository) Create(o *User) (err error) {
	if o.Id ==*new(bson.ObjectId) {
		o.Id = bson.NewObjectId()
	}
	err = r.C.Insert(o)
	return
}

func (r *UserRepository) Update(o User) (err error) {
	err = r.C.Update(bson.M{"_id": o.Id},o)
	return
}


func (r *UserRepository) GetById(id string) (User,error) {
	var user User
	err := r.C.FindId(bson.ObjectIdHex(id)).One(&user)
	return user,err
}


func (r *UserRepository) GetByNameAndPassword(username string,password string ) (User,error) {
	var user User
	err := r.C.Find(bson.M{"username":username,"hashed_password":password}).Select(bson.M{"_id":1}).One(&user)
	return user,err
}


func (r *UserRepository) getByQuery(query *mgo.Query) []Order {
	iter := query.Iter()
	var result Order
	var orders []Order
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Id)
		orders = append(orders, result)
	}
	return orders
}