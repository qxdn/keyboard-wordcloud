package mask

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/psykhi/wordclouds"
)

// embed default image
//go:embed mask.png
var f embed.FS

func loadDefaultPng() (image.Image, error) {
	mask, err := f.ReadFile("mask.png")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(mask)
	return png.Decode(reader)
}

func LoadDefaulMask(width int, height int, exclude color.RGBA) []*wordclouds.Box {
	res := make([]*wordclouds.Box, 0)
	img, err := loadDefaultPng()
	if err != nil {
		panic(err)
	}

	// scale
	imgw := img.Bounds().Dx()
	imgh := img.Bounds().Dy()

	wr := float64(width) / float64(imgw)
	wh := float64(height) / float64(imgh)
	scalingRatio := math.Min(wr, wh)
	// center
	xoffset := 0.0
	yoffset := 0.0
	if scalingRatio*float64(imgw) < float64(width) {
		xoffset = (float64(width) - scalingRatio*float64(imgw)) / 2
		res = append(res, &wordclouds.Box{
			Top:    float64(height),
			Left:   0.0,
			Right:  xoffset,
			Bottom: 0,
		})
		res = append(res, &wordclouds.Box{
			Top:    float64(height),
			Left:   float64(width) - xoffset,
			Right:  float64(width),
			Bottom: 0,
		})
	}

	if scalingRatio*float64(imgh) < float64(height) {
		yoffset = (float64(height) - scalingRatio*float64(imgh)) / 2
		res = append(res, &wordclouds.Box{
			Top:    yoffset,
			Left:   0.0,
			Right:  float64(width),
			Bottom: 0,
		})
		res = append(res, &wordclouds.Box{
			Top:    float64(height),
			Left:   0.0,
			Right:  float64(width),
			Bottom: float64(height) - yoffset,
		})
	}
	step := 3
	bounds := img.Bounds()
	for i := bounds.Min.X; i < bounds.Max.X; i = i + step {
		for j := bounds.Min.Y; j < bounds.Max.Y; j = j + step {
			r, g, b, a := img.At(i, j).RGBA()
			er, eg, eb, ea := exclude.RGBA()

			if r == er && g == eg && b == eb && a == ea {
				b := &wordclouds.Box{
					Top:    math.Min(float64(j+step)*scalingRatio+yoffset, float64(height)),
					Left:   float64(i)*scalingRatio + xoffset,
					Right:  math.Min(float64(i+step)*scalingRatio+xoffset, float64(width)),
					Bottom: float64(j)*scalingRatio + yoffset,
				}
				res = append(res, b)
			}
		}
	}

	return res
}

func LoadMask(path string, width int, height int, exclude color.RGBA) []*wordclouds.Box {
	if path == "" {
		return LoadDefaulMask(width, height, exclude)
	}
	return wordclouds.Mask(path, width, height, exclude)
}
