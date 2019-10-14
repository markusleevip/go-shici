package gofish

import "time"

const(
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:69.0) Gecko/20100101 Firefox/69.0"
	Qps = 50
)

var rateLimiter = time.Tick(
	time.Second/Qps)

type GoFish struct{
	Request *Request
}

func NewGoFish() *GoFish{
	return &GoFish{}
}

func (g *GoFish) Visit() error {
	return  g.Request.Do()
}