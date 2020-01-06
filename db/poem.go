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


func (p *Poem) HasRow() bool{
	sql := "select id from poem where title=? and author=? and dynasty=? "

	stmtOut, err := db.Prepare(sql)
	if checkError(err) {
		return false
	}
	rows, err := stmtOut.Query(p.Title,p.Author,p.Dynasty)
	defer rows.Close()
	if checkError(err) {
		return false
	}
	if rows.Next(){
		err = rows.Scan(&p.Id)
		if err != nil {
			return false
		}
		//p.UpdateLevel()
		return true
	}
	return false

}

func QueryPoemsHasWord(field string,value string)(poems []Poem,err error){
	sql :="select id,title,author,dynasty,content from poem where 1=1  "
	sql += fmt.Sprintf(	" and %s like ?", field)
	fmt.Println(sql,value)
	stmtOut, err := db.Prepare(sql)
	if checkError(err) {
		return nil, err
	}
	rows, err := stmtOut.Query(value)
	defer rows.Close()
	if checkError(err) {
		return nil, err
	}
	for rows.Next() {
		p :=Poem{}
		err = rows.Scan(&p.Id,&p.Title,&p.Author,&p.Dynasty,&p.Content)
		if err != nil {
			return nil, err
		}
		poems =append(poems,p)
	}
	return poems, nil
}

func (p *Poem) UpdateContent(newContent string) bool{
	if p.Id==0{
		return false
	}
	stmtUpdate, err := db.Prepare("update poem set content = ? where id = ?")
	if checkError(err) {
		return false
	}

	_, err = stmtUpdate.Exec(newContent,p.Id)
	if checkError(err) {
		return false
	}

	return true

}


func QueryPoemsByAuthor(author string) (poems []Poem,err error){

	return queryPoems("author",author)
}


func queryPoems(field string,value string)(poems []Poem,err error){
	sqlStr :="select id,title,author,dynasty,content from poem where 1=1  "
	sqlStr += fmt.Sprintf(" and %s = ?", field)
	fmt.Println(sqlStr,value)
	stmtOut, err := db.Prepare(sqlStr)
	if checkError(err) {
		return nil, err
	}
	rows, err := stmtOut.Query(value)
	defer rows.Close()
	if checkError(err) {
		return nil, err
	}
	for rows.Next() {
		p :=Poem{}
		err = rows.Scan(&p.Id,&p.Title,&p.Author,&p.Dynasty,&p.Content)
		if err != nil {
			return nil, err
		}
		poems =append(poems,p)
	}
	return poems, nil
}
