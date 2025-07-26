package main

import (
	"github.com/GiffE/gokol/glue"
	"github.com/GiffE/gokol/samples/math"
	sapp "github.com/GiffE/gokol/sapp"
	"github.com/GiffE/gokol/sg"
)

type State struct {
	rx, ry float32
	Pip    sg.Pipeline
	Bind   sg.Bindings
}

var state State

var vs_source_glcore = `
#version 410

uniform vec4 vs_params[4];
layout(location = 0) in vec4 pos;
layout(location = 0) out vec4 color;
layout(location = 1) in vec4 color0;
layout(location = 1) out vec2 uv;
layout(location = 2) in vec2 texcoord0;

void main()
{
    gl_Position = mat4(vs_params[0], vs_params[1], vs_params[2], vs_params[3]) * pos;
    color = color0;
    uv = texcoord0 * 5.0;
}`

var fs_source_glcore = `
#version 410

uniform sampler2D tex_smp;

layout(location = 0) out vec4 frag_color;
layout(location = 1) in vec2 uv;
layout(location = 0) in vec4 color;

void main()
{
    frag_color = texture(tex_smp, uv) * color;
}`

type vertex struct {
	x, y, z float32
	color   uint32
	u, v    int16
}

func Init() {
	sg.Setup(&sg.Desc{
		Environment: glue.Environment(),
	})

	// Cube vertex buffer with packed vertex formats for color and texture coords.
	// Note that a vertex format which must be portable across all
	// backends must only use the normalized integer formats
	// (BYTE4N, UBYTE4N, SHORT2N, SHORT4N), which can be converted
	// to floating point formats in the vertex shader inputs.
	//
	// The reason is that D3D11 cannot convert from non-normalized
	// formats to floating point inputs (only to integer inputs),
	// and WebGL2 / GLES2 don't support integer vertex shader inputs.
	vertices := []vertex{
		/* pos                  color       uvs */
		{-1.0, -1.0, -1.0, 0xFF0000FF, 0, 0},
		{1.0, -1.0, -1.0, 0xFF0000FF, 32767, 0},
		{1.0, 1.0, -1.0, 0xFF0000FF, 32767, 32767},
		{-1.0, 1.0, -1.0, 0xFF0000FF, 0, 32767},

		{-1.0, -1.0, 1.0, 0xFF00FF00, 0, 0},
		{1.0, -1.0, 1.0, 0xFF00FF00, 32767, 0},
		{1.0, 1.0, 1.0, 0xFF00FF00, 32767, 32767},
		{-1.0, 1.0, 1.0, 0xFF00FF00, 0, 32767},

		{-1.0, -1.0, -1.0, 0xFFFF0000, 0, 0},
		{-1.0, 1.0, -1.0, 0xFFFF0000, 32767, 0},
		{-1.0, 1.0, 1.0, 0xFFFF0000, 32767, 32767},
		{-1.0, -1.0, 1.0, 0xFFFF0000, 0, 32767},

		{1.0, -1.0, -1.0, 0xFFFF007F, 0, 0},
		{1.0, 1.0, -1.0, 0xFFFF007F, 32767, 0},
		{1.0, 1.0, 1.0, 0xFFFF007F, 32767, 32767},
		{1.0, -1.0, 1.0, 0xFFFF007F, 0, 32767},

		{-1.0, -1.0, -1.0, 0xFFFF7F00, 0, 0},
		{-1.0, -1.0, 1.0, 0xFFFF7F00, 32767, 0},
		{1.0, -1.0, 1.0, 0xFFFF7F00, 32767, 32767},
		{1.0, -1.0, -1.0, 0xFFFF7F00, 0, 32767},

		{-1.0, 1.0, -1.0, 0xFF007FFF, 0, 0},
		{-1.0, 1.0, 1.0, 0xFF007FFF, 32767, 0},
		{1.0, 1.0, 1.0, 0xFF007FFF, 32767, 32767},
		{1.0, 1.0, -1.0, 0xFF007FFF, 0, 32767},
	}

	state.Bind.VertexBuffers[0] = sg.MakeBuffer(&sg.BufferDesc[vertex]{
		Data:  vertices,
		Label: "texcube-vertices",
	})
	indicies := []uint16{
		0, 1, 2, 0, 2, 3,
		6, 5, 4, 7, 6, 4,
		8, 9, 10, 8, 10, 11,
		14, 13, 12, 15, 14, 12,
		16, 17, 18, 16, 18, 19,
		22, 21, 20, 23, 22, 20,
	}
	state.Bind.IndexBuffer = sg.MakeBuffer(&sg.BufferDesc[uint16]{
		Usage: sg.Usage{IndexBuffer: true},
		Data:  indicies,
		Label: "cube-indicies",
	})

	// create a checkerboard texture
	pixels := [16]uint32{
		0xFFFFFFFF, 0xFF000000, 0xFFFFFFFF, 0xFF000000,
		0xFF000000, 0xFFFFFFFF, 0xFF000000, 0xFFFFFFFF,
		0xFFFFFFFF, 0xFF000000, 0xFFFFFFFF, 0xFF000000,
		0xFF000000, 0xFFFFFFFF, 0xFF000000, 0xFFFFFFFF,
	}

	state.Bind.Images[0] = sg.MakeImage(&sg.ImageDesc[uint32]{
		Width:  4,
		Height: 4,
		Data: sg.ImageData[uint32]{
			Subimage: [6][16][]uint32{
				{pixels[:]},
			},
		},
		Label: "texcube-texture",
	})

	// create a sampler object with default attributes
	state.Bind.Samplers[0] = sg.MakeSampler(&sg.SamplerDesc{
		Label: "texcube-sampler",
	})

	// create shader
	shd := sg.MakeShader(&sg.ShaderDesc{
		Label: "cube_shader",
		VertexFunc: sg.ShaderFunction{
			Source: vs_source_glcore,
			Entry:  "main",
		},
		FragmentFunc: sg.ShaderFunction{
			Source: fs_source_glcore,
			Entry:  "main",
		},
		Attrs: [16]sg.ShaderVertexAttr{
			{BaseType: sg.ShaderAttrBaseTypeFloat, GlslName: "pos"},
			{BaseType: sg.ShaderAttrBaseTypeFloat, GlslName: "color0"},
			{BaseType: sg.ShaderAttrBaseTypeFloat, GlslName: "texcoord0"},
		},
		UniformBlocks: [8]sg.ShaderUniformBlock{
			{Stage: sg.ShaderStageVertex,
				Layout: sg.UniformLayoutStd140,
				Size:   64,
				GlslUniforms: [16]sg.ShaderUniform{
					{Type: sg.UniformTypeFloat4,
						ArrayCount: 4,
						GlslName:   "vs_params"},
				},
			},
		},
		Images: [16]sg.ShaderImage{
			{Stage: sg.ShaderStageFragment,
				ImageType:    sg.ImageType2d,
				SampleType:   sg.ImageSampleTypeFloat,
				Multisampled: false,
			},
		},
		Samplers: [16]sg.ShaderSampler{
			{Stage: sg.ShaderStageFragment,
				SamplerType: sg.SamplerTypeFiltering},
		},
		ImageSamplerPairs: [16]sg.ShaderImageSamplerPair{
			{Stage: sg.ShaderStageFragment,
				ImageSlot:   0,
				SamplerSlot: 0,
				GlslName:    "tex_smp",
			},
		},
	})

	state.Pip = sg.MakePipeline(&sg.PipelineDesc{
		Layout: sg.VertexLayoutState{
			Attrs: [16]sg.VertexAttrState{
				{Format: sg.VertexFormatFloat3},
				{Format: sg.VertexFormatUbyte4N},
				{Format: sg.VertexFormatShort2N},
			},
		},
		Shader:    shd,
		IndexType: sg.IndexTypeUint16,
		CullMode:  sg.CullModeBack,
		Depth: sg.DepthState{
			WriteEnabled: true,
			Compare:      sg.CompareFuncLessEqual,
		},
		Label: "texcube-pipeline",
	})
}

func Frame() {
	w := sapp.Widthf()
	h := sapp.Heightf()
	t := float32(sapp.FrameDuration() * 60)

	proj := math.Perspective(60, w/h, 0.01, 10)
	view := math.LookAt(math.Vec3f{0, 1.5, 6}, math.Vec3f{}, math.Vec3f{0, 1, 0})
	view_proj := proj.Mul4(view)
	state.rx += 1
	state.ry += 2 * t
	rxm := math.Rotate(state.rx, math.Vec3f{1, 0, 0})
	rym := math.Rotate(state.ry, math.Vec3f{0, 1, 0})
	model := rxm.Mul4(rym)
	mvp := view_proj.Mul4(model)

	sg.BeginPass(&sg.Pass{
		Action: sg.PassAction{
			Colors: [sg.MaxColorAttachments]sg.ColorAttachmentAction{
				{
					LoadAction: sg.LoadActionClear,
					ClearValue: sg.Color{R: 0.25, G: 0.5, B: 0.75, A: 1.0},
				},
			},
		},
		Swapchain: glue.Swapchain(),
	})
	sg.ApplyPipeline(state.Pip)
	sg.ApplyBindings(&state.Bind)
	sg.ApplyUniform(0, &mvp)
	sg.Draw(0, 36, 1)
	sg.EndPass()
	sg.Commit()
}

func Cleanup() {
	sg.Shutdown()
}

func main() {
	sapp.Run(&sapp.AppDesc{
		Width:       800,
		Height:      600,
		Init:        Init,
		Cleanup:     Cleanup,
		Frame:       Frame,
		SampleCount: 4,
		WindowTitle: "Textured Cube (sokol-app)",
	})
}
