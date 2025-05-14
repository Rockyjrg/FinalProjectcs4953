package main

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Note struct {
	X, Y   float32 //position
	Speed  float32 //vertical speed
	Active bool    //is the note current active and falling?
	Key    int32   //the key associated with this note's lane
	Lane   int     //the index of the lane this note belongs to
}

// constants for game behavior
const (
	LaneWidth     float32 = 100 //width allocated per lane
	LanePadding   float32 = 50  //padding between lanes or from screen edge
	TargetY       float32 = 450 //y position where notes should be hit
	NoteRadius    float32 = 20  //radius of the note circle
	HitTolerance  float32 = 30  //distance from targetY to consider a hit
	NoteSpeed     float32 = 300 //default speed of falling notes
	NoteSpawnRate float32 = 0.5 //how often notes spawn
)

// array to hold lane-specific date
var (
	//LaneXPositions stores the X coord for the center of each lane
	LaneXPositions = []float32{
		LanePadding + NoteRadius,               //lane 0 (D)
		LanePadding + LaneWidth + NoteRadius,   //lane 1 (F)
		LanePadding + 2*LaneWidth + NoteRadius, //lane 2 (J)
		LanePadding + 3*LaneWidth + NoteRadius, //lane 3 (K)
	}
	//LaneKeys stores the rl key code associated with each lane
	LaneKeys = []int32{
		rl.KeyD,
		rl.KeyF,
		rl.KeyJ,
		rl.KeyK,
	}
)

// global slice to hold all active and inactive notes
var notes []Note

// noteSpawnTimer tracks time for spawning notes
var noteSpawnTimer float32

// InitNotes initialize or resets the notes slice and spawning timer
func InitNotes() {
	notes = []Note{}   //clear notes slice
	noteSpawnTimer = 0 //reset timer
}

// SpawnNote creates a new note in a lane
func SpawnNote(laneIndex int) {
	//ensure lane index is valid
	if laneIndex < 0 || laneIndex >= len(LaneXPositions) {
		return
	}
	notes = append(notes, Note{
		X:      LaneXPositions[laneIndex], //use x position for the given lane
		Y:      -NoteRadius * 2,           //start notes above the screen
		Speed:  NoteSpeed,
		Active: true,                //note is active when spawned
		Key:    LaneKeys[laneIndex], //associate the note the lane's key
		Lane:   laneIndex,           //store the lane index
	})
}

// UpdateNotes moves active notes down the screen and checks for misses
func UpdateNotes() {
	noteSpawnTimer += rl.GetFrameTime()
	//random spawning accross lanes for demo
	if noteSpawnTimer >= NoteSpawnRate {
		noteSpawnTimer = 0
		//spawn a note in a random lane for now
		randomLane := rand.IntN(len(LaneXPositions))
		SpawnNote(randomLane)
	}

	// Create a new slice to hold the notes that should be kept.
	updatedNotes := []Note{}
	screenHeight := float32(rl.GetScreenHeight())

	//iterate through notes to update their positions and check for misses
	for i := range notes {
		note := &notes[i]

		if note.Active {
			//move the note down
			note.Y += note.Speed * rl.GetFrameTime()

			//check if not passed target without being hit
			if note.Y > TargetY+HitTolerance {
				note.Active = false //deactivate note
				ResetCombo()        //break the current combo
				//TODO: Add a "Miss" sound and feedback
				fmt.Printf("Missed note in lane %d! (Y: %.2f)\n", note.Lane, note.Y)
			}
		}
		if note.Active || note.Y < screenHeight+NoteRadius { // Keep if active OR not yet off-screen
			updatedNotes = append(updatedNotes, *note) // Append the (potentially modified) note value
		}
	}
	notes = updatedNotes
}

// renders active notes on screen
func DrawNotes() {
	for _, note := range notes {
		if note.Active {
			rl.DrawCircle(int32(note.X), int32(note.Y), NoteRadius, rl.White)
		}
	}
}

// DrawGame draws the static parts of the game interface
func DrawGameArea() {
	//draw the target circles at the bottom of each lane
	for i := range LaneXPositions {
		//draw a hollow circle for the zone where a player should click the circle
		rl.DrawCircleLines(int32(LaneXPositions[i]), int32(TargetY), NoteRadius+5, rl.DarkGray)

	}

	//draw lines to seperate lanes
	screenHeight := float32(rl.GetScreenHeight())
	for i := range LaneXPositions {
		rl.DrawLine(int32(LaneXPositions[i]-LaneWidth/2), 0, int32(LaneXPositions[i]-LaneWidth/2), int32(screenHeight), rl.DarkGray)
	}
	rl.DrawLine(int32(LaneXPositions[len(LaneXPositions)-1]+LaneWidth/2), 0, int32(LaneXPositions[len(LaneXPositions)-1]+LaneWidth/2), int32(screenHeight), rl.DarkGray)
}
