# add text to image

```go
type text struct{
	value string
	fontSize int
	fontType font.Face
	padding  int[]
}

func TextToPng(text string, backgroud Image, textList []*text)(targetFile string, err error)

func TextToJpeg(text, backgroud string)(targetFile string, err error)
```

# 参考链接

https://learnku.com/articles/44827

