package service

import (
"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type CatalogRepository struct {
	C *mgo.Collection
}

func NewCatalogRepository(c *mgo.Collection) *CatalogRepository {
	return &CatalogRepository{c}
}

func (r *CatalogRepository) FindByText(search string,page Page) []struct{Id bson.ObjectId;Title string}{
	var resp []struct{Id bson.ObjectId;Title string}
	//err:=r.C.Find(bson.M{"$text":bson.M{"$search":search}}).Skip(page.skip()).Limit(page.limit()).Select(bson.M{"title":1}).All(&resp)
	err:=r.C.Find(bson.M{"$text":bson.M{"$search":search}}).Skip(page.skip()).Limit(page.limit()).Select(bson.M{"title":1}).All(&resp)

	if err!=nil {
		log.Fatalf("Error %s" , err.Error())
	}
	return resp
}

func (r *CatalogRepository) CountByText(search string) int{
	resp,err:=r.C.Find(bson.M{"$text":bson.M{"$search":search}}).Count()

	if err!=nil {
		log.Fatalf("Error %s" , err.Error())
	}
	return resp
}




