//go:build darwin && !SOKOL_METAL

package gokol

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore -framework OpenGL
#define SOKOL_GLCORE

#define SOKOL_NO_ENTRY
#define SOKOL_APP_IMPL
#include "sokol_app.h"
*/
import "C"
