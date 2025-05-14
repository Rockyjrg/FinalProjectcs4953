package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1920/2, 1080/2, "Rhythm Game")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	//disable default esc key behavior
	rl.SetExitKey(rl.KeyNull)

	InitNotes()
	InitScore()

	//set initial game state
	currentState = StateMainMenu

	//main game loop
	for !rl.WindowShouldClose() {
		//handle logic and drawing based on the current game state
		switch currentState {
		case StateMainMenu:
			UpdateMainMenu()
			DrawMainMenu()
		case StatePlaying:
			//check for pause
			if rl.IsKeyPressed(rl.KeyEscape) {
				currentState = StatePaused
			} else {
				UpdateGame()
				DrawGame()
			}
		case StatePaused:
			UpdatePauseMenu()
			DrawPauseMenu()
		}
	}
	//rl.CloseWindow()
}

func UpdateGame() {
	UpdateNotes()
	HandleInput()
}

func DrawGame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	DrawGameArea() //clear and reset notes
	DrawNotes()    //draw falling notes
	DrawScore()    //draw score, combo, multiplier

	rl.EndDrawing()
}

// InitGame sets up the initial state for a new game
func InitGame() {
	InitNotes()
	InitScore()
}
