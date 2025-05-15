package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1920/2, 1080/2, "Rhythm Game")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	defer UnloadGameSounds() //function to unload sounds on exit

	defer UnloadKeyLights()

	rl.SetTargetFPS(60)

	//disable default esc key behavior
	rl.SetExitKey(rl.KeyNull)

	InitNotes() //includes music and sound loading
	InitScore()

	//set initial game state
	currentState = StateMainMenu

	//main game loop
	for !rl.WindowShouldClose() {
		UpdateGameMusic()

		//handle logic and drawing based on the current game state
		switch currentState {
		case StateMainMenu:
			UpdateMainMenu()
			DrawMainMenu()
		case StatePlaying:
			//check for pause
			if rl.IsKeyPressed(rl.KeyEscape) {
				currentState = StatePaused
				PauseGameMusic()
			} else {
				UpdateGame()
				DrawGame()
			}
		case StatePaused:
			UpdatePauseMenu() //handle resuming music
			DrawPauseMenu()
		}
	}
}

func UpdateGame() {
	UpdateNotes() //handles note movement and spawning (random for now)
	HandleInput() //handles hit accuracy, bad presses, and key press sounds
	UpdateScore() //handles feedback timer and maybe other score updates
	UpdateKeyLights()
}

func UpdateKeyLights() {
	for i := range keyFlashTimer {
		if keyFlashTimer[i] > 0 {
			keyFlashTimer[i] -= rl.GetFrameTime()
			if keyFlashTimer[i] <= 0 {
				keyPressedState[i] = false
			}
		}
	}
}

func DrawGame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	DrawGameArea()    //lanes and target circles
	DrawNotes()       //falling notes
	DrawScore()       //score, combo, multiplier
	DrawHitFeedback() //temporary text feedback for hits/misses/bad presses
	DrawKeyLights()

	rl.EndDrawing()
}

// InitGame sets up the initial state for a new game
func InitGame() {
	InitNotes() //re-initializes notes, resets random spawner, loads music and sounds
	InitScore() //resets score/combo
	PlayGameMusic()
}

func DrawKeyLights() {
	scale := 2.8
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	baseX := screenWidth - int(float32(keyLightTextures[0].Width)*float32(scale)*float32(len(keyLightTextures))) - 50
	y := screenHeight - int(float32(keyLightTextures[0].Height)*float32(scale)) - 50 // position near bottom

	spacing := int(float32(keyLightTextures[0].Width)*float32(scale)) + 10 // 10 px gap between each sprite

	for i, tex := range keyLightTextures {
		color := rl.Fade(rl.White, 0.4)
		if keyPressedState[i] {
			color = rl.White
		}

		x := baseX + spacing*int(i)

		rl.DrawTextureEx(tex, rl.Vector2{X: float32(x), Y: float32(y)}, 0, float32(scale), color)
	}
}
