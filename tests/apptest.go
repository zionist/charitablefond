package tests

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/controllers"
	"github.com/zionist/charitablefond/app/models"
	"labix.org/v2/mgo/bson"
)

type ApplicationTest struct {
	revel.TestSuite
}

func (t ApplicationTest) Before() {
	println("Set up")
}

/*
func (t ApplicationTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html")
}
*/

func (t ApplicationTest) TestConfig() {
	p := controllers.MongoDbController{}
	_, f := p.GetConfig("mongodb.host", "")
	t.AssertEqual(f, true)
	_, f = p.GetConfig("mongodb.port", "")
	t.AssertEqual(f, true)

}

func (t ApplicationTest) TestConnectToDb() {
	p := controllers.MongoDbController{}
	p.Connect()
	type Person struct {
		Name  string
		Phone string
	}
	var err error
	c := controllers.Session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		panic(err)
	}
	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	t.AssertEqual(err, nil)

	t.AssertEqual(result.Phone, "+55 53 8116 9639")
}

func (t ApplicationTest) TestGetConnectionUrl() {
	p := controllers.MongoDbController{}
	p.GetConnectionUrl()
	t.AssertEqual(p.Url, "mongodb://127.0.0.1:27017")
}

func (t ApplicationTest) After() {
	println("Tear down")
}

//models tests
type ModelsTest struct {
	revel.TestSuite
}

func (t ModelsTest) TestModels() {
	fmt.Printf("%v", models.Page{})

}
