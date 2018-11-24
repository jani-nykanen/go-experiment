// Scene type
// (c) Jani Nykänen

package main

// Scene interface
type scene interface {
	init(g *graphics, trans *transition, evMan *eventManager, ass *assetPack) error
	update(input *inputManager, tm float32)
	draw(g *graphics)
	destroy()
	onChange()
}
