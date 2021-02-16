package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"github.com/markwinter/chip8-go/pkg/chipeight"
	"image/color"
	"log"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

var (
	file = flag.String("file", "", "Path to the Chip8 ROM to load")
)

type Game struct {
	c8          *chipeight.Chipeight
	pixel       *ebiten.Image
	frameBuffer [2048]uint8
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.c8.Step()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.c8.ShouldDraw() {
		g.frameBuffer = g.c8.GetScreen()
	}

	for row := 0; row < screenHeight; row++ {
		for col := 0; col < screenWidth; col++ {
			if g.frameBuffer[(row*screenWidth)+col] == 0 {
				continue
			}

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(col), float64(row))
			err := screen.DrawImage(g.pixel, opts)
			if err != nil {
				log.Printf("DrawImage error: %v", err)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (width, height int) {
	return screenWidth, screenHeight
}

func main() {
	flag.Parse()

	if *file == "" {
		log.Fatalf("Must provide path to Chip8 ROM")
	}

	c8 := chipeight.NewChipeight()

	loadErr := c8.LoadROM(*file)
	if loadErr != nil {
		log.Fatalf("failed to load file: %v", loadErr)
	}

	pixel, _ := ebiten.NewImage(1, 1, ebiten.FilterNearest)
	err := pixel.Fill(color.White)
	if err != nil {
		log.Fatalf("Failed to create pixel: %v", err)
	}

	game := &Game{}
	game.c8 = c8
	game.pixel = pixel

	ebiten.SetWindowSize(screenWidth*10*2, screenHeight*10*2)
	ebiten.SetWindowTitle("Chip8 Emulator")
	ebiten.SetMaxTPS(60)
	ebiten.SetVsyncEnabled(true)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
