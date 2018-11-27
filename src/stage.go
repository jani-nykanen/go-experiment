// Stage
// (c) Jani Nykänen

package main

import (
	"strconv"
)

// Constants
const (
	stageYOff = 8
)

// Stage type
type stage struct {
	baseMap        *tilemap
	solidMap       []int
	width          int32
	height         int32
	index          int
	bmpFont        *bitmap
	bmpBorders     *bitmap
	bmpWall        *bitmap
	bmpGremlin     *bitmap
	sMove          *sample
	sTransform     *sample
	xpos           int32
	ypos           int32
	gremlins       []*gremlin
	stars          []*star
	anyMoving      bool
	oldMovingState bool
	moves          int
	gameRef        *game
}

// Add a gremlin
func (s *stage) addGremlin(x, y, color int32, sleeping bool) {

	s.gremlins = append(s.gremlins, createGremlin(x, y, color, s, sleeping))
}

// Add a star
func (s *stage) addStar(x, y, color int32) {

	s.stars = append(s.stars, createStar(x, y, color, s))
}

// Parse objects
func (s *stage) parseObjects() {

	// Go through tiles
	var tileID int32
	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			// Get tileID
			tileID = int32(s.baseMap.getTile(x, y))
			if tileID <= 0 {
				continue
			}

			// (We are not using switch here because
			// it does work as I would like to in Go)

			// If gremlin
			if tileID >= 2 && tileID <= 4 {
				s.addGremlin(x, y, tileID-2, false)

				// If star
			} else if tileID >= 5 && tileID <= 7 {
				s.addStar(x, y, tileID-5)

				// If sleeping
			} else if tileID >= 8 && tileID <= 10 {
				s.addGremlin(x, y, tileID-8, true)

				// If boulder
			} else if tileID == 11 {
				s.addGremlin(x, y, 3, false)

			}
		}
	}
}

// Update stage
func (s *stage) update(input *inputManager, audio *audioManager, tm float32) {

	// Check if something is moving
	// and check star collisions
	s.oldMovingState = s.anyMoving
	s.anyMoving = false
	gremlinCount := 0
	transforming := false
	for i := 0; i < len(s.gremlins); i++ {

		// Count existing gremlins
		if s.gremlins[i].exist && s.gremlins[i].color != 3 {
			gremlinCount++
		}

		// Check stars collisions before updating
		// the gremlin itself
		for i2 := 0; i2 < len(s.stars); i2++ {

			s.gremlins[i].getStarCollision(s.stars[i2], s)
		}

		// Check if active
		if !s.anyMoving && s.gremlins[i].isActive() {

			s.anyMoving = true
			// Check if starting to transform
			if s.gremlins[i].dying && !s.gremlins[i].transformPlayed {
				transforming = true
				s.gremlins[i].transformPlayed = true
			}
		}
	}

	// Play transformation sample
	if transforming {

		audio.playSample(s.sTransform, 0.30)
	}

	// If no gremlins left, victory
	if gremlinCount == 0 {

		s.gameRef.showInfoBox(true)
		return
	}

	// Update gremlins
	for i := 0; i < len(s.gremlins); i++ {

		s.gremlins[i].update(input, s, tm)
	}

	// If something moved, reduce moves
	if s.oldMovingState != s.anyMoving && s.anyMoving {

		// Play sample
		audio.playSample(s.sMove, 0.30)

		s.moves--
		// If negative moves, restart
		if s.moves < 0 {
			s.gameRef.showInfoBox(false)
			return
		}
	}

	// Update stars
	for i := 0; i < len(s.stars); i++ {

		s.stars[i].update(input, s, tm)
	}
}

// Draw borders
func (s *stage) drawBorders(g *graphics) {

	ypos := int32(s.ypos - 8)
	xpos := int32(s.xpos - 8)
	xjump := s.width*16 + 8
	yjump := s.height*16 + 8
	shadowJump := int32(12)

	// "Shadow"
	g.setGlobalColor(0, 72, 184, 255)
	g.fillRect(xpos+shadowJump, ypos+shadowJump, xjump+8+2, yjump+8+2)

	// Draw white outline
	g.setGlobalColor(255, 255, 255, 255)
	g.fillRect(xpos-1, ypos-1, xjump+8+2, yjump+8+2)

	// Horizontal
	for x := 0; x < int(s.width)*2; x++ {

		// Top
		g.drawBitmapRegion(s.bmpBorders, 8, 0, 8, 8,
			s.xpos+int32(x)*8, ypos, flipNone)
		// Bottom
		g.drawBitmapRegion(s.bmpBorders, 8, 16, 8, 8,
			s.xpos+int32(x)*8, ypos+yjump, flipNone)
	}

	// Vertical
	for y := 0; y < int(s.height)*2; y++ {

		// Left
		g.drawBitmapRegion(s.bmpBorders, 0, 8, 8, 8,
			xpos, s.ypos+int32(y)*8, flipNone)
		// Right
		g.drawBitmapRegion(s.bmpBorders, 16, 8, 8, 8,
			xpos+xjump, s.ypos+int32(y)*8, flipNone)
	}

	// Corners
	g.drawBitmapRegion(s.bmpBorders, 0, 0, 8, 8,
		xpos, ypos, flipNone)
	g.drawBitmapRegion(s.bmpBorders, 16, 0, 8, 8,
		xpos+xjump, ypos, flipNone)
	g.drawBitmapRegion(s.bmpBorders, 0, 16, 8, 8,
		xpos, ypos+yjump, flipNone)
	g.drawBitmapRegion(s.bmpBorders, 16, 16, 8, 8,
		xpos+xjump, ypos+yjump, flipNone)
}

// Draw walls
func (s *stage) drawWalls(g *graphics) {

	// Draw tiles (temp)
	var tileID int32
	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			// Get tileID
			tileID = int32(s.baseMap.getTile(x, y))
			if tileID <= 0 {
				continue
			}

			// If wall
			if tileID == 1 {

				g.drawBitmapRegion(s.bmpWall, 0, 0, 16, 16, (x * 16), (y * 16), flipNone)
			}
		}
	}
}

// Draw objects
func (s *stage) drawObjects(g *graphics) {

	// Draw gremlins
	for i := 0; i < len(s.gremlins); i++ {

		s.gremlins[i].draw(s.bmpGremlin, g)
	}

	// Draw stars
	for i := 0; i < len(s.stars); i++ {

		s.stars[i].draw(s.bmpGremlin, g)
	}
}

// Draw the background
func (s *stage) drawBackground(g *graphics) {

	// Clear screen
	g.clearScreen(30, 160, 248)

	// Draw borders
	s.drawBorders(g)

	// Background
	g.setGlobalColor(0, 0, 0, 255)
	g.fillRect(s.xpos, s.ypos, s.width*16, s.height*16)

}

// Draw info
func (s *stage) drawInfo(g *graphics) {

	stageIndexY := int32(8)
	nameY := int32(24)
	xoff := int32(-6)
	starXoff := int32(-3)
	bottomY := int32(16)
	bottomXOff := int32(2)
	difMinusX := int32(-4)

	// Draw stage index
	g.drawText(s.bmpFont, "Stage "+strconv.Itoa(s.index),
		128, stageIndexY, xoff, 0, true)
	// Draw stage name
	g.drawText(s.bmpFont, "\""+s.baseMap.name+"\"",
		128, nameY, xoff, 0, true)

	// Draw difficulty text
	str := "Difficulty: "
	g.drawText(s.bmpFont, "Difficulty: ",
		bottomXOff, 240-bottomY, xoff, 0, false)
	// Draw difficulty
	g.drawText(s.bmpFont, getDifficultyString(s.baseMap.difficulty),
		bottomXOff+int32(len(str)*10)+difMinusX, 240-bottomY, starXoff, 0, false)

	// Draw moves
	str = "Moves: " + strconv.Itoa(s.moves)
	g.drawText(s.bmpFont, str,
		256-int32(len(str)+1)*10+bottomXOff,
		240-bottomY, xoff, 0, false)
}

// Draw stage
func (s *stage) draw(g *graphics) {

	// Draw background
	s.drawBackground(g)

	// Draw gremlins
	g.translate(s.xpos, s.ypos)

	// Draw walls
	s.drawWalls(g)
	// Draw objects
	s.drawObjects(g)

	g.translate(0, 0)

	// Draw info
	s.drawInfo(g)
}

// Check solid data
func (s *stage) isTileSolid(x, y int) int {

	if x < 0 || y < 0 || x >= s.baseMap.width || y >= s.baseMap.height {
		return 1
	}

	return s.solidMap[y*s.baseMap.width+x]
}

// Update solid data
func (s *stage) updateSolid(x, y int, value int) {

	if x < 0 || y < 0 || x >= s.baseMap.width || y >= s.baseMap.height {
		return
	}

	s.solidMap[y*s.baseMap.width+x] = value
}

// Create a new stage
func createStage(index int, ass *assetPack, gameRef *game) *stage {

	s := new(stage)

	// Load base map
	s.baseMap = ass.getTilemap(strconv.Itoa(index))
	// Create solid map
	s.solidMap = make([]int, s.baseMap.width*s.baseMap.height)
	for i := 0; i < len(s.solidMap); i++ {

		// Check walls
		if s.baseMap.data[i] == 1 {
			s.solidMap[i] = 1
		} else {
			s.solidMap[i] = 0
		}
	}

	// Get assets
	s.bmpFont = ass.getBitmap("font")
	s.bmpBorders = ass.getBitmap("borders")
	s.bmpWall = ass.getBitmap("wall")
	s.bmpGremlin = ass.getBitmap("gremlin")
	s.sMove = ass.getSample("move")
	s.sTransform = ass.getSample("transform")

	// Get data
	s.width = int32(s.baseMap.width)
	s.height = int32(s.baseMap.height)
	s.moves = s.baseMap.moveLimit

	// Calculate position
	s.xpos = 128 - s.width*16/2
	s.ypos = stageYOff + (240-stageYOff)/2 - s.height*16/2

	// Set misc variables to default
	s.anyMoving = false
	s.oldMovingState = false
	s.index = index
	s.gameRef = gameRef

	// Create an empty object lists
	s.gremlins = make([]*gremlin, 0)

	// Parse objects
	s.parseObjects()

	return s
}
