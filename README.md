[![Go Reference](https://pkg.go.dev/badge/github.com/greenthepear/egriden.svg)](https://pkg.go.dev/github.com/greenthepear/egriden)
[![Go Report Card](https://goreportcard.com/badge/github.com/greenthepear/egriden)](https://goreportcard.com/report/github.com/greenthepear/egriden)

**egriden** is a framework for the [Ebitengine](https://ebitengine.org/) game engine, perfect for creating simple grid-based puzzle or strategy games. 

Current features:
- A **grid layer** system
- **Levels** (aka scenes) and conventional **free layers**
- "**Gobjects**" with custom draw and update scripts
- Animatable **sprite** system with easy loading from PNG files

Instead of the common component approach seen in engines like Unity, it is more akin to GameMaker with how it handles object interactions.

It is an evolution of the messy code base created for [TacZ](https://github.com/greenthepear/TacZ) and is currently used by me to create a word-based action puzzle game. *[Click here to follow development.](https://greenthepear.com)*

***Currently unstable!*** Contributions of any kind are welcome. Check [the changelog](CHANGELOG.md) for updates.

# Install

This is meant to be used **with** (rather than instead of) Ebitengine so follow the [installation guide for Ebitengine](https://ebitengine.org/en/documents/install.html).

In your go project's directory:
```
$ go get github.com/greenthepear/egriden
```
or to get the latest development version of this repo:
```
$ go get github.com/greenthepear/egriden@main
```

# Quick boilerplate tutorial

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
    g.DrawAllGridLayers(screen) //Draw layers according to their Z order

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
    g.InitEgridenAssets()

    layer0 := g.CreateSimpleGridLayerOnTop(
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
    //structure for the Gobject interface. You can define update, draw
    //functions, sprite drawing options and other fun stuff.
    layer0.AddGobject(goPlayer.Build(), 1, 5)

    //Run the game
    if err = ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}
```

# License
BSD-3-Clause license

*Ebitengine by Hajime Hoshi is licensed under the Apache License 2.0*
