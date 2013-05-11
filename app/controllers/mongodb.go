/*
Package mongodb implements simple mongodb connection wrapper 
*/
package controllers

import (
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"time"
)

const connection_duration = "3s"
const port = "27017"
const host = "127.0.0.1"
const base = "fond"

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
		duration, _ := time.ParseDuration(connection_duration)
		var err error
		Session, err = mgo.DialWithTimeout(c.Url, duration)
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

// Create connection url from conf or constants
// TODO: move constants to one file
func (p *MongoDbController) GetConnectionUrl() {
	p.Host, _ = p.GetConfig("mongodb.host", host)
	p.Port, _ = p.GetConfig("mongodb.port", port)
	p.User, _ = p.GetConfig("mongodb.user", "")
	p.Base, _ = p.GetConfig("mongodb.base", base)
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
