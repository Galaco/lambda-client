package sky

// language=glsl
var Fragment = `
    #version 410

	in vec3 UV;

    out vec4 frag_colour;

	uniform samplerCube cubemapSample;

    void main() {
		// Output color = color of the texture at the specified UV
		frag_colour = texture( cubemapSample, UV );
    }
` + "\x00"
