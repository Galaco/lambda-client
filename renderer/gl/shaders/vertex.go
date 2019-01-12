package shaders

// language=glsl
var Vertex = `
    #version 410

	uniform mat4 projection;
	uniform mat4 view;
	uniform mat4 model;

    layout(location = 0) in vec3 vertexPosition;
	layout(location = 1) in vec2 vertexUV;
	layout(location = 2) in vec3 vertexNormal;
	layout(location = 3) in vec4 vertexTangent;
	layout(location = 4) in vec2 lightmapUV;

	// Output data ; will be interpolated for each fragment.
	out vec2 UV;
	out vec2 LightmapUV;

    void main() {
		gl_Position = projection * view * model * vec4(vertexPosition, 1.0);

    	// UV of the vertex. No special space for this one.
    	UV = vertexUV;
		LightmapUV = lightmapUV;
    }
` + "\x00"
