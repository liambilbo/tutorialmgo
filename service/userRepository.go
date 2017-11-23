package service

import (
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

func (r *UserRepository) GetByLastName(lastName string,page Page) []User {
	var users []User
	r.C.Find(bson.M{"last_name":lastName}).Skip(page.skip()).Limit(page.limit()).All(&users)
	return users
}

func (r *UserRepository) GetByLastNamePattern(lastName string,page Page) []User {
	var users []User
	r.C.Find(bson.M{"last_name":bson.M{"$regex":bson.RegEx{lastName, ""}}}).Skip(page.skip()).Limit(page.limit()).All(&users)
	return users
}

func (r *UserRepository) GetByZip(ziplow int,ziphigh int,page Page) []User {
	var users []User
	r.C.Find(bson.M{"addresses.zip":bson.M{"$gt":ziplow,"$lt":ziphigh}}).Skip(page.skip()).Limit(page.limit()).All(&users)
	return users
}

func (r *UserRepository) GetByFirstAddressState(state string,page Page) []User {
	var users []User
	r.C.Find(bson.M{"addresses.0.state":state}).Skip(page.skip()).Limit(page.limit()).All(&users)
	return users
}

func (r *UserRepository) GetByAddressElem(elem bson.M,page Page) []User {
	var users []User
	r.C.Find(bson.M{"addresses":bson.M{"$elemMatch":elem}}).Skip(page.skip()).Limit(page.limit()).All(&users)
	return users
}

