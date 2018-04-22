package framebuffer

// FixScreenInfo  Device independent unchangeable information about a frame buffer device and a specific video mode
type FixScreenInfo struct {
	ID         [16]byte  // Identification string eg "TT Builtin"
	SmemStart  uintptr   // Start of frame buffer mem
	SmemLen    uint32    // Length of frame buffer mem
	Type       uint32    // FB_TYPE_
	TypeAux    uint32    // Interleave for interleaved Planes
	Visual     uint32    // FB_VISUAL_
	Xpanstep   uint16    // Zero if no hardware panning
	Ypanstep   uint16    // Zero if no hardware panning
	Ywrapstep  uint16    // Zero if no hardware ywrap
	LineLength uint32    // Length of a line in bytes
	MmioStart  uintptr   // Start of Memory Mapped I/O (physical address)
	MmioLen    uint32    // Length of Memory Mapped I/O
	Accel      uint32    // Type of acceleration available
	Reserved   [3]uint16 // Reserved for future compatibility
}

// BitField for the color
type BitField struct {
	Offset   uint32 // Beginning of bitfield
	Length   uint32 // Length of bitfield
	MsbRight uint32 // != 0 : Most significant bit is right
}

// VarScreenInfo contains device independent changeable information about a frame buffer device and a specific video mode.
type VarScreenInfo struct {
	Xres, Yres,
	XresVirtual, YresVirtual,
	Xoffset, Yoffset,
	BitsPerPixel, Grayscale uint32
	Red, Green, Blue, Alpha BitField
	Nonstd, Activate,
	Height, Width,
	AccelFlags, Pixclock,
	LeftMargin, RightMargin, UpperMargin, LowerMargin,
	HsyncLen, VsyncLen, Sync,
	Vmode, Rotate, Colorspace uint32
	Reserved [4]uint32
}

// PixelFormat of the framebuffer
type PixelFormat int

const (
	// UnknownPixelFormat for when detection is not (yet) implemented
	UnknownPixelFormat PixelFormat = iota
	// BGR565 is a 16bit pixelformat
	BGR565
)

// PixelFormat returns the detected pixelformat
func (s *VarScreenInfo) PixelFormat() PixelFormat {
	switch s.BitsPerPixel {
	case 16:
		if s.Blue.Offset == 0 &&
			s.Blue.Length == 5 &&
			s.Green.Offset == 5 &&
			s.Green.Length == 6 &&
			s.Red.Offset == 11 &&
			s.Red.Length == 5 &&
			s.Alpha.Length == 0 {
			return BGR565
		}
	}
	return UnknownPixelFormat
}
