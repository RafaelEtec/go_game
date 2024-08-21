package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

type Vision struct {
	CanSee     bool
	BlinkEvery uint
}

type Player struct {
	*Sprite
	*Vision
	Velocity float64
	Health   uint
}

type Enemy struct {
	*Sprite
	Velocity      float64
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	HealAmount uint
}

type Game struct {
	player      *Player
	enemies     []*Enemy
	potions     []*Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func movePlayer(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X += g.player.Velocity
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X -= g.player.Velocity
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y -= g.player.Velocity
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y += g.player.Velocity
	}
}

func followPlayer(g *Game) {
	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += enemy.Velocity
			} else if enemy.X > g.player.X {
				enemy.X -= enemy.Velocity
			}
			if enemy.Y < g.player.Y {
				enemy.Y += enemy.Velocity
			} else if enemy.Y > g.player.Y {
				enemy.Y -= enemy.Velocity
			}
		}
	}
}

func updatePlayerVelocity(g *Game) {
	g.player.Velocity += 0.02
}

func handlePotion(g *Game) {
	for _, potion := range g.potions {
		if (g.player.X < potion.X+10 && g.player.X > potion.X-10) && (g.player.Y < potion.Y+10 && g.player.Y > potion.Y-10) {
			g.player.Health += potion.HealAmount
			updatePlayerVelocity(g)
			fmt.Printf("Picked up potion! Health: %d\n", g.player.Health)
		}
	}
}

func handleClick(g *Game) {

}

func showTiles(g *Game, opts ebiten.DrawImageOptions, screen *ebiten.Image) {
	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
			opts.GeoM.Reset()
		}
	}
}

func darken() {

}

func (g *Game) Update() error {
	movePlayer(g)
	followPlayer(g)
	handlePotion(g)
	handleClick(g)
	//darken(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{55, 38, 35, 255})
	opts := ebiten.DrawImageOptions{}

	//showTiles(g, opts, screen)
	showPlayer(g, opts, screen)
	showPotions(g, opts, screen)
	showEnemies(g, opts, screen)
}

func showPlayer(g *Game, opts ebiten.DrawImageOptions, screen *ebiten.Image) {
	opts.GeoM.Translate(g.player.X, g.player.Y)

	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()
}

func showEnemies(g *Game, opts ebiten.DrawImageOptions, screen *ebiten.Image) {
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
}

func showPotions(g *Game, opts ebiten.DrawImageOptions, screen *ebiten.Image) {
	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Give me a name!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	enemyImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/misc/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Vision: &Vision{
				CanSee:     true,
				BlinkEvery: 2,
			},
			Velocity: 2,
			Health:   2,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: enemyImg,
					X:   50.0,
					Y:   150.0,
				},
				0.5,
				true,
			},
			{
				&Sprite{
					Img: enemyImg,
					X:   150.0,
					Y:   50.0,
				},
				0.5,
				true,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					X:   210.0,
					Y:   100.0,
				},
				1.0,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
