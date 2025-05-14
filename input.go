package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// HandleInput processes player input
func HandleInput() {
	// if rl.IsKeyPressed(rl.KeyEscape) {
	// 	currentState = StatePaused
	// 	return
	// }

	//iterate through the keys assigned for gameplay
	for _, key := range LaneKeys {
		//if a gameplay key is pressed
		if rl.IsKeyPressed(key) {
			CheckHit(key)
		}
	}
}

// CheckHit checks if a pressed key successfully hit an active note
// within the hit window
func CheckHit(pressedKey int32) {
	//find the lane index for the pressed key
	laneIndex := -1
	for i, key := range LaneKeys {
		if key == pressedKey {
			laneIndex = i
			break
		}
	}

	//if the pressed key doesn't correspond to a lane, do nothing
	if laneIndex == -1 {
		return
	}

	//iterate through all active notes to find a potential hit in the correct lane
	for i := range notes {
		if notes[i].Active && notes[i].Key == pressedKey { //check to see if note is active and in correct lane
			//calculate vertical distance from note's center to the target Y pos
			distance := math.Abs(float64(notes[i].Y - TargetY))

			//check if note is within the acceptable hit tolerance
			if distance < float64(HitTolerance) {
				notes[i].Active = false //mark note as hit
				AddScore()
				return
			}
		}
	}
}
