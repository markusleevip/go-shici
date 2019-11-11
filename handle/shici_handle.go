package handle

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/markusleevip/go-shici/gofish"
	"io"
)

var baseUrl ="https://so.gushiwen.org"
type AuthorHandle struct{}
func (h *AuthorHandle) Worker(body io.Reader,url string) {
	doc, err:=goquery.NewDocumentFromReader(body)
	if err !=nil{
		fmt.Println("doc.err.",err)
	}
	doc.Find(".sons").Find(".cont").Find("a").Each(func(i int, s *goquery.Selection) {
		author:= s.Text()
		fmt.Printf("%d author=%s\n",i,author)
		link,_ := s.Attr("href")
		fmt.Printf("%dlink=%s\n",i,link)

		h := PoemHomeHandle{}
		fish:=gofish.NewGoFish()
		request,err := gofish.NewRequest("GET",baseUrl+link,gofish.UserAgent,&h,nil)
		if err!=nil{
			fmt.Println(err)
			return
		}
		fish.Request = request
		fish.Visit()
	})
}

type PoemHomeHandle struct{}
func (h *PoemHomeHandle) Worker(body io.Reader,url string) {
	doc, err:=goquery.NewDocumentFromReader(body)
	if err !=nil{
		fmt.Println("doc.err.",err)
	}
	doc.Find(".sonspic").Find(".cont").Find("p").Find("a").Each(func(i int, s *goquery.Selection) {

		link,_ := s.Attr("href")
		fmt.Println("作品主页=",baseUrl+link)

		h := PoemInfoHandle{}
		fish:=gofish.NewGoFish()
		request,err := gofish.NewRequest("GET",baseUrl+link,gofish.UserAgent,&h,nil)
		if err!=nil{
			fmt.Println(err)
			return
		}
		fish.Request = request
		fish.Visit()
	})
}