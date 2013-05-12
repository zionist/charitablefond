package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/constants"
	"github.com/zionist/charitablefond/app/models"
	"labix.org/v2/mgo/bson"
	"time"
)

type Application struct {
	MongoDbController
}

//Any method on a controller that is exported and returns a revel.Result may be treated as an Action.
func (c Application) Index() revel.Result {
	collection := Session.DB(c.Base).C("test")
	type Person struct {
		Name  string
		Phone string
		Time  time.Time
	}
	collection.Insert(&Person{"Ale", "+55 53 8116 9639", time.Now()},
		&Person{"Cla", "+55 53 8402 8510", time.Now()})
	hello := "Hello world"
	revel.INFO.Println("started")
	return c.Render(hello)
}

//GET page
func (c Application) Page(url string) revel.Result {
	err, found := c.CheckPageExists(url)
	if err != nil {
		c.RenderError(err)
	}
	if found == false {
		return c.NotFound("Страница не найдена")
	}
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	result := models.Page{}
	if err = collection.Find(bson.M{"url": url}).One(&result); err != nil {
		c.RenderError(err)
	}
	fmt.Println(result.Content)
	c.RenderArgs["page_header"] = result.Header
	c.RenderArgs["page_content"] = result.Content
	return c.RenderTemplate("Application/Page.html")
}

//Admin pages
//Show list of pages
func (c Application) AdminListPages() revel.Result {
	return nil
}

//Create plain page
func (c Application) AdminCreatePage() revel.Result {
	return c.RenderTemplate("Application/AdminCreate.html")
}

//Update plain page
func (c Application) AdminUpdatePage() revel.Result {
	return nil
}

//Delete plain page
func (c Application) AdminDeletePage() revel.Result {
	return nil
}

//POST pages
func (c Application) CreatePage(page_header, page_content, page_url string) revel.Result {
	c.Validation.MinSize(page_header, 1).Message("Требуется заголовок")
	c.Validation.MinSize(page_url, 1).Message("Требуется ссылка на страницу")
	c.Validation.MinSize(page_content, 1).Message("Требуется контент")
	if c.Validation.HasErrors() {
		revel.INFO.Printf("CreatePage validation errors %v", c.Validation.Errors)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Application.AdminCreatePage)
	}
	//TODO: Add permission (sessison check)
	//Check page exists
	err, found := c.CheckPageExists(page_url)
	if err != nil {
		c.RenderError(err)
	}
	if found == true {
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{fmt.Sprintf("Страница со ссылкой %s уже  создана", page_url), ""})
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Application.AdminCreatePage)
	}
	//Save page
	p := models.Page{Header: page_header, Url: page_url, Content: page_content}
	if err := c.SavePage(p); err != nil {
		c.RenderError(err)
	}
	return c.RenderTemplate("Application/AdminCreated.html")
}

func (c Application) CheckPageExists(url string) (err error, found bool) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	result := models.Page{}
	empty := models.Page{}
	err = collection.Find(bson.M{"url": url}).One(&result)
	fmt.Println(result)
	if result != empty {
		found = true
	} else {
		found = false
	}
	return
}

func (c Application) SavePage(p models.Page) (err error) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	//collection := Session.DB("fond").C("pages")
	fmt.Println("%+v", Session)
	err = collection.Insert(&p)
	revel.INFO.Printf("Page %s saved", p.Url)
	return
}
