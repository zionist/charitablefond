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
}

//Front page
func (c PageController) Index() revel.Result {
	revel.INFO.Println("Index page")
	return c.Redirect(constants.FrontPage)
}

//GET page
func (c PageController) Page(url string) revel.Result {
	revel.INFO.Println("Page.Page started")
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
	return c.RenderTemplate("Page/Page.html")
}

//Admin pages
//List of pages
func (c PageController) AdminListPages() revel.Result {
	result := []models.Page{}
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	if err := collection.Find(bson.M{}).All(&result); err != nil {
		c.RenderError(err)
	}
	//TODO: add sorting
	//Cut content to 32 
	new_result := []models.Page{}
	for _, v := range result {
		if len(v.Content) > 64 {
			v.Content = v.Content[0:64]
		}
		new_result = append(new_result, v)
	}
	c.RenderArgs["pages"] = new_result
	return c.RenderTemplate("Page/AdminListPages.html")
}

//Delete plain page
//TODO: add permissions check for deleting
func (c PageController) AdminDelete(url string) revel.Result {
	if err := c.DelPages(url); err != nil {
		return c.RenderError(err)
	}
	c.RenderArgs["page_content"] = "Удалено"
	return c.RenderTemplate("Page/Page.html")
}

//Create creation page
func (c PageController) AdminCreatePage() revel.Result {
	return c.RenderTemplate("Page/AdminCreatePage.html")
}

//POST create plain pages 
func (c PageController) CreatePage(page_header, page_content, page_url string) revel.Result {
	revel.INFO.Println("Page.CreatePage started")
	c.Validation.MinSize(page_header, 1).Message("Требуется заголовок")
	c.Validation.MinSize(page_url, 1).Message("Требуется ссылка на страницу")
	c.Validation.MinSize(page_content, 1).Message("Требуется контент")
	if c.Validation.HasErrors() {
		revel.INFO.Printf("CreatePage validation errors %v", c.Validation.Errors)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(PageController.AdminCreatePage)
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
		return c.Redirect(PageController.AdminCreatePage)
	}
	//Save page
	p := models.Page{Header: page_header, Url: page_url, Content: page_content}
	if err := c.SavePage(p); err != nil {
		c.RenderError(err)
	}
	return c.RenderTemplate("Page/AdminPageCreated.html")
}

func (c PageController) CheckPageExists(url string) (err error, found bool) {
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
func (c PageController) AdminUpdatePage(url string) revel.Result {
	revel.INFO.Println("Page.UpdatePage started")
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
	return c.RenderTemplate("Page/AdminUpdatePage.html")
}

//POST update plain pages 
func (c PageController) UpdatePage(page_header, page_content, page_url string) revel.Result {
	revel.INFO.Println("Page.UpdatePage started")
	c.Validation.MinSize(page_header, 1).Message("Требуется заголовок")
	c.Validation.MinSize(page_url, 1).Message("Требуется ссылка на страницу")
	c.Validation.MinSize(page_content, 1).Message("Требуется контент")
	if c.Validation.HasErrors() {
		revel.INFO.Printf("CreatePage validation errors %v", c.Validation.Errors)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(PageController.AdminCreatePage)
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

func (c PageController) DelPages(url string) (err error) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	err = collection.Remove(bson.M{"url": url})
	revel.INFO.Printf("Pages with url %s removed", url)
	return
}

func (c PageController) SavePage(p models.Page) (err error) {
	collection := Session.DB(c.Base).C(constants.PageCollectionName)
	err = collection.Insert(&p)
	revel.INFO.Printf("Page %s saved", p.Url)
	return
}
