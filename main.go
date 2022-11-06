package main

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"path"
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
	"gopkg.in/natefinch/lumberjack.v2"
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

type CheckPoint struct {
	Date       string         `json:"timestamp"`
	Checkpoint map[string]int `json:"checkpoint"`
}

var (
	ImageSavePath  string
	checkpointPath string
	WordCounts     = make(map[string]int, 255) // count keyboard counts
	DefaultColors  = []color.RGBA{
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
		FontMinSize:     100,
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
		Debug:        false,
		SizeFunction: "linear",
	}
)

func _generateWordClouds(wordcount map[string]int, conf *WordCloudConf) (image.Image, error) {
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

func isNewDay(t time.Time) bool {
	begin := time.Now().Truncate(24 * time.Hour)
	return t.Before(begin)
}

func init() {
	log.SetLevel(log.InfoLevel)
	// rolling log
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/log.log",
		MaxSize:    350, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	baseDir, _ := os.UserHomeDir()
	savePath := path.Join(baseDir, "Pictures", "wordclouds")
	// ignore exist
	os.MkdirAll(savePath, os.ModePerm)
	ImageSavePath = savePath

	// word dir
	wordDir, _ := os.Getwd()
	checkpointPath = path.Join(wordDir, "checkpoint.json")
	checkpoint := CheckPoint{}
	// load checkpoint
	if _, err := os.Stat(checkpointPath); err == nil {
		log.Info("checkpoint file exist")
		// checkpoint file exist
		content, err := ioutil.ReadFile(checkpointPath)
		if err != nil {
			log.Panic(err)
		}
		json.Unmarshal(content, &checkpoint)
		WordCounts = checkpoint.Checkpoint
		date, err := time.Parse(time.RFC3339, checkpoint.Date)
		if err != nil {
			log.Panic(err)
		}
		if isNewDay(date) {
			generateWordClouds(date)
		}
	}
}

func countKeyBoard(ch <-chan types.KeyboardEvent, lock *sync.RWMutex) {
	for {
		k := <-ch
		code := k.VKCode.String()[3:]
		log.Debugf("Received %v %s", k.Message, code)
		if k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN {
			// TODO: 限制上限
			if WordCounts[code] == math.MaxInt {
				log.Warnf("Key:%v count max", k.VKCode)
				continue
			}
			lock.Lock()
			WordCounts[code]++
			lock.Unlock()
		}
	}
}

func generateWordClouds(date time.Time) {
	defer saveCheckPoint()
	// fotmat to YYYY-MM-DD
	filename := date.Format("2006-01-02.png")
	filename = path.Join(ImageSavePath, filename)
	log.Infof("generate date = %v wordcloud", date.Format("2006-01-02"))
	img, err := _generateWordClouds(WordCounts, nil)
	for i := 0; i < 255; i++ {
		for k := range WordCounts {
			WordCounts[k] = 0
		}
	}
	if err != nil {
		log.Errorf("generate wordcloud fail: %v", err.Error())
		return
	}
	gg.SavePNG(filename, img)
}

func generateWordCloudsCron(lock *sync.RWMutex) {
	yesterday := time.Now().AddDate(0, 0, -1)
	lock.Lock()
	defer lock.Unlock()
	generateWordClouds(yesterday)
}

func saveCheckPoint() {
	checkpoint := CheckPoint{
		Date:       time.Now().Format(time.RFC3339),
		Checkpoint: WordCounts,
	}
	content, err := json.Marshal(checkpoint)
	if err != nil {
		log.Error(err)
	}
	ioutil.WriteFile(checkpointPath, content, os.ModePerm)
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
	log.Info("installed keyboard hook")
	defer keyboard.Uninstall()

	go countKeyBoard(ch, rwlock)
	go func(lock *sync.RWMutex) {
		for {
			lock.RLock()
			saveCheckPoint()
			lock.RUnlock()
			// save per 5 minute
			time.Sleep(60 * time.Second)
		}
	}(rwlock)
	// generate file at every day 00:00am
	gocron.Every(1).Days().At("00:00").Do(generateWordCloudsCron, rwlock)
	//gocron.Every(60).Seconds().Do(generateWordCloudsCron, rwlock)

	gocron.Start()

	// os signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// recieve stop signal
	s := <-signalChan
	log.Debugf("received os signal: %v", s)
}
