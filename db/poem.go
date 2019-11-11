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

func (p *Poem) Save() {
	data , _ := json.Marshal(p)
	fmt.Println(string(data))
}
