package glue

import (
	sapp "github.com/GiffE/gokol/sapp"
	sg "github.com/GiffE/gokol/sg"
)

func Environment() sg.Environment {
	return sg.Environment{
		Defaults: sg.EnvironmentDefaults{
			ColorFormat: sg.PixelFormat(sapp.ColorFormat()),
			DepthFormat: sg.PixelFormat(sapp.DepthFormat()),
			SampleCount: sapp.SampleCount(),
		},
		Metal: sg.MetalEnvironment{
			Device: sapp.MetalGetDevice(),
		},
		//   env.d3d11.device = sapp_d3d11_get_device();
		//   env.d3d11.device_context = sapp_d3d11_get_device_context();
		//   env.wgpu.device = sapp_wgpu_get_device();
	}

}

func Swapchain() sg.Swapchain {
	return sg.Swapchain{
		Width:       sapp.Width(),
		Height:      sapp.Height(),
		SampleCount: sapp.SampleCount(),
		ColorFormat: sg.PixelFormat(sapp.ColorFormat()),
		DepthFormat: sg.PixelFormat(sapp.DepthFormat()),

		// We could optimize these CGo call by using build
		// versions.  Metal is obviously apple only.
		Metal: sg.MetalSwapchain{
			CurrentDrawable:     sapp.MetalGetCurrentDrawable(),
			DepthStencilTexture: sapp.MetalGetDepthStencilTexture(),
			MSAAColorTexture:    sapp.MetalGetMSAAColorTexture(),
		},

		//  swapchain.d3d11.render_view = sapp_d3d11_get_render_view();
		//  swapchain.d3d11.resolve_view = sapp_d3d11_get_resolve_view();
		//  swapchain.d3d11.depth_stencil_view = sapp_d3d11_get_depth_stencil_view();
		//  swapchain.wgpu.render_view = sapp_wgpu_get_render_view();
		//  swapchain.wgpu.resolve_view = sapp_wgpu_get_resolve_view();
		//  swapchain.wgpu.depth_stencil_view = sapp_wgpu_get_depth_stencil_view();
		//  swapchain.gl.framebuffer = sapp_gl_get_framebuffer();
	}
}
