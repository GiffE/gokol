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
layout(location = 0) in vec4 position;
layout(location = 0) out vec4 color;
layout(location = 1) in vec4 color0;

void main()
{
    gl_Position = mat4(vs_params[0], vs_params[1], vs_params[2], vs_params[3]) * position;
    color = color0;
}`

var fs_source_glcore = `
#version 410

layout(location = 0) out vec4 frag_color;
layout(location = 0) in vec4 color;

void main()
{
    frag_color = color;
}`

func Init() {
	sg.Setup(&sg.Desc{
		Environment: glue.Environment(),
	})

	// cube vertex buffer
	vertices := []float32{
		-1.0, -1.0, -1.0, 1.0, 0.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0, 0.0, 1.0,

		-1.0, -1.0, 1.0, 0.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 0.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 0.0, 1.0, 0.0, 1.0,

		-1.0, -1.0, -1.0, 0.0, 0.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 0.0, 0.0, 1.0, 1.0,
		-1.0, 1.0, 1.0, 0.0, 0.0, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 0.0, 1.0, 1.0,

		1.0, -1.0, -1.0, 1.0, 0.5, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.5, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 0.5, 0.0, 1.0,
		1.0, -1.0, 1.0, 1.0, 0.5, 0.0, 1.0,

		-1.0, -1.0, -1.0, 0.0, 0.5, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 0.5, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.5, 1.0, 1.0,
		1.0, -1.0, -1.0, 0.0, 0.5, 1.0, 1.0,

		-1.0, 1.0, -1.0, 1.0, 0.0, 0.5, 1.0,
		-1.0, 1.0, 1.0, 1.0, 0.0, 0.5, 1.0,
		1.0, 1.0, 1.0, 1.0, 0.0, 0.5, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0, 0.5, 1.0,
	}
	vbuf := sg.MakeBuffer(&sg.BufferDesc[float32]{
		Data:  vertices,
		Label: "cube-verticies",
	})
	indicies := []uint16{
		0, 1, 2, 0, 2, 3,
		6, 5, 4, 7, 6, 4,
		8, 9, 10, 8, 10, 11,
		14, 13, 12, 15, 14, 12,
		16, 17, 18, 16, 18, 19,
		22, 21, 20, 23, 22, 20,
	}
	ibuf := sg.MakeBuffer(&sg.BufferDesc[uint16]{
		Usage: sg.Usage{IndexBuffer: true},
		Data:  indicies,
		Label: "cube-indicies",
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
			{BaseType: sg.ShaderAttrBaseTypeFloat, GlslName: "position"},
			{BaseType: sg.ShaderAttrBaseTypeFloat, GlslName: "color0"},
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
	})

	state.Pip = sg.MakePipeline(&sg.PipelineDesc{
		Layout: sg.VertexLayoutState{
			Buffers: [8]sg.VertexBufferLayoutState{
				{Stride: 28},
			},
			Attrs: [16]sg.VertexAttrState{
				{Format: sg.VertexFormatFloat3},
				{Format: sg.VertexFormatFloat4},
			},
		},
		Shader:    shd,
		IndexType: sg.IndexTypeUint16,
		CullMode:  sg.CullModeBack,
		Depth: sg.DepthState{
			WriteEnabled: true,
			Compare:      sg.CompareFuncLessEqual,
		},
		Label: "cube-pipeline",
	})

	// setup resource bindings
	state.Bind = sg.Bindings{
		VertexBuffers: [8]sg.Buffer{vbuf},
		IndexBuffer:   ibuf,
	}
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
		WindowTitle: "Cube (sokol-app)",
	})
}
