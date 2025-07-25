package gokol

/*
#include "sokol_app.h"
#include <stdlib.h>
typedef void (*_sokolgo_void_callback)(uintptr_t userdata);
typedef void (*_sokolgo_event_callback)(const sapp_event* event, uintptr_t userdata);
extern void cb_sokol_Init(uintptr_t userdata);
extern void cb_sokol_Frame(uintptr_t userdata);
extern void cb_sokol_Cleanup(uintptr_t userdata);
extern void cb_sokol_Event(sapp_event* event, uintptr_t userdata);

#cgo nocallback sapp_isvalid
#cgo nocallback sapp_widthf
#cgo nocallback sapp_heightf
#cgo nocallback sapp_color_format
#cgo nocallback sapp_depth_format
#cgo nocallback sapp_sample_count
#cgo nocallback sapp_high_dpi
#cgo nocallback sapp_dpi_scale
#cgo nocallback sapp_show_keyboard
#cgo nocallback sapp_keyboard_shown
#cgo nocallback sapp_is_fullscreen
#cgo nocallback sapp_toggle_fullscreen
#cgo nocallback sapp_show_mouse
#cgo nocallback sapp_mouse_shown
#cgo nocallback sapp_lock_mouse
#cgo nocallback sapp_mouse_locked
#cgo nocallback sapp_set_mouse_cursor
#cgo nocallback sapp_get_mouse_cursor
#cgo nocallback sapp_query_desc
#cgo nocallback sapp_request_quit
#cgo nocallback sapp_cancel_quit
#cgo nocallback sapp_quit
#cgo nocallback sapp_consume_event
#cgo nocallback sapp_frame_count
#cgo nocallback sapp_frame_duration
#cgo nocallback sapp_set_clipboard_string
#cgo nocallback sapp_get_clipboard_string
#cgo nocallback sapp_set_window_title
#cgo nocallback sapp_get_num_dropped_files
#cgo nocallback sapp_get_dropped_file_path
#cgo nocallback sapp_metal_get_device
#cgo nocallback sapp_metal_get_current_drawable
#cgo nocallback sapp_metal_get_depth_stencil_texture
#cgo nocallback sapp_metal_get_msaa_color_texture

#cgo noescape sapp_isvalid
#cgo noescape sapp_widthf
#cgo noescape sapp_heightf
#cgo noescape sapp_color_format
#cgo noescape sapp_depth_format
#cgo noescape sapp_sample_count
#cgo noescape sapp_high_dpi
#cgo noescape sapp_dpi_scale
#cgo noescape sapp_show_keyboard
#cgo noescape sapp_keyboard_shown
#cgo noescape sapp_is_fullscreen
#cgo noescape sapp_toggle_fullscreen
#cgo noescape sapp_show_mouse
#cgo noescape sapp_mouse_shown
#cgo noescape sapp_lock_mouse
#cgo noescape sapp_mouse_locked
#cgo noescape sapp_set_mouse_cursor
#cgo noescape sapp_get_mouse_cursor
#cgo noescape sapp_query_desc
#cgo noescape sapp_request_quit
#cgo noescape sapp_cancel_quit
#cgo noescape sapp_quit
#cgo noescape sapp_consume_event
#cgo noescape sapp_frame_count
#cgo noescape sapp_frame_duration
#cgo noescape sapp_set_clipboard_string
#cgo noescape sapp_get_clipboard_string
#cgo noescape sapp_set_window_title
#cgo noescape sapp_get_num_dropped_files
#cgo noescape sapp_get_dropped_file_path
#cgo noescape sapp_metal_get_device
#cgo noescape sapp_metal_get_current_drawable
#cgo noescape sapp_metal_get_depth_stencil_texture
#cgo noescape sapp_metal_get_msaa_color_texture
*/
import "C"
import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

type InitCallback func()
type FrameCallback func()
type CleanupCallback func()
type EventCallback func(ev Event)

type AppDesc struct {
	WindowTitle string
	// the preferred width of the window / canvas
	Width int
	// the preferred height of the window / canvas
	Height int
	// whether the rendering canvas is full-resolution on HighDPI displays
	HighDPI bool
	// whether the window should be created in fullscreen mode
	Fullscreen bool
	// MSAA sample count
	SampleCount int
	// the preferred swap interval (ignored on some platforms)
	SampleInterval int
	// whether the framebuffer should have an alpha channel (ignored on some platforms)
	Alpha bool
	// enable clipboard access, default is false
	EnableClipboard bool
	// max size of clipboard content in bytes
	ClipboardSize int
	// enable file dropping (drag'n'drop), default is false
	EnableDragnDrop bool
	// max number of dropped files to process (default: 1)
	MaxDroppedFiles int
	// max length in bytes of a dropped UTF-8 file path (default: 2048)
	MaxDroppedFilePathLength int

	Init    InitCallback
	Frame   FrameCallback
	Cleanup CleanupCallback
	Event   EventCallback
	// not supported
	// allocator

	// User data callbacks (should be easy, just needs to be added to the wrapped user data)
	// the initial window icon to set
	// Icon                        any // TODO: Add Icon support
	// logging callback override (default: NO LOGGING!)
	// Logger                      any // TODO: Add Logger support
	// GL version
	// html5 props
}

/*
EventType

The type of event that's passed to the event handler callback
in the sapp_event.type field. These are not just "traditional"
input events, but also notify the application about state changes
or other user-invoked actions.
*/
type EventType uint32

const (
	EventTypeInvalid EventType = iota
	EventTypeKeyDown
	EventTypeKeyUp
	EventTypeChar
	EventTypeMouseDown
	EventTypeMouseUp
	EventTypeMouseScroll
	EventTypeMouseMove
	EventTypeMouseEnter
	EventTypeMouseLeave
	EventTypeTouchesBegan
	EventTypeTouchesMoved
	EventTypeTouchesEnded
	EventTypeTouchesCancelled
	EventTypeResized
	EventTypeIconified
	EventTypeRestored
	EventTypeFocused
	EventTypeUnfocused
	EventTypeSuspended
	EventTypeResumed
	EventTypeQuitRequested
	EventTypeClipboardPasted
	EventTypeFilesDropped
)

/*
MouseButton

The currently pressed mouse button in the events MOUSE_DOWN
and MOUSE_UP, stored in the struct field sapp_event.mouse_button.
*/
type MouseButton uint32

const (
	MouseButtonLeft    MouseButton = 0x0
	MouseButtonRight   MouseButton = 0x1
	MouseButtonMiddle  MouseButton = 0x2
	MouseButtonInvalid MouseButton = 0x100
)

/*
MouseCursor

Predefined cursor image definitions, set with sapp_set_mouse_cursor(sapp_mouse_cursor cursor)
*/
type MouseCursor uint32

const (
	MouseCursorDefault = iota
	MouseCursorArrow
	MouseCursorIbeam
	MouseCursorCrosshair
	MouseCursorPointingHand
	MouseCursorResizeEW
	MouseCursorResizeNS
	MouseCursorResizeNWSE
	MouseCursorResizeNESW
	MouseCursorResizeAll
	MouseCursorNotAllowed
)

/*
Keycode

The 'virtual keycode' of a KEY_DOWN or KEY_UP event in the
struct field sapp_event.key_code.

Note that the keycode values are identical with GLFW.
*/
type Keycode uint32

const (
	KeycodeInvalid Keycode = iota
	KeycodeSpace
	KeycodeApostrophe // '
	KeycodeComma      // ,
	KeycodeMinus      // -
	KeycodePeriod     // .
	KeycodeSlash      // /
	Keycode0
	Keycode1
	Keycode2
	Keycode3
	Keycode4
	Keycode5
	Keycode6
	Keycode7
	Keycode8
	Keycode9
	KeycodeSemicolon // ;
	KeycodeEqual     // =
	KeycodeA
	KeycodeB
	KeycodeC
	KeycodeD
	KeycodeE
	KeycodeF
	KeycodeG
	KeycodeH
	KeycodeI
	KeycodeJ
	KeycodeK
	KeycodeL
	KeycodeM
	KeycodeN
	KeycodeO
	KeycodeP
	KeycodeQ
	KeycodeR
	KeycodeS
	KeycodeT
	KeycodeU
	KeycodeV
	KeycodeW
	KeycodeX
	KeycodeY
	KeycodeZ
	KeycodeLeftBracket  // [
	KeycodeBackslash    // \
	KeycodeRightBracket // ]
	KeycodeGraveAccent  // `
	KeycodeWorld1       // non-US #1
	KeycodeWorld2       // non-US #2
	KeycodeEscape
	KeycodeEnter
	KeycodeTab
	KeycodeBackspace
	KeycodeInsert
	KeycodeDelete
	KeycodeRight
	KeycodeLeft
	KeycodeDown
	KeycodeUp
	KeycodePageUp
	KeycodePageDown
	KeycodeHome
	KeycodeEnd
	KeycodeCapsLock
	KeycodeScrollLock
	KeycodeNumLock
	KeycodePrintScreen
	KeycodePause
	KeycodeF1
	KeycodeF2
	KeycodeF3
	KeycodeF4
	KeycodeF5
	KeycodeF6
	KeycodeF7
	KeycodeF8
	KeycodeF9
	KeycodeF10
	KeycodeF11
	KeycodeF12
	KeycodeF13
	KeycodeF14
	KeycodeF15
	KeycodeF16
	KeycodeF17
	KeycodeF18
	KeycodeF19
	KeycodeF20
	KeycodeF21
	KeycodeF22
	KeycodeF23
	KeycodeF24
	KeycodeF25
	KeycodeKp0
	KeycodeKp1
	KeycodeKp2
	KeycodeKp3
	KeycodeKp4
	KeycodeKp5
	KeycodeKp6
	KeycodeKp7
	KeycodeKp8
	KeycodeKp9
	KeycodeKpDecimal
	KeycodeKpDivide
	KeycodeKpMultiply
	KeycodeKpSubtract
	KeycodeKpAdd
	KeycodeKpEnter
	KeycodeKpEqual
	KeycodeLeftShift
	KeycodeLeftControl
	KeycodeLeftAlt
	KeycodeLeftSuper
	KeycodeRightShift
	KeycodeRightControl
	KeycodeRightAlt
	KeycodeRightSuper
	KeycodeMenu
)

type AndroidToolType uint8

const (
	AndroidToolTypeUnknown AndroidToolType = 0
	AndroidToolTypeFinger  AndroidToolType = 1
	AndroidToolTypeStylus  AndroidToolType = 2
	AndroidToolTypeMouse   AndroidToolType = 3
)

type Touchpoint struct {
	Identifier      uint64
	PosX, PosY      float32
	AndroidToolType AndroidToolType
	Changed         bool
}

/*
This is an all-in-one event struct passed to the event handler
user callback function. Note that it depends on the event
type what struct fields actually contain useful values, so you
should first check the event type before reading other struct
fields.
*/
type Event struct {
	internal *C.sapp_event
}

// current frame counter, always valid, useful for checking if two events were issued in the same frame
func (e Event) FrameCount() uint64 { return uint64(e.internal.frame_count) }

// the event type, always valid
func (e Event) Type() EventType { return EventType(e.internal._type) }

// the virtual key code, only valid in KEY_UP, KEY_DOWN
func (e Event) KeyCode() Keycode { return Keycode(e.internal.key_code) }

// the UTF-32 character code, only valid in CHAR events
func (e Event) CharCode() rune { return rune(e.internal.char_code) }

// true if this is a key-repeat event, valid in KEY_UP, KEY_DOWN and CHAR
func (e Event) KeyRepeat() bool { return bool(e.internal.key_repeat) }

// current modifier keys, valid in all key-, char- and mouse-events
func (e Event) Modifiers() uint32 { return uint32(e.internal.modifiers) }

// mouse button that was pressed or released, valid in MOUSE_DOWN, MOUSE_UP
func (e Event) MouseButton() MouseButton { return MouseButton(e.internal.mouse_button) }

// current horizontal mouse position in pixels, always valid except during mouse lock
func (e Event) MouseX() float32 { return float32(e.internal.mouse_x) }

// current vertical mouse position in pixels, always valid except during mouse lock
func (e Event) MouseY() float32 { return float32(e.internal.mouse_y) }

// relative horizontal mouse movement since last frame, always valid
func (e Event) MouseDX() float32 { return float32(e.internal.mouse_dx) }

// relative vertical mouse movement since last frame, always valid
func (e Event) MouseDY() float32 { return float32(e.internal.mouse_dy) }

// horizontal mouse wheel scroll distance, valid in MOUSE_SCROLL events
func (e Event) ScollX() float32 { return float32(e.internal.scroll_x) }

// vertical mouse wheel scroll distance, valid in MOUSE_SCROLL events
func (e Event) ScrollY() float32 { return float32(e.internal.scroll_y) }

// current window width and framebuffer sizes in pixels, always valid
func (e Event) WindowWidth() int { return int(e.internal.window_width) }

// current window height and framebuffer sizes in pixels, always valid
func (e Event) WindowHeight() int { return int(e.internal.window_height) }

// = window_width * dpi_scale
func (e Event) FramebufferWidth() int { return int(e.internal.framebuffer_width) }

// = window_height * dpi_scale
func (e Event) FramebufferHeight() int { return int(e.internal.framebuffer_height) }

// current touch points, valid in TOUCHES_BEGIN, TOUCHES_MOVED, TOUCHES_ENDED
func (e Event) Touches() []Touchpoint {
	panic("not implemented")
	// TODO: Implement touches, this should be a fairly straight forward copy.
}

type wrappedUserData struct {
	init    InitCallback
	frame   FrameCallback
	cleanup CleanupCallback
	event   EventCallback
}

//export cb_sokol_Init
func cb_sokol_Init(userdata uintptr) {
	cb := cgo.Handle(userdata).Value().(*wrappedUserData).init
	if cb != nil {
		cb()
	}
}

//export cb_sokol_Frame
func cb_sokol_Frame(userdata uintptr) {
	cb := cgo.Handle(userdata).Value().(*wrappedUserData).frame
	if cb != nil {
		cb()
	}
}

//export cb_sokol_Cleanup
func cb_sokol_Cleanup(userdata uintptr) {
	cb := cgo.Handle(userdata).Value().(*wrappedUserData).cleanup
	if cb != nil {
		cb()
	}
}

//export cb_sokol_Event
func cb_sokol_Event(evt *C.sapp_event, userdata uintptr) {
	cb := cgo.Handle(userdata).Value().(*wrappedUserData).event
	if cb != nil {
		cb(Event{internal: evt})
	}
}

func Run(opts *AppDesc) {
	runtime.LockOSThread()
	var titleStr *C.char
	if opts.WindowTitle != "" {
		titleStr = C.CString(opts.WindowTitle)
		defer C.free(unsafe.Pointer(titleStr))
	}

	var userDataObj *wrappedUserData
	if opts.Init != nil || opts.Frame != nil || opts.Cleanup != nil {
		userDataObj = &wrappedUserData{
			init:    opts.Init,
			frame:   opts.Frame,
			cleanup: opts.Cleanup,
			event:   opts.Event,
		}
	}

	var udata cgo.Handle
	if userDataObj != nil {
		udata = cgo.NewHandle(userDataObj)
		defer udata.Delete()
	}

	desc := C.sapp_desc{
		window_title:                 titleStr,
		width:                        C.int(opts.Width),
		height:                       C.int(opts.Height),
		alpha:                        C.bool(opts.Alpha),
		sample_count:                 C.int(opts.SampleCount),
		swap_interval:                C.int(opts.SampleInterval),
		high_dpi:                     C.bool(opts.HighDPI),
		fullscreen:                   C.bool(opts.Fullscreen),
		enable_clipboard:             C.bool(opts.EnableClipboard),
		clipboard_size:               C.int(opts.ClipboardSize),
		enable_dragndrop:             C.bool(opts.EnableDragnDrop),
		max_dropped_files:            C.int(opts.MaxDroppedFiles),
		max_dropped_file_path_length: C.int(opts.MaxDroppedFilePathLength),
		user_data:                    unsafe.Pointer(udata),
	}
	if opts.Init != nil {
		desc.init_userdata_cb = (C._sokolgo_void_callback)(C.cb_sokol_Init)
	}
	if opts.Frame != nil {
		desc.frame_userdata_cb = (C._sokolgo_void_callback)(C.cb_sokol_Frame)
	}
	if opts.Cleanup != nil {
		desc.cleanup_userdata_cb = (C._sokolgo_void_callback)(C.cb_sokol_Cleanup)
	}
	if opts.Event != nil {
		desc.event_userdata_cb = (C._sokolgo_event_callback)(C.cb_sokol_Event)
	}
	C.sapp_run(&desc)
}

// returns true after sokol-app has been initialized
func IsValid() bool { return bool(C.sapp_isvalid()) }

// returns the current framebuffer width in pixels
func Width() int { return int(C.sapp_width()) }

// returns the current framebuffer width in pixels
func Widthf() float32 { return float32(C.sapp_widthf()) }

// returns the current framebuffer height in pixels
func Height() int { return int(C.sapp_width()) }

// returns the current framebuffer height in pixels
func Heightf() float32 { return float32(C.sapp_heightf()) }

// get default framebuffer color pixel format
func ColorFormat() int { return int(C.sapp_color_format()) }

// get default framebuffer depth pixel format
func DepthFormat() int { return int(C.sapp_depth_format()) }

// get default framebuffer sample count
func SampleCount() int { return int(C.sapp_sample_count()) }

// returns true when high_dpi was requested and actually running in a high-dpi scenario
func HighDPI() bool { return bool(C.sapp_high_dpi()) }

// returns the dpi scaling factor (window pixels to framebuffer pixels)
func DPIScale() float32 { return float32(C.sapp_dpi_scale()) }

// show or hide the mobile device onscreen keyboard
func ShowKeyboard(show bool) { C.sapp_show_keyboard(C.bool(show)) }

// return true if the mobile device onscreen keyboard is currently shown
func KeyboardShown() bool { return bool(C.sapp_keyboard_shown()) }

// query fullscreen mode
func IsFullscreen() bool { return bool(C.sapp_is_fullscreen()) }

// toggle fullscreen mode
func ToggleFullscreen() { C.sapp_toggle_fullscreen() }

// show or hide the mouse cursor
func ShowMouse(show bool) { C.sapp_show_mouse(C.bool(show)) }

// show or hide the mouse cursor
func MouseShown() bool { return bool(C.sapp_mouse_shown()) }

// enable/disable mouse-pointer-lock mode
func LockMouse(lock bool) { C.sapp_lock_mouse(C.bool(lock)) }

// return true if in mouse-pointer-lock mode (this may toggle a few frames later)
func MouseLocked() bool { return bool(C.sapp_mouse_locked()) }

// set mouse cursor type
func SetMouseCursor(cursor MouseCursor) { C.sapp_set_mouse_cursor(C.sapp_mouse_cursor(cursor)) }

// get current mouse cursor type
func GetMouseCursor() MouseCursor { return MouseCursor(C.sapp_get_mouse_cursor()) }

// return the userdata pointer optionally provided in sapp_desc
// func UserData() unsafe.Pointer { return unsafe.Pointer(C.sapp_userdata()) }

// return a copy of the sapp_desc structure
func QueryDesc() AppDesc {
	desc := C.sapp_query_desc()
	return AppDesc{
		WindowTitle:              C.GoString(desc.window_title),
		Width:                    int(desc.width),
		Height:                   int(desc.height),
		HighDPI:                  bool(desc.high_dpi),
		Fullscreen:               bool(desc.fullscreen),
		SampleCount:              int(desc.sample_count),
		SampleInterval:           int(desc.swap_interval),
		Alpha:                    bool(desc.alpha),
		EnableClipboard:          bool(desc.enable_clipboard),
		ClipboardSize:            int(desc.clipboard_size),
		EnableDragnDrop:          bool(desc.enable_dragndrop),
		MaxDroppedFiles:          int(desc.max_dropped_files),
		MaxDroppedFilePathLength: int(desc.max_dropped_file_path_length),
	}
}

// initiate a "soft quit" (sends SAPP_EVENTTYPE_QUIT_REQUESTED)
func RequestQuit() { C.sapp_request_quit() }

// cancel a pending quit (when SAPP_EVENTTYPE_QUIT_REQUESTED has been received)
func CancelQuit() { C.sapp_cancel_quit() }

// initiate a "hard quit" (quit application without sending SAPP_EVENTTYPE_QUIT_REQUESTED)
func Quit() { C.sapp_quit() }

// call from inside event callback to consume the current event (don't forward to platform)
func ConsumeEvent() { C.sapp_consume_event() }

// get the current frame counter (for comparison with sapp_event.frame_count)
func FrameCount() uint64 { return uint64(C.sapp_frame_count()) }

// get an averaged/smoothed frame duration in seconds
func FrameDuration() float64 { return float64(C.sapp_frame_duration()) }

// write string into clipboard
func SetClipboardString(str string) { C.sapp_set_clipboard_string(tmpstring(str)) }

// read string from clipboard (usually during SAPP_EVENTTYPE_CLIPBOARD_PASTED)
func GetClipboardString() string { return C.GoString(C.sapp_get_clipboard_string()) }

// set the window title (only on desktop platforms)
func SetWindowTitle(str string) { C.sapp_set_window_title(tmpstring(str)) }

// set the window icon (only on Windows and Linux)
// func SetIcon(iconDesc *IconDesc) { C.sapp_set_icon((*C.sapp_icon_desc)(unsafe.Pointer(iconDesc))) }

// gets the total number of dropped files (after an SAPP_EVENTTYPE_FILES_DROPPED event)
func GetNumDroppedFiles() int { return int(C.sapp_get_num_dropped_files()) }

// gets the dropped file paths
func GetDroppedFilePath(index int) string {
	return C.GoString(C.sapp_get_dropped_file_path(C.int(index)))
}

// Metal: get bridged pointer to Metal device object
func MetalGetDevice() unsafe.Pointer { return C.sapp_metal_get_device() }

// Metal: get bridged pointer to MTKView's current drawable of type CAMetalDrawable
func MetalGetCurrentDrawable() unsafe.Pointer {
	return C.sapp_metal_get_current_drawable()
}

// Metal: get bridged pointer to MTKView's depth-stencil texture of type MTLTexture
func MetalGetDepthStencilTexture() unsafe.Pointer {
	return C.sapp_metal_get_depth_stencil_texture()
}

// Metal: get bridged pointer to MTKView's msaa-color-texture of type MTLTexture (may be null)
func MetalGetMSAAColorTexture() unsafe.Pointer {
	return C.sapp_metal_get_msaa_color_texture()
}

func tmpstring(s string) *C.char {
	if s == "" {
		return nil
	}
	p := make([]byte, len(s)+1)
	copy(p, s)
	return (*C.char)(unsafe.Pointer(&p[0]))
}
