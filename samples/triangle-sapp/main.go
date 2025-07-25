package main

import (
	glue "github.com/GiffE/gokol/glue"
	sapp "github.com/GiffE/gokol/sapp"
	sg "github.com/GiffE/gokol/sg"
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

func Init() {
	sg.Setup(&sg.Desc{
		Environment: glue.Environment(),
	})

	vertices := []float32{
		// positions  		// colors
		0.0, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0,
	}

	state.Bind.VertexBuffers[0] = sg.MakeBuffer(&sg.BufferDesc[float32]{
		Data:  vertices,
		Label: "triangle-vertices",
	})

	// create shader from code-generated sg_shader_desc
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

	// create a pipeline object (default render states are fine for triangle)
	state.Pip = sg.MakePipeline(&sg.PipelineDesc{
		Shader: shd,
		Label:  "triangle-pipeline",
		Layout: sg.VertexLayoutState{
			Attrs: [16]sg.VertexAttrState{
				// if the vertex layout doesn't have gaps, don't need to provide strides and offsets
				{Format: sg.VertexFormatFloat3},
				{Format: sg.VertexFormatFloat4},
			},
		},
	})

	state.PassAction = sg.PassAction{
		Colors: [sg.MaxColorAttachments]sg.ColorAttachmentAction{
			{
				LoadAction: sg.LoadActionClear,
				ClearValue: sg.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
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
	sg.ApplyBindings(&state.Bind)
	sg.Draw(0, 3, 1)
	sg.EndPass()
	sg.Commit()
}

func Cleanup() {
	sg.Shutdown()
}

func main() {
	sapp.Run(&sapp.AppDesc{
		Width:       400,
		Height:      400,
		Init:        Init,
		Cleanup:     Cleanup,
		Frame:       Frame,
		WindowTitle: "Triangle (sokol-app)",
		Event:       func(e sapp.Event) {},
	})
}
