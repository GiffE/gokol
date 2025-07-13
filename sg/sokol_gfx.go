package sg

/*
#include "sokol_gfx.h"

#cgo nocallback sg_setup
#cgo noescape sg_setup
#cgo nocallback sg_begin_pass
#cgo noescape sg_begin_pass
#cgo nocallback sg_end_pass
#cgo noescape sg_end_pass
#cgo nocallback sg_shutdown
#cgo noescape sg_shutdown
#cgo nocallback sg_commit
#cgo noescape sg_commit
#cgo nocallback sg_isvalid
#cgo noescape sg_isvalid
*/
import "C"

import "unsafe"

func Setup(desc *Desc) {
	C.sg_setup(&C.sg_desc{
		buffer_pool_size:                                C.int(desc.BufferPoolSize),
		image_pool_size:                                 C.int(desc.ImagePoolSize),
		sampler_pool_size:                               C.int(desc.SamplerPoolSize),
		shader_pool_size:                                C.int(desc.ShaderPoolSize),
		pipeline_pool_size:                              C.int(desc.PipelinePoolSize),
		attachments_pool_size:                           C.int(desc.AttachmentsPoolSize),
		uniform_buffer_size:                             C.int(desc.BufferPoolSize),
		max_dispatch_calls_per_pass:                     C.int(desc.MaxDispatchCallsPerPass),
		max_commit_listeners:                            C.int(desc.MaxCommitListeners),
		disable_validation:                              C.bool(desc.DisableValidation),
		d3d11_shader_debugging:                          C.bool(desc.D3D11ShaderDebugging),
		mtl_force_managed_storage_mode:                  C.bool(desc.MtlForceManagedStorageMode),
		mtl_use_command_buffer_with_retained_references: C.bool(desc.MtlUseCommandBufferWithRetainedReferences),
		wgpu_disable_bindgroups_cache:                   C.bool(desc.WgpuDisableBindgroupsCache),
		wgpu_bindgroups_cache_size:                      C.int(desc.WgpuBindgroupsCacheSize),
		environment: C.sg_environment{
			defaults: C.sg_environment_defaults{
				color_format: C.sg_pixel_format(desc.Environment.Defaults.ColorFormat),
				depth_format: C.sg_pixel_format(desc.Environment.Defaults.DepthFormat),
				sample_count: C.int(desc.Environment.Defaults.SampleCount),
			},
			metal: C.sg_metal_environment{device: desc.Environment.Metal.Device},
			d3d11: C.sg_d3d11_environment{
				device:         desc.Environment.D3d11.Device,
				device_context: desc.Environment.D3d11.DeviceContext,
			},
			wgpu: C.sg_wgpu_environment{device: desc.Environment.Wgpu.Device},
		},
		// not implemented
		// allocator,
		// logger,
	})
}

func BeginPass(pass *Pass) {
	sgPass := C.sg_pass{
		compute: C.bool(pass.Compute),
		action: C.sg_pass_action{
			colors: [C.SG_MAX_COLOR_ATTACHMENTS]C.sg_color_attachment_action{},
			depth: C.sg_depth_attachment_action{
				load_action:  C.sg_load_action(pass.Action.Depth.LoadAction),
				store_action: C.sg_store_action(pass.Action.Depth.StoreAction),
				clear_value:  C.float(pass.Action.Depth.ClearValue),
			},
			stencil: C.sg_stencil_attachment_action{
				load_action:  C.sg_load_action(pass.Action.Stencil.LoadAction),
				store_action: C.sg_store_action(pass.Action.Stencil.StoreAction),
				clear_value:  C.uint8_t(pass.Action.Stencil.ClearValue),
			},
		},
		attachments: C.sg_attachments{id: C.uint32_t(pass.Attachments.Id)},
		swapchain: C.sg_swapchain{
			width:        C.int(pass.Swapchain.Width),
			height:       C.int(pass.Swapchain.Height),
			sample_count: C.int(pass.Swapchain.SampleCount),
			color_format: C.sg_pixel_format(pass.Swapchain.ColorFormat),
			depth_format: C.sg_pixel_format(pass.Swapchain.DepthFormat),
			metal: C.sg_metal_swapchain{
				current_drawable:      pass.Swapchain.Metal.CurrentDrawable,
				depth_stencil_texture: pass.Swapchain.Metal.DepthStencilTexture,
				msaa_color_texture:    pass.Swapchain.Metal.MSAAColorTexture,
			},
			d3d11: C.sg_d3d11_swapchain{
				render_view:        pass.Swapchain.D3D11.RenderView,
				resolve_view:       pass.Swapchain.D3D11.ResolveView,
				depth_stencil_view: pass.Swapchain.D3D11.DepthStencilView,
			},
			wgpu: C.sg_wgpu_swapchain{
				render_view:        pass.Swapchain.Wgpu.RenderView,
				resolve_view:       pass.Swapchain.Wgpu.ResolveView,
				depth_stencil_view: pass.Swapchain.Wgpu.DepthStencilView,
			},
			gl: C.sg_gl_swapchain{framebuffer: C.uint32_t(pass.Swapchain.Gl.Framebuffer)},
		},
		label: tmpstring(pass.Label),
	}
	for i, t := range pass.Action.Colors {
		sgPass.action.colors[i] = C.sg_color_attachment_action{
			load_action:  C.sg_load_action(t.LoadAction),
			store_action: C.sg_store_action(t.StoreAction),
			clear_value: C.sg_color{
				r: C.float(t.ClearValue.R),
				g: C.float(t.ClearValue.G),
				b: C.float(t.ClearValue.B),
				a: C.float(t.ClearValue.A),
			},
		}
	}

	C.sg_begin_pass(&sgPass)
}

func EndPass()      { C.sg_end_pass() }
func Shutdown()     { C.sg_shutdown() }
func Commit()       { C.sg_commit() }
func IsValid() bool { return bool(C.sg_isvalid()) }

type Desc struct {
	BufferPoolSize      int
	ImagePoolSize       int
	SamplerPoolSize     int
	ShaderPoolSize      int
	PipelinePoolSize    int
	AttachmentsPoolSize int
	UniformBufferSize   int
	// max expected number of dispatch calls per pass (default: 1024)
	MaxDispatchCallsPerPass int
	MaxCommitListeners      int
	// disable validation layer even in debug mode, useful for tests
	DisableValidation bool
	// if true, HLSL shaders are compiled with D3DCOMPILE_DEBUG | D3DCOMPILE_SKIP_OPTIMIZATION
	D3D11ShaderDebugging bool
	// for debugging: use Metal managed storage mode for resources even with UMA
	MtlForceManagedStorageMode bool
	// Metal: use a managed MTLCommandBuffer which ref-counts used resources
	MtlUseCommandBufferWithRetainedReferences bool
	// set to true to disable the WebGPU backend BindGroup cache
	WgpuDisableBindgroupsCache bool
	// number of slots in the WebGPU bindgroup cache (must be 2^N)
	WgpuBindgroupsCacheSize int
	// Allocator                                 Allocator
	// Logger                                    Logger // optional log function override
	Environment Environment
}

const (
	InvalidId                 = 0
	NumInflightFrames         = 2
	MaxColorAttachments       = 4
	MaxUniformBlockMembers    = 16
	MaxVertexAttributes       = 16
	MaxMipmaps                = 16
	MaxTexturearrayLayers     = 128
	MaxUniformBlockBindSlots  = 8
	MaxVertexBufferBindSlots  = 8
	MaxImageBindSlots         = 16
	MaxSamplerBindSlots       = 16
	MaxStorageBufferBindSlots = 8
	MaxImageSamplerPairs      = 16
)

type (
	Environment struct {
		Defaults EnvironmentDefaults
		Metal    MetalEnvironment
		D3d11    D3D11Environment
		Wgpu     WgpuEnviornment
	}

	EnvironmentDefaults struct {
		ColorFormat PixelFormat
		DepthFormat PixelFormat
		SampleCount int
	}
	MetalEnvironment struct {
		Device unsafe.Pointer
	}
	D3D11Environment struct {
		Device        unsafe.Pointer
		DeviceContext unsafe.Pointer
	}
	WgpuEnviornment struct {
		Device unsafe.Pointer
	}
)

type PixelFormat uint32

const (
	PixelFormatDefault PixelFormat = iota // value 0 reserved for default-init
	PixelFormatNone

	PixelFormatR8
	PixelFormatR8SN
	PixelFormatR8UI
	PixelFormatR8SI

	PixelFormatR16
	PixelFormatR16SN
	PixelFormatR16UI
	PixelFormatR16SI
	PixelFormatR16F
	PixelFormatRG8
	PixelFormatRG8SN
	PixelFormatRG8UI
	PixelFormatRG8SI

	PixelFormatR32UI
	PixelFormatR32SI
	PixelFormatR32F
	PixelFormatRG16
	PixelFormatRG16SN
	PixelFormatRG16UI
	PixelFormatRG16SI
	PixelFormatRG16F
	PixelFormatRGBA8
	PixelFormatSRGB8A8
	PixelFormatRGBA8SN
	PixelFormatRGBA8UI
	PixelFormatRGBA8SI
	PixelFormatBGRA8
	PixelFormatRGB10A2
	PixelFormatRG11B10F
	PixelFormatRGB9E5

	PixelFormatRG32UI
	PixelFormatRG32SI
	PixelFormatRG32F
	PixelFormatRGBA16
	PixelFormatRGBA16SN
	PixelFormatRGBA16UI
	PixelFormatRGBA16SI
	PixelFormatRGBA16F

	PixelFormatRGBA32UI
	PixelFormatRGBA32SI
	PixelFormatRGBA32F

	PixelFormatDepth
	PixelFormatDepthStencil

	PixelFormatBC1_RGBA
	PixelFormatBC2_RGBA
	PixelFormatBC3_RGBA
	PixelFormatBC3_SRGBA
	PixelFormatBC4_R
	PixelFormatBC4_RSN
	PixelFormatBC5_RG
	PixelFormatBC5_RGSN
	PixelFormatBC6H_RGBF
	PixelFormatBC6H_RGBUF
	PixelFormatBC7_RGBA
	PixelFormatBC7_SRGBA
	PixelFormatETC2_RGB8
	PixelFormatETC2_SRGB8
	PixelFormatETC2_RGB8A1
	PixelFormatETC2_RGBA8
	PixelFormatETC2_SRGB8A8
	PixelFormatEAC_R11
	PixelFormatEAC_R11SN
	PixelFormatEAC_RG11
	PixelFormatEAC_RG11SN

	PixelFormatASTC_4x4_RGBA
	PixelFormatASTC_4x4_SRGBA

	PixelFormatNum
)

type LoadAction uint32

const (
	LoadActionDefault LoadAction = iota
	LoadActionClear
	LoadActionLoad
	LoadActionDontCare
)

type StoreAction uint32

const (
	StoreactionDefault StoreAction = iota
	StoreactionStore
	StoreactionDontcare
)

type Color struct {
	R, G, B, A float32
}

type (
	PassAction struct {
		Colors  [MaxColorAttachments]ColorAttachmentAction
		Depth   DepthAttachmentAction
		Stencil StencilAttachmentAction
	}
	ColorAttachmentAction struct {
		LoadAction  LoadAction  // default: SG_LOADACTION_CLEAR
		StoreAction StoreAction // default: SG_STOREACTION_STORE
		ClearValue  Color       // default: { 0.5f, 0.5f, 0.5f, 1.0f }
	}
	DepthAttachmentAction struct {
		LoadAction  LoadAction  // default: SG_LOADACTION_CLEAR
		StoreAction StoreAction // default: SG_STOREACTION_DONTCARE
		ClearValue  float32     // default: 1.0
	}
	StencilAttachmentAction struct {
		LoadAction  LoadAction  // default: SG_LOADACTION_CLEAR
		StoreAction StoreAction // default: SG_STOREACTION_DONTCARE
		ClearValue  uint8       // default: 0
	}
)

type (
	Buffer      struct{ Id uint32 }
	Image       struct{ Id uint32 }
	Sampler     struct{ Id uint32 }
	Shader      struct{ Id uint32 }
	Pipeline    struct{ Id uint32 }
	Attachments struct{ Id uint32 }
)

type Pass struct {
	Compute     bool
	Action      PassAction
	Attachments Attachments
	Swapchain   Swapchain
	Label       string
}

type (
	Swapchain struct {
		Width       int
		Height      int
		SampleCount int
		ColorFormat PixelFormat
		DepthFormat PixelFormat
		Metal       MetalSwapchain
		D3D11       D3D11Swapchain
		Wgpu        WgpuSwapchain
		Gl          GlSwapchain
	}

	MetalSwapchain struct {
		CurrentDrawable     unsafe.Pointer // CAMetalDrawable (NOT MTLDrawable!!!)
		DepthStencilTexture unsafe.Pointer // MTLTexture
		MSAAColorTexture    unsafe.Pointer // MTLTexture
	}
	D3D11Swapchain struct {
		RenderView       unsafe.Pointer // ID3D11RenderTargetView
		ResolveView      unsafe.Pointer // ID3D11RenderTargetView
		DepthStencilView unsafe.Pointer // ID3D11DepthStencilView
	}
	WgpuSwapchain struct {
		RenderView       unsafe.Pointer // WGPUTextureView
		ResolveView      unsafe.Pointer // WGPUTextureView
		DepthStencilView unsafe.Pointer // WGPUTextureView
	}
	GlSwapchain struct {
		Framebuffer uint32 // GL framebuffer object
	}
)

func tmpstring(s string) *C.char {
	if s == "" {
		return nil
	}
	p := make([]byte, len(s)+1)
	copy(p, s)
	return (*C.char)(unsafe.Pointer(&p[0]))
}
