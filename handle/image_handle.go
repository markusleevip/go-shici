package handle

import (
	"bufio"
	"fmt"
	"github.com/golang/freetype"
	"github.com/markusleevip/go-shici/bean"
	"github.com/markusleevip/go-shici/db"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const contentReStr =".*?[，。、！：？；]"

var (
	dpi      =  float64(72)
	fontFile = "../data/kaiti.TTF"
	qrcodeFile = "../data/qrcode.jpg"
	hinting  = "none"
	size     = float64(44)
	width	 = 443
	height   = 959
	spacing  = float64(1.5)
	leftSpace= 100
	topSpace = 50
	titleLineSize = 8
	wonb     = false
	pageSize = 11
	calcSize = float64(13)
	contentRes = regexp.MustCompile(contentReStr)
	imgFileName = "/data/poems/%s.png"
	imgFileNameQrCode= "/data/poems/qr/%s.jpg"
	imgFileNameSave= "%s.jpg"
)

func calcImage(poem db.Poem) (imgBean bean.ImageBean ,err error){
	imgBean.Height=height
	imgBean.Spacing = spacing
	imgBean.Size = size
	imgBean.LeftSpace= leftSpace
	imgBean.TopSpace = topSpace
	content:=contentRes.FindAllString(strings.Replace(poem.Content," ","",-1),-1)


	resultContent := make([]string,0)
	if utf8.RuneCountInString(poem.Title)>titleLineSize{
		resultTitle := SubStringTitle(poem.Title)
		resultContent = append(resultContent,resultTitle...)
	}else{
		resultContent = append(resultContent,poem.Title)
	}

	resultContent = append(resultContent,"")
	resultContent = append(resultContent,content...)
	resultContent = append(resultContent,poem.Author+" "+poem.Dynasty)
	imgBean.Lines=len(resultContent)
	for index:= range resultContent{
		if utf8.RuneCountInString(resultContent[index])>imgBean.MaxLen {
			imgBean.MaxLen = utf8.RuneCountInString(resultContent[index])
		}
	}
	if imgBean.MaxLen>7 {
		imgBean.LeftSpace=50
	}
	if imgBean.MaxLen>10 {
		imgBean.Size=32
	}
	if imgBean.Lines>pageSize{
		lensF:= float64(imgBean.Lines) / calcSize
		if imgBean.Size<36{
			lensF= float64(imgBean.Lines) / float64(pageSize)
		}
		imgBean.Height = int(lensF*float64(height))+80
	}
	if imgBean.Lines<9 {
		imgBean.TopSpace=200
	}

	imgBean.Content = resultContent
	fmt.Println("imgBean.MaxLen=",imgBean.MaxLen)
	fmt.Println("imgBean.Height=",imgBean.Height)
	fmt.Println("imgBean.Size=",imgBean.Size)
	fmt.Println("imgBean.LeftSpace=",imgBean.LeftSpace)
	fmt.Println("imgBean.TopSpace=",imgBean.TopSpace)
	fmt.Println("imgBean.Lines=",imgBean.Lines)
	fmt.Println("imgBean.Spacing=",imgBean.Spacing)
	fmt.Println("imgBean.Content=",imgBean.Content)
	return imgBean, nil
}

func CreateShiImage(poem db.Poem) {

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	imgBean,_ := calcImage(poem)
	resultContent:=imgBean.Content
	// Initialize the context.
	fg, bg := image.Black, image.NewUniform(color.RGBA{189, 153, 95, 0xff})
	if wonb {
		fg, bg = image.White, image.Black
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, imgBean.Height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(imgBean.Size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}


	// Draw the text.
	pt := freetype.Pt(imgBean.LeftSpace, imgBean.TopSpace+int(c.PointToFixed(imgBean.Size)>>6))
	for _, s := range resultContent {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(imgBean.Size * imgBean.Spacing)
	}

	// Save that RGBA image to disk.
	fileName :=poem.Author+"_"+strings.Replace(poem.Title,"/","",-1)
	outFile, err := os.Create(fmt.Sprintf(imgFileName,fileName))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	CreateQrCodeImg(poem,fmt.Sprintf(imgFileName,fileName),fmt.Sprintf(imgFileNameQrCode,fileName+"_qr"))
	fmt.Println("Wrote out.png OK.")
}

func CreateQrCodeImg(poem db.Poem,fileName string,qrFileName string){
	imgb,_ := os.Open(fileName)
	pngPic ,_ :=png.Decode(imgb)

	wmb, _ := os.Open(qrcodeFile)
	watermark, _ := jpeg.Decode(wmb)
	defer wmb.Close()

	offset := image.Pt(pngPic.Bounds().Dx()-watermark.Bounds().Dx()-10, pngPic.Bounds().Dy()-watermark.Bounds().Dy()-10)
	b := pngPic.Bounds()
	m := image.NewNRGBA(b)

	draw.Draw(m,b,pngPic,image.Point{},draw.Src)
	draw.Draw(m,watermark.Bounds().Add(offset),watermark,image.ZP,draw.Over)

	imgw, _ := os.Create(qrFileName)
	jpeg.Encode(imgw,m ,&jpeg.Options{100})
	defer imgw.Close()
}

func SubStringTitle(title string) (titles []string){
	titleSize :=utf8.RuneCountInString(title)
	num := titleSize/titleLineSize
	//fmt.Println(num)
	titles =make([]string,0)
	titleRune := []rune(title)
	for i:=0;i<num;i++{
		titles = append(titles,string(titleRune[i*titleLineSize:(i+1)*titleLineSize]))
	}
	if titleSize>titleLineSize*num{
		titles = append(titles,string(titleRune[num*titleLineSize:]))
	}
	return titles
}
