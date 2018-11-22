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
	baseMap      *tilemap
	width        int32
	height       int32
	index        int
	bmpTestTiles *bitmap
	bmpFont      *bitmap
	bmpBorders   *bitmap
	bmpWall      *bitmap
	xpos         int32
	ypos         int32
	objects      []object
}

// Get difficulty string
func (s *stage) getDifficultyString() string {

	dif := s.baseMap.difficulty
	ret := ""

	// Full stars
	for i := 0; i < int(dif/2); i++ {
		ret += "#"
	}

	// Half stars
	if dif%2 == 1 {
		ret += "$"
	}

	return ret
}

// Update stage
func (s *stage) update(input *inputManager, tm float32) {

	// Update objects
	for i := 0; i < len(s.objects); i++ {

		s.objects[i].update(input, tm)
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
	g.setGlobalColor(0, 8, 120, 255)
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

// Draw map
func (s *stage) drawMap(g *graphics) {

	// Clear screen
	g.clearScreen(0, 72, 184)

	// Draw borders
	s.drawBorders(g)

	// Background
	g.setGlobalColor(0, 0, 0, 255)
	g.fillRect(s.xpos, s.ypos, s.width*16, s.height*16)

	// Draw tiles (temp)
	var tileID, sx, sy int32
	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			// Get tileID
			tileID = int32(s.baseMap.getTile(x, y))
			if tileID <= 0 {
				continue
			}

			tileID--
			// If wall
			if tileID == 0 {

				g.drawBitmapRegion(s.bmpWall, 0, 0, 16, 16, int32(s.xpos+x*16), int32(s.ypos+y*16), flipNone)

			} else {

				// Draw tile
				sx = tileID % 16
				sy = tileID / 16
				g.drawBitmapRegion(s.bmpTestTiles, sx*16, sy*16, 16, 16,
					int32(s.xpos+x*16), int32(s.ypos+y*16), flipNone)
			}
		}
	}

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
	g.drawText(s.bmpFont, s.getDifficultyString(),
		bottomXOff+int32(len(str)*10)+difMinusX, 240-bottomY, starXoff, 0, false)

	// Draw moves
	str = "Moves: " + strconv.Itoa(s.baseMap.moveLimit)
	g.drawText(s.bmpFont, str,
		256-int32(len(str)+1)*10+bottomXOff,
		240-bottomY, xoff, 0, false)
}

// Draw stage
func (s *stage) draw(g *graphics) {

	// Draw map
	s.drawMap(g)

	// Draw objects
	for i := 0; i < len(s.objects); i++ {

		s.objects[i].draw(g)
	}

	// Draw info
	s.drawInfo(g)
}

// Add an object
func (s *stage) addObject(o object) {
	s.objects = append(s.objects, o)
}

// Create a new stage
func createStage(index int, ass *assetPack) *stage {

	s := new(stage)

	// Load base map
	s.baseMap = ass.getTilemap(strconv.Itoa(index))
	// Get assets
	s.bmpTestTiles = ass.getBitmap("testTiles")
	s.bmpFont = ass.getBitmap("font")
	s.bmpBorders = ass.getBitmap("borders")
	s.bmpWall = ass.getBitmap("wall")
	// Get data
	s.width = int32(s.baseMap.width)
	s.height = int32(s.baseMap.height)

	// Calculate position
	s.xpos = 128 - s.width*16/2
	s.ypos = stageYOff + (240-stageYOff)/2 - s.height*16/2

	// Create an empty object list
	s.objects = make([]object, 0)

	s.index = index

	return s
}
