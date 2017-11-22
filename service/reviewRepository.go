package service

import (
"fmt"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
)

type ReviewRepository struct {
	C *mgo.Collection
}

func NewReviewRepository(c *mgo.Collection) *ReviewRepository {
	return &ReviewRepository{c}
}

func (r *ReviewRepository) Create(o *Review) (err error) {
	if o.Id ==*new(bson.ObjectId) {
		o.Id = bson.NewObjectId()
	}
	err = r.C.Insert(o)
	return
}

func (r *ReviewRepository) Update(o Review) (err error) {
	err = r.C.Update(bson.M{"_id": o.Id},o)
	return
}


func (r *ReviewRepository) GetById(id string) (Review,error) {
	var review Review
	err := r.C.FindId(bson.ObjectIdHex(id)).One(&review)
	return review,err
}


func (r *ReviewRepository) GetByProductId(id string,page Page) []Review {
	query:=r.C.Find(bson.M{"product_id":bson.ObjectIdHex(id)}).Sort("-helpful_votes").Limit(page.limit()).Skip(page.skip())
	return r.getByQuery(query)
}


func (r *ReviewRepository) getByQuery(query *mgo.Query) []Review {
	iter := query.Iter()
	var result Review
	var reviews []Review
	for iter.Next(&result) {
		fmt.Printf("Result: %v\n", result.Id)
		reviews = append(reviews, result)
	}
	return reviews
}