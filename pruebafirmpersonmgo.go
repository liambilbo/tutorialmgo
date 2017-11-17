package main

import (
	"encoding/json"
	"fmt"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

const (
	firmKey     = "firm"
	physicalKey = "physical"
)

type omit *struct{}

type PersonRecord struct {
	RecordID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	PersonID string        `bson:"person_id",json:"person_id"`
	Type     string        `bson:"type",json:"type"`
	Name     string        `bson:"name",json:"name"`
}

type FirmRecord struct {
	PersonRecord `bson:",inline""`
}

type PhysicalRecord struct {
	PersonRecord  `bson:",inline""`
	Surname       string `bson:"surname",json:"surname"`
	SecondSurname string `bson:"secondsurname",json:"secondsurname"`
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("persons")

	firm := &FirmRecord{}
	firm.PersonID = uuid.New()

	firm.Name = "Lorea Engine S.A."

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(firm)

	err = c.Insert(firm)

	if err != nil {
		log.Fatal(err)
	}

	result := FirmRecord{}
	err = c.Find(bson.M{"name": "Lorea Engine S.A."}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name:%s , Id %s ", result.Name, result.PersonID)
}
