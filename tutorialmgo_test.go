package main

import (
	. "github.com/liambilbo/tutorialmgo/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

var productRepository *ProductRepository
var categoryRepository *CategoryRepository
var orderRepository *OrderRepository
var userRepository *UserRepository
var reviewRepository *ReviewRepository

var idProduct string
var idCategory string

func init() {

	session, err := mgo.Dial("localhost")

	session.SetMode(mgo.Monotonic, false)

	if err != nil {
		log.Fatalf("[MongoDB session] ;: %s \n", err)
	}

	colproducts := session.DB("tutorial").C("products")
	colproducts.RemoveAll(nil)
	productRepository = NewProductRepository(colproducts)

	colcategories := session.DB("tutorial").C("categories")
	colcategories.RemoveAll(nil)
	categoryRepository = NewCategoryRepository(colcategories)

	colusers := session.DB("tutorial").C("users")
	colusers.RemoveAll(nil)
	userRepository = NewUserRepository(colusers)

	colorders := session.DB("tutorial").C("orders")
	colorders.RemoveAll(nil)
	orderRepository = NewOrderRepository(colorders)


	colreviews := session.DB("tutorial").C("reviews")
	colreviews.RemoveAll(nil)
	reviewRepository = NewReviewRepository(colreviews)
}


func createUpdateProduct() Product {

	product := Product{
		Slug:        "wheelbarrow-9092",
		Name:        "Extra Large Wheelbarrow",
		Sku:         "90923",
		Description: "Heavy duty wheelbarrow...",
		Details: Details{
			Weight:       47,
			WeightUnits:  "lbs",
			ModelNum:     4039283402,
			Manufacterer: "Acme",
			Color:        "Green",
		},
		TotalReviews:  4,
		AverageReview: 4.5,
		Price: Price{
			Retail: 589700,
			Sale:   489700,
		},
		HistoricPrices: []Price{
			Price{Sale: 429700, Retail: 529700, RangeDate: RangeDate{Start: date("2010-4-1"), End: date("2010-4-8")}},
			Price{Sale: 529700, Retail: 529700, RangeDate: RangeDate{Start: date("2010-4-9"), End: date("2010-4-6")}},
		},
		PrimaryCategory: bson.ObjectIdHex("6a5b1476238d3b4dd5000048"),
		CategoryIds:     []bson.ObjectId{bson.ObjectIdHex("6a5b1476238d3b4dd5000048"), bson.ObjectIdHex("6a5b1476238d3b4dd5000049")},
		MainCatId:       bson.ObjectIdHex("6a5b1476238d3b4dd5000048"),
		Tags:            []string{"tools", "gardening", "soil"},
	}


	if err := productRepository.Create(&product); err != nil {
		log.Fatalf("[Create Product] : %s \n", err)
	}



	idProduct = product.Id.Hex()

	fmt.Printf("[Create Product] id %s \n", idProduct)


	/*
		product.Tags = append(product.Tags, "ligth")

		if err := productRepository.Update(product); err != nil {
			log.Fatalf("[Update Product] %s ", err)
		}


	*/


	return product
}

func TestAll(t *testing.T) {
	createUpdateProduct()
	createProductIndexSlug()
	createProductIndexTags()
	createUpdateCategory()
	createUpdateReview()
	createReviewIndexProductId()
	createUpdateOrder()
	createUpdateUser()
    t.Run("P=1",FindProducts)
	t.Run("C=1",FindCategories)
	t.Run("R=1",FindReviews)
	t.Run("U=1",FindReviews)
}

func FindProducts(t *testing.T) {
	products:=productRepository.GetByColorAndManufacturer("Green","Acme",Page{1,1})
	assert.NotZerof(t,len(products),"[GetByColorAndManufacturer] Product not founded")

	products=productRepository.GetByCategoryId("6a5b1476238d3b4dd5000048")
	assert.NotZerof(t,len(products),"[GetByCategoryId] Product not founded")

	products=productRepository.GetByColorExists(true,Page{1,1})
	assert.NotZerof(t,len(products),"[GetByColorExists] Product not founded")

	products=productRepository.GetByColorExists(false,Page{1,1})
	assert.Zerof(t,len(products),"[GetByColorExists] Product founded")

	products=productRepository.GetByFirstTag("tools",Page{1,1})
	assert.NotZerof(t,len(products),"[GetByFirstTag] Product not founded")

	products=productRepository.GetByFirstTag("soil",Page{1,1})
	assert.Zerof(t,len(products),"[GetByFirstTag] Product founded")
}

func FindCategories(t *testing.T) {
	product,err:=productRepository.GetBySlug("wheelbarrow-9092")
	if err != nil {
		t.Fail()
	}

	categories:=categoryRepository.GetCategoByIds(convertObjectIdToString(product.CategoryIds))
	assert.NotZerof(t,len(categories),"[GetCategoByIds] Categories not founded")
}

func FindReviews(t *testing.T) {
	product,err:=productRepository.GetBySlug("wheelbarrow-9092")
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
	assert.Equal(t,1,count,"[CountByProduct] Count not equal 1")
	assert.Equal(t,float64(4),average,"[CountByProduct] Average not equal 4")
}

func FindUsers(t *testing.T) {

	user,err:=userRepository.GetByNameAndPassword("kbanker","bd1cfa194c3a603e7186780824b04419")

	if err != nil {
		t.Fail()
	}
	assert.Equal(t,user.Id.Hex(),"4c4b1476238d3b4dd5000001","Error GetByNameAndPassword")

	users:=userRepository.GetByLastName("Banker",Page{1,1})
	assert.NotZerof(t,len(users),"[GetByLastName] Users not founded")
	users=userRepository.GetByLastNamePattern("^Ba",Page{1,1})
	assert.NotZerof(t,len(users),"[GetByLastNamePattern] Users not founded")
	users=userRepository.GetByZip(10019,10040,Page{1,1})
	assert.NotZerof(t,len(users),"[GetByZip] Users not founded")

	users=userRepository.GetByFirstAddressState("NY",Page{1,1})
	assert.NotZerof(t,len(users),"[GetByFirstAddressState] Users not founded")

	users=userRepository.GetByAddressElem(bson.M{"name":"home","state":"NY"},Page{1,1})
	assert.NotZerof(t,len(users),"[GetByAddressElem] Users not founded")

	users=userRepository.GetByAddressSize(2,Page{1,1})
	assert.NotZerof(t,len(users),"[GetByAddressSize] Users not founded")

}


func getProductById(id string) Product {
	product,_:=productRepository.GetById(id)
	return product
}





func findReviewsOfAProduct(){
	product,err:=productRepository.GetBySlug("wheelbarrow-9092")
	if err!=nil {
		log.Fatalf("[GetBySlug] %s \n",err)
	}
	reviews :=reviewRepository.GetByProductId(product.Id.Hex(),Page{20,0})
	for _,v:=range reviews{
		fmt.Printf("[Reviews of product %s] %s \n",product.Slug,v.Id.Hex())
	}



}


func convertObjectIdToString (ids []bson.ObjectId) []string {
	var result []string
	for  _,v:=range ids{
		result=append(result,v.Hex())
	}
	return result
}







func date(d string) time.Time {
	t, _:=time.Parse("2017-01-20",d)
	return t
}

func createProductIndexSlug() {

	index := mgo.Index{
		Key: []string{"slug"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productRepository.C.EnsureIndex(index)
}

func createProductIndexTags() {

	index := mgo.Index{
		Key: []string{"tags"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productRepository.C.EnsureIndex(index)
}

func createReviewIndexProductId() {

	index := mgo.Index{
		Key: []string{"product_id"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	reviewRepository.C.EnsureIndex(index)
}


func createUserIndexAddressState() {

	index := mgo.Index{
		Key: []string{"addresses.state"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productRepository.C.EnsureIndex(index)
}

func createUpdateReview(){
	review:=Review{Id:bson.ObjectIdHex("4c4b1476238d3b4dd5000041"),
		ProductId:bson.ObjectIdHex(idProduct),
		Date:date("2010-5-7"),
		Title:"Amazing",
		Text:"Has a squeaky wheel, but still a darn good wheelbarrow.",
		Rating:4,
		UserId:bson.ObjectIdHex("4c4b1476238d3b4dd5000042"),
		UserName:"dgreenthumb",
		HelpfulVotes:3,
		VoterIds:[]bson.ObjectId{
			bson.ObjectIdHex("4c4b1476238d3b4dd5000033"),
			bson.ObjectIdHex("7a4f0376238d3b4dd5000003"),
			bson.ObjectIdHex("92c21476238d3b4dd5000032"),
		},
	}

	if err:=reviewRepository.Create(&review); err!=nil {
		log.Fatalf("[Creating Review] %s",err)
	}

}

func createUpdateCategory(){
	category:=Category{
		CategoryId:CategoryId{Id:bson.ObjectIdHex("6a5b1476238d3b4dd5000048"),
			Name:"Gardening Tools",
			Slug:"gardening-tools",},
		Description:"Gardening gadgets galore!",
		ParentId:bson.ObjectIdHex("55804822812cb336b78728f9"),
		Ancestors:[]CategoryId{CategoryId{Id:bson.ObjectIdHex("558048f0812cb336b78728fa"),Name:"Home",Slug:"home"},
			CategoryId{Id:bson.ObjectIdHex("55804822812cb336b78728f9"),Name:"Outdoors",Slug:"outdoors"},},
	}


	if err:=categoryRepository.Create(&category);err!=nil{
		log.Fatalf("[Creating Category] %s",err)
	}

	idCategory=category.Id.Hex()
	fmt.Printf("[Create Category] id %s \n", idProduct)
}

func createUpdateOrder(){
	order:=Order{Id:bson.ObjectIdHex("6a5b1476238d3b4dd5000048"),
		UserId:bson.ObjectIdHex("4c4b1476238d3b4dd5000001"),
		State:"CART",
		LineItems:[]LineItem{LineItem{Id:bson.ObjectIdHex("4c4b1476238d3b4dd5003981"),
			Name:"Extra Large Wheelbarrow",
			Price:Price{Sale:4897,Retail:5897,},
			Sku:"9092",
			Quantity:1,},
			{Id:bson.ObjectIdHex("4c4b1476238d3b4dd5003982"),
				Name:"Rubberized Work Glove, Black",
				Price:Price{Sale:1299,Retail:1499,},
				Sku:"10027",
				Quantity:2,},
		},
		ShippingAddress:Address{State:"NY",
			Street: "588 5th Street",
			City:"Brooklyn",
			Zip:11215,},
		Subtotal:1028,
	}

	if err:=orderRepository.Create(&order); err!=nil {
		log.Fatalf("[Creating Order] %s",err)
	}

}


func createUpdateUser(){
	user:=User{Id:bson.ObjectIdHex("4c4b1476238d3b4dd5000001"),
		UserName:"kbanker",
		Email:"kylebanker@gmail.com",
		FirstName:"Kyle",
		LastName:"Banker",
		HashedPassword:"bd1cfa194c3a603e7186780824b04419",
		Addresses:[]Address{Address{Name:"home",
			Street:"588 5th Street",
			City:"Brooklyn",
			State:"NY",
			Zip: 11215},
			Address{Name:"work",
				Street:"1 E. 23rd Street",
				City:"New York",
				State:"NY",
				Zip: 10010},
		},
	}

	if err:=userRepository.Create(&user); err!=nil {
		log.Fatalf("[Creating User] %s",err)
	}
}
