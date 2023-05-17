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
