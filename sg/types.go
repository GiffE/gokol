package sg

import "unsafe"

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

/*
sg_image_desc

Creation parameters for sg_image objects, used in the sg_make_image() call.

The default configuration is:

.type               SG_IMAGETYPE_2D
.usage              .immutable = true
.width              0 (must be set to >0)
.height             0 (must be set to >0)
.num_slices         1 (3D textures: depth; array textures: number of layers)
.num_mipmaps        1
.pixel_format       SG_PIXELFORMAT_RGBA8 for textures, or sg_desc.environment.defaults.color_format for render targets
.sample_count       1 for textures, or sg_desc.environment.defaults.sample_count for render targets
.data               an sg_image_data struct to define the initial content
.label              0 (optional string label for trace hooks)

Q: Why is the default sample_count for render targets identical with the
"default sample count" from sg_desc.environment.defaults.sample_count?

A: So that it matches the default sample count in pipeline objects. Even
though it is a bit strange/confusing that offscreen render targets by default
get the same sample count as 'default swapchains', but it's better that
an offscreen render target created with default parameters matches
a pipeline object created with default parameters.

NOTE:

Regular (non-attachment) images with usage.immutable must be fully initialized by
providing a valid .data member which points to initialization data.

Images with usage.render_attachment or usage.storage_attachment must
*not* be created with initial content. Be aware that the initial
content of render- and storage-attachment images is undefined.

ADVANCED TOPIC: Injecting native 3D-API textures:

The following struct members allow to inject your own GL, Metal or D3D11
textures into sokol_gfx:

.gl_textures[SG_NUM_INFLIGHT_FRAMES]
.mtl_textures[SG_NUM_INFLIGHT_FRAMES]
.d3d11_texture
.d3d11_shader_resource_view
.wgpu_texture
.wgpu_texture_view

For GL, you can also specify the texture target or leave it empty to use
the default texture target for the image type (GL_TEXTURE_2D for
SG_IMAGETYPE_2D etc)

For D3D11 and WebGPU, either only provide a texture, or both a texture and
shader-resource-view / texture-view object. If you want to use access the
injected texture in a shader you *must* provide a shader-resource-view.

The same rules apply as for injecting native buffers (see sg_buffer_desc
documentation for more details).
*/
type ImageDesc[T any] struct {
	Type                    ImageType
	Usage                   ImageUsage
	Width                   int
	Height                  int
	NumSlices               int
	NumMipmaps              int
	PixelFormat             PixelFormat
	SampleCount             int
	Data                    ImageData[T]
	Label                   string
	GlTextures              [NumInflightFrames]uint32
	GlTextureTarget         uint32
	MtlTextures             [NumInflightFrames]unsafe.Pointer
	D3D11Texture            unsafe.Pointer
	D3D11ShaderResourceView unsafe.Pointer
	WgpuTexture             unsafe.Pointer
	WgpuTextureView         unsafe.Pointer
}

/*
sg_image_usage

Describes how the image object is going to be used:
render_attachment (default: false)

the image object is used as color-, resolve- or depth-stencil-
attachment in a render pass

storage_attachment (default: false)

the image object is used as storage-attachment in a
compute pass (to be written to by compute shaders)

immutable (default: true)

the image content cannot be updated from the CPU side
(but may be updated by the GPU in a render- or compute-pass)

dynamic_update (default: false)

the image content is updated infrequently by the CPU

stream_update (default: false)

the image content is updated each frame by the CPU via

Note that the usage as texture binding is implicit and always allowed.
*/
type ImageUsage struct {
	RenderAttachment  bool
	StorageAttachment bool
	Immutable         bool
	DynamicUpdate     bool
	StreamUpdate      bool
}

/*
sg_cube_face

The cubemap faces. Use these as indices in the sg_image_desc.content
array.
*/
type CubeFace uint32

const (
	CubeFacePosX CubeFace = iota
	CubeFaceNegX
	CubeFacePosY
	CubeFaceNegY
	CubeFacePosZ
	CubeFaceNegZ
	CubeFaceNum
)

/*
sg_image_data

Defines the content of an image through a 2D array of sg_range structs.
The first array dimension is the cubemap face, and the second array
dimension the mipmap level.
*/
type ImageData[T any] struct {
	Subimage [CubeFaceNum][MaxMipmaps][]T
}

/*
sg_sampler_desc

Creation parameters for sg_sampler objects, used in the sg_make_sampler() call

	.min_filter:        SG_FILTER_NEAREST
	.mag_filter:        SG_FILTER_NEAREST
	.mipmap_filter      SG_FILTER_NEAREST
	.wrap_u:            SG_WRAP_REPEAT
	.wrap_v:            SG_WRAP_REPEAT
	.wrap_w:            SG_WRAP_REPEAT (only SG_IMAGETYPE_3D)
	.min_lod            0.0f
	.max_lod            FLT_MAX
	.border_color       SG_BORDERCOLOR_OPAQUE_BLACK
	.compare            SG_COMPAREFUNC_NEVER
	.max_anisotropy     1 (must be 1..16)
*/
type SamplerDesc struct {
	MinFilter     Filter
	MagFilter     Filter
	MipmapFilter  Filter
	WrapU         Wrap
	WrapV         Wrap
	WrapW         Wrap
	MinLod        float32
	MaxLod        float32
	BorderColor   BorderColor
	Compare       CompareFunc
	MaxAnisotropy uint32
	Label         string
	GlSampler     uint32
	MtlSampler    unsafe.Pointer
	D3D11Sampler  unsafe.Pointer
	WgpuSampler   unsafe.Pointer
}

/*
sg_filter

The filtering mode when sampling a texture image. This is
used in the sg_sampler_desc.min_filter, sg_sampler_desc.mag_filter
and sg_sampler_desc.mipmap_filter members when creating a sampler object.

For the default is SG_FILTER_NEAREST.
*/
type Filter uint32

const (
	FilterDefault Filter = iota
	FilterNearest
	FilterLinear
	FilterNum
)

/*
sg_wrap

The texture coordinates wrapping mode when sampling a texture
image. This is used in the sg_image_desc.wrap_u, .wrap_v
and .wrap_w members when creating an image.

The default wrap mode is SG_WRAP_REPEAT.

NOTE: SG_WRAP_CLAMP_TO_BORDER is not supported on all backends
and platforms. To check for support, call sg_query_features()
and check the "clamp_to_border" boolean in the returned
sg_features struct.

Platforms which don't support SG_WRAP_CLAMP_TO_BORDER will silently fall back
to SG_WRAP_CLAMP_TO_EDGE without a validation error.
*/
type Wrap uint32

const (
	WrapDefault Wrap = iota // value 0 reserved for default-init
	WrapRepeat
	WrapClampToEdge
	WrapClampToBorder
	WrapMirroredRepeat
	WrapNum
)

/*
sg_border_color

The border color to use when sampling a texture, and the UV wrap
mode is SG_WRAP_CLAMP_TO_BORDER.

The default border color is SG_BORDERCOLOR_OPAQUE_BLACK
*/
type BorderColor uint32

const (
	BorderColorDefault BorderColor = iota // value 0 reserved for default-init
	BorderColorTransparentBlack
	BorderColorOpaqueBlack
	BorderColorOpaqueWhite
	BorderColorNum
)
