package main

import (
	_ "embed"
	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed life.frag
var FRAGMENT_CODE string

func LoadRaylib() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "raylib stuff")
	rl.SetTargetFPS(60)

	rl.SetWindowPosition(0, 0)
	rl.ToggleFullscreen()
}

func LoadInitialImage() *rl.Image {
	// return rl.GenImagePerlinNoise(TEXTURE_WIDTH, TEXTURE_HEIGHT, 0, 0, 4)
	// return rl.GenImageWhiteNoise(TEXTURE_WIDTH, TEXTURE_HEIGHT, 0.5)
	return rl.GenImageColor(TEXTURE_WIDTH, TEXTURE_HEIGHT, rl.Black)
}

func LoadBuffers(image *rl.Image) (
	buffers []rl.RenderTexture2D,
	sourceBuffer int,
) {
	buffers = []rl.RenderTexture2D{
		rl.LoadRenderTexture(TEXTURE_WIDTH, TEXTURE_HEIGHT),
		rl.LoadRenderTexture(TEXTURE_WIDTH, TEXTURE_HEIGHT),
	}

	sourceBuffer = 0

	// load initial image into texture
	rl.UpdateTexture(
		buffers[sourceBuffer].Texture,
		rl.LoadImageColors(image),
	)

	return buffers, sourceBuffer
}

func LoadShader() rl.Shader {
	shader := rl.LoadShaderFromMemory("", FRAGMENT_CODE)
	textureResolutionLoc := rl.GetShaderLocation(shader, "textureResolution")
	rl.SetShaderValue(
		shader,
		textureResolutionLoc,
		[]float32{TEXTURE_WIDTH, TEXTURE_HEIGHT},
		rl.ShaderUniformVec2,
	)

	return shader
}
