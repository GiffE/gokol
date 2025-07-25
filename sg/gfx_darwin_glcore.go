//go:build darwin && !SOKOL_METAL

package sg

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework OpenGL
#define SOKOL_GFX_IMPL
#define SOKOL_GLCORE
#include "sokol_gfx.h"
*/
import "C"
