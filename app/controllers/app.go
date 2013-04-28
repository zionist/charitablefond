package controllers

import (
    "github.com/robfig/revel"
)

type Application struct {
	*revel.Controller
}

//Any method on a controller that is exported and returns a revel.Result may be treated as an Action.
func (c Application) Index() revel.Result {
    hello := "Hello world"
    revel.INFO.Println("started")
	return c.Render(hello)
}
