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

func (r *CatalogRepository) FindByText(search string,page Page) []struct{Id bson.ObjectId;Title string;Score float64}{
	var resp []struct{Id bson.ObjectId;Title string;Score float64}
	//err:=r.C.Find(bson.M{"$text":bson.M{"$search":search}}).Skip(page.skip()).Limit(page.limit()).Select(bson.M{"title":1}).All(&resp)
	err:=r.C.Find(bson.M{"$text":bson.M{"$search":search}}).Skip(page.skip()).Limit(page.limit()).
		Select(bson.M{"_id":0,"title":1,"score":bson.M{"$meta":"textScore"}}).
				Sort("$textScore:score").All(&resp)

	if err!=nil {
		log.Fatalf("Error %s" , err.Error())
	}
	return resp
}

func (r *CatalogRepository) FindByTextAggregation(search string,page Page) []struct{Id bson.ObjectId;Title string;Score float64}{
	var resp []struct{Id bson.ObjectId;Title string;Score float64}

	err:=r.C.Pipe([]bson.M{
		bson.M{"$match":bson.M{"$text":bson.M{"$search":search}}},
		bson.M{"$sort":bson.M{"score":bson.M{"$meta":"textScore"}}},
		bson.M{"$project":bson.M{"title":1,"score":bson.M{"$meta":"textScore"}}},
		bson.M{"$skip":page.skip()},
		bson.M{"$limit":page.limit()},
		},
	).All(&resp)


	if err!=nil {
		log.Fatalf("Error %s" , err.Error())
	}
	return resp
}

func (r *CatalogRepository) FindByTextAggregation2(search string,page Page) []struct{Id bson.ObjectId;Title string;Score float64}{
	var resp []struct{Id bson.ObjectId;Title string;Score float64}

	err:=r.C.Pipe([]bson.M{
		bson.M{"$match":bson.M{"$text":bson.M{"$search":search}}},
		bson.M{"$project":bson.M{"title":1,"score":bson.M{"$meta":"textScore"}}},
		bson.M{"$sort":bson.M{"score":-1}},
		bson.M{"$skip":page.skip()},
		bson.M{"$limit":page.limit()},
	},
	).All(&resp)


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

func (r *CatalogRepository) FindByTextAndStatus(search string,status string,page Page) []struct{Title string;Status string}{
	var resp []struct{Title string;Status string}

	err:=r.C.Find(bson.M{"$text":bson.M{"$search":search},"status":status}).Skip(page.skip()).Limit(page.limit()).Select(bson.M{"_id":0,"title":1,"status":1}).All(&resp)
	if err!=nil {
		log.Fatalf("Error %s" , err.Error())
	}
	return resp
}

func (r *CatalogRepository) FindByTextSmartScore(search string,language string,page Page) []struct{Title string;Score float64;Multiplier float64;AdjScore float64}{
	var resp []struct{Title string;Score float64;Multiplier float64;AdjScore float64}

	err:=r.C.Pipe([]bson.M{
		bson.M{"$match":bson.M{"$text":bson.M{"$search":search,"$language":language}}},
		bson.M{"$project":bson.M{"title":1,
		                         "score":bson.M{"$meta":"textScore"},
		"multiplier":bson.M{"$cond":[]interface{}{"$longDescription",1.0,3.0}}}},
		bson.M{"$project":bson.M{"_id":0,
		                         "title":1,
			                     "score":1,
								  "multiplier":1,
			                     "adjScore":bson.M{"$multiply":[]interface{}{"$score","$multiplier"}}},
		                     },
		bson.M{"$project":bson.M{"_id":0,
			"title":1,
			"score":1,
			"multiplier":1,
			"adjScore":1,},
		},
		bson.M{"$sort":bson.M{"adjScore":-1}},
		bson.M{"$skip":page.skip()},
		bson.M{"$limit":page.limit()},
	},
	).All(&resp)


	if err!=nil {
		log.Fatalf("Error %s" , err.Error())
	}
	return resp
}




