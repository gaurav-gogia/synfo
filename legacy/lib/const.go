package lib

// Global Constants
const (
	IMAGE = iota + 1
	VIDEO
	AUDIO
	ARCHIVE
)

// Internal constants
const (
	mountinfoPath = "/proc/self/mountinfo"
	partfile      = ".part"
)
