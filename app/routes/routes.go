// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/robfig/revel"


type tMongoDbController struct {}
var MongoDbController tMongoDbController


func (p tMongoDbController) Connect(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("MongoDbController.Connect", args).Url
}

func (p tMongoDbController) Disconnect(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("MongoDbController.Disconnect", args).Url
}


type tStatic struct {}
var Static tStatic


func (p tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (p tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (p tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (p tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (p tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


type tApplication struct {}
var Application tApplication


func (p tApplication) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Index", args).Url
}

func (p tApplication) Page(
		url string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "url", url)
	return revel.MainRouter.Reverse("Application.Page", args).Url
}


