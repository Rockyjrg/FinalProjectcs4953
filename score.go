package main

import (
	"fmt"
	"strconv"

	// Needed for timing feedback display
	rl "github.com/gen2brain/raylib-go/raylib"
)

var score int
var combo int
var multiplier int

const ComboBreakSoundMultiplierThreshold = 2 // Play sound if multiplier is this or higher when combo breaks

// Score points per accuracy (remains the same)
var pointsPerAccuracy = map[string]int{
	"Perfect!": 500,
	"Great!":   300,
	"Good!":    100,
	"OK!":      50,
}

// Hit Feedback Display (remains the same)
var (
	feedbackText     string
	feedbackColor    rl.Color
	feedbackTimer    float32
	feedbackDuration float32 = 0.8
)

// InitScore initializes or resets the scoring variables and feedback display.
func InitScore() {
	score = 0
	combo = 0
	multiplier = 1
	feedbackText = ""
	feedbackTimer = 0
}

// ResetCombo resets combo count and multiplier.
func ResetCombo() {
	if comboBreakSound.FrameCount > 0 && multiplier >= ComboBreakSoundMultiplierThreshold {
		rl.PlaySound(comboBreakSound)
		fmt.Printf("Combo broken! Multiplier was %d. Playing combo break sound.\n", multiplier) // Debug
	} else if multiplier < ComboBreakSoundMultiplierThreshold && combo > 0 {
		// Debugging for combo breaks below the threshold
		fmt.Printf("Combo broken. Multiplier was %d (below threshold).\n", multiplier)
	}

	combo = 0
	multiplier = 1
}

// AddScoreWithAccuracy increases score based on the accuracy of the hit.
func AddScoreWithAccuracy(accuracy string) {
	points, ok := pointsPerAccuracy[accuracy]
	if !ok {
		points = 0
	}

	// Only increase combo and multiplier for successful hits (those with points > 0)
	if points > 0 {
		combo++

		if combo > 0 && combo%10 == 0 {
			// Store previous multiplier to check if it increased
			multiplier++
		}
		score += points * multiplier
	}
}

// UpdateScore handles any time-based score logic, like the feedback timer.
func UpdateScore() {
	if feedbackTimer > 0 {
		feedbackTimer -= rl.GetFrameTime()
		if feedbackTimer <= 0 {
			feedbackText = ""
		}
	}
}

// DrawScore renders the current score, combo, and multiplier on the screen.
func DrawScore() {
	rl.DrawText("Score: "+strconv.Itoa(score), 10, 10, 20, rl.White)
	rl.DrawText("Combo: "+strconv.Itoa(combo), 10, 40, 20, rl.White)
	rl.DrawText("x"+strconv.Itoa(multiplier), 10, 70, 20, rl.White)
}

// ShowHitFeedback sets the text and color for temporary hit/miss feedback.
func ShowHitFeedback(text string, color rl.Color) {
	feedbackText = text
	feedbackColor = color
	feedbackTimer = feedbackDuration
}

// DrawHitFeedback renders the temporary hit/miss feedback text.
func DrawHitFeedback() {
	if feedbackText != "" && feedbackTimer > 0 {
		screenWidth := int32(rl.GetScreenWidth())
		textY := int32(TargetY + NoteRadius + 20)
		textX := screenWidth/2 - rl.MeasureText(feedbackText, 30)/2
		rl.DrawText(feedbackText, textX, textY, 30, feedbackColor)
	}
}

// ResetGame calls the initialization functions for game components
// to reset the state for a new game.
func ResetGame() {
	StopGameMusic()
	InitNotes() // Reloads music and sounds
	InitScore()
}
