**egriden** is a framework for the [Ebitengine](https://ebitengine.org/) game engine, perfect for creating simple grid-based puzzle or strategy games.

Currently far from stable or well-documented. Contributions welcome.

# Boilerplate tutorial

```go
package main

import (
    "github.com/greenthepear/egriden"
    "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    egriden.EgridenAssets //Add assets needed for Egriden
    
    //Anything else you want here.
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.DrawAllLayers(screen) //Draw layers according to their Z order

    // or do it your way with g.GridLayer(z).Draw()
}

func (g *Game) Update() error {
    // ... your game logic here
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 320, 320 //Define screen layout
}

func main(){
    //Initialize
    g := &Game{}
	g.InitEgridenComponents()

    layer0 := g.CreateGridLayerOnTop(
        "Background", //Name
        16, //Pixel size of a tile in the grid
        10, 10, //Width and height of the grid
        egriden.Sparce, //Draw mode
        0, 0) //Draw offset from the top left corner

    //Create an image sequence from all the PNGs in a folder
    //Image sequences are made of frames, controlled by the frame index
    seq, err := egriden.CreateImageSequenceFromFolder("idle", "./Graphics/player/idle/")
    
	if err != nil {
		log.Fatal(err)
	}

    //Create SpritePack with the sequence. A sprite pack can have multiple sequences,
    //which can be switched using their names (keys)
    playerSprites := egriden.NewSpritePackWithSequence(seq)

    //Create Gobject (short for grid object or go object or game object or whatever
    //you like) with the ImagePack
    goPlayer := egriden.NewBaseGobject("player", playerSprites)

    //Add to layer, Build() method needed for a baseGobject, otherwise create your own
    //structure for the Gobject interface. You can define create, update, draw
    //functions and other fun stuff.
    layer0.AddGobject(goPlayer.Build(), 1, 5)

    //Run the game
    if err = ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}
```
