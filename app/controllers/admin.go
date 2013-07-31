package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/constants"
	"github.com/zionist/charitablefond/app/models"
	"labix.org/v2/mgo/bson"
)

type AdminController struct {
	*revel.Controller
	MongoDbController
	UserController
}

func (c AdminController) CheckContentExists(page_url string, collection string) (count int, err error) {
	coll := Session.DB(Base).C(collection)
	query := coll.Find(bson.M{"url": page_url})
	count, err = query.Count()
	return
}

func (c AdminController) GetAdminListPages() revel.Result {
	collection := Session.DB(Base).C(constants.PageCollectionName)
	result := []models.Page{}
	if err := collection.Find(bson.M{}).All(&result); err != nil {
		c.RenderError(err)
	}
	new_result := []models.Page{}
	// Cut content to 120
	for _, v := range result {
		if len(v.Content) > 120 {
			v.Content = v.Content[0:120]
		}
		new_result = append(new_result, v)
	}
	c.RenderArgs["content"] = new_result
	c.RenderArgs["title"] = "Страницы сайта"
	c.RenderArgs["content_type"] = "page"
	return c.RenderTemplate("Admin/ListPage.html")
}

//Admin pages
//List of content types
func (c AdminController) GetAdminListContent(content_type string) revel.Result {
	revel.INFO.Println(content_type)
	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
	//TODO: Make type cast
	if content_type == "page" {
		return c.GetAdminListPages()
	} else if content_type == "block" {

		//result := []models.Block{}
		//collection := Session.DB(Base).C(constants.BlockCollectionName)
	}
	c.Response.Status = 500
	revel.ERROR.Println("Wrong admin list type")
	return c.RenderText("internal_server_error")
	//TODO: add sorting
}

//TODO: add permissions check for deleting
func (c AdminController) GetAdminDeleteContent(content_type string, url string) revel.Result {

	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
	if content_type == "page" {
		if err := c.DelPages(url); err != nil {
			return c.RenderError(err)
		}
		c.RenderArgs["page_content"] = c.Message("deleted")
		c.RenderArgs["logged"] = "true"
		return c.RenderTemplate("Page/Page.html")
	}
	c.Response.Status = 500
	revel.ERROR.Println("Wrong content type")
	return c.RenderText("internal_server_error")
}

//GET any of Admin Page
func (c AdminController) GetAdminCreateContent(content_type string) revel.Result {
	revel.INFO.Printf("GetAdminContentCreate started for %s type", content_type)
	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
	revel.INFO.Println(content_type)
	if content_type == "page" {
		return c.RenderTemplate("Admin/CreatePage.html")
	}
	c.Response.Status = 500
	revel.ERROR.Println("Wrong content type")
	return c.RenderText("internal_server_error")
}

func (c AdminController) PostAdminCreatePage() revel.Result {
	revel.INFO.Println("Page.PostAdminCreatePage started")
	// params map
	args := c.Request.Request.Form
	page_header := args["page_header"][0]
	page_content := args["page_content"][0]
	page_url := args["page_url"][0]
	c.Validation.MinSize(page_header, 1).Message(c.Message("header_required"))
	c.Validation.MinSize(page_url, 1).Message(c.Message("url_required"))
	c.Validation.MinSize(page_content, 1).Message(c.Message("content_required"))
	if c.Validation.HasErrors() {
		revel.INFO.Printf("CreatePage validation errors %v", c.Validation.Errors)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/admin/create/page")
	}
	// TODO: Add permission (sessison check)
	count, err := c.CheckContentExists(page_url, constants.PageCollectionName)
	if err != nil {
		c.RenderError(err)
	}
	if count != 0 {
		c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{fmt.Sprintf("%s %s", c.Message("already_created"), page_url), ""})
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/admin/create/page")
	}
	// Save page
	p := models.Page{Header: page_header, Url: page_url, Content: page_content}
	if err := c.SavePage(p); err != nil {
		c.RenderError(err)
	}

	return c.Redirect("/admin/update/page/%s", page_url)
}

//POST create plain content
func (c AdminController) PostAdminCreateContent(content_type string) revel.Result {
	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
	return c.PostAdminCreatePage()
}

func (c AdminController) GetAdminUpdateContent(content_type, url string) revel.Result {
	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
    if content_type == "page" {
	    return c.GetAdminUpdatePage(url)
    }
    c.Response.Status = 500
    revel.ERROR.Printf("Wrong content type %s", content_type)
	return c.RenderText("internal_server_error")
}

//Create update page
func (c AdminController) GetAdminUpdatePage(url string) revel.Result {
	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
	revel.INFO.Println("Page.UpdatePage started")
	count, err := c.CheckContentExists(url, constants.PageCollectionName)
	if err != nil {
		c.RenderError(err)
	}
	if count <= 0 {
		return c.NotFound(c.Message("not_found"))
	} else if count > 1 {
		revel.ERROR.Println("There more than one page accesed by url")
		c.Response.Status = 500
		return c.RenderText("internal_server_error")
	}
	collection := Session.DB(Base).C(constants.PageCollectionName)
	result := models.Page{}
	if err = collection.Find(bson.M{"url": url}).One(&result); err != nil {
		c.RenderError(err)
	}
	c.RenderArgs["page_header"] = result.Header
	c.RenderArgs["page_content"] = result.Content
	c.RenderArgs["page_url"] = result.Url
	return c.RenderTemplate("Admin/UpdatePage.html")
}

func (c AdminController) PostAdminUpdateContent(content_type string) revel.Result {
	if !c.LoggedIn() {
		return c.Forbidden(c.Message("forbidden"))
	}
    if content_type == "page" {
        return c.PostAdminUpdatePage()
    }
    c.Response.Status = 500
    revel.ERROR.Println("Wrong content type")
	return c.RenderText("internal_server_error")
}

//POST update plain pages
func (c AdminController) PostAdminUpdatePage() revel.Result {
	revel.INFO.Println("Page.UpdatePage started")
	args := c.Request.Request.Form
	page_header := args["page_header"][0]
	page_content := args["page_content"][0]
	page_url := args["page_url"][0]
	c.Validation.MinSize(page_header, 1).Message(c.Message("header_required"))
	c.Validation.MinSize(page_url, 1).Message(c.Message("url_required"))
	c.Validation.MinSize(page_content, 1).Message(c.Message("content_required"))
	if c.Validation.HasErrors() {
		revel.INFO.Printf("CreatePage validation errors %v", c.Validation.Errors)
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(AdminController.GetAdminCreateContent(c, "page"))
	}
	//TODO: Add permission (session check)
	//Get page by url
	//TODO: Remove security hole (user can delete all pages with same url using hidden value)
	//Remove all pages with same url
	//TODO: Do real update not delete
	if err := c.DelPages(page_url); err != nil {
		return c.RenderError(err)
	}
	//Save page
	p := models.Page{Header: page_header, Url: page_url, Content: page_content}
	if err := c.SavePage(p); err != nil {
		return c.RenderError(err)
	}
	return c.Redirect("/admin/update/page/%s", page_url)
	//return c.Redirect(AdminController.GetAdminUpdatePage(c, page_url))
}

func (c AdminController) DelPages(url string) (err error) {
	collection := Session.DB(Base).C(constants.PageCollectionName)
	err = collection.Remove(bson.M{"url": url})
	revel.INFO.Printf("Pages with url %s removed", url)
	return
}

func (c AdminController) SavePage(p models.Page) (err error) {
	collection := Session.DB(Base).C(constants.PageCollectionName)
	err = collection.Insert(&p)
	revel.INFO.Printf("Page %s saved", p.Url)
	return
}
