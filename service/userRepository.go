package service

import (
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
	"log"
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


func (r *UserRepository) GetByAddressSize(size int,page Page) []User {
	var users []User
	r.C.Find(bson.M{"addresses":bson.M{"$size":size}}).Skip(page.skip()).Limit(page.limit()).All(&users)
	return users
}


func (r *UserRepository) UpdateAddress(userId string,address Address) (User,error) {
	var user User

	_,err:=r.C.Upsert(bson.M{"_id":bson.ObjectIdHex(userId),"addresses.name":address.Name},bson.M{"$set":bson.M{"addresses.$":address}})
	if err!=nil{
		_,err=r.C.Upsert(bson.M{"_id":bson.ObjectIdHex(userId)},bson.M{"$push":bson.M{"addresses":address}})
	}

	if err!=nil {
		log.Fatalf("Error %s",err.Error())
	}

	user,err=r.GetById(userId)
    return user,err
}


func (r *UserRepository) RemoveAddress(userId string,address Address) (User,error) {
	var user User

	_,err:=r.C.Upsert(bson.M{"_id":bson.ObjectIdHex(userId),"addresses.name":address.Name},bson.M{"$pull":bson.M{"addresses":bson.M{"name":address.Name}}})

	if err!=nil {
		log.Fatalf("Error %s",err.Error())
	}

	user,err=r.GetById(userId)
	return user,err
}


func (r *UserRepository) UpdateAddressFindModify(userId string,address Address) (User,error) {
	var user User

	change:=mgo.Change{Update:bson.M{"$set":bson.M{"addresses.$":address}},Remove:false,ReturnNew:true,Upsert:true}
	_, err:=r.C.Find(bson.M{"_id":bson.ObjectIdHex(userId),"addresses.name":address.Name}).Apply(change,&user)

	if err!=nil{
		change:=mgo.Change{Update:bson.M{"$push":bson.M{"addresses":address}},Remove:false,ReturnNew:true,Upsert:true}
		_, err=r.C.Find(bson.M{"_id":bson.ObjectIdHex(userId)}).Apply(change,&user)
	}

	if err!=nil {
		log.Fatalf("Error %s",err.Error())
	}


	return user,err
}


