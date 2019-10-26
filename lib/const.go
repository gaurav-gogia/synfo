package lib

// Global Constants
const (
	IMAGE = iota + 1
	VIDEO
	AUDIO
	ARCHIVE
)

const (
	defaultBuffer = 10 * 1024
	defaultModel  = "hog"
	mountinfoPath = "/proc/self/mountinfo"
	partfile      = ".part"
)

const (
	AUTOCMD = "AUTO"
	EXTCMD  = "EXTRACT"
)
