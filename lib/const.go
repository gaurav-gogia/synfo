package lib

// Global Constants
const (
	IMAGE   = "image"
	VIDEO   = "video"
	AUDIO   = "audio"
	ARCHIVE = "archive"
)

const (
	defaultBuffer    = 10 * 1024
	defaultModel     = "hog"
	defaultFt        = "image"
	defaultDiskImage = "evi.iso"
	mountinfoPath    = "/proc/self/mountinfo"
)

const (
	extcmduse = "Extracts specified type of files from target device."
	apdcmduse = "Runs an HoG/CNN based automated PoI Identification module."
	awdcmduse = "Runs a CNN based automated Weapon Detection module."

	srcflaghelp   = "Source path for special device block file."
	dstflaghelp   = "Destination directory for disk image."
	bsflaghelp    = "Buffer Size in bytes to be used during disk imaging."
	poiflaghelp   = "Directory of images with known faces."
	modelflaghelp = "ML Model type to be used for face detection {hog | cnn}."
	ftflaghelp    = "Type of file(s) to be extracted {image | audio | video | archive}."

	helpusageflag     = "Shows this help message."
	exampleusageflag  = "Shows example usage."
	flashusageflag    = "Boosts disk imaging speed."
	imageanalysisflag = "Runs analysis on existing disk image instead of cloning it."
)

// Global command name constants
const (
	APDCMD = "apd"
	AWDCMD = "awd"
	EXTCMD = "ext"
)

// Constant error strings
const (
	ErrUnableToMount = "unable to attach image, make sure that it's not already mounted"
)

// File headers and footers
var (
	PNG_HEADER = []byte{0x89, 0x50, 0x4e, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	PNG_FOOTER = []byte{0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

	JPG_HEADER = []byte{0xff, 0xd8, 0xff}
	JPG_FOOTER = []byte{0xff, 0xd9}

	GIF87A_HEADER = []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}
	GIF98A_HEADER = []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}
	GIF_FOOTER    = []byte{0x00, 0x3b}
)
