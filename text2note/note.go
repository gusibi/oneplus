package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/image/font"

	"github.com/fogleman/gg"
)

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

type Text struct {
	Text     string
	FontSize int
	Padding  *Padding // 上，右，下，左
}

type background struct {
	ImageFile string
	ImageData string
	Width     int
	Height    int
	Padding   int
	FontSize  float64
}

var (
	NormalBackGroundHeader = &background{
		ImageFile: "./images/note_header.png",
		ImageData: "",
		Padding:   30,
		Width:     990,
		Height:    133,
	}
	NormalBackGroundBody = &background{
		ImageFile: "./images/note_body.png",
		ImageData: "",
		Padding:   80,
		Width:     990,
		Height:    100,
	}
	NormalBackGroundFooter = &background{
		ImageFile: "./images/note_footer_without_watermark.png",
		ImageData: "",
		Padding:   80,
		Width:     990,
		Height:    133,
		FontSize:  20,
	}
)

type NoteImageGenerator interface {
	CreatePngNote(ctx context.Context) (string, error)
	CreateJpegNote(ctx context.Context) (string, error)
}

type NoteImageGenerate struct {
	TargetPath  string
	Text        string
	WaterMark   string
	FontSize    float64
	LineSpacing float64
	fontFace    font.Face

	backgroundHeader *background
	backgroundBody   *background
	backgroundFooter *background
	backgroundHeight int
}

func NewNoteImageGenerator(text string, fontSize, lineSpacing float64, waterMark, targetPath string) *NoteImageGenerate {
	return &NoteImageGenerate{
		Text: text, FontSize: fontSize,
		WaterMark:        waterMark,
		LineSpacing:      lineSpacing,
		TargetPath:       targetPath,
		backgroundBody:   NormalBackGroundBody,
		backgroundFooter: NormalBackGroundFooter,
		backgroundHeader: NormalBackGroundHeader,
	}
}

func (ng *NoteImageGenerate) CreateNote(ctx context.Context, savePng bool) (targetFile string, err error) {
	lines, err := ng.loadText(ctx)
	if err != nil {
		return "", err
	}
	textHeight := ng.getTextNoteHeight(ctx, lines)
	dc, err := ng.loadBackground(ctx, textHeight)
	if err != nil {
		return "", err
	}
	dc, err = ng.drawText(ctx, dc, lines)
	if err != nil {
		return "", err
	}
	if ng.WaterMark != "" {
		dc, err = ng.drawWaterMark(ctx, dc)
	}
	if savePng {
		dc.SavePNG("note.png")
	}
	return "", nil
}

func (ng *NoteImageGenerate) getTextNoteHeight(ctx context.Context, lines []*Text) int {
	count := 0
	for _, line := range lines {
		count += 1
		if line.Padding == nil {
			continue
		}
		if line.Padding.Bottom > 0 {
			count += line.Padding.Bottom
		}
		if line.Padding.Top > 0 {
			count += line.Padding.Top
		}
	}
	height := int(float64(count*int(ng.FontSize)) * ng.LineSpacing)
	if height > ng.backgroundBody.Height {
		return height
	} else {
		return ng.backgroundBody.Height
	}
}

func (ng *NoteImageGenerate) CreatePngNote(ctx context.Context) (string, error) {
	return ng.CreateNote(ctx, true)
}

func (ng *NoteImageGenerate) loadBackground(ctx context.Context, textHeight int) (*gg.Context, error) {
	bth, err := gg.LoadImage(ng.backgroundHeader.ImageFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	btb, err := gg.LoadImage(ng.backgroundBody.ImageFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	btf, err := gg.LoadImage(ng.backgroundFooter.ImageFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	backgroundHeight := textHeight + ng.backgroundHeader.Height + ng.backgroundFooter.Height
	ng.backgroundHeight = backgroundHeight
	dc := gg.NewContext(ng.backgroundHeader.Width, backgroundHeight)
	dc.DrawRoundedRectangle(0, 0, float64(ng.backgroundHeader.Width), float64(backgroundHeight), 0)
	dc.Clip()
	// 写入 header
	dc.DrawImage(bth, 0, 0)
	// 写入 body
	for y := ng.backgroundHeader.Height; y < textHeight+ng.backgroundHeader.Height; y += ng.backgroundBody.Height {
		dc.DrawImage(btb, 0, y)
	}
	// 写入 footer
	dc.DrawImage(btf, 0, textHeight+ng.backgroundHeader.Height)
	//dc.SavePNG("temp_back.png")
	return dc, nil
}

func (ng *NoteImageGenerate) CreateJpegNote(ctx context.Context) (string, error) {
	return ng.CreateNote(ctx, false)
}

func (ng *NoteImageGenerate) loadFont(ctx context.Context) error {
	FontPath := "/Users/gs/github/oneplus/text2note/font/FangZhengFangSongJianTi.ttf"
	if fontFace, err := gg.LoadFontFace(FontPath, ng.FontSize); err != nil {
		return err
	} else {
		ng.fontFace = fontFace
	}
	return nil
}

func (ng *NoteImageGenerate) isNeedSep(a, b rune) bool {
	if (unicode.IsLetter(a) || unicode.IsDigit(a)) && (unicode.IsLetter(b) || unicode.IsDigit(b)) {
		return true
	}
	return false
}

func (ng *NoteImageGenerate) paragraphToLines(paragraph string) ([]*Text, error) {
	//length := len(paragraph)
	textWidth := ng.backgroundBody.Width - ng.backgroundBody.Padding*2
	var lines []*Text
	cWidth, startIdx, lastIdx := 0, 0, 0
	var lastChar rune
	for i, c := range paragraph {
		fmt.Printf("last char:%s, current char:%s \n", string(lastChar), string(c))
		if awidth, ok := ng.fontFace.GlyphAdvance(c); ok {
			width := int(float64(awidth) / 64)
			if cWidth+width > textWidth {
				if i == 0 {
					return nil, errors.New("font size too large")
				}
				lines = append(lines, &Text{
					Text: paragraph[startIdx:lastIdx],
				})
				//// 判断是否需要连字符，只有英文需要
				//if i-lastIdx == 1 && ng.isNeedSep(lastChar, c) {
				//	lines = append(lines, &Text{
				//		Text: fmt.Sprintf("%s-", paragraph[startIdx:lastIdx]),
				//	})
				//} else {
				//	lines = append(lines, &Text{
				//		Text: paragraph[startIdx:lastIdx],
				//	})
				//}
				startIdx = lastIdx
				cWidth = width
			} else {
				cWidth += width
			}
		} else {
			return nil, fmt.Errorf("Get char: '%s' width fail ", string(c))
		}
		lastIdx = i
		lastChar = c
	}
	lines = append(lines, &Text{
		Text: paragraph[startIdx:],
	})
	if len(lines) > 0 {
		lines[len(lines)-1].Padding = &Padding{0, 0, 1, 0}
	}
	return lines, nil
}

func (ng *NoteImageGenerate) loadText(ctx context.Context) (results []*Text, err error) {
	paragraphList := strings.Split(ng.Text, "\n")
	for _, paragraph := range paragraphList {
		lines, err := ng.paragraphToLines(paragraph)
		if err != nil {
			return nil, err
		}
		results = append(results, lines...)
	}
	if len(results) > 0 {
		results[len(results)-1].Padding = nil
	}
	return results, nil
}

func (ng *NoteImageGenerate) drawText(ctx context.Context, dc *gg.Context, textList []*Text) (*gg.Context, error) {
	x := ng.backgroundBody.Padding
	y := float64(ng.backgroundHeader.Height) + ng.FontSize
	//textWidth := ng.backgroundBody.Width - ng.backgroundBody.Padding*2
	FontPath := "/Users/gs/github/oneplus/text2note/font/FangZhengFangSongJianTi.ttf"
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(FontPath, ng.FontSize); err != nil {
		panic(err)
	}
	for _, text := range textList {
		dc.DrawString(text.Text, float64(x), y)
		if text.Padding != nil && text.Padding.Bottom > 0 {
			y = y + ng.FontSize*float64(text.Padding.Bottom)
		}
		y = y + ng.FontSize*ng.LineSpacing
	}
	return dc, nil
}

func (ng *NoteImageGenerate) drawWaterMark(ctx context.Context, dc *gg.Context) (*gg.Context, error) {
	FontPath := "/Users/gs/github/oneplus/text2note/font/FangZhengFangSongJianTi.ttf"
	dc.SetRGB(0, 90, 0)
	if err := dc.LoadFontFace(FontPath, ng.backgroundFooter.FontSize); err != nil {
		panic(err)
	}
	y := ng.backgroundHeight - int(ng.backgroundFooter.FontSize)
	dc.DrawString(ng.WaterMark, float64(ng.backgroundFooter.Padding), float64(y))
	return dc, nil
}

func Text2Note(text string, fontSize, lineSpacing float64, waterMark, targetPath string) string {
	ng := NewNoteImageGenerator(text, fontSize, lineSpacing, waterMark, targetPath)
	ng.loadFont(nil)
	ng.CreatePngNote(nil)
	return targetPath
}

func main() {
	text := "Hello, world! \n微信(weixin|wechat) Python SDK 支持开放平台和公众平台 支持微信小程序云开发"
	//text = "PrintRanges defines the set of printable characters according to Go. ASCII space, U+0020, is handled separately.PrintRanges defines the set of printable characters according to Go. ASCII space, U+0020, is handled separately.\n\n"
	Text2Note(text, 40.0, 1.75, "by 公号「四月」（hiiiapril）", "")
}
