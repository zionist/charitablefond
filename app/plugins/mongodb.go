/*
Package mongodb implements simple mongodb connection wrapper 
*/
package mongodb

import (
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"time"
)

const connection_duration = "3s"
const port = "27017"
const host = "127.0.0.1"

var required_configs = [...]string{"mongodb.host", "mongodb.port"}

type MongoDbPlugin struct {
	revel.EmptyPlugin
	Host     string
	Port     string
	User     string
	Password string
	Url      string
	Session  *mgo.Session
}

//revel.Config.String wrapper
func (p *MongoDbPlugin) GetConfig(config_string string, defaul string) (value string, found bool) {
	value, found = revel.Config.String(config_string)
	if !found {
		value = defaul
	}
	return
}

// Connect to mongodb
func (p *MongoDbPlugin) Connect() {
	if p.Url == "" {
		p.GetConnectionUrl()
	}
	duration, _ := time.ParseDuration(connection_duration)
	var err error
	p.Session, err = mgo.DialWithTimeout(p.Url, duration)
	//p.Session.SetMode(mgo.Strong, true)
	if err != nil {
		revel.ERROR.Fatal(err)
	}
}

//Create connection url from conf or constants
// TODO: move constants to one file
func (p *MongoDbPlugin) GetConnectionUrl() {
	p.Host, _ = p.GetConfig("mongodb.host", host)
	p.Port, _ = p.GetConfig("mongodb.port", port)
	p.User, _ = p.GetConfig("mongodb.user", "")
	if p.Password == "" {
		p.Password, _ = p.GetConfig("mongodb.pass", "")
	}
	if p.User != "" && p.Password != "" {
		p.Url = fmt.Sprintf("mongodb://%s:%s@%s:%s", p.User, p.Password, p.Host, p.Port)
	} else {
		p.Url = fmt.Sprintf("mongodb://%s:%s", p.Host, p.Port)
	}
}

func (p *MongoDbPlugin) OnAppStart() {
	p.Connect()
}
