package main

import (
	"image"
	"image/color"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fogleman/gg"
	"github.com/jasonlvhit/gocron"
	"github.com/psykhi/wordclouds"
	"github.com/qxdn/keyboard-wordcloud/modules/font"
	"github.com/qxdn/keyboard-wordcloud/modules/hook/keyboard"
	"github.com/qxdn/keyboard-wordcloud/modules/hook/types"
	"github.com/qxdn/keyboard-wordcloud/modules/mask"
	log "github.com/sirupsen/logrus"
)

type WordCloudConf struct {
	FontMaxSize     int
	FontMinSize     int
	RandomPlacement bool
	FontFile        string
	Colors          []color.RGBA
	BackgroundColor color.RGBA
	Width           int
	Height          int
	Mask            []*wordclouds.Box
	SizeFunction    string
	Debug           bool
}

var (
	WordCounts    = make(map[string]int, 255) // count keyboard counts
	DefaultColors = []color.RGBA{
		{0xa7, 0x1b, 0x1b, 0xff},
		{0x48, 0x48, 0x4B, 0xff},
		{0x59, 0x3a, 0xee, 0xff},
		{0x65, 0xCD, 0xFA, 0xff},
		{0x70, 0xD6, 0xBF, 0xff},
	}
	FontPath             = font.LoadFontPath("arial.ttf")
	Width                = 4096
	Height               = 4096
	DefaultWordCloudConf = WordCloudConf{
		FontMaxSize:     400,
		FontMinSize:     10,
		RandomPlacement: false,
		FontFile:        FontPath,
		Colors:          DefaultColors,
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		Width:           Width,
		Height:          Height,
		Mask: mask.LoadDefaulMask(Width, Height, color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}),
		Debug:        true,
		SizeFunction: "linear",
	}
)

func generateWordCloud(wordcount map[string]int, conf *WordCloudConf) (image.Image, error) {
	if conf == nil {
		conf = &DefaultWordCloudConf
	}
	colors := make([]color.Color, 0)
	for _, c := range conf.Colors {
		colors = append(colors, c)
	}
	oarr := []wordclouds.Option{wordclouds.FontFile(conf.FontFile),
		wordclouds.FontMaxSize(conf.FontMaxSize),
		wordclouds.FontMinSize(conf.FontMinSize),
		wordclouds.Colors(colors),
		wordclouds.MaskBoxes(conf.Mask),
		wordclouds.Height(conf.Height),
		wordclouds.Width(conf.Width),
		wordclouds.WordSizeFunction(conf.SizeFunction),
		wordclouds.RandomPlacement(conf.RandomPlacement),
		wordclouds.BackgroundColor(conf.BackgroundColor)}
	if conf.Debug {
		oarr = append(oarr, wordclouds.Debug())
	}
	w, err := wordclouds.NewWordcloud(wordcount, oarr...)
	if err != nil {
		return nil, err
	}
	return w.Draw(), nil
}

func init() {
	log.SetLevel(log.DebugLevel)
	/**
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/log.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	*/
}

func countKeyBoard(ch <-chan types.KeyboardEvent, lock *sync.RWMutex) {
	for {
		k := <-ch
		log.Debugf("Received %v %s", k.Message, k.VKCode.String())
		if k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN {
			// TODO: 限制上限
			lock.Lock()
			WordCounts[k.VKCode.String()[3:]]++
			lock.Unlock()
		}
	}
}

func generateAndSaveWordCloud(lock *sync.RWMutex) {
	yesterday := time.Now().AddDate(0, 0, -1)
	// fotmat to YYYY-MM-DD
	outName := yesterday.Format("./record/2006-01-02.png")
	log.Infof("generate date = %v wordcloud", yesterday.Format("2006-01-02"))
	lock.Lock()
	img, err := generateWordCloud(WordCounts, nil)
	for i := 0; i < 255; i++ {
		for k := range WordCounts {
			WordCounts[k] = 0
		}
	}
	lock.Unlock()
	if err != nil {
		log.Errorf("generate wordcloud fail: %v", err.Error())
		return
	}
	gg.SavePNG(outName, img)
}

func main() {
	// keyboard event chan
	ch := make(chan types.KeyboardEvent, 100)
	rwlock := &sync.RWMutex{}

	// install hook
	if err := keyboard.Install(ch); err != nil {
		log.Fatal(err)
		return
	}
	defer keyboard.Uninstall()

	go countKeyBoard(ch, rwlock)

	// generate file at every day 00:00am
	gocron.Every(1).Days().At("00:00").Do(generateAndSaveWordCloud, rwlock)

	gocron.Start()

	// os signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// recieve stop signal
	s := <-signalChan
	log.Debugf("received os signal: %v", s)

}
