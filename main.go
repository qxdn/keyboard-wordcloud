package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/qxdn/keyboard-wordcloud/modules/hook/keyboard"
	"github.com/qxdn/keyboard-wordcloud/modules/hook/types"
)

/*
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
*/

func main() {
	ch := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(ch); err != nil {
		fmt.Println(err)
		return
	}
	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")

	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("Received timeout signal")
			return
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return
		case k := <-ch:
			fmt.Printf("Received %v %v\n", k.Message, k.VKCode)
			continue
		}
	}

}
