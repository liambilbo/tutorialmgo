package main

import (
	. "github.com/liambilbo/tutorialmgo/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"fmt"
)

var productRepository *ProductRepository
var categoryRepository *CategoryRepository
var orderRepository *OrderRepository
var userRepository *UserRepository

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


}

func main() {
	createUpdateProduct()
	createUpdateCategory()
	findProduct()
	findCategoriesOfAProduct()
}

func date(d string) time.Time {
	/*
		t, _:=time.Parse("2017-01-20",d)
		return t*/
	return time.Now()
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

	createIndexSlug()

	if err := productRepository.Create(&product); err != nil {
		log.Fatalf("[Create Product] : %s \n", err)
	}

	/*

	idProduct = product.Id.Hex()

	fmt.Printf("[Create Product] id %s \n", idProduct)

	product.Tags = append(product.Tags, "ligth")

	if err := productRepository.Update(product); err != nil {
		log.Fatalf("[Update Product] %s ", err)
	}


*/


	return product
}

func getProductById(id string) Product {
	product,_:=productRepository.GetById(id)
	return product
}


func createIndexSlug() {

	index := mgo.Index{
		Key: []string{"slug"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	productRepository.C.EnsureIndex(index)
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

	orderRepository.Create(&order)



}



func findProduct(){
	products:=productRepository.GetByCategoryId("6a5b1476238d3b4dd5000048")
	for _,p:=range products{
		fmt.Println("Product Id : %s \n",p.Id)

	}
}

func findCategoriesOfAProduct(){
	product,err:=productRepository.GetBySlug("wheelbarrow-9092")
	if err!=nil {
		log.Fatalf("[GetBySlug] %s \n",err)
	}

	categories:=categoryRepository.GetCategoByIds(convertObjectIdToString(product.CategoryIds))

	for _,v:=range categories{
		fmt.Printf("[Category of %s] %s \n",product.Slug,v.Id.Hex())

	}

}

func convertObjectIdToString (ids []bson.ObjectId) []string {
	var result []string
	for  _,v:=range ids{
		result=append(result,v.Hex())
	}
	return result
}
