package handle

import (
	"github.com/markusleevip/go-shici/db"
	"testing"
)

func Test_CreateImg(t *testing.T) {
	poems, err := db.QueryPoemsByAuthor("王安石")
	if err != nil {
		return
	}
	for index:= range poems{
		CreateShiImage(poems[index])
	}
}
