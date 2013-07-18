package controllers

import (
	"crypto/sha512"
	"fmt"
	"github.com/robfig/revel"
	"github.com/zionist/charitablefond/app/constants"
	"github.com/zionist/charitablefond/app/models"
	"io"
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
	return c.RenderTemplate("Page/Login.html")
}

func deny_login(c UserController) revel.Result {
	c.Validation.Errors = append(c.Validation.Errors, &revel.ValidationError{c.Message("wrong_login_or_pass"), ""})
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
			revel.INFO.Printf("user %s has logged in", username)
			c.Session["user"] = username
		}
	}
	return c.Redirect("/admin/pages/list")
}

func (c UserController) Logout() revel.Result {
	revel.INFO.Println("User.Logout started")
	if c.LoggedIn() {
		revel.INFO.Printf("user %s has logged out", c.Session["user"])
		c.Session["user"] = ""
		return c.Redirect(PageController.Index)
	}
	return c.Forbidden(c.Message("forbidden"))
}

//Check is user logged in
func (c UserController) LoggedIn() (hasper bool) {
	if len(c.Session["user"]) != 0 {
		hasper = true
	}
  return
}
