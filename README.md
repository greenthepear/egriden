**egriden** is a framework for the [Ebitengine](https://ebitengine.org/) game engine, perfect for creating simple grid-based puzzle or strategy games. Instead of the common component approach seen in engines like Unity, it is more akin to GameMaker with how it handles object interaction.

Current features:
- A **grid layer** system
- **Levels** (aka scenes) and conventional **free layers**
- "**Gobjects**" with custom draw and update scripts
- Animatable **sprite** system with easy loading from PNG files

It is an evolution of the messy code base created for [TacZ](https://github.com/greenthepear/TacZ) and is currently used by me to create a word-based action puzzle game. *[Click here to follow development.](https://greenthepear.com)*

***Currently unstable!*** Contributions of any kind are welcome. Check changelog.md for updates.

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

# License
MIT
