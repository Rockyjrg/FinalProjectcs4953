package main

import (
	"fmt"
	"math/rand/v2"
	"os"

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
	NoteSpeed     float32 = 400 //default speed of falling notes
	NoteSpawnRate float32 = 0.5 //how often notes spawn
)

const (
	PerfectTolerance float32 = 10
	GreatTolerance   float32 = 20
	GoodTolerance    float32 = 30
	OKTolerance      float32 = HitTolerance
)

// array to hold lane-specific date
var (
	//LaneXPositions stores the X coord for the center of each lane
	LaneXPositions = []float32{
		LanePadding + NoteRadius,
		LanePadding + LaneWidth + NoteRadius,
		LanePadding + 2*LaneWidth + NoteRadius,
		LanePadding + 3*LaneWidth + NoteRadius,
	}
	//LaneKeys stores the rl key code associated with each lane
	LaneKeys = []int32{
		rl.KeyD,
		rl.KeyF,
		rl.KeyJ,
		rl.KeyK,
	}
)

var (
	keyLightTextures [4]rl.Texture2D
	keyPressedState  [4]bool
	keyFlashTimer    [4]float32
	keyFlashDuration float32 = 0.1 // how long the light stays lit in seconds
)

// global slice to hold all active and inactive notes
var notes []Note

// noteSpawnTimer tracks time for spawning notes
var noteSpawnTimer float32

var gameMusic rl.Music

var (
	keyPressSound   rl.Sound // Sound played when a gameplay key is pressed
	comboBreakSound rl.Sound // Sound played when a high enough multiplier is broken
)

// InitNotes initializes or resets the notes slice, spawning timer, and loads assets (music, sounds).
func InitNotes() {
	notes = []Note{}   //clear notes slice
	noteSpawnTimer = 0 //reset timer

	// Load assets
	LoadGameMusic("assets/song.ogg")
	LoadGameSounds()

	keyLightTextures[0] = rl.LoadTexture("assets/Kita.png")
	keyLightTextures[1] = rl.LoadTexture("assets/Nijika.png")
	keyLightTextures[2] = rl.LoadTexture("assets/Bocchi.png")
	keyLightTextures[3] = rl.LoadTexture("assets/Ryo.png")
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

func UnloadKeyLights() {
	for _, tex := range keyLightTextures {
		rl.UnloadTexture(tex)
	}
}

// Load sound effects
func LoadGameSounds() {
	fmt.Println("Loading sound effects...")
	keyPressSound = rl.LoadSound("assets/333041__christopherderp__videogame-menu-button-clicking-sound-10.wav")
	comboBreakSound = rl.LoadSound("assets/416638__funwithsound__baby-gibberish-grunt-3.wav")

	if keyPressSound.FrameCount == 0 {
		fmt.Println("Warning: Failed to load key_press.wav")
	}
	if comboBreakSound.FrameCount == 0 {
		fmt.Println("Warning: Failed to load combo_break.wav")
	}

	//Set volume for sounds
	rl.SetSoundVolume(keyPressSound, 0.1)
	rl.SetSoundVolume(comboBreakSound, 1.5)
}

func UnloadGameSounds() {
	fmt.Println("Unloading sound effects...")
	// Check if sounds were loaded before unloading to avoid errors
	if keyPressSound.FrameCount > 0 {
		rl.UnloadSound(keyPressSound)
	}
	if comboBreakSound.FrameCount > 0 {
		rl.UnloadSound(comboBreakSound)
	}
}

// LoadGameMusic loads the background music file.
func LoadGameMusic(filePath string) {
	fmt.Printf("Loading music from %s...\n", filePath) // Example print
	// Stop and unload any previously loaded music before loading a new one
	StopGameMusic()

	// Check if file exists before trying to load
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Music file not found: %s. Make sure it exists in the 'assets' folder.\n", filePath)
		return
	}

	gameMusic = rl.LoadMusicStream(filePath)
	if gameMusic.FrameCount == 0 {
		fmt.Printf("Failed to load music: %s\n", filePath)
	} else {
		fmt.Printf("Music loaded successfully: %s\n", filePath)
		//set volume or pitch here if needed
		rl.SetMusicVolume(gameMusic, 0.07)
	}
}

// PlayGameMusic starts or resumes the loaded music.
func PlayGameMusic() {
	if gameMusic.FrameCount > 0 {
		rl.PlayMusicStream(gameMusic)
		fmt.Println("Music Started.")
	}
}

// PauseGameMusic pauses the loaded music.
func PauseGameMusic() {
	if gameMusic.FrameCount > 0 && rl.IsMusicStreamPlaying(gameMusic) {
		rl.PauseMusicStream(gameMusic)
		fmt.Println("Music Paused.")
	}
}

// ResumeGameMusic resumes the loaded music.
func ResumeGameMusic() {
	if gameMusic.FrameCount > 0 && !rl.IsMusicStreamPlaying(gameMusic) { // Only resume if not already playing
		rl.ResumeMusicStream(gameMusic)
		fmt.Println("Music Resumed.")
	}
}

// StopGameMusic stops and unloads the loaded music.
// Called when the game ends or returns to menu.
func StopGameMusic() {
	if gameMusic.FrameCount > 0 {
		if gameMusic.CtxData != nil {
			rl.StopMusicStream(gameMusic)
			rl.UnloadMusicStream(gameMusic)
			fmt.Println("Music Stopped and Unloaded.")
		}
	}
}

// UpdateGameMusic processes the music stream buffers.
// Must be called every frame in the main loop.
func UpdateGameMusic() {
	rl.UpdateMusicStream(gameMusic)
}

// UpdateNotes moves active notes down the screen and checks for misses
func UpdateNotes() {
	noteSpawnTimer += rl.GetFrameTime()
	//random spawning across lanes for demo
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

			//check if not passed target without being hit (a miss)
			// A note is missed if it passes the OK tolerance threshold.
			if note.Y > TargetY+OKTolerance { // Using OKTolerance as the miss threshold
				note.Active = false
				ResetCombo() // This will now potentially play a sound
				ShowHitFeedback("Miss!", GetAccuracyColor("Miss!"))
				fmt.Printf("Missed note in lane %d! (Y: %.2f)\n", note.Lane, note.Y)
			}
		}

		// Keep the note if it's still active
		// OR if it's inactive but not yet far off-screen (useful for hit/miss feedback)
		if note.Active || note.Y < screenHeight+NoteRadius*3 {
			updatedNotes = append(updatedNotes, *note)
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

// DrawGameArea draws the static parts of the game interface
func DrawGameArea() {
	for i := range LaneXPositions {
		rl.DrawCircleLines(int32(LaneXPositions[i]), int32(TargetY), NoteRadius+5, rl.DarkGray)
	}

	screenHeight := float32(rl.GetScreenHeight())
	for i := range LaneXPositions {
		rl.DrawLine(int32(LaneXPositions[i]-LaneWidth/2), 0, int32(LaneXPositions[i]-LaneWidth/2), int32(screenHeight), rl.DarkGray)
	}
	rl.DrawLine(int32(LaneXPositions[len(LaneXPositions)-1]+LaneWidth/2), 0, int32(LaneXPositions[len(LaneXPositions)-1]+LaneWidth/2), int32(screenHeight), rl.DarkGray)
}
