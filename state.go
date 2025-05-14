package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// constants to represent the different states of the game
const (
	StateMainMenu = iota //0
	StatePlaying         //1
	StatePaused          //2
)

var currentState int

// updateMainMenu handles input and logic for the main menu
func UpdateMainMenu() {
	//if enter is pressed, transistion to the playing state and initialize the game
	if rl.IsKeyPressed(rl.KeyEnter) {
		InitGame()
		currentState = StatePlaying
	} else if rl.IsKeyPressed(rl.KeyQ) {
		rl.CloseWindow()
	}
}

// renders main menu screen
func DrawMainMenu() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	//draw menu text elements
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())
	rl.DrawText("Bocchi the Rock! Rhythm Game!", screenWidth/2-rl.MeasureText("Bocchi the Rock! Rhythm Game!", 40)/2, screenHeight/4, 40, rl.White)
	rl.DrawText("Press ENTER to Start", screenWidth/2-rl.MeasureText("Press ENTER to Start", 20)/2, screenHeight/4+100, 20, rl.LightGray)
	rl.DrawText("Press Q to Quit", screenWidth/2-rl.MeasureText("Press Q to Quit", 20)/2, screenHeight/4+150, 20, rl.LightGray)

	rl.EndDrawing()
}

// handles input and logic for the pause menu
func UpdatePauseMenu() {
	if rl.IsKeyReleased(rl.KeyEscape) {
		currentState = StatePlaying
	} else if rl.IsKeyPressed(rl.KeyM) {
		//if M is pressed, reset game and go back to menu
		ResetGame()
		currentState = StateMainMenu
	} else if rl.IsKeyPressed(rl.KeyQ) {
		//close the app if Q is pressed
		rl.CloseWindow()
	}
}

// renders pause menu overlay
func DrawPauseMenu() {
	DrawGame() //show game in the background

	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())
	panelWidth := int32(400)
	panelHeight := int32(200)
	panelX := screenWidth/2 - panelWidth/2
	panelY := screenHeight/2 - panelHeight/2

	//draw semi transparent dark gray panel
	rl.DrawRectangle(panelX, panelY, panelWidth, panelHeight, rl.Fade(rl.DarkGray, 0.8))

	//draw pause menu text centered within the panel
	textX := panelX + panelWidth/2
	textY := panelY + 20
	rl.DrawText("Paused", textX-rl.MeasureText("Paused", 30)/2, textY, 30, rl.White)
	rl.DrawText("Press ESC to Resume", textX-rl.MeasureText("Press Esc to Resume", 20)/2, textY+50, 20, rl.LightGray)
	rl.DrawText("Press M for Main Menu", textX-rl.MeasureText("Press M for Main Menu", 20)/2, textY+90, 20, rl.LightGray)
	rl.DrawText("Press Q to Quit", textX-rl.MeasureText("Press Q to Quit", 20)/2, textY+130, 20, rl.LightGray)
}
