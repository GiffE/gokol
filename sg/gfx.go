package sg

/*
#include "sokol_gfx.h"
#include <stdlib.h>

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
#cgo nocallback sg_make_buffer
#cgo noescape sg_make_buffer
#cgo nocallback sg_make_shader
#cgo noescape sg_make_shader
#cgo nocallback sg_setup
#cgo noescape sg_setup
#cgo nocallback sg_apply_pipeline
#cgo noescape sg_apply_pipeline
#cgo nocallback sg_apply_bindings
#cgo noescape sg_apply_bindings
#cgo nocallback sg_apply_uniforms
#cgo noescape sg_apply_uniforms
#cgo nocallback sg_draw
#cgo noescape sg_draw
#cgo nocallback sg_query_buffer_overflow
#cgo noescape sg_query_buffer_overflow
*/
import "C"

import (
	"reflect"
	"runtime"
	"unsafe"
)

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
	}
	if pass.Label != "" {
		sgPass.label = Str(pass.Label)
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

func MakeBuffer[T any](desc *BufferDesc[T]) Buffer {
	label, free := Strs(desc.Label)
	defer free()
	sgDesc := C.sg_buffer_desc{
		label: label[0],
		size:  C.size_t(desc.Size),
		usage: C.sg_buffer_usage{
			vertex_buffer:  C.bool(desc.Usage.VertexBuffer),
			index_buffer:   C.bool(desc.Usage.IndexBuffer),
			storage_buffer: C.bool(desc.Usage.StorageBuffer),
			immutable:      C.bool(desc.Usage.Immutable),
			dynamic_update: C.bool(desc.Usage.DynamicUpdate),
			stream_update:  C.bool(desc.Usage.StreamUpdate),
		},
		data: C.sg_range{
			ptr:  unsafe.Pointer(&desc.Data[0]),
			size: C.size_t(unsafe.Sizeof(desc.Data[0]) * uintptr(len(desc.Data))),
		},
		// optionally inject backend-specific resources
	}
	return Buffer{Id: uint32(C.sg_make_buffer(&sgDesc).id)}
}

func MakeShader(desc *ShaderDesc) Shader {
	strs, free := Strs(desc.Label, // 0
		desc.VertexFunc.Entry,           // 1
		desc.VertexFunc.Source,          // 2
		desc.VertexFunc.D3D11Filepath,   // 3
		desc.VertexFunc.D3D11Target,     // 4
		desc.FragmentFunc.Entry,         // 5
		desc.FragmentFunc.Source,        // 6
		desc.FragmentFunc.D3D11Filepath, // 7
		desc.FragmentFunc.D3D11Target,   // 8
		desc.ComputeFunc.Entry,          // 9
		desc.ComputeFunc.D3D11Filepath,  // 10
		desc.ComputeFunc.D3D11Target,    // 11
	)
	defer free()
	sgShaderDesc := C.sg_shader_desc{
		label: strs[0],
		mtl_threads_per_threadgroup: C.sg_mtl_shader_threads_per_threadgroup{
			x: C.int(desc.MtlTreadsPerThreadGroup.X),
			y: C.int(desc.MtlTreadsPerThreadGroup.Y),
			z: C.int(desc.MtlTreadsPerThreadGroup.Z),
		},
		vertex_func: C.sg_shader_function{
			entry:          strs[1],
			source:         strs[2],
			bytecode:       rangeRef(desc.VertexFunc.Bytecode),
			d3d11_filepath: strs[3],
			d3d11_target:   strs[4], // default: "vs_4_0" or "ps_4_0"
		},
		fragment_func: C.sg_shader_function{
			entry:          strs[5],
			source:         strs[6],
			bytecode:       rangeRef(desc.FragmentFunc.Bytecode),
			d3d11_filepath: strs[7],
			d3d11_target:   strs[8], // default: "vs_4_0" or "ps_4_0"
		},
	}
	// Convert vertex attributes
	for i, attr := range desc.Attrs {
		attrStr, attrFree := Strs(attr.GlslName, attr.HlslSemName)
		defer attrFree()
		sgShaderDesc.attrs[i] = C.sg_shader_vertex_attr{
			base_type:      C.sg_shader_attr_base_type(attr.BaseType),
			hlsl_sem_index: C.uint8_t(attr.HlslSemIndex),
			glsl_name:      attrStr[0],
			// TODO: yield these into the above alloc to reuse
			hlsl_sem_name: attrStr[1],
		}
	}

	// Convert uniform blocks
	for i, block := range desc.UniformBlocks {
		sgUniformBlock := C.sg_shader_uniform_block{
			stage:                 C.sg_shader_stage(block.Stage),
			size:                  C.uint32_t(block.Size),
			hlsl_register_b_n:     C.uint8_t(block.HlslRegisterBN),
			msl_buffer_n:          C.uint8_t(block.MslBufferN),
			wgsl_group0_binding_n: C.uint8_t(block.WgslGroup0BindingN),
			layout:                C.sg_uniform_layout(block.Layout),
		}
		for j, uniform := range block.GlslUniforms {
			attrStr, uniformFree := Strs(uniform.GlslName)
			defer uniformFree()
			sgUniformBlock.glsl_uniforms[j] = C.sg_glsl_shader_uniform{
				_type:       C.sg_uniform_type(uniform.Type),
				array_count: C.uint16_t(uniform.ArrayCount),
				glsl_name:   attrStr[0],
			}
		}
		sgShaderDesc.uniform_blocks[i] = sgUniformBlock
	}

	// Convert storage buffers
	// for i, buffer := range desc.StorageBuffers {
	// 	sgShaderDesc.storage_buffers[i] = C.sg_shader_storage_buffer{
	// 		stage:                 C.sg_shader_stage(buffer.Stage),
	// 		readonly:              C.bool(buffer.Readonly),
	// 		hlsl_register_t_n:     C.uint8_t(buffer.HlslRegisterTN),
	// 		hlsl_register_u_n:     C.uint8_t(buffer.HlslRegisterUN),
	// 		msl_buffer_n:          C.uint8_t(buffer.MslBufferN),
	// 		wgsl_group1_binding_n: C.uint8_t(buffer.WgslGroup1BindingN),
	// 		glsl_binding_n:        C.uint8_t(buffer.GlslBindingN),
	// 	}
	// }

	// Convert images
	// for i, image := range desc.Images {
	// 	sgShaderDesc.images[i] = C.sg_shader_image{
	// 		stage:                 C.sg_shader_stage(image.Stage),
	// 		image_type:            C.sg_image_type(image.ImageType),
	// 		sample_type:           C.sg_image_sample_type(image.SampleType),
	// 		multisampled:          C.bool(image.Multisampled),
	// 		hlsl_register_t_n:     C.uint8_t(image.HlslRegisterTN),
	// 		msl_texture_n:         C.uint8_t(image.MslTextureN),
	// 		wgsl_group1_binding_n: C.uint8_t(image.WgslGroup1BindingN),
	// 	}
	// }

	// Convert samplers
	// for i, sampler := range desc.Samplers {
	// 	sgShaderDesc.samplers[i] = C.sg_shader_sampler{
	// 		stage:                 C.sg_shader_stage(sampler.Stage),
	// 		sampler_type:          C.sg_sampler_type(sampler.SamplerType),
	// 		hlsl_register_s_n:     C.uint8_t(sampler.HlslRegisterSN),
	// 		msl_sampler_n:         C.uint8_t(sampler.MslSamplerN),
	// 		wgsl_group1_binding_n: C.uint8_t(sampler.WgslGroup1BindingN),
	// 	}
	// }

	// Convert image sampler pairs
	// for i, pair := range desc.ImageSamplerPairs {
	// 	sgShaderDesc.image_sampler_pairs[i] = C.sg_shader_image_sampler_pair{
	// 		stage:        C.sg_shader_stage(pair.Stage),
	// 		image_slot:   C.uint8_t(pair.ImageSlot),
	// 		sampler_slot: C.uint8_t(pair.SamplerSlot),
	// 		glsl_name:    tmpstring(pair.GlslName),
	// 	}
	// }
	return Shader{Id: uint32(C.sg_make_shader(&sgShaderDesc).id)}
}

func MakePipeline(desc *PipelineDesc) Pipeline {
	strs, free := Strs(desc.Label)
	defer free()
	sgDesc := C.sg_pipeline_desc{
		label:                     strs[0],
		compute:                   C.bool(desc.Compute),
		shader:                    C.sg_shader{C.uint32_t(desc.Shader.Id)},
		color_count:               C.int(desc.ColorCount),
		primitive_type:            C.sg_primitive_type(desc.PrimitiveType),
		index_type:                C.sg_index_type(desc.IndexType),
		cull_mode:                 C.sg_cull_mode(desc.CullMode),
		face_winding:              C.sg_face_winding(desc.FaceWinding),
		sample_count:              C.int(desc.SampleCount),
		alpha_to_coverage_enabled: C.bool(desc.AlphaToCoverageEnabled),
		depth: C.sg_depth_state{
			write_enabled: C.bool(desc.Depth.WriteEnabled),
			compare:       C.sg_compare_func(desc.Depth.Bias),
		},
		layout: C.sg_vertex_layout_state{},
		//  sg_color_target_state colors[SG_MAX_COLOR_ATTACHMENTS];
		// stencil:     C.sg_stencil_state{},
		//  sg_color blend_color;
	}
	for i, buf := range desc.Layout.Buffers {
		sgDesc.layout.buffers[i] = C.sg_vertex_buffer_layout_state{
			stride:    C.int(buf.Stride),
			step_func: C.sg_vertex_step(buf.StepFunc),
			step_rate: C.int(buf.StepRate),
		}
	}
	for i, attr := range desc.Layout.Attrs {
		sgDesc.layout.attrs[i] = C.sg_vertex_attr_state{
			buffer_index: C.int(attr.BufferIndex),
			offset:       C.int(attr.Offset),
			format:       C.sg_vertex_format(attr.Format),
		}
	}
	return Pipeline{Id: uint32(C.sg_make_pipeline(&sgDesc).id)}
}

func ApplyBindings(bind *Bindings) {
	sgBindings := C.sg_bindings{}
	for i, buffer := range bind.VertexBuffers {
		sgBindings.vertex_buffers[i] = C.sg_buffer{id: C.uint32_t(buffer.Id)}
	}
	for i, offset := range bind.VertexBufferOffsets {
		sgBindings.vertex_buffer_offsets[i] = C.int(offset)
	}
	if bind.IndexBuffer.Id != 0 {
		sgBindings.index_buffer = C.sg_buffer{id: C.uint32_t(bind.IndexBuffer.Id)}
	}
	sgBindings.index_buffer_offset = C.int(bind.IndexBufferOffset)
	for i, image := range bind.Images {
		sgBindings.images[i] = C.sg_image{id: C.uint32_t(image.Id)}
	}
	for i, sampler := range bind.Samplers {
		sgBindings.samplers[i] = C.sg_sampler{id: C.uint32_t(sampler.Id)}
	}
	for i, buffer := range bind.StorageBuffers {
		sgBindings.storage_buffers[i] = C.sg_buffer{id: C.uint32_t(buffer.Id)}
	}
	C.sg_apply_bindings(&sgBindings)
}

func ApplyPipeline(pip Pipeline) { C.sg_apply_pipeline(C.sg_pipeline{id: C.uint32_t(pip.Id)}) }
func ApplyUniform[T any](ubSlot int, data *T) {
	toRange := C.sg_range{}
	pinner := runtime.Pinner{}
	pinner.Pin(data)
	defer pinner.Unpin()
	// Use reflection to check if T is a slice
	v := reflect.ValueOf(data).Elem() // Dereference the pointer
	if v.Kind() == reflect.Slice {
		// Handle slice: get pointer to first element and total data size
		if v.Len() > 0 {
			toRange.ptr = unsafe.Pointer(v.Index(0).Addr().Pointer())
			elemSize := v.Type().Elem().Size()
			toRange.size = C.size_t(uintptr(v.Len()) * elemSize)
		} else {
			toRange.ptr = nil
			toRange.size = 0
		}
	} else {
		// Handle regular types: get pointer and size
		toRange.ptr = unsafe.Pointer(data)
		toRange.size = C.size_t(unsafe.Sizeof(*data))
	}
	C.sg_apply_uniforms(C.int(ubSlot), &toRange)
}
func EndPass()      { C.sg_end_pass() }
func Shutdown()     { C.sg_shutdown() }
func Commit()       { C.sg_commit() }
func IsValid() bool { return bool(C.sg_isvalid()) }
func Draw(baseElement, numElements, numInstances int) {
	C.sg_draw(C.int(baseElement), C.int(numElements), C.int(numInstances))
}
func QueryBufferOverflow(buf Buffer) bool {
	return bool(C.sg_query_buffer_overflow(C.sg_buffer{id: C.uint32_t(buf.Id)}))
}

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
	// optional log callback
	Logger      LogCallback
	Environment Environment
}

type LogCallback func()

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

type Usage struct {
	VertexBuffer  bool
	IndexBuffer   bool
	StorageBuffer bool
	Immutable     bool
	DynamicUpdate bool
	StreamUpdate  bool
}

type (
	ShaderDesc struct {
		VertexFunc              ShaderFunction
		FragmentFunc            ShaderFunction
		ComputeFunc             ShaderFunction
		Attrs                   [MaxVertexAttributes]ShaderVertexAttr
		UniformBlocks           [MaxUniformBlockBindSlots]ShaderUniformBlock
		StorageBuffers          [MaxStorageBufferBindSlots]ShaderStorageBuffer
		Images                  [MaxImageBindSlots]ShaderImage
		Samplers                [MaxSamplerBindSlots]ShaderSampler
		ImageSamplerPairs       [MaxImageSamplerPairs]ShaderImageSamplerPair
		MtlTreadsPerThreadGroup MtlShaderThreadsPerThreadGroup
		Label                   string
	}
	ShaderFunction struct {
		Source        string
		Bytecode      []byte
		Entry         string
		D3D11Target   string // default: "vs_4_0" or "ps_4_0"
		D3D11Filepath string
	}
	ShaderVertexAttr struct {
		BaseType     ShaderAttrBaseType // default: UNDEFINED (disables validation)
		GlslName     string             // [optional] GLSL attribute name
		HlslSemName  string             // HLSL semantic name
		HlslSemIndex uint8              // HLSL semantic index
	}
	ShaderUniformBlock struct {
		Stage              ShaderStage
		Size               uint32
		HlslRegisterBN     uint8 // HLSL register(bn)
		MslBufferN         uint8 // MSL [[buffer(n)]]
		WgslGroup0BindingN uint8 // WGSL @group(0) @binding(n)
		Layout             UniformLayout
		GlslUniforms       [MaxUniformBlockMembers]ShaderUniform
	}
	ShaderUniform struct {
		Type       UniformType
		ArrayCount uint16 // 0 or 1 for scalars, >1 for arrays
		GlslName   string // glsl name binding is required on GL 4.1 and WebGL2
	}
	ShaderStorageBuffer struct {
		Stage              ShaderStage
		Readonly           bool
		HlslRegisterTN     uint8 // HLSL register(tn) bind slot (for readonly access)
		HlslRegisterUN     uint8 // HLSL register(un) bind slot (for read/write access)
		MslBufferN         uint8 // MSL [[buffer(n)]] bind slot
		WgslGroup1BindingN uint8 // WGSL @group(1) @binding(n) bind slot
		GlslBindingN       uint8 // GLSL layout(binding=n)
	}
	ShaderImage struct {
		Stage              ShaderStage
		ImageType          ImageType
		SampleType         ImageSampleType
		Multisampled       bool
		HlslRegisterTN     uint8 // HLSL register(tn) bind slot
		MslTextureN        uint8 // MSL [[texture(n)]] bind slot
		WgslGroup1BindingN uint8 // WGSL @group(1) @binding(n) bind slot
	}
	ShaderSampler struct {
		Stage              ShaderStage
		SamplerType        SamplerType
		HlslRegisterSN     uint8 // HLSL register(sn) bind slot
		MslSamplerN        uint8 // MSL [[sampler(n)]] bind slot
		WgslGroup1BindingN uint8 // WGSL @group(1) @binding(n) bind slot
	}
	ShaderImageSamplerPair struct {
		Stage       ShaderStage
		ImageSlot   uint8
		SamplerSlot uint8
		GlslName    string // glsl name binding required because of GL 4.1 and WebGL2
	}
	MtlShaderThreadsPerThreadGroup struct {
		X, Y, Z int
	}
)

type ShaderAttrBaseType uint32

const (
	ShaderAttrBaseTypeUndefined ShaderAttrBaseType = iota
	ShaderAttrBaseTypeFloat
	ShaderAttrBaseTypeSint
	ShaderAttrBaseTypeUint
)

type ShaderStage uint32

const (
	ShaderStageNone ShaderStage = iota
	ShaderStageVertex
	ShaderStageFragment
	ShaderStageCompute
)

type UniformType uint32

const (
	UniformTypeInvalid UniformType = iota
	UniformTypeFloat
	UniformTypeFloat2
	UniformTypeFloat3
	UniformTypeFloat4
	UniformTypeInt
	UniformTypeInt2
	UniformTypeInt3
	UniformTypeInt4
	UniformTypeMat4
)

type UniformLayout uint32

const (
	UniformLayoutDefault UniformLayout = iota // value 0 reserved for default-init
	UniformLayoutNative                       // default: layout depends on currently active backend
	UniformLayoutStd140                       // std140: memory layout according to std140
)

type ImageType uint32

const (
	ImageTypeDefault ImageType = iota // value 0 reserved for default-init
	ImageType2d
	ImageTypeCube
	ImageType3d
	ImageTypeArray
)

type ImageSampleType uint32

const (
	ImageSampleTypeDefault ImageSampleType = iota // value 0 reserved for default-init
	ImageSampleTypeFloat
	ImageSampleTypeDepth
	ImageSampleTypeSint
	ImageSampleTypeUint
	ImageSampleTypeUnfilterableFloat
)

type SamplerType uint32

const (
	SamplerTypeDefault SamplerType = iota
	SamplerTypeFiltering
	SamplerTypeNonfiltering
	SamplerTypeComparison
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
	PipelineDesc struct {
		Compute                bool
		Shader                 Shader
		Layout                 VertexLayoutState
		Depth                  DepthState
		Stencil                StencilState
		ColorCount             int
		Colors                 [MaxColorAttachments]ColorTargetState
		PrimitiveType          PrimitiveType
		IndexType              IndexType
		CullMode               CullMode
		FaceWinding            FaceWinding
		SampleCount            int
		BlendColor             Color
		AlphaToCoverageEnabled bool
		Label                  string
	}
	VertexLayoutState struct {
		Buffers [MaxVertexBufferBindSlots]VertexBufferLayoutState
		Attrs   [MaxVertexAttributes]VertexAttrState
	}
	VertexBufferLayoutState struct {
		Stride   int
		StepFunc VertexStep
		StepRate int
	}
	VertexAttrState struct {
		BufferIndex int
		Offset      int
		Format      VertexFormat
	}
	DepthState struct {
		PixelFormat    PixelFormat
		Compare        CompareFunc
		WriteEnabled   bool
		Bias           float32
		BiasSlopeScale float32
		BiasClamp      float32
	}
	StencilState struct {
		Enabled   bool
		Front     StencilFaceState
		Back      StencilFaceState
		ReadMask  uint8
		WriteMask uint8
		Ref       uint8
	}
	StencilFaceState struct {
		Compare     CompareFunc
		FailOp      StencilOp
		DepthFailOp StencilOp
		PassOp      StencilOp
	}
	ColorTargetState struct {
		PixelFormat PixelFormat
		WriteMask   ColorMask
		Blend       BlendState
	}
	BlendState struct {
		Enabled        bool
		SrcFactorRGB   BlendFactor
		DstFactorRgb   BlendFactor
		OpRgb          BlendOp
		SrcFactorAlpha BlendFactor
		DstFactorAlpha BlendFactor
		OpAlpha        BlendOp
	}
)

type IndexType uint32

const (
	IndexTypeDefault IndexType = iota // value 0 reserved for default-init
	IndexTypeNone
	IndexTypeUint16
	IndexTypeUint32
)

type VertexStep uint32

const (
	VertexStepDefault VertexStep = iota // value 0 reserved for default-init
	VertexStepPerVertex
	VertexStepPerInstance
)

type VertexFormat uint32

const (
	VertexFormatInvalid VertexFormat = iota
	VertexFormatFloat
	VertexFormatFloat2
	VertexFormatFloat3
	VertexFormatFloat4
	VertexFormatInt
	VertexFormatInt2
	VertexFormatInt3
	VertexFormatInt4
	VertexFormatUint
	VertexFormatUint2
	VertexFormatUint3
	VertexFormatUint4
	VertexFormatByte4
	VertexFormatByte4N
	VertexFormatUbyte4
	VertexFormatUbyte4N
	VertexFormatShort2
	VertexFormatShort2N
	VertexFormatUshort2
	VertexFormatUshort2N
	VertexFormatShort4
	VertexFormatShort4N
	VertexFormatUshort4
	VertexFormatUshort4N
	VertexFormatUint10N2
	VertexFormatHalf2
	VertexFormatHalf4
	VertexFormatNum
)

type CompareFunc uint32

const (
	CompareFuncDefault CompareFunc = iota // value 0 reserved for default-init
	CompareFuncNever
	CompareFuncLess
	CompareFuncEqual
	CompareFuncLessEqual
	CompareFuncGreater
	CompareFuncNotEqual
	CompareFuncGreaterEqual
	CompareFuncAlways
)

type StencilOp uint32

const (
	StencilOpDefault StencilOp = iota // value 0 reserved for default-init
	StencilOpKeep
	StencilOpZero
	StencilOpReplace
	StencilOpIncrClamp
	StencilOpDecrClamp
	StencilOpInvert
	StencilOpIncrWrap
	StencilOpDecrWrap
)

type ColorMask uint32

const (
	ColorMaskDefault = ColorMask(0x0)  // value 0 reserved for default-init
	ColorMaskNone    = ColorMask(0x10) // special value for 'all channels disabled
	ColorMaskR       = ColorMask(0x1)
	ColorMaskG       = ColorMask(0x2)
	ColorMaskRg      = ColorMask(0x3)
	ColorMaskB       = ColorMask(0x4)
	ColorMaskRb      = ColorMask(0x5)
	ColorMaskGb      = ColorMask(0x6)
	ColorMaskRgb     = ColorMask(0x7)
	ColorMaskA       = ColorMask(0x8)
	ColorMaskRa      = ColorMask(0x9)
	ColorMaskGa      = ColorMask(0xA)
	ColorMaskRga     = ColorMask(0xB)
	ColorMaskBa      = ColorMask(0xC)
	ColorMaskRba     = ColorMask(0xD)
	ColorMaskGba     = ColorMask(0xE)
	ColorMaskRgba    = ColorMask(0xF)
)

type BlendFactor uint32

const (
	BlendFactorDefault BlendFactor = iota // value 0 reserved for default-init
	BlendFactorZero
	BlendFactorOne
	BlendFactorSrcColor
	BlendFactorOneMinusSrcColor
	BlendFactorSrcAlpha
	BlendFactorOneMinusSrcAlpha
	BlendFactorDstColor
	BlendFactorOneMinusDstColor
	BlendFactorDstAlpha
	BlendFactorOneMinusDstAlpha
	BlendFactorSrcAlphaSaturated
	BlendFactorBlendColor
	BlendFactorOneMinusBlendColor
	BlendFactorBlendAlpha
	BlendFactorOneMinusBlendAlpha
)

type BlendOp uint32

const (
	BlendOpDefault BlendOp = iota // value 0 reserved for default-init
	BlendOpAdd
	BlendOpSubtract
	BlendOpReverseSubtract
	BlendOpMin
	BlendOpMax
)

type PrimitiveType uint32

const (
	PrimitiveTypeDefault PrimitiveType = iota // value 0 reserved for default-init
	PrimitiveTypePoints
	PrimitiveTypeLines
	PrimitiveTypeLineStrip
	PrimitiveTypeTriangles
	PrimitiveTypeTriangleStrip
)

type CullMode uint32

const (
	CullModeDefault CullMode = iota // value 0 reserved for default-init
	CullModeNone
	CullModeFront
	CullModeBack
)

type FaceWinding uint32

const (
	FaceWindingDefault FaceWinding = iota // value 0 reserved for default-init
	FaceWindingCcw
	FaceWindingCw
)

type Bindings struct {
	VertexBuffers       [MaxVertexBufferBindSlots]Buffer
	VertexBufferOffsets [MaxVertexBufferBindSlots]int
	IndexBuffer         Buffer
	IndexBufferOffset   int
	Images              [MaxImageBindSlots]Image
	Samplers            [MaxSamplerBindSlots]Sampler
	StorageBuffers      [MaxStorageBufferBindSlots]Buffer
}

type BufferDesc[T any] struct {
	Size  uint64
	Usage Usage
	Data  []T
	Label string

	// optionally inject backend-specific resources
	// GlBuffers   [NumInflightFrames]uint32
	// MtlBuffers  [NumInflightFrames]unsafe.Pointer
	// D3D11Buffer unsafe.Pointer
	// WgpuBuffer  unsafe.Pointer
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

func rangeRef[T any](p []T) C.sg_range {
	if p == nil {
		return C.sg_range{}
	}
	return C.sg_range{
		ptr:  unsafe.Pointer(&p[0]),
		size: C.size_t(len(p)),
	}
}
