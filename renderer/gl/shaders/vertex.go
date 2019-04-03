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

	out vec2 UV;
	out vec3 EyeDirection;
	out vec3 LightDirection;
	
	// temporary
	uniform vec3 lightPos = vec3(0.0, 0.0, 100.0);

	void calculateEyePosition() {
		// View space vertex position
		vec4 P = view * model * vec4(vertexPosition, 1.0);
		// Normal vector
		vec3 N = normalize(mat3(view * model) * vertexNormal);
		// Tangent vector
		vec3 T = normalize(mat3(view * model) * vertexTangent.xyz);
		// Bitangent vector
		vec3 B = cross(N, T);
		// Vector from target to viewer
		vec3 V = -P.xyz;

		EyeDirection = normalize(vec3(dot(V, T), dot(V, B), dot(V, N)));

		vec3 L = lightPos - P.xyz;
		LightDirection = normalize(vec3(dot(L, T), dot(L, B), dot(L, N)));
	}

    void main() {
		gl_Position = projection * view * model * vec4(vertexPosition, 1.0);

    	UV = vertexUV;

		// bump + specular related
		calculateEyePosition();
    }
` + "\x00"
