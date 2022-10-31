package main

import (
	"github.com/mazznoer/colorgrad"
	"github.com/veandco/go-sdl2/sdl"
)

type rgb struct {
	red   uint8
	green uint8
	blue  uint8
}

const max_iterations = 1000

var color_table [max_iterations + 1]rgb
var largeur int32 = 1024
var hauteur int32 = 768
var window *sdl.Window
var renderer *sdl.Renderer

// Dessine la fractale
func draw(start_x, start_y, zoom float64) (pos_r, pos_i float64) {
	var x int32
	var y int32
	var i int
	var new_r, new_i, old_r, old_i float64

	for x = 0; x < int32(largeur); x++ {
		pos_r = float64((x-largeur/2.0))/(zoom*float64(largeur)) + start_x
		for y = 0; y < int32(hauteur); y++ {
			pos_i = float64(y-hauteur/2.0)/(zoom*float64(largeur)) + start_y
			new_i = 0
			new_r = 0
			old_i = 0
			old_r = 0
			i = 0

			for ((new_r*new_r + new_i*new_i) < 4.0) && (i < max_iterations) {
				old_r = new_r
				old_i = new_i
				new_r = old_r*old_r - old_i*old_i + pos_r
				new_i = 2.0*old_r*old_i + pos_i
				i = i + 1
			}

			r := color_table[i].red
			g := color_table[i].green
			b := color_table[i].blue
			renderer.SetDrawColor(r, g, b, 0)
			renderer.DrawPoint(x, y)

		}
	}
	return
}

func main() {
	var pcx, pcy float64

	// Initialise SDL
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, renderer, err = sdl.CreateWindowAndRenderer(largeur, hauteur, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	window.SetTitle("Fractale de Mandelbrot")

	err = renderer.SetDrawColor(0, 0, 0, 0)
	if err != nil {
		panic(err)
	}

	err = renderer.Clear()
	if err != nil {
		panic(err)
	}

	// Initialise le degrade
	grad, _ := colorgrad.NewGradient().
		HtmlColors("#FF4400", "#FFCC00", "#871E13", "#A34608", "#0066FF", "#FF4400").Domain(0, 10, 35, 55, max_iterations*0.85, max_iterations).Interpolation(colorgrad.InterpolationCatmullRom).Build()

	for dummy := 0; dummy < max_iterations; dummy++ {
		color_table[dummy].red, color_table[dummy].green, color_table[dummy].blue = grad.At(float64(dummy)).RGB255()

	}
	color_table[max_iterations].red = 0
	color_table[max_iterations].green = 0
	color_table[max_iterations].blue = 0

	// Position de depart
	start_x := -0.60
	start_y := 0.0
	var zoom float64 = 0.25

	pos_r, pos_i := draw(start_x, start_y, zoom)

	renderer.Present()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					// Si clic alors on ajuste le centrage et on zoome / dezoome
					pcx = (float64(t.X) - (float64(largeur) / 2)) / (float64(largeur) / 2)
					pcy = (float64(t.Y) - (float64(hauteur) / 2)) / (float64(hauteur) / 2)
					start_x = start_x + (((pos_r - start_x) / 2) * pcx)
					start_y = start_y + (((pos_i - start_y) / 2) * pcy)
					if t.Button == 1 {
						zoom = zoom * 2
					} else {
						zoom = zoom / 2
					}

					pos_r, pos_i = draw(start_x, start_y, zoom)
					renderer.Present()

				}
			}
		}

		sdl.Delay(10)
	}
}
