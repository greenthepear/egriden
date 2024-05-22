# Changelog

**egriden** is not yet stable, hence the v0.x.x and will introduce breaking changes until v1.

## v0.2.0 - WIP
- Added **levels**, an interface for them and methods for `BaseLevel`. `g.InitEgridenComponents()` now creates a default level.
- Added **free layers** which allow you to place and draw objects anywhere on the screen according to their XY coordinates.
- Added methods for deleting and moving Gobjects within a grid layer: `l.DeleteAt()`, `l.MoveGobjectTo()`.
- Added `DrawOptions` field and methods for SpritePacks, allowing you to customize rendering of sprites using `ebiten.DrawImageOptions`. You can now also apply offsets for sprites with the new `o.SetDrawOffets()` method.
- Added `l.IsXYwithinBounds()`, `l.IsScreenXYwithinBounds()`, `ScreenXYtoGrid()`, `SnapScreenXYtoCellAnchor()` to make interactions between the screen (cursor) and grid layers easier.
- Removed `baseGobjectWithoutScripts`, you can use the new `OnDrawFunc` and `OnUpdateFunc` fields to assign scripts, which are nil by default.
- Made layer selection methods return nil instead of panicing if z is out of bounds.
- Renamed `l.AddObject` to `l.AddGobject`.
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