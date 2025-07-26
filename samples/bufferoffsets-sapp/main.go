package main

import (
	"unsafe"

	"github.com/GiffE/gokol/glue"
	sapp "github.com/GiffE/gokol/sapp"
	"github.com/GiffE/gokol/sg"
)

type State struct {
	Pip        sg.Pipeline
	Bind       sg.Bindings
	PassAction sg.PassAction
}

var state State

var vs_source_glcore = `
#version 410
layout(location = 0) in vec4 position;
layout(location = 0) out vec4 color;
layout(location = 1) in vec4 color0;

void main()
{
    gl_Position = position;
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

type vertex struct {
	x, y, r, g, b float32
}

func Init() {
	sg.Setup(&sg.Desc{
		Environment: glue.Environment(),
	})

	// a 2D triangle and quad in 1 vertex buffer and 1 index buffer
	vertices := []vertex{
		/* triangle */
		{0.0, 0.55, 1.0, 0.0, 0.0},
		{0.25, 0.05, 0.0, 1.0, 0.0},
		{-0.25, 0.05, 0.0, 0.0, 1.0},

		/* quad */
		{-0.25, -0.05, 0.0, 0.0, 1.0},
		{0.25, -0.05, 0.0, 1.0, 0.0},
		{0.25, -0.55, 1.0, 0.0, 0.0},
		{-0.25, -0.55, 1.0, 1.0, 0.0},
	}

	state.Bind.VertexBuffers[0] = sg.MakeBuffer(&sg.BufferDesc[vertex]{
		Data:  vertices,
		Label: "vertex-buffer",
	})

	// an index buffer with 2 triangles
	indices := []uint16{
		0, 1, 2,
		0, 1, 2, 0, 2, 3,
	}
	state.Bind.IndexBuffer = sg.MakeBuffer(&sg.BufferDesc[uint16]{
		Usage: sg.Usage{IndexBuffer: true},
		Data:  indices,
		Label: "index-buffer",
	})

	shd := sg.MakeShader(&sg.ShaderDesc{
		Label: "triangle_shader",
		VertexFunc: sg.ShaderFunction{
			Source: vs_source_glcore,
			Entry:  "main",
		},
		FragmentFunc: sg.ShaderFunction{
			Source: fs_source_glcore,
			Entry:  "main",
		},
		Attrs: [16]sg.ShaderVertexAttr{
			{GlslName: "position", BaseType: sg.ShaderAttrBaseTypeFloat},
			{GlslName: "color0", BaseType: sg.ShaderAttrBaseTypeFloat},
		},
	})

	// a pipeline state object
	state.Pip = sg.MakePipeline(&sg.PipelineDesc{
		Shader:    shd,
		IndexType: sg.IndexTypeUint16,
		Layout: sg.VertexLayoutState{
			Attrs: [16]sg.VertexAttrState{
				{Format: sg.VertexFormatFloat2},
				{Format: sg.VertexFormatFloat3},
			},
		},
		Label: "pipeline",
	})

	state.PassAction = sg.PassAction{
		Colors: [sg.MaxColorAttachments]sg.ColorAttachmentAction{
			{
				LoadAction: sg.LoadActionClear,
				ClearValue: sg.Color{R: 0.5, G: 0.5, B: 1.0, A: 1.0},
			},
		},
	}
}

func Frame() {
	sg.BeginPass(&sg.Pass{
		Action:    state.PassAction,
		Swapchain: glue.Swapchain(),
	})
	sg.ApplyPipeline(state.Pip)
	// render the triangle
	state.Bind.VertexBufferOffsets[0] = 0
	state.Bind.IndexBufferOffset = 0
	sg.ApplyBindings(&state.Bind)
	sg.Draw(0, 3, 1)
	// render the quad
	var x vertex
	var y uint16
	state.Bind.VertexBufferOffsets[0] = int(3 * unsafe.Sizeof(x))
	state.Bind.IndexBufferOffset = int(3 * unsafe.Sizeof(y))
	sg.ApplyBindings(&state.Bind)
	sg.Draw(0, 6, 1)
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
		WindowTitle: "Buffer Offsets (sokol-app)",
	})
}
