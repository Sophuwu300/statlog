package units

const (
	B   = 1
	KiB = 1024 * B
	MiB = 1024 * KiB
	GiB = 1024 * MiB

	KB = 1000 * B
	MB = 1000 * KB
	GB = 1000 * MB
)

var Labels = map[int]string{
	B:   "B",
	KiB: "KiB",
	MiB: "MiB",
	GiB: "GiB",
	KB:  "KB",
	MB:  "MB",
	GB:  "GB",
}
