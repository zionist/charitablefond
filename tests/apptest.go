package tests

import "github.com/robfig/revel"
import "github.com/zionist/charitablefond/app/plugins"

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

func (t ApplicationTest) TestGetConfig() {
    p := mongodb.MongoDbPlugin{}
    p.TestMode = true
    //should be wriiten from conf file
    var found bool
    _, found = p.GetConfig("mongodb.host")
    t.Assert(found == true)
    _, found = p.GetConfig("mongodb.port")
    t.Assert(found == true)

    _, found = p.GetConfig("noconfig.noconfig")
    t.Assert(found == false)
}

//func (t ApplicationTest) TestThatSessionWorks() {
//    p := mongodb.MongoDbPlugin{}
//    p.OnAppStart()
//}

func (t ApplicationTest) After() {
	println("Tear down")
}
