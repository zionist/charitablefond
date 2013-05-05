package unittests

import (
	. "launchpad.net/gocheck"
	"testing"
	"github.com/zionist/charitablefond/app/plugins"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MongoDbSuite struct{}

var _ = Suite(&MongoDbSuite{})

func (s *MongoDbSuite) TestHelloWorld(c *C) {
	c.Check(42, Equals, 42)
}

/* func (t ApplicationTest) TestGetConfig() {
    p := mongodb.MongoDbPlugin{}
    p.TestMode = true
    //should be wriiten from conf file
    var found bool
    _, found = p.GetConfig("mongodb.host")
    t.Assert(found == true)
    _, found = p.GetConfig("mongodb.port")
    t.Assert(found == true)

    _, found = p.GetConfig("noconfig.noconfig")
    t.Assert(found == false)
}*/

