# watermark
读照片EXIF信息，给相照片加上时间的水印

![Go](https://github.com/hellojukay/watermark/workflows/Go/badge.svg?branch=master)

我会经常打印一些生活的照片，那些照片都我的回忆，我希望能给照片加上时间的水印，这样在我看到照片的时候就知道是什么时候拍摄的。我的拍摄器材是尼康 D750 无法自动添加时间水印, 所以我写了这个程序。

# build
```shell
go get -u github.com/jteeuwen/go-bindata/...
go generate && go build
```

原图
![2020-07-19_071154.JPG](2020-07-19_071154.JPG)
```shell
./watermark -i 2020-07-19_071154.JPG
```
加上水印之后的图片
![watermark_2020-07-19_071154.JPG](watermark_2020-07-19_071154.JPG)
