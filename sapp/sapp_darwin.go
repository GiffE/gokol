//go:build darwin

package gokol

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore -framework Metal -framework MetalKit
#define SOKOL_METAL

#define SOKOL_NO_ENTRY
#define SOKOL_APP_IMPL
#include "sokol_app.h"
*/
import "C"
