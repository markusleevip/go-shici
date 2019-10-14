package handle

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
)

type AuthorHandle struct{}
func (h *AuthorHandle) Worker(body io.Reader,url string) {
	doc, err:=goquery.NewDocumentFromReader(body)
	if err !=nil{
		fmt.Errorf("doc.err.",err)
	}
	doc.Find(".sons").Find(".cont").Find("a").Each(func(i int, s *goquery.Selection) {
		author:= s.Text()
		fmt.Printf("%d author=%s\n",i,author)
		link,_ := s.Attr("href")
		fmt.Printf("%dlink=%s\n",i,link)
	})
}
