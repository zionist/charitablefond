package controllers

import (
	"fmt"
  "io"
	"github.com/robfig/revel"
	"crypto/sha512"
	"github.com/zionist/charitablefond/app/constants"
	"github.com/zionist/charitablefond/app/models"
	"labix.org/v2/mgo/bson"
)

/*
  userController uses for user validation and auth
*/
type UserController struct {
	*revel.Controller
	MongoDbController
}

//GET auth page
func (c UserController) GetLoginPage() revel.Result {
	revel.INFO.Println("UserController.Auth started")
	/*
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
	*/
	return c.RenderTemplate("Page/Login.html")
}

func deny_login(c UserController) revel.Result {
  c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{"Не верный пользователь или пароль", ""})
  c.Validation.Keep()
  c.FlashParams()
  return c.Redirect(UserController.Login)
}

func (c UserController) Login(username, password string) revel.Result {
	revel.INFO.Println("User.Login started")
	collection := Session.DB(Base).C(constants.UsersCollectionName)
	result := models.User{}
	if err := collection.Find(bson.M{"username": username}).One(&result); err != nil {
		c.RenderError(err)
	}
	if result.Username != username {
      return deny_login(c)
	} else {
    hash := sha512.New()
    io.WriteString(hash, password)
    sum := fmt.Sprintf("%x", hash.Sum(nil))
    if sum != result.Password {
      return deny_login(c)
    } else {
	    fmt.Println(c.Session["user"])
	    c.Session["user"] = username
    }
  }
	return c.RenderTemplate("Page/Login.html")
}
