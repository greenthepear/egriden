# Changelog

**egriden** is not yet stable, hence the v0.x.x and will introduce breaking changes until v1.

## v0.1.1

- Removed `Gobject.DoesDrawScriptOverwriteSprite()` and introduced `Gobject.DrawSprite()`. Use the latter inside the `Gobject.OnDraw()` for combine sprite and custom drawing.
    - `Gobject.OnDraw()` returned now function needs a layer pointer to achieve this. This also allows for custom drawing to be interconnected with the layer fields.
- Removed the non-functional `OnCreate()` as its too much dependent on how you want to implement custom Gobjects and how they are created. The creator or build function should replace this.

## v0.1.0

Initial version, introduces:
- Gobjects
- Grid Layers
- Sprite Packs and Image Sequences