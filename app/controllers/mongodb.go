/*
Package mongodb implements simple mongodb connection wrapper 
*/
package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"github.com/zionist/charitablefond/app/constants"
)


var required_configs = [...]string{"mongodb.host", "mongodb.port"}

var (
	Session   *mgo.Session
)

type MongoDbController struct {
	*revel.Controller
	Host      string
	Port      string
	User      string
	Password  string
	Base      string
	Url       string
}

//revel.Config.String wrapper
func (p *MongoDbController) GetConfig(config_string string, defaul string) (value string, found bool) {
	value, found = revel.Config.String(config_string)
	if !found {
		value = defaul
	}
	return
}

// Connect to mongodb
func (c *MongoDbController) Connect() revel.Result {
	if Session != nil {
		revel.INFO.Println("Already connected")

	} else {
		revel.INFO.Println("Connect to database")
		if c.Url == "" {
			c.GetConnectionUrl()
		}
		var err error
		Session, err = mgo.Dial(c.Url)
		//p.Session.SetMode(mgo.Strong, true)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
	}
	return nil
}

// Connect to mongodb
func (c *MongoDbController) Disconnect() revel.Result {
	if Session == nil {
		revel.INFO.Println("can't dissonnect. Not connected")

	} else {
		revel.INFO.Println("Disconnect from database")
		Session.Close()
	}
	return nil
}

// TODO: move constants to one file
func (p *MongoDbController) GetConnectionUrl() {
	p.Host, _ = p.GetConfig("mongodb.host", constants.MongoHost)
	p.Port, _ = p.GetConfig("mongodb.port", constants.MongoPort)
	p.User, _ = p.GetConfig("mongodb.user", "")
	p.Base, _ = p.GetConfig("mongodb.base", constants.MongoBase)
	if p.Password == "" {
		p.Password, _ = p.GetConfig("mongodb.pass", "")
	}
	if p.User != "" && p.Password != "" {
		p.Url = fmt.Sprintf("mongodb://%s:%s@%s:%s", p.User, p.Password, p.Host, p.Port)
	} else {
		p.Url = fmt.Sprintf("mongodb://%s:%s", p.Host, p.Port)
	}
}

func init() {
	revel.InterceptMethod((*MongoDbController).Connect, revel.BEFORE)
}

/*
func init() {
	revel.InterceptMethod((*MongoDbController).Disconnect, revel.AFTER)
}
*/
