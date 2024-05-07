# Changelog

**egriden** is not yet stable, hence the v0.x.x and will introduce breaking changes until v1.

## v0.2.0 - WIP
- Added **levels**
- Gotta make **free layers**
- Renamed `l.AddObject` to `l.AddGobject`
- Updated to Ebitengine v2.7.3
- Updated localization files

## v0.1.1 - 2024-04-29

- Removed `Gobject.DoesDrawScriptOverwriteSprite()` and introduced `Gobject.DrawSprite()`. Use the latter inside the `Gobject.OnDraw()` for combine sprite and custom drawing.
    - `Gobject.OnDraw()` returned now function needs a layer pointer to achieve this. This also allows for custom drawing to be interconnected with the layer fields.
- Removed the non-functional `OnCreate()` as its too much dependent on how you want to implement custom Gobjects and how they are created. The creator/build function should replace this.
- Updated to go v1.22.2 and Ebitengine v2.7.2

## v0.1.0

Initial version, introduces:
- Gobjects
- Grid Layers
- Sprite Packs and Image Sequences