package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WINDOW_WIDTH  = 1920
	WINDOW_HEIGHT = WINDOW_WIDTH * 9 / 16

	TEXTURE_SCALE  = 10
	TEXTURE_WIDTH  = WINDOW_WIDTH / TEXTURE_SCALE
	TEXTURE_HEIGHT = TEXTURE_WIDTH * 9 / 16
)

var (
	FPS int32 = 30

	//go:embed fragment.glsl
	FRAGMENT_CODE string
)

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "raylib stuff")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// rl.SetWindowPosition(0, 0)
	rl.ToggleFullscreen()

	// image := rl.GenImagePerlinNoise(TEXTURE_WIDTH, TEXTURE_HEIGHT, 0, 0, 4)
	// image := rl.GenImageWhiteNoise(TEXTURE_WIDTH, TEXTURE_HEIGHT, 0.5)
	image := rl.GenImageColor(TEXTURE_WIDTH, TEXTURE_HEIGHT, rl.Black)

	buffers := []rl.RenderTexture2D{
		rl.LoadRenderTexture(TEXTURE_WIDTH, TEXTURE_HEIGHT),
		rl.LoadRenderTexture(TEXTURE_WIDTH, TEXTURE_HEIGHT),
	}
	sourceBuffer := 0

	// load initial image into texture
	rl.UpdateTexture(
		buffers[sourceBuffer].Texture,
		rl.LoadImageColors(image),
	)

	// load shader with texture size input
	shader := rl.LoadShaderFromMemory("", FRAGMENT_CODE)
	textureResolutionLoc := rl.GetShaderLocation(shader, "textureResolution")
	rl.SetShaderValue(
		shader,
		textureResolutionLoc,
		[]float32{TEXTURE_WIDTH, TEXTURE_HEIGHT},
		rl.ShaderUniformVec2,
	)

	paused := true

	// periodic FPS log
	go func() {
		for {
			time.Sleep(time.Second)
			if !paused {
				log.Print("Current FPS: ", rl.GetFPS())
			}
		}
	}()

	drawCell := func(c color.RGBA) {
		TextureMode(buffers[sourceBuffer], func() {
			pos := rl.GetMousePosition()
			rl.DrawPixel(
				int32(pos.X)/TEXTURE_SCALE,
				TEXTURE_HEIGHT-int32(pos.Y)/TEXTURE_SCALE-1,
				c,
			)
		})
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			paused = !paused
			if paused {
				rl.SetTargetFPS(60)
			} else {
				rl.SetTargetFPS(FPS)
			}
		}

		for char := rl.GetCharPressed(); char != 0; {
			log.Print(char)
			switch char {
			case '+':
				FPS *= 2
				rl.SetTargetFPS(FPS)
			case '-':
				if FPS <= 2 {
					break
				}

				FPS /= 2
				rl.SetTargetFPS(FPS)
			}

			break
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
				fmt.Sprintf("FPS: %d", FPS),
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

func DrawingMode(f func()) {
	rl.BeginDrawing()
	defer rl.EndDrawing()
	f()
}

func ShaderMode(shader rl.Shader, f func()) {
	rl.BeginShaderMode(shader)
	defer rl.EndShaderMode()
	f()
}

func TextureMode(texture rl.RenderTexture2D, f func()) {
	rl.BeginTextureMode(texture)
	defer rl.EndTextureMode()
	f()
}
