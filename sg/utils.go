package sg

import (
	"strings"
	"unsafe"
)

// #include <stdlib.h>
import "C"

// Str takes a null-terminated Go string and returns its sokol-compatible address.
// This function reaches into Go string storage in an unsafe way so the caller
// must ensure the string is not garbage collected.
func Str(str string) *C.char {
	if !strings.HasSuffix(str, "\x00") {
		panic("str argument missing null terminator: " + str)
	}
	return (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
}

// GoStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func GoStr(cstr *C.char) string {
	return C.GoString((*C.char)((unsafe.Pointer(cstr))))
}

// Strs takes a list of Go strings (with or without null-termination) and
// returns their C counterpart.
//
// The returned free function must be called once you are done using the strings
// in order to free the memory.
//
// If no strings are provided as a parameter this function will panic.
func Strs(strs ...string) (cstrs []*C.char, free func()) {
	if len(strs) == 0 {
		panic("Strs: expected at least 1 string")
	}

	// Allocate a contiguous array large enough to hold all the strings' contents.
	n := 0
	for i := range strs {
		n += len(strs[i])
	}
	if n == 0 {
		return make([]*C.char, len(strs)), func() {}
	}
	n += len(strs)
	data := C.malloc(C.size_t(n))

	// Copy all the strings into data.
	dataSlice := (*[1 << 30]byte)(data)[:n]
	css := make([]*C.char, len(strs)) // Populated with pointers to each string.
	offset := 0
	for i := range strs {
		if len(strs[i]) == 0 {
			css[i] = nil
			dataSlice[offset] = 0
			offset += 1
			continue
		}
		copy(dataSlice[offset:offset+len(strs[i])], strs[i][:]) // Copy strs[i] into proper data location.
		dataSlice[offset+len(strs[i])] = 0
		css[i] = (*C.char)(unsafe.Pointer(&dataSlice[offset])) // Set a pointer to it.
		offset += len(strs[i]) + 1
	}

	return css, func() { C.free(data) }
}
