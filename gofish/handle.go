package gofish

import "io"

type Handle interface {
	Worker(body io.Reader,url string)
}
