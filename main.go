package main

import (
	"github.com/ojrac/opensimplex-go"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"

	"image"
	"image/color"
	"math/rand"
	"math"
	"fmt"
	"time"
)

const WSIZE = 512
var world [WSIZE][WSIZE]float64
var m *image.RGBA

var scale float64 = 0.07

func run() {
	fmt.Println()
	cfg := pixelgl.WindowConfig{
		Title: "Island",
		Bounds: pixel.R(0, 0, 512, 512),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	win.Clear(colornames.Skyblue)

	p := pixel.PictureDataFromImage(m)
	pixel.NewSprite(p, p.Bounds()).
		Draw(win, pixel.IM.Moved(c).Scaled(c, 512/WSIZE))

	for !win.Closed() {
	win.Update()
	p := pixel.PictureDataFromImage(m)
	pixel.NewSprite(p, p.Bounds()).
		Draw(win, pixel.IM.Moved(c).Scaled(c, 512/WSIZE))
		if win.JustPressed(pixelgl.KeySpace) {
			gen(int64(rand.Float64()*1000000))
		}
	}
}

func gen(seed int64) {
	fmt.Println(seed)
	noise := opensimplex.NewNormalized(seed)
	m = image.NewRGBA(image.Rect(0, 0, WSIZE, WSIZE))
	for i := 0; i< WSIZE; i++ {
		for j := 0; j<WSIZE; j++ {
			//n := octSum(noise, 8, float64(i), float64(j), 0.6, 0.007)
			n := octSum(noise, 8, float64(i), float64(j), 0.55, 0.015)
		    n = subGrad(i, j, n)	
			world[i][j] = n
			m.Set(i, j, detColor(n))
			//m.Set(i, j, pixel.RGB(n, n, n))
		}
	}
}

func detColor(n float64) color.RGBA {
	n *= 255
	switch {
		case n > 255 * 0.65:
			return colornames.White
		case n > 255 * 0.57:
			return colornames.Lightgrey
		case n > 255 * 0.5:
			return colornames.Grey
		case n > 255 * 0.4:
			return colornames.Green
		case n > 255 * 0.3:
			return colornames.Forestgreen
		case n > 255 * 0.27:
			return colornames.Lightgoldenrodyellow
		default:
			return colornames.Skyblue
	}
}

func subGrad(x int, y int, n float64) float64 {
	tdX := 256.0 * math.Cos((-0.7 * float64(x)) * (math.Pi/180)) + 256
	tdY := 256.0 * math.Cos((-0.7 * float64(y)) * (math.Pi/180)) + 256
	sub := (tdX+tdY)/2 / WSIZE
	return n-sub
}

func octSum(noise opensimplex.Noise, octaves int, x float64, y float64, persistence float64, fr float64) float64 {
	maxAmp := 0.0
	amp := 1.0
	freq := fr
	res := 0.0

	for i:=0; i<octaves; i++ {
		res += noise.Eval2(x*freq, y*freq) * amp
		maxAmp += amp
		amp *= persistence
		freq *= 2
	}
	res /= maxAmp

	return res
}

func main() {
	rand.Seed(time.Now().UnixNano())
	m = image.NewRGBA(image.Rect(0, 0, WSIZE, WSIZE))
	gen(int64(rand.Float64()*1000000))
	pixelgl.Run(run)
}

