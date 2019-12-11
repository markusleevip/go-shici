package db

import (
	"encoding/json"
	"fmt"
)

type Poem struct {
	Id 		int
	Title 	string
	Author 	string
	Dynasty string
	Content string
}

func (p *Poem) Insert() bool {
	stmtInsert, err := db.Prepare(" insert into poem (title,author,dynasty,content) values (?,?,?,?) ")
	if checkError(err) {
		return false
	}
	_ ,err = stmtInsert.Exec(&p.Title,&p.Author,&p.Dynasty,&p.Content)
	if checkError(err) {
		return false
	}
	return true

}
func (p *Poem) Save() {
	data , _ := json.Marshal(p)
	fmt.Println(string(data))
	p.Insert()

}
