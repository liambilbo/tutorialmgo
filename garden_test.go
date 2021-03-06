package main

import (
	. "github.com/liambilbo/tutorialmgo/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

var productGRepository *ProductRepository
var categoryGRepository *CategoryRepository
var orderGRepository *OrderRepository
var userGRepository *UserRepository
var reviewGRepository *ReviewRepository


func init() {

	session, err := mgo.Dial("localhost")

	session.SetMode(mgo.Monotonic, false)

	if err != nil {
		log.Fatalf("[MongoDB session] ;: %s \n", err)
	}

	colproducts := session.DB("garden").C("products")
	productGRepository = NewProductRepository(colproducts)

	colcategories := session.DB("garden").C("categories")
	categoryGRepository = NewCategoryRepository(colcategories)

	colusers := session.DB("garden").C("users")
	userGRepository = NewUserRepository(colusers)

	colorders := session.DB("garden").C("orders")
	orderGRepository = NewOrderRepository(colorders)


	colreviews := session.DB("garden").C("reviews")
	reviewGRepository = NewReviewRepository(colreviews)


	//createGardenProductIndexSlug()
	createGardenProductIndexText()
	createGardenProductIndexTags()
	createGardenReviewIndexProductId()
	createGardenUserIndexAddressState()
}



func TestGardenAll(t *testing.T) {
    //t.Run("GP=1",FindGardenProducts)
	//t.Run("GC=1",FindGardenCategories)
	//t.Run("GR=1",FindGardenReviews)
	t.Run("GU=1",FindGardenUsers)
	//t.Run("GO=1",FindGardenOrders)
}

func FindGardenProducts(t *testing.T) {
	products:=productGRepository.GetByColorAndManufacturer("Green","Acme",Page{1,1})
	assert.NotZerof(t,len(products),"[GetByColorAndManufacturer] Product not founded")

	products=productGRepository.GetByCategoryId("6a5b1476238d3b4dd5000048")
	assert.NotZerof(t,len(products),"[GetByCategoryId] Product not founded")

	products=productGRepository.GetByColorExists(true,Page{1,1})
	assert.NotZerof(t,len(products),"[GetByColorExists] Product not founded")

	products=productGRepository.GetByColorExists(false,Page{1,1})
	assert.Zerof(t,len(products),"[GetByColorExists] Product founded")

	products=productGRepository.GetByFirstTag("tools",Page{1,1})
	assert.NotZerof(t,len(products),"[GetByFirstTag] Product not founded")
	products=productGRepository.GetByFirstTag("soil",Page{1,1})
	assert.Zerof(t,len(products),"[GetByFirstTag] Product founded")

	productGRepository.LoadProductsByCategory()
	preport:=productGRepository.FindProductByText("gardens")
	assert.Zerof(t,len(preport),"[FindProductByText] Product founded")


}

func FindGardenCategories(t *testing.T) {
	product,err:=productGRepository.GetBySlug("wheel-barrow-9092")
	if err != nil {
		t.Fail()
	}

	categories:=categoryGRepository.GetCategoByIds(convertObjectIdToString(product.CategoryIds))
	assert.NotZerof(t,len(categories),"[GetCategoByIds] Categories not founded")
}

func FindGardenReviews(t *testing.T) {
	product,err:=productGRepository.GetBySlug("wheel-barrow-9092")
	if err != nil {
		t.Fail()
	}

	reviews :=reviewGRepository.GetByProductId(product.Id.Hex(),Page{20,1})
	assert.NotZerof(t,len(reviews),"[GetByProductId] Reviews not founded")

	reviews =reviewGRepository.GetByWhere("function() { return this.helpful_votes > 2; }",Page{20,1})
	assert.NotZerof(t,len(reviews),"[GetByWhere] Reviews not founded")

	reviews =reviewGRepository.GetByWhere("function() { return this.helpful_votes > 2; }",Page{20,1})
	assert.NotZerof(t,len(reviews),"[GetByWhere] Reviews not founded")

	reviews =reviewGRepository.GetByWhere("(this.helpful_votes) > 2",Page{20,1})
	assert.NotZerof(t,len(reviews),"[GetByWhere] Reviews not founded")

	reviews =reviewGRepository.GetByText(bson.RegEx{"Wheel|worst"," i"},Page{20,1})
	assert.NotZerof(t,len(reviews),"[GetByText] Reviews not founded")

	average ,count:=reviewGRepository.CountByProductId(product.Id.Hex())
	assert.Equal(t,3,count,"[CountByProduct] Count not equal 1")
	assert.Equal(t,float64(4.333333333333333),average,"[CountByProduct] Average not equal 4")

	ratings:=reviewGRepository.CountRatingByProductId(product.Id.Hex())
	assert.Equal(t,ratings[5],1,"[CountRatingByProductId] Count not equal 1")

}

func FindGardenUsers(t *testing.T) {

	user,err:=userGRepository.GetByNameAndPassword("kbanker","bd1cfa194c3a603e7186780824b04419")

	if err != nil {
		t.Fail()
	}
	assert.Equal(t,user.Id.Hex(),"4c4b1476238d3b4dd5000001","Error GetByNameAndPassword")

	address:=Address{Name:"mycasa",Zip:12003,Street:"Verca suena",City:"Larrabeyu",State:"BK"}
	usernew,err:=userGRepository.UpdateAddress(user.Id.Hex(),address)


	if err != nil {
		t.Fail()
	}
	assert.Equal(t,usernew.Id,user.Id,"[UpdateAddress] Users not founded")


	address=Address{Name:"mycasa",Zip:12003,Street:"Completan de rozas",City:"Larrabeyu",State:"BK"}
	usernew,err=userGRepository.UpdateAddressFindModify(user.Id.Hex(),address)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t,usernew.Id,user.Id,"[UpdateAddress] Users not founded")



	usernew,err=userGRepository.RemoveAddress(user.Id.Hex(),address)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t,usernew.Id,user.Id,"[RemoveAddress] Users not founded")

	users:=userGRepository.GetByLastName("Banker",Page{1,1})
	assert.NotZerof(t,len(users),"[GetByLastName] Users not founded")
	users=userGRepository.GetByLastNamePattern("^Ba",Page{1,1})
	assert.NotZerof(t,len(users),"[GetByLastNamePattern] Users not founded")
	users=userGRepository.GetByZip(10019,10040,Page{1,1})
	assert.NotZerof(t,len(users),"[GetByZip] Users not founded")

	users=userGRepository.GetByFirstAddressState("NY",Page{1,1})
	assert.NotZerof(t,len(users),"[GetByFirstAddressState] Users not founded")

	users=userGRepository.GetByAddressElem(bson.M{"name":"home","state":"NY"},Page{1,1})
	assert.NotZerof(t,len(users),"[GetByAddressElem] Users not founded")

	users=userGRepository.GetByAddressSize(2,Page{1,1})
	assert.NotZerof(t,len(users),"[GetByAddressSize] Users not founded")




}


func FindGardenOrders(t *testing.T) {

	report:=orderGRepository.GetReportByMonthAfter(date("2017-01-20"))

	reportline:=report["2014-8"]

	assert.Equal(t,reportline.Number,2,"Error GetReportByMonthAfter")
	assert.Equal(t,reportline.Subtotal,11093,"Error GetReportByMonthAfter")


	user:=orderGRepository.GetBestManhatanCustomer(date("2010-07-01"))
	assert.Equal(t,user.UserId,"4c4b1476238d3b4dd5000002","Error GetBestManhatanCustomer")

	orderGRepository.ReportCustomerOrders(date("2010-07-01"))

}



func createGardenProductIndexSlug() {

	index := mgo.Index{
		Key: []string{"slug"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productGRepository.C.EnsureIndex(index)
}


func createGardenProductIndexText() {

	index := mgo.Index{
		Key: []string{"$text:name","$text:description","$text:tags"},
	}

	productGRepository.C.EnsureIndex(index)
}

func createGardenProductIndexTags() {

	index := mgo.Index{
		Key: []string{"tags"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productGRepository.C.EnsureIndex(index)
}

func createGardenReviewIndexProductId() {

	index := mgo.Index{
		Key: []string{"product_id"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	reviewGRepository.C.EnsureIndex(index)
}


func createGardenUserIndexAddressState() {

	index := mgo.Index{
		Key: []string{"addresses.state"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productGRepository.C.EnsureIndex(index)
}
