package main

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/psykhi/wordclouds"
	"github.com/qxdn/keyboard-wordcloud/modules/font"
	"github.com/qxdn/keyboard-wordcloud/modules/mask"
)

var (
	wcColor = []color.RGBA{
		{0xa7, 0x1b, 0x1b, 0xff},
		{0x48, 0x48, 0x4B, 0xff},
		{0x59, 0x3a, 0xee, 0xff},
		{0x65, 0xCD, 0xFA, 0xff},
		{0x70, 0xD6, 0xBF, 0xff},
	}
)

func main() {
	wordCounts := map[string]int{"important": 10, "notebook": 10, "mesh": 10, "abc": 10, "mass": 20}
	colors := make([]color.Color, 0)
	for _, c := range wcColor {
		colors = append(colors, c)
	}
	fontPath := font.LoadFontPath("arial.ttf")
	fmt.Println(fontPath)
	boxes := mask.LoadMask("", 2048, 2048, color.RGBA{0, 0, 0, 0})
	w := wordclouds.NewWordcloud(
		wordCounts,
		wordclouds.Colors(colors),
		wordclouds.FontMaxSize(200),
		wordclouds.FontFile(fontPath),

		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.MaskBoxes(boxes),
		wordclouds.Debug(),
	)
	gg.SavePNG("temp.png", w.Draw())
}
