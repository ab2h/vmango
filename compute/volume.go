package compute

type VolumeFormat int

const (
	FormatUnknown = VolumeFormat(0)
	FormatRaw     = VolumeFormat(1)
	FormatQcow2   = VolumeFormat(2)
)

func (format VolumeFormat) String() string {
	switch format {
	default:
		return "unknown"
	case FormatRaw:
		return "raw"
	case FormatQcow2:
		return "qcow2"
	}
}

type Volume struct {
	Type       string
	Path       string
	Size       uint64 // MiB
	Pool       string
	Format     VolumeFormat
	AttachedTo string
}

type VolumePool struct {
	Name string
	Size uint64 // MiB
	Used uint64 // MiB
	Free uint64 // MiB
}

func (pool *VolumePool) UsagePercent() int {
	return int(100 * pool.Used / pool.Free)
}

func (pool *VolumePool) FreeGB() uint64 {
	return pool.Free / 1024
}

type VolumeRepository interface {
	Get(path string) (*Volume, error)
	Create(pool, name string, format VolumeFormat, size uint64) (*Volume, error)
	Delete(path string) error
	List() ([]*Volume, error)
	Pools() ([]*VolumePool, error)
}
