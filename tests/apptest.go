package tests

import (
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/plugins"
	"labix.org/v2/mgo/bson"
)

type ApplicationTest struct {
	revel.TestSuite
}

func (t ApplicationTest) Before() {
	println("Set up")
}

func (t ApplicationTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html")
}

func (t ApplicationTest) TestConfig() {
	p := mongodb.MongoDbPlugin{}
	_, f := p.GetConfig("mongodb.host", "")
	t.AssertEqual(f, true)
	_, f = p.GetConfig("mongodb.port", "")
	t.AssertEqual(f, true)

}

func (t ApplicationTest) TestConnectToDb() {

	p := mongodb.MongoDbPlugin{}
	p.Connect()
	type Person struct {
		Name  string
		Phone string
	}
	var err error
	c := p.Session.DB("test").C("people")
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
	p := mongodb.MongoDbPlugin{}
	p.GetConnectionUrl()
	t.AssertEqual(p.Url, "mongodb://127.0.0.1:27017")
}

func (t ApplicationTest) After() {
	println("Tear down")
}
