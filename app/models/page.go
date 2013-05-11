package models

/*
import (
	"fmt"
	"github.com/zionist/charitablefond/app/controllers"
)
*/

type Page struct {
	PageId  int
	Content string
	Header  string
	Footer  string
}

/*
func (p *Page) Save() {
	d := mongodb.MongoDbPlugin{}
	fmt.Println(d)
}
*/
