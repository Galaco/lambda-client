package shaders

//language=glsl
var Fragment = `
    #version 410

	uniform int useLightmap;

	uniform sampler2D albedoSampler;
	uniform sampler2D normalSampler;
	uniform sampler2D lightmapTextureSampler;

	in vec2 UV;
	in vec3 EyeDirection;
	in vec3 LightDirection;

    out vec4 frag_colour;

	vec4 GetAlbedo(in sampler2D sampler, in vec2 uv) 
	{
		return texture(sampler, uv).rgba;
	}

	float CalculateNormalFactor(in sampler2D sampler, vec2 uv)
	{
		vec3 L = normalize(LightDirection);
		vec3 N = normalize(texture(sampler, uv).xyz * 2.0 - 1.0);	// transform coordinates to -1-1 from 0-1

		return max(dot(N,L), 0.0);
	}

	vec4 GetSpecular(in sampler2D normalSampler, vec2 uv) {
		vec3 N = normalize(texture(normalSampler, uv).xyz * 2.0 - 1.0);	// transform coordinates to -1-1 from 0-1
		vec3 L = normalize(LightDirection);
		vec3 R = reflect(-L, N);
		vec3 V = normalize(EyeDirection);	

		// replace with texture where relevant
		vec3 specularSample = vec3(1.0);

		return max(pow(dot(R, V), 5.0), 0.0) * vec4(specularSample.xyz, 1.0);
	}

	// Lightmaps the face
	// Does nothing if lightmap was not defined
	vec4 GetLightmap(in sampler2D lightmap, in vec2 uv) 
	{
		return texture(lightmap, uv).rgba;
	}

    void main() 
	{
		float bumpFactor = CalculateNormalFactor(normalSampler, UV);
		vec4 diffuse = GetAlbedo(albedoSampler, UV);

		vec4 specular = GetSpecular(normalSampler, UV);

		frag_colour = diffuse + specular;
    }
` + "\x00"
