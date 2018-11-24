// Star
// (c) Jani Nykänen

package main

// Star type
type star struct {
	x, y   int32
	vx, vy float32
	spr    sprite
	color  int32
}

// Animate
func (st *star) animate(tm float32) {

	animSpeed := float32(8.0)

	st.spr.animate(st.color*5+2, 0, 3, animSpeed, tm)
}

// Update
func (st *star) update(input *inputManager, s *stage, tm float32) {

	// Update solid
	s.updateSolid(int(st.x), int(st.y), 1)

	// Animate
	st.animate(tm)
}

// Draw
func (st *star) draw(bmp *bitmap, g *graphics) {

	// Draw sprite
	st.spr.draw(g, bmp, int32(st.x)*16, int32(st.y)*16, flipNone)
}

// Create a star
func createStar(x, y, color int32, s *stage) *star {

	st := new(star)

	// Set position
	st.x = x
	st.y = y
	// Update solid
	s.updateSolid(int(st.x), int(st.y), 1)

	// Create sprite
	st.spr = createSprite(16, 16)
	st.spr.row = color*5 + 2

	// Set color
	st.color = color

	return st
}
