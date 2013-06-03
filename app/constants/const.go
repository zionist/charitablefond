package constants

import (
	"regexp"
)

const MongoConnection_duration = "3s"
const MongoPort = "27017"
const MongoHost = "127.0.0.1"
const MongoBase = "fond"

const PageCollectionName = "pages"
const UsersCollectionName = "users"
const FrontPage = "/page/index"

// used for correct icon display
var IconTypesRegex []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("^web"),
	regexp.MustCompile("^cloud"),
	regexp.MustCompile("^ip"),
	regexp.MustCompile("^it"),
}

const DefaultIcon = "cloud"
