# Changelog

**egriden** is not yet stable, hence the v0.x.x and will introduce breaking changes until v1.

## v0.3.0 - 2025-08-03
- GridLayer cells can now be **rectangles** of any side lengths, not only squares.
- Added **padding** which creates gaps between cells.
    - Added `examples/non-square-grid` example to showcase the above changes.
- XY of Gobjects have been split into [`(*BaseGobject).GridPos`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseGobject.GridPos) and [`(*BaseGobject).ScreenPos`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseGobject.ScreenPos), the latter using floats, like ebitengine does for drawing anyway.
- Added a **new [`Cell`](https://pkg.go.dev/github.com/greenthepear/egriden#Cell) type** and associated methods to make GridLayer interactions clearer. Uses the [imggg](https://github.com/greenthepear/imggg) reimplementation of the standard `image` package to use data structures such as `Point` and `Rectangle`.
    - `ScreenXYtoGrid()` and `SnapScreenXYtoCellAnchor()` have been replaced with [`(*GridLayer).CellAtScreenPos`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.CellAtScreenPos) and [`Cell.Anchor`](https://pkg.go.dev/github.com/greenthepear/egriden#Cell.Anchor).
- Added **neighbor** utilities to easily get slices/sets of neighboring cells: [`(Cell).GetNeighbors`](https://pkg.go.dev/github.com/greenthepear/egriden#Cell.GetNeighbors), [`(Cell).GetNeighborsFunc`](https://pkg.go.dev/github.com/greenthepear/egriden#Cell.GetNeighborsFunc), [`(Cell).GetNeighborsSet`](https://pkg.go.dev/github.com/greenthepear/egriden#Cell.GetNeighborsSet), [`(Cell).GetNeighborsSetFunc`](https://pkg.go.dev/github.com/greenthepear/egriden#Cell.GetNeighborsSetFunc). Also contains predefined neighborhoods.
- Made multiple changes for initializing and managing SpritePacks.
    - Added [`CreateSpritePacksFromYaml`](https://pkg.go.dev/github.com/greenthepear/egriden#CreateSpritePacksFromYaml) which allows for **creating SpritePacks from yaml data and embedded files** (`embed.FS`).
    - Added [`CreateImageSequenceFromImages`](https://pkg.go.dev/github.com/greenthepear/egriden#CreateImageSequenceFromImages) to allow making image sequences from `image.Image`.
    - Deprecated [`CreateImageSequenceFromFolder`](https://pkg.go.dev/github.com/greenthepear/egriden#CreateImageSequenceFromFolder), replacing it with a more general [`CreateImageSequenceFromGlob`](https://pkg.go.dev/github.com/greenthepear/egriden#CreateImageSequenceFromGlob).
    - Made the ImageSequence fields (`Name` and `Frames`) public so you can create them however you want, if you don't like the above options.
    - Added [`(SpritePack).FrameAt`](https://pkg.go.dev/github.com/greenthepear/egriden#SpritePack.FrameAt) to quickly get a frame of a specific ImageSequence.
- [`(*BaseLevel).CreateGridLayerOnTop`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseLevel.CreateGridLayerOnTop) now takes in a new [`GridLayerParameters`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayerParameters) structure so the signature isn't so long and unreadable.
    - Deprecated [`(*BaseLevel).CreateSimpleGridLayerOnTop`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseLevel.CreateSimpleGridLayerOnTop) for this reason.
- Added [`(*GridLayer).SwapGobjectsAt`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.SwapGobjectsAt) and [`(*GridLayer).SwapGobjectsAtCells`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.SwapGobjectsAtCells) which lets you easily swap objects in grid layers.
- Replaced most instances of `image.Point` with `imggg.Point`.
- Added [`(*GridLayer).DebugDrawCheckerBoard`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.DebugDrawCheckerBoard) which draws a checkerboard pattern of the grid's cells for debugging purposes.
- Updated Ebitengine to v2.8.8.
- Changed to the Apache License 2.0 to be uniform with Ebitengine.
- Updated localization files.

## v0.2.1 - 2024-06-15
- Changed signature of [`(Gobject).OnDrawFunc`](https://pkg.go.dev/github.com/greenthepear/egriden#Gobject.OnDrawFunc) and [`(Gobject).OnUpdateFunc`](https://pkg.go.dev/github.com/greenthepear/egriden#Gobject.OnUpdateFunc) functions by adding a `self Gobject` parameter. This removes confusing pointer stuff that happens when you want to make a `BaseGobject` with custom functions like so:
    ```go
    o := NewBaseGobject(...)
    o.OnDrawFunc = func (i *ebiten.Image, l Layer) {
        o.NextFrame() //Won't do anything!
        //Just changes the frame of the original gobject 'template' which will
        //not be running in the game, since that's what o points to.
        o.DrawSprite(i, l)
    }
    l.AddGobject(o.Build(), 0, 0)
    ```
    - Also added a `l Layer` parameter to `(Gobject).OnUpdateFunc` for convenience. Getting the object's layer pointer into the `OnUpdate` function was tricky beforehand.
- Fixed the small issue that `(Gobject).OnUpdate` would never run ever.
- Added a new simple example [gopher-party](./examples/gopher-party/) to showcase many things added in v0.2.0 and this patch.
- Updated Ebitengine to v2.7.5.
- Updated localization files.

## v0.2.0 "Liberty" - 2024-06-05
- Added [**free layers**](https://pkg.go.dev/github.com/greenthepear/egriden#FreeLayer) which allow you to place and draw objects anywhere on the screen according to their XY coordinates.
- Added **levels**, an interface for them and methods for [`BaseLevel`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseLevel). [`(*EgridenAssets).InitEgridenAssets()`](https://pkg.go.dev/github.com/greenthepear/egriden#EgridenAssets.InitEgridenAssets) now creates a default level.
    - `(*EgridenAssets).InitEgridenAssets()` was renamed from `(*EgridenAssets).InitEgridenComponents()`.
- Added methods for deleting and moving Gobjects within a grid layer: [`(*GridLayer).DeleteAt()`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.DeleteAt), [`(*GridLayer).MoveGobjectTo()`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.MoveGobjectTo).
- Added [`(SpritePack).DrawOptions`](https://pkg.go.dev/github.com/greenthepear/egriden#SpritePack.DrawOptions) field and [`(Gobject).SetDrawOptions`](https://pkg.go.dev/github.com/greenthepear/egriden#Gobject.SetDrawOptions) method, allowing you to customize rendering of sprites using `ebiten.DrawImageOptions`. You can now also apply offsets for sprites with the new [`(Gobject).SetDrawOffsets()`](https://pkg.go.dev/github.com/greenthepear/egriden#Gobject.SetDrawOffsets) method.
- Added [`(GridLayer).IsXYwithinBounds()`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.IsXYwithinBounds), [`(GridLayer).IsScreenXYwithinBounds()`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.IsScreenXYwithinBounds), [`ScreenXYtoGrid()`](https://pkg.go.dev/github.com/greenthepear/egriden#ScreenXYtoGrid) and [`SnapScreenXYtoCellAnchor()`](https://pkg.go.dev/github.com/greenthepear/egriden#SnapScreenXYtoCellAnchor) to make interactions between the screen (cursor) and grid layers easier.
- Removed `baseGobjectWithoutScripts`, you can use the new [`(BaseGobject).OnDrawFunc`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseGobject.OnDrawFunc) and [`(BaseGobject).OnUpdateFunc`](https://pkg.go.dev/github.com/greenthepear/egriden#BaseGobject.OnUpdateFunc) fields to assign scripts, which are nil by default.
- Made layer selection methods return nil instead of panicking if z is out of bounds.
- Renamed `(*GridLayer).AddObject` to [`(*GridLayer).AddGobject`](https://pkg.go.dev/github.com/greenthepear/egriden#GridLayer.AddGobject).
- More tests and bug fixes.
- "Updated" to go v1.22.3 and Ebitengine v2.7.3.
- Updated localization files.

## v0.1.1 - 2024-04-29

- Removed `Gobject.DoesDrawScriptOverwriteSprite()` and introduced `Gobject.DrawSprite()`. Use the latter inside the `Gobject.OnDraw()` for combine sprite and custom drawing.
    - `Gobject.OnDraw()` returned function now needs a layer pointer to achieve this. This also allows for custom drawing to be interconnected with the layer fields.
- Removed the non-functional `OnCreate()` as its too much dependent on how you want to implement custom Gobjects and how they are created. The creator/build function should replace this.
- Updated to go v1.22.2 and Ebitengine v2.7.2.

## v0.1.0

Initial version, introduces:
- Gobjects
- Grid Layers
- Sprite Packs and Image Sequences
