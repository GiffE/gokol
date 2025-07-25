//go:build darwin && SOKOL_METAL

package sg

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Metal
#define SOKOL_GFX_IMPL
#define SOKOL_METAL
#include "sokol_gfx.h"
*/
import "C"
