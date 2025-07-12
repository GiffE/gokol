<p align="center">
    <img src="assets/logo_full_large.png" style="width: 60%" /><br/><br/>Simple
    cross-platform bindings for the
    <a href="https://github.com/floooh/sokol">SOKOL</a>
     libraries in Go.<br/><br/>
</p>

# Gokol

This package provides Go bindings for the Sokol C libraries ([https://github.com/floooh/sokol](https://github.com/floooh/sokol)). Sokol is a cross-platform, C-based library for game development, offering graphics, audio, and input functionality.

## Status

This project is currently under development. While parts may be functional, it's still in an early stage and may contain bugs or experience breaking changes. The bindings currently include headers from the 6-29-2025 release. Not all Sokol functions are yet implemented. Feel free to submit issues or PR's for missing features.

### Progress

- `sokol_app.h`
  - 🚧 macOS
  - iOS
  - Windows
  - Android
  - HTML5
  - Linux
- `sokol_gfx.h`

## Installation

```bash
go get -u github.com/GiffE/gokol
```
