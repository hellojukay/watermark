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
	"strings"
	"sync"
	"time"

	_ "embed"
	"flag"

	"github.com/golang/freetype"
	"github.com/rwcarlsen/goexif/exif"
)

//go:embed Monaco_Linux.ttf
var bf []byte
var (
	file       string
	fontFile   string
	fontBuffer []byte
	output     string
	fontSize   = 34
	wg         sync.WaitGroup
)

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
	file = strings.TrimPrefix(file, "./")
	if output == "" {
		output = fmt.Sprintf("watermark_%s", strings.TrimPrefix(file, ".\\"))
	}
	if fontFile == "" {
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
	var begin = time.Now()
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
	log.Printf("[%10s] read file finish\n", time.Now().Sub(begin))
	// 去读日期信息
	v, err := x.Get("DateTimeOriginal")
	if err != nil {
		log.Fatalf("照片%s读日期信息失败了,%s", fname, err)
	}
	t, err := time.Parse(`"2006:01:02 15:04:05"`, v.String())
	if err != nil {
		log.Fatalf("照片%s无法解析日期 %s\n", fname, v.String())
	}
	log.Printf("[%10s] read time from image finish\n", time.Now().Sub(begin))

	var txt = t.Format("2006-01-02 15:04:05")
	font, err := freetype.ParseFont(fontBuffer)
	if err != nil {
		log.Fatalf("解析字体文件 %s 文件失败，这可能不是一个合法的字体文件,%s\n", fontFile, err)
	}
	fh.Seek(0, 0)
	log.Printf("[%10s] encode image begin\n", time.Now().Sub(begin))

	jpgimg, err := jpeg.Decode(fh)
	if err != nil {
		log.Fatalf("解析图片文件失败，请检查文件是否合法 %s\n", err)
	}
	log.Printf("[%10s] encode image finish\n", time.Now().Sub(begin))

	img := image.NewNRGBA(jpgimg.Bounds())
	var maxY = img.Bounds().Dy()
	var maxX = img.Bounds().Dx()
	log.Printf("[%10s] set color begin\n", time.Now().Sub(begin))

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				img.Set(x, y, jpgimg.At(x, y))
			}(x, y)
		}
	}
	wg.Wait()
	log.Printf("[%10s] draw text begin\n", time.Now().Sub(begin))

	f := freetype.NewContext()
	f.SetFont(font)
	f.SetDPI(400)
	f.SetFontSize(float64(fontSize))
	f.SetClip(jpgimg.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
	pt := freetype.Pt(img.Bounds().Dx()-2300, img.Bounds().Dy()-50)
	_, err = f.DrawString(txt, pt)
	if err != nil {
		log.Fatalf("写入水印失败 %s\n", err.Error())
	}
	log.Printf("[%10s] draw text finish\n", time.Now().Sub(begin))

	// 保存到新的文件中
	newfile, _ := os.Create(output)
	defer newfile.Close()

	err = jpeg.Encode(newfile, img, &jpeg.Options{
		Quality: 100,
	})
	log.Printf("[%10s] encode to file finished\n", time.Now().Sub(begin))

	if err != nil {
		log.Fatalf("保存水印文件%s失败%s", output, err)
		os.Exit(1)
	}
	log.Printf("[%10s] save to file finish\n", time.Now().Sub(begin))

}
