package shaders

//language=glsl
var Fragment = `
    #version 410

	uniform int useLightmap;

	uniform sampler2D albedoSampler;
	uniform sampler2D normalSampler;
	uniform sampler2D lightmapTextureSampler;


	in vec2 UV;
	in vec2 LightmapUV;

    out vec4 frag_colour;

	// Basetexture
	// Nothing is renderable without a base texture
	void AddAlbedo(inout vec4 fragColour, in sampler2D sampler, in vec2 uv) 
	{
		fragColour = texture(sampler, uv).rgba;
	}

//	vec3 CalculateNormal(in sampler2D sampler, vec2 uv)    // Calculate new normal based off Normal Texture and TBN matrix
//	{
//    	vec3 BumpMapNormal = texture(sampler, uv).xyz;
//    	BumpMapNormal = normalize(BumpMapNormal * 2.0 - 1.0);	// transform coordinates to -1-1 from 0-1
//    	vec3 NewNormal = normalize(TBN * BumpMapNormal);	// Tangent Space Conversion
//    	return NewNormal;
//	}

	// Lightmaps the face
	// Does nothing if lightmap was not defined
	void AddLightmap(inout vec4 fragColour, in sampler2D lightmap, in vec2 uv) 
	{
		fragColour = fragColour * texture(lightmap, uv).rgba;
	}

    void main() 
	{
		AddAlbedo(frag_colour, albedoSampler, UV);

//		bumpNormal = CalculateNormal(normalSampler, UV);

		//if (useLightmap == 1) {
		//	AddLightmap(frag_colour, lightmapTextureSampler, LightmapUV);
		//}
    }
` + "\x00"
