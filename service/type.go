package service

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Details struct {
	Weight       int `bson:"weight" json:"weight"`
	WeightUnits  string `bson:"weight_units" json:"weight_units"`
	ModelNum     int `bson:"model_num" json:"model_num"`
	Manufacterer string `bson:"manufacturer" json:"manufacturer"`
	Color        string `bson:"color" json:"color"`
}

type Price struct {
	Retail int `bson:"retail" json:"retail"`
	Sale   int `bson:"sale" json:"sale"`
	RangeDate `bson:",inline"`
}

type RangeDate struct {
	Start time.Time `bson:"start,omitempty" json:"start"`
	End   time.Time `bson:"end,omitempty" json:"end"`
}

type Product struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	Slug            string  `bson:"slug" json:"slug"`
	Sku             string  `bson:"sku" json:"sku"`
	Name            string  `bson:"name" json:"name"`
	Description     string  `bson:"description" json:"description"`
	Details         Details  `bson:"details" json:"details"`
	TotalReviews    int      `bson:"total_reviews" json:"total_reviews"`
	AverageReview   float64  `bson:"average_review" json:"average_review"`
	Price           Price    `bson:"price" json:"price"`
	HistoricPrices  []Price  `bson:"price_history" json:"price_history"`
	PrimaryCategory bson.ObjectId `bson:"primary_category" json:"primary_category"`
	CategoryIds     []bson.ObjectId `bson:"category_ids" json:"category_ids"`
	MainCatId       bson.ObjectId `bson:"main_cat_id" json:"main_cat_id"`
	Tags            []string `bson:"tags" json:"tags"`
}

type CategoryId struct {
	Id bson.ObjectId `bson:"_id" json:"_id"`
	Slug string  `bson:"slug" json:"slug"`
	Name string `bson:"name" json:"name"`
}
type Category struct {
	CategoryId `bson:",inline"`
	Description string `bson:"description" json:"description"`
	ParentId bson.ObjectId `bson:"parent_id" json:"parent_id"`
	Ancestors []CategoryId `bson:"ancestors" json:"ancestors"`
}

type LineItem struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	Sku string `bson:"sku" json:"sku"`
	Name string `bson:"name" json:"name"`
	Quantity int `bson:"quantity" json:"quantity"`
	Price Price `bson:"price" json:"price"`
}

type Address struct {
	Street string `bson:"street" json:"street"`
	City string `bson:"city" json:"city"`
	State string `bson:"state" json:"state"`
	Zip int `bson:"zip" json:"zip"`
	Name string `bson:"name,omitempty" json:"name"`
	}

type PaymentMethod struct {
	Name string `bson:"name" json:"name"`
	PaymentToken string `bson:"payment_token" json:"payment_token"`
}


type Order struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	UserId bson.ObjectId `bson:"user_id" json:"user_id"`
	State string `bson:"state" json:"state"`
	LineItems []LineItem `bson:"line_item" json:"line_item"`
	ShippingAddress Address `bson:"shipping_address" json:"shipping_address"`
	Subtotal int `bson:"subtotal" json:"subtotal"`
}

type User struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	UserName string `bson:"username" json:"username"`
	Email string`bson:"email" json:"email"`
	FirstName string`bson:"first_name" json:"first_name"`
	LastName string`bson:"last_name" json:"last_name"`
	HashedPassword string `bson:"hashed_password" json:"hashed_password"`
	Addresses []Address `bson:"addresses" json:"addresses"`
	PaymentMethods PaymentMethod `bson:"payment_methods" json:"payment_methods"`
}


type Review struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	ProductId bson.ObjectId `bson:"product_id" json:"product_id"`
	Date time.Time`bson:"date" json:"date"`
	Title string`bson:"title" json:"title"`
	Text string`bson:"text" json:"text"`
	Rating  int `bson:"rating" json:"rating"`
	UserId bson.ObjectId `bson:"user_id" json:"user_id"`
	UserName string `bson:"user_name" json:"user_name"`
	HelpfulVotes  int   `bson:"helpful_votes" json:"helpful_votes"`
	VoterIds []bson.ObjectId `bson:"voter_ids" json:"voter_ids"`
}

type Book struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
	Title string `bson:"title" json:"title"`
	Isbn string `bson:"isbn" json:"isbn"`
	PageCount int `bson:"pageCount" json:"pageCount"`
	PublishedDate time.Time `bson:"publishedDate" json:"publishedDate"`
	TumnnailUrl string `bson:"tumnnailUrl" json:"tumnnailUrl"`
	LongDescription string `bson:"longDescription" json:"longDescription"`
	Status string `bson:"longDescription" json:"longDescription"`
	Authors []string `bson:"authors" json:"authors"`
	Categories []string `bson:"categories" json:"categories"`
}



type Page struct{
	Size int
	Number int
}

func (p *Page) add(pos int) Page {
	return Page{p.Size,p.Number + pos}
}

func (p *Page) skip() int {
	return p.Size * (p.Number - 1)
}

func (p *Page) limit() int {
	return p.Size
}