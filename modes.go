package main

import rl "github.com/gen2brain/raylib-go/raylib"

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
