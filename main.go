package main

import (
	"log"
	"os"
	"time"

	"flag"
	"github.com/rwcarlsen/goexif/exif"
)

var file string

func init() {
	flag.StringVar(&file, "img", "", "image file path")
	flag.Parse()
	// 必须指定图片
	if file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
func main() {
	fname := file

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	x, err := exif.Decode(f)
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
	println(txt)
}
