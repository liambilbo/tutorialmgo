package service

import (
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
	var review []Review
	r.C.Find(bson.M{"product_id":bson.ObjectIdHex(id)}).Sort("-helpful_votes").Limit(page.limit()).Skip(page.skip()).All(&review)
	return review
}


func (r *ReviewRepository) GetByWhere(where string,page Page) []Review {
	var review []Review
	r.C.Find(bson.M{"$where":where}).Sort("-helpful_votes").Limit(page.limit()).Skip(page.skip()).All(&review)
	return review
}


func (r *ReviewRepository) GetByText(regexp bson.RegEx,page Page) []Review {
	var review []Review
	r.C.Find(bson.M{"text":regexp}).Sort("-helpful_votes").Limit(page.limit()).Skip(page.skip()).All(&review)
	return review
}


func  (r *ReviewRepository) CountByProductId(productId string) (average float64,count int) {
	//resp := []bson.M{}
	resp := bson.M{}
	pipe :=r.C.Pipe([]bson.M{
		bson.M{"$match":bson.M{"product_id":bson.ObjectIdHex(productId)}},
		bson.M{"$group":bson.M{"_id":"$product_id",
		                       "average":bson.M{"$avg":"$rating"},
		                       "count":bson.M{"$sum":1}}},
	})

	//pipe.All(&resp)

	pipe.One(&resp)

	if c,okc:=resp["count"];okc {count= c.(int)}
	if a,okr:=resp["average"];okr {average= a.(float64)}

	return
}