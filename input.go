package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// HandleInput processes player input
func HandleInput() {
	//iterate through the keys assigned for gameplay
	for i, key := range LaneKeys {
		//if a gameplay key is pressed
		if rl.IsKeyPressed(key) {
			if keyPressSound.FrameCount > 0 {
				rl.PlaySound(keyPressSound)
			}
			keyPressedState[i] = true
			keyFlashTimer[i] = keyFlashDuration
			// Check if a hit occurred or if it was a bad press.
			CheckHit(key)
		}
	}
}

// CheckHit checks if a pressed key successfully hit an active note
// within the hit window and determines the accuracy.
func CheckHit(pressedKey int32) {
	//find the lane index for the pressed key.
	laneIndex := -1
	for i, key := range LaneKeys {
		if key == pressedKey {
			laneIndex = i
			break
		}
	}

	//if the pressed key doesn't correspond to a lane, do nothing.
	if laneIndex == -1 {
		return
	}

	//find closest hittable note in the pressed lane within the OVERALL HitTolerance
	closestNoteIndex := -1

	// Iterate through active notes in the pressed lane
	for i := range notes {
		// Only consider notes that are active and in the correct lane
		if notes[i].Active && notes[i].Key == pressedKey {
			distance := float32(math.Abs(float64(notes[i].Y - TargetY)))
			// Check if this note is within the overall hittable window
			if distance < HitTolerance {
				// Found a hittable note. Let's take the first one found in the slice.
				closestNoteIndex = i
				break // No need to check other notes in this lane for this key press
			}
		}
	}

	if closestNoteIndex != -1 {
		//hittable note was found and hit!
		note := &notes[closestNoteIndex] // Get the note using the found index

		note.Active = false // Mark the note as hit (inactive)

		// Calculate the exact distance for accuracy
		distance := float32(math.Abs(float64(note.Y - TargetY)))

		// Determine accuracy based on distance and add score
		accuracy := DetermineHitAccuracy(distance)
		AddScoreWithAccuracy(accuracy) // Call scoring function with accuracy

		// Show temporary feedback text
		ShowHitFeedback(accuracy, GetAccuracyColor(accuracy)) // Show hit feedback

		fmt.Printf("Hit note in lane %d! Distance: %.2f, Accuracy: %s\n", note.Lane, distance, accuracy)

	} else {
		// --- No hittable note was found in the pressed lane within the HitTolerance
		// This is a "bad press". Reset the combo.
		// Only reset combo if the game is actually playing, not in menus etc.
		if currentState == StatePlaying {
			ResetCombo()                                                  // This will now potentially play the combo break sound
			ShowHitFeedback("Bad Press!", GetAccuracyColor("Bad Press!")) // Show bad press feedback
			fmt.Println("Bad Press!")
		}
	}
}

// DetermineHitAccuracy calculates the accuracy string based on distance from target.
func DetermineHitAccuracy(distance float32) string {
	if distance < PerfectTolerance {
		return "Perfect!"
	} else if distance < GreatTolerance {
		return "Great!"
	} else if distance < GoodTolerance {
		return "Good!"
	} else if distance < OKTolerance {
		return "OK!"
	}
	return "" //should not happen if distance is < HitTolerance
}

// GetAccuracyColor returns a color for the given accuracy string.
func GetAccuracyColor(accuracy string) rl.Color {
	switch accuracy {
	case "Perfect!":
		return rl.Yellow
	case "Great!":
		return rl.Green
	case "Good!":
		return rl.Blue
	case "OK!":
		return rl.Purple
	case "Miss!": // Handled in UpdateNotes, but useful to have a color here
		return rl.Red
	case "Bad Press!":
		return rl.DarkGray
	default:
		return rl.White // Default color
	}
}
