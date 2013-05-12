package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/constants"
	"github.com/zionist/charitablefond/app/models"
	"labix.org/v2/mgo/bson"
)

type Application struct {
	MongoDbController
}

//Any method on a controller that is exported and returns a revel.Result may be treated as an Action.
func (c Application) Index() revel.Result {
	hello := "Hello world"
	revel.INFO.Println("started")
	return c.Render(hello)
}

//GET page
func (c Application) Page(url string) revel.Result {
	revel.INFO.Println("Application.Page started")
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
	c.RenderArgs["page_header"] = result.Header
	c.RenderArgs["page_content"] = result.Content
	return c.RenderTemplate("Application/Page.html")
}

//Admin pages
//Show list of pages
func (c Application) AdminListPages() revel.Result {
	return nil
}

//Delete plain page
func (c Application) AdminDeletePage() revel.Result {
	return nil
}

//Create creation page
func (c Application) AdminCreatePage() revel.Result {
	return c.RenderTemplate("Application/AdminCreatePage.html")
}

//POST create plain pages 
func (c Application) CreatePage(page_header, page_content, page_url string) revel.Result {
	revel.INFO.Println("Application.CreatePage started")
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
	return c.RenderTemplate("Application/AdminPageCreated.html")
}

func (c Application) CheckPageExists(url string) (err error, found bool) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	result := models.Page{}
	empty := models.Page{}
	err = collection.Find(bson.M{"url": url}).One(&result)
	if result != empty {
		found = true
	} else {
		found = false
	}
	return
}

//Create update page
func (c Application) AdminUpdatePage(url string) revel.Result {
	revel.INFO.Println("Application.UpdatePage started")
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
	c.RenderArgs["page_header"] = result.Header
	c.RenderArgs["page_content"] = result.Content
	c.RenderArgs["page_url"] = result.Url
	return c.RenderTemplate("Application/AdminUpdatePage.html")
}

//POST update plain pages 
func (c Application) UpdatePage(page_header, page_content, page_url string) revel.Result {
	revel.INFO.Println("Application.UpdatePage started")
	c.Validation.MinSize(page_header, 1).Message("Требуется заголовок")
	c.Validation.MinSize(page_url, 1).Message("Требуется ссылка на страницу")
	c.Validation.MinSize(page_content, 1).Message("Требуется контент")
	if c.Validation.HasErrors() {
		revel.INFO.Printf("CreatePage validation errors %v", c.Validation.Errors)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Application.AdminCreatePage)
	}
	//TODO: Add permission (session check)
	//Get page by url
	//TODO: Remove security hole (user can delete all pages with same url using hidden value)
	//Remove all pages with same url
	if err := c.DelPages(page_url); err != nil {
		c.RenderError(err)
	}
	//Save page
	p := models.Page{Header: page_header, Url: page_url, Content: page_content}
	if err := c.SavePage(p); err != nil {
		c.RenderError(err)
	}
	return c.Redirect("/admin/update/%s", page_url)
}

func (c Application) DelPages(url string) (err error) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	err = collection.Remove(bson.M{"url": url})
	revel.INFO.Printf("Pages with url %s removed", url)
	return
}

func (c Application) SavePage(p models.Page) (err error) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	err = collection.Insert(&p)
	revel.INFO.Printf("Page %s saved", p.Url)
	return
}
