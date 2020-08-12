//go:generate go-bindata -o=asset.go -pkg=main ./Monaco_Linux.ttf
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"time"

	"flag"

	"github.com/golang/freetype"
	"github.com/rwcarlsen/goexif/exif"
)

var file string
var fontFile string
var fontBuffer []byte
var output string
var fontSize = 34

func init() {
	flag.StringVar(&file, "i", "", "image file path")
	flag.StringVar(&fontFile, "f", "", "font file path")
	flag.StringVar(&output, "w", "", "output file,if not set , watermark prefix will add")
	flag.Parse()
	// 必须指定图片
	if file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if output == "" {
		output = fmt.Sprintf("watermark_%s", file)
	}
	if fontFile == "" {
		bf, err := Asset("Monaco_Linux.ttf")
		if err != nil {
			log.Fatalf("去读默认字体 Monaco_Linux.ttf 失败,%s", err)
		}
		fontBuffer = bf
	} else {
		bf, err := ioutil.ReadFile(fontFile)
		if err != nil {
			log.Fatalf("读取字体文件 %s 失败", fontFile)
		}
		fontBuffer = bf
	}
}
func main() {
	fname := file

	fh, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	x, err := exif.Decode(fh)
	if err != nil {
		log.Fatal(err)
	}

	// 去读日期信息
	v, err := x.Get("DateTimeOriginal")
	if err != nil {
		log.Fatalf("去读日期信息失败了", err)
	}
	t, err := time.Parse(`"2006:01:02 15:04:05"`, v.String())
	if err != nil {
		log.Fatalf("无法解析日期 %s\n", v.String())
	}
	var txt = t.Format("2006-01-02 15:04:05")
	font, err := freetype.ParseFont(fontBuffer)
	if err != nil {
		log.Fatalf("解析字体文件 %s 文件失败，这可能不是一个合法的字体文件,%s\n", fontFile, err)
	}
	fh.Seek(0, 0)
	jpgimg, err := jpeg.Decode(fh)
	if err != nil {
		log.Fatalf("解析图片文件失败，请检查文件是否合法 %s\n", err)
	}

	img := image.NewNRGBA(jpgimg.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, jpgimg.At(x, y))
		}
	}
	f := freetype.NewContext()
	f.SetFont(font)
	f.SetDPI(400)
	f.SetFontSize(float64(fontSize))
	f.SetClip(jpgimg.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
	pt := freetype.Pt(img.Bounds().Dx()-2300, img.Bounds().Dy()-50)
	_, err = f.DrawString(txt, pt)
	// 保存到新的文件中
	newfile, _ := os.Create(output)
	defer newfile.Close()

	err = jpeg.Encode(newfile, img, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		log.Fatalf("保存水印文件%s失败%s", output, err)
		os.Exit(1)
	}
}
