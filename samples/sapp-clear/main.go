package main

import (
	glue "github.com/GiffE/gokol/glue"
	sapp "github.com/GiffE/gokol/sapp"
	sg "github.com/GiffE/gokol/sg"
)

var passAction sg.PassAction

func Init() {
	sg.Setup(&sg.Desc{
		Environment: glue.Environment(),
	})

	passAction = sg.PassAction{
		Colors: [sg.MaxColorAttachments]sg.ColorAttachmentAction{
			sg.ColorAttachmentAction{
				LoadAction: sg.LoadActionClear,
				ClearValue: sg.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
			},
		},
	}
}

func Frame() {
	sg.BeginPass(&sg.Pass{
		Action:    passAction,
		Swapchain: glue.Swapchain(),
	})
	sg.EndPass()
	sg.Commit()
}

func Cleanup() {
	sg.Shutdown()
}

func main() {
	sapp.Run(&sapp.AppDesc{
		Width:       400,
		Height:      300,
		WindowTitle: "Clear (sokol app)",
		Init:        Init,
		Cleanup:     Cleanup,
		Frame:       Frame,
	})
}
