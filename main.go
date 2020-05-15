package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli"
)

var widthMax float64
var heightMax float64
var seed int64
var destination string

var stepMax float64
var stepChange float64

func DrawAndSave() {

	rand.Seed(seed)

	// starting conditions
	height := (rand.Float64() * heightMax)
	slope := (rand.Float64()*stepMax)*2 - stepMax

	dest := image.NewRGBA(image.Rect(0, 0, int(widthMax), int(heightMax)))
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)

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
		if height < 0 {
			height = 0
			slope *= -1
		}
		gc.BeginPath()
		gc.MoveTo(i, heightMax)
		gc.LineTo(i, height)
		gc.Close()
		gc.FillStroke()
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
			Name:  "width",
			Value: 1280,
			Aliases: []string{"w"},
			Usage: "the `WIDTH` of the image",
		},
		&cli.Float64Flag{
			Name:  "height",
			Value: 720,
			Aliases: []string{"ht"},
			Usage: "the `HEIGHT` of the image",
		},
		&cli.Int64Flag{
			Name:  "seed",
			Value: time.Now().UTC().UnixNano(),
			Aliases: []string{"s"},
			Usage: "the given `SEED` for the mountains",
		},
		&cli.StringFlag{
			Name:  "destination",
			Value: "mountains.png",
			Aliases: []string{"d"},
			Usage: "the destination `PATH` of the generated image",
		},
		&cli.Float64Flag{
			Name:  "change",
			Value: 1.0,
			Aliases: []string{"c"},
			Usage: "the `AMOUNT` of change each step",
		},
		&cli.Float64Flag{
			Name:  "step",
			Value: 3.0,
			Aliases: []string{"st"},
			Usage: "the `HEIGHT` of the next step",
		},
	}
	app.Action = func(c *cli.Context) error {
		widthMax = c.Float64("width")
		heightMax = c.Float64("height")
		seed = c.Int64("seed")
		stepChange = c.Float64("change")
		stepMax = c.Float64("step")
		destination = c.String("destination")

		DrawAndSave()
		return nil
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
