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

	report:=catalogRepository.FindByTextAndStatus("mongodb","MEAP",Page{Size:10,Number:1})
	assert.Equal(t,1,len(report),"Error FindByTextAndStatus")




}

func ExampleFindBooksByText() {

	books:=catalogRepository.FindByText("actions",Page{5,1})

	for _,v :=range books {
		fmt.Printf("Title : %s , Score : %f\n" , v.Title, v.Score)
	}
	// Output:
	//Title : Spring Batch in Action , Score : 11.666667
	//Title : Hadoop in Action , Score : 9.249910
	//Title : HTML5 in Action , Score : 9.052156
	//Title : Jess in Action , Score : 8.780835
	//Title : MongoDB in Action , Score : 8.772244

}

func ExampleFindByTextAggregation() {

	books:=catalogRepository.FindByTextAggregation2("actions",Page{5,1})

	for _,v :=range books {
		fmt.Printf("Title : %s , Score : %f\n" , v.Title, v.Score)
	}
	// Output:
	//Title : Spring Batch in Action , Score : 11.666667
	//Title : Hadoop in Action , Score : 9.249910
	//Title : HTML5 in Action , Score : 9.052156
	//Title : Jess in Action , Score : 8.780835
	//Title : MongoDB in Action , Score : 8.772244

}

func ExampleFindByTextSmartScore() {

	books:=catalogRepository.FindByTextSmartScore("\"mongodb\" in action","english",Page{5,1})

	for _,v :=range books {
		fmt.Printf("Title : %s , Score : %f , Multiplier : %f , AdjScore : %f\n" , v.Title, v.Score , v.Multiplier, v.AdjScore)
	}
	// Output:
    // Title : MongoDB in Action, Second Edition , Score : 12.500000 , Multiplier : 3.000000 , AdjScore : 0.000000
	// Title : MongoDB in Action , Score : 18.448653 , Multiplier : 1.000000 , AdjScore : 0.000000

}






