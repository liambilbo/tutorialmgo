package main

import (
	. "github.com/liambilbo/tutorialmgo/service"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

var catalogRepository *CatalogRepository


func init() {

	session, err := mgo.Dial("localhost")

	session.SetMode(mgo.Monotonic, false)

	if err != nil {
		log.Fatalf("[MongoDB session] ;: %s \n", err)
	}

	books := session.DB("catalog").C("books")
	catalogRepository = NewCatalogRepository(books)

	createBookIndexText()


}



func TestCatalogAll(t *testing.T) {
    t.Run("CB=1",FindBooks)
}


func createBookIndexText() {

	weights:=make(map[string]int)

	weights["title"]=10
	weights["categories"]=5

	index := mgo.Index{
		Key: []string{/*"$text:'$**'",*/"$text:title","$text:shortDescription","$text:longDescription","$text:authors","$text:categories"},
		Name:"itext",
		Weights: weights,
	}

	err:=catalogRepository.C.EnsureIndex(index)

	if err!=nil {
		log.Fatalf("%s",err.Error())
	}
}

func FindBooks(t *testing.T) {

	numberbooks1:=catalogRepository.CountByText("\"mongodb\" in action")
	assert.Equal(t,2,numberbooks1,"Error CountByText")
	numberbooks2:=catalogRepository.CountByText("mongodb in action")
	assert.Equal(t,188,numberbooks2,"Error CountByText")
	numberbooks3:=catalogRepository.CountByText("mongodb")
	assert.Equal(t,2,numberbooks3,"Error CountByText")
	numberbooks4:=catalogRepository.CountByText("mongodb -second")
	assert.Equal(t,1,numberbooks4,"Error CountByText")
	numberbooks5:=catalogRepository.CountByText("mongodb -\"second edition\"")
	assert.Equal(t,1,numberbooks5,"Error CountByText")

}

func ExampleFindBooksByText() {

	books:=catalogRepository.FindByText("actions",Page{5,1})

	for _,v :=range books {
		fmt.Printf("Title : %s\n" , v.Title)
	}
	// Output:
	// Title : SQR in PeopleSoft and Other Applications
	// Title : Algorithms of the Intelligent Web
	// Title : Java Persistence with Hibernate
	// Title : SOA Security
	// Title : Ext JS in Action, Second Edition

}