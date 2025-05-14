package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var score int
var combo int
var multiplier int

func InitScore() {
	score = 0
	combo = 0
	multiplier = 1
}

// ResetCombo resets combo count and multiplier
// called when a note is missed
func ResetCombo() {
	combo = 0
	multiplier = 1
	//TODO: think about adding some kind of visual or sound cue
}

// AddScore increases score, updates the combo, and increments multiplier
func AddScore() {
	combo++

	//increase multi for every 10 successful hits in a row
	if combo > 0 && combo%10 == 0 {
		multiplier++
	}
	score += 100 * multiplier
}

func DrawScore() {
	rl.DrawText("Score: "+strconv.Itoa(score), 10, 10, 20, rl.White)
	rl.DrawText("Combo: "+strconv.Itoa(combo), 10, 40, 20, rl.White)
	rl.DrawText("x"+strconv.Itoa(multiplier), 10, 70, 20, rl.White)
}

func ResetGame() {
	InitNotes()
	InitScore()
}
