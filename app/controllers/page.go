package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/constants"
	"github.com/zionist/charitablefond/app/models"
	"labix.org/v2/mgo/bson"
)

type PageController struct {
	*revel.Controller
	MongoDbController
	UserController
}

//Front page
func (c PageController) Index() revel.Result {
	revel.INFO.Println("Index page")
	return c.Redirect(constants.FrontPage)
}

//GET page
func (c PageController) GetPage(url string) revel.Result {
	revel.INFO.Printf("Page.Page with url %s started", url)
	collection := Session.DB(Base).C(constants.PageCollectionName)
	query := collection.Find(bson.M{"url": url})
	count, err := query.Count()
	if err != nil {
		c.RenderError(err)
	}
	fmt.Println(count)
	if count <= 0 {
		return c.NotFound(c.Message("page_not_found"))
	} else if count > 1 {
		revel.ERROR.Println("There more than one page accesed by url")
		c.Response.Status = 500
		return c.RenderText("internal_server_error")
	}

	var result = models.Page{}
	query.One(&result)

	c.RenderArgs["page_header"] = result.Header
	c.RenderArgs["page_content"] = result.Content
	c.RenderArgs["page_content"] = result.Content
	if c.LoggedIn() == true {
		c.RenderArgs["logged"] = "true"
		c.RenderArgs["url"] = url
	}
	return c.RenderTemplate("Page/Page.html")
}
