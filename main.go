package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_WIDTH  = 1920
	WINDOW_HEIGHT = WINDOW_WIDTH * 9 / 16

	TEXTURE_SCALE  = 10
	TEXTURE_WIDTH  = WINDOW_WIDTH / TEXTURE_SCALE
	TEXTURE_HEIGHT = TEXTURE_WIDTH * 9 / 16
)

func main() {
	LoadRaylib()

	image := LoadInitialImage()
	buffers, sourceBuffer := LoadBuffers(image)
	shader := LoadShader()

	paused := true
	targetFPS := int32(60)

	drawCell := func(c color.RGBA) {
		TextureMode(buffers[sourceBuffer], func() {
			pos := rl.GetMousePosition()
			rl.DrawPixel(
				int32(pos.X)/TEXTURE_SCALE,
				// flip Y since texture is also flipped
				// also offset by one to compensate for inaccuracies
				TEXTURE_HEIGHT-int32(pos.Y)/TEXTURE_SCALE-1,
				c,
			)
		})
	}

	for !rl.WindowShouldClose() {
		// toggle pause
		if rl.IsKeyPressed(rl.KeySpace) {
			paused = !paused
			if paused {
				rl.SetTargetFPS(60)
			} else {
				rl.SetTargetFPS(targetFPS)
			}
		}

		// handle other inputs
		for char := rl.GetCharPressed(); char != 0; {
			log.Print(char)
			switch char {
			case '+':
				targetFPS *= 2
				rl.SetTargetFPS(targetFPS)

			case '-':
				// don't allow too low FPS
				if targetFPS <= 2 {
					break
				}

				targetFPS /= 2
				rl.SetTargetFPS(targetFPS)
			}

			break // TODO: GetCharPressed continuously returns the last char,
			// never 0, therefore we can't really loop here - figure this out!
		}

		// update target buffer if unpaused or stepped
		if !paused || rl.IsKeyPressed(rl.KeyEnter) {
			TextureMode(buffers[1-sourceBuffer], func() {
				ShaderMode(shader, func() {
					t := &buffers[sourceBuffer].Texture

					// draw texture flipped cuz raylib flips original data
					// when loading normal textures
					rl.DrawTextureRec(
						*t,
						rl.Rectangle{
							X:      0,
							Y:      0,
							Width:  float32(t.Width),
							Height: -float32(t.Height),
						},
						rl.Vector2{X: 0, Y: 0},
						rl.White,
					)
				})
			})

			sourceBuffer = 1 - sourceBuffer // flip buffers
		}

		// draw on source buffer
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			drawCell(rl.White)
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			drawCell(rl.Black)
		}

		// show source bufer
		DrawingMode(func() {
			rl.DrawTextureEx(
				buffers[sourceBuffer].Texture,
				rl.Vector2{X: 0, Y: 0},
				0,
				TEXTURE_SCALE,
				rl.White,
			)

			rl.DrawText(
				fmt.Sprintf("FPS: %d / %d", rl.GetFPS(), targetFPS),
				20,
				WINDOW_HEIGHT-50,
				50,
				rl.Red,
			)

			if paused {
				rl.DrawText(
					"Paused!",
					20,
					20,
					50,
					rl.Red,
				)
			}
		})
	}
}
