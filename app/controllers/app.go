package controllers

import (
	"github.com/robfig/revel"
)

type Application struct {
	PageController
}

//Test application
func (c Application) Index() revel.Result {
	revel.INFO.Println("started")
	return nil
}
