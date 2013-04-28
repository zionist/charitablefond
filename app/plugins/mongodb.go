/*
Package mongodb implements simple mongodb connection wrapper 
*/
package mongodb

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"time"
	//"labix.org/v2/mgo/bson"
)

const connection_duration = "3s"
var required_configs = [...]string{"mongodb.host", "mongodb.port"}

type MongoDbPlugin struct {
	revel.EmptyPlugin
    TestMode bool
	Host    string
	Port    string
	Session *mgo.Session

}

//CheckConfig checks app.conf for mongodb connection credentials 
func (p MongoDbPlugin) GetConfig(config_string string) (value string, found bool) {
	if value, found = revel.Config.String(config_string); !found {
        if !p.TestMode {
		    revel.ERROR.Fatal("No %s in config", config_string)
        }
	}
	return
}

//func (p MongoDbPlugin) CheckConfigs() () {

func (p MongoDbPlugin) OnAppStart() {
	duration, _ := time.ParseDuration(connection_duration)
	var err error
	p.Session, err = mgo.DialWithTimeout(p.Host, duration)
	if err != nil {
		revel.ERROR.Fatal(err)
	}
}
