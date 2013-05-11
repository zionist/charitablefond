package controllers

import (
	"time"
	"github.com/robfig/revel"
)

type Application struct {
	MongoDbController
}

//Any method on a controller that is exported and returns a revel.Result may be treated as an Action.
func (c Application) Index() revel.Result {
	//fmt.Println("%+v", Session)
	collection := Session.DB(base).C("test")
	type Person struct {
		Name  string
		Phone string
		Time time.Time
	}
	collection.Insert(&Person{"Ale", "+55 53 8116 9639", time.Now()},
		&Person{"Cla", "+55 53 8402 8510", time.Now()})
	hello := "Hello world"
	revel.INFO.Println("started")
	//fmt.Printf("%+v", c)
	return c.Render(hello)
}

func (c Application) Page(url string) revel.Result {
	//revel.INFO.Printf("plain page %s load", url)
	c.RenderArgs["hello"] = url
	return c.RenderTemplate("Application/Index.html")
}
