package main

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

var (
	widthMax           float64
	heightMax          float64
	seed               int64
	destination        string
	startingColor      string
	endingColor        string
	backgroundColor    string
	gradientBackground string
	topToBottom        bool
	layers             int
	colorA             [3]float64
	colorB             [3]float64
	colorC             [3]float64
	colorD             [3]float64
)

var stepMax float64
var stepChange float64

func ColorizeBackground(img *image.RGBA) *image.RGBA {
	for x := 0.00; x < widthMax; x++ {
		for y := 0.00; y < heightMax; y++ {
			r, g, b := linearGradient(x, y, colorC, colorD)
			img.Set(int(x), int(y), color.RGBA{r, g, b, 255})
		}
	}
	return img
}

func Hex2RGB(hex string) [3]float64 {
	values, err := strconv.ParseUint(string(hex), 16, 32)
	if err != nil {
		return [3]float64{0, 0, 0}
	}

	return [3]float64{float64(values >> 16), float64((values >> 8) & 0xFF), float64(values & 0xFF)}
}

func linearGradient(x, y float64, color1, color2 [3]float64) (uint8, uint8, uint8) {
	d := x / widthMax
	if topToBottom {
		d = y / widthMax
	}
	r := color1[0] + d*(color2[0]-color1[0])
	g := color1[1] + d*(color2[1]-color1[1])
	b := color1[2] + d*(color2[2]-color1[2])
	return uint8(r), uint8(g), uint8(b)
}

func DrawAndSave() {
	rand.Seed(seed)

	colorA = Hex2RGB(startingColor)
	colorB = Hex2RGB(endingColor)
	colorC = Hex2RGB(backgroundColor)
	colorD = Hex2RGB(gradientBackground)

	dest := image.NewRGBA(image.Rect(0, 0, int(widthMax), int(heightMax)))
	dest = ColorizeBackground(dest)

	cache := make(map[float64]float64)
	for l := 1; l <= layers; l++ {
		// starting conditions
		height := (rand.Float64() * heightMax)
		slope := (rand.Float64()*stepMax)*2 - stepMax

		gc := draw2dimg.NewGraphicContext(dest)
		gc.SetLineWidth(5)
		if l > 1 {
			lighten := (rand.Float64()*25 - 25)
			colorA = [3]float64{colorA[0] - lighten, colorA[1] - lighten, colorA[2] - lighten}
			colorB = [3]float64{colorB[0] - lighten, colorB[1] - lighten, colorB[2] - lighten}
		}
		for i := 0.0; i < widthMax; i++ {
			height += slope
			slope += (rand.Float64()*stepChange)*2 - stepChange

			if slope > stepMax {
				slope = stepMax
			}
			if slope < -stepMax {
				slope = -stepMax
			}
			if height > heightMax {
				height = heightMax
				slope *= -1
			}

			if height >= cache[i] && l > 1 {
				height -= 25
				slope *= -1
			}

			if height < 0 {
				height = 0
				slope *= -1
			}
			if cache[i] <= height {
				cache[i] = height
			}

			r, g, b := linearGradient(i, height, colorA, colorB)
			gc.SetFillColor(color.RGBA{r, g, b, 255})
			gc.SetStrokeColor(color.RGBA{r, g, b, 255})
			gc.BeginPath()
			gc.MoveTo(i, heightMax)
			gc.LineTo(i, height)
			gc.Close()
			gc.FillStroke()
		}
	}
	err := draw2dimg.SaveToPngFile(destination, dest)
	if err != nil {
		panic(err)
	}
}

func main() {
	widthMax = 1280.0
	heightMax = 720.0
	app := cli.NewApp()
	app.Name = "mountains"
	app.Usage = "Commandline tool to procedurally generate a mountain range as png"
	app.UsageText = "mountains [options]"
	app.Flags = []cli.Flag{
		&cli.Float64Flag{
			Name:    "width",
			Value:   1280,
			Aliases: []string{"w"},
			Usage:   "the `WIDTH` of the image",
		},
		&cli.Float64Flag{
			Name:    "height",
			Value:   720,
			Aliases: []string{"ht"},
			Usage:   "the `HEIGHT` of the image",
		},
		&cli.Int64Flag{
			Name:    "seed",
			Value:   time.Now().UTC().UnixNano(),
			Aliases: []string{"s"},
			Usage:   "the given `SEED` for the mountains",
		},
		&cli.StringFlag{
			Name:    "destination",
			Value:   "mountains.png",
			Aliases: []string{"d"},
			Usage:   "the destination `PATH` of the generated image",
		},
		&cli.StringFlag{
			Name:    "starting-color",
			Value:   "000000",
			Aliases: []string{"sc"},
			Usage:   "the starting color for the gradient",
		},
		&cli.StringFlag{
			Name:    "ending-color",
			Value:   "000000",
			Aliases: []string{"ec"},
			Usage:   "the ending color for the gradient",
		},
		&cli.StringFlag{
			Name:    "background-color",
			Value:   "FFFFFF",
			Aliases: []string{"bg"},
			Usage:   "the background color of the image",
		},
		&cli.StringFlag{
			Name:    "gradient-color",
			Value:   "FFFFFF",
			Aliases: []string{"gc"},
			Usage:   "the gradient target for the background",
		},
		&cli.Float64Flag{
			Name:    "change",
			Value:   1.0,
			Aliases: []string{"c"},
			Usage:   "the `AMOUNT` of change each step",
		},
		&cli.Float64Flag{
			Name:    "step",
			Value:   3.0,
			Aliases: []string{"st"},
			Usage:   "the `HEIGHT` of the next step",
		},
		&cli.BoolFlag{
			Name:    "top-to-bottom",
			Aliases: []string{"ttb"},
			Usage:   "set the background gradient to 'left-to-right'",
		},
		&cli.IntFlag{
			Name:    "layers",
			Value:   1,
			Aliases: []string{"l"},
			Usage:   "set the number of layers",
		},
	}
	app.Action = func(c *cli.Context) error {
		widthMax = c.Float64("width")
		heightMax = c.Float64("height")
		seed = c.Int64("seed")
		stepChange = c.Float64("change")
		stepMax = c.Float64("step")
		destination = c.String("destination")
		startingColor = c.String("starting-color")
		endingColor = c.String("ending-color")
		backgroundColor = c.String("background-color")
		gradientBackground = c.String("gradient-color")
		topToBottom = c.Bool("top-to-bottom")
		layers = c.Int("layers")
		fmt.Println(c.Int("layers"), layers)

		DrawAndSave()
		return nil
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
