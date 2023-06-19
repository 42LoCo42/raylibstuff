#version 330

uniform sampler2D texture0;
uniform vec2 textureResolution;

in vec2 fragTexCoord;

out vec4 finalColor;

vec2 uv(vec2 pixel) { return pixel / textureResolution; }

int isOn(vec2 pos) { return int(texture(texture0, pos).r > 0.5); }

void main() {
	// clang-format off
	int sum =
		// moore neighbourhood
		isOn(fragTexCoord + uv(vec2(-1, -1))) +
		isOn(fragTexCoord + uv(vec2(-1,  0))) +
		isOn(fragTexCoord + uv(vec2(-1,  1))) +
		isOn(fragTexCoord + uv(vec2( 0, -1))) +
		isOn(fragTexCoord + uv(vec2( 0,  1))) +
		isOn(fragTexCoord + uv(vec2( 1, -1))) +
		isOn(fragTexCoord + uv(vec2( 1,  0))) +
		isOn(fragTexCoord + uv(vec2( 1,  1)));
	// clang-format on

	// CHANGE HERE
	bool lives =
		sum == 3 || (sum == 2 && isOn(fragTexCoord) == 1); // Game of Life
	// sum == 2 && isOn(fragTexCoord) == 0; // Seeds

	float on = lives ? 1.0 : 0.0;
	finalColor = vec4(on, on, on, 1);
}
