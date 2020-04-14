package compute

type VirtualMachineRepository interface {
	List() ([]*VirtualMachine, error)
	Get(id, node string) (*VirtualMachine, error)
	Create(id, node string, arch Arch, vcpus int, memory Size, volumes []*VirtualMachineAttachedVolume, interfaces []*VirtualMachineAttachedInterface, config *VirtualMachineConfig) (*VirtualMachine, error)
	Delete(id, node string) error
	Update(id, node string, params VirtualMachineUpdateParams) error
	AttachVolume(machineId, node string, attachedVolume *VirtualMachineAttachedVolume) error
	DetachVolume(machineId, node, attachmentDeviceName string) error
	AttachInterface(id, node string, iface *VirtualMachineAttachedInterface) error
	DetachInterface(id, node, mac string) error
	GetConsoleStream(id, node string) (VirtualMachineConsoleStream, error)
	GetGraphicStream(id, node string) (VirtualMachineGraphicStream, error)
	Poweroff(id, node string) error
	Reboot(id, node string) error
	Start(id, node string) error
}

type VirtualMachineConsoleStream interface {
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	Close() error
}

type VirtualMachineGraphicStream VirtualMachineConsoleStream

type VirtualMachineState int

const (
	StateUnknown = VirtualMachineState(0)
	StateStopped = VirtualMachineState(1)
	StateRunning = VirtualMachineState(2)
)

func (state VirtualMachineState) String() string {
	switch state {
	default:
		return "unknown"
	case StateStopped:
		return "stopped"
	case StateRunning:
		return "running"
	}
}

type VirtualMachineCpuPin struct {
	Vcpus    map[uint][]uint
	Emulator []uint
}

type VirtualMachineConfig struct {
	Hostname string
	Keys     []*Key
	Userdata []byte
}

type VirtualMachineGraphic struct {
	Type   GraphicType
	Listen string
	Port   int
}

func (g VirtualMachineGraphic) Vnc() bool {
	return g.Type == GraphicTypeVnc
}

type VirtualMachine struct {
	Id         string
	NodeId     string
	VCpus      int
	Arch       Arch
	State      VirtualMachineState
	Memory     Size
	Interfaces []*VirtualMachineAttachedInterface
	Volumes    []*VirtualMachineAttachedVolume
	Config     *VirtualMachineConfig
	Cpupin     *VirtualMachineCpuPin
	GuestAgent bool
	Autostart  bool
	Graphic    VirtualMachineGraphic
}

func (vm *VirtualMachine) AttachmentInfo(path string) *VirtualMachineAttachedVolume {
	for _, attachedVolume := range vm.Volumes {
		if attachedVolume.Path == path {
			return attachedVolume
		}
	}
	return nil
}

func (vm *VirtualMachine) IpAddressList() []string {
	iplist := []string{}
	for _, iface := range vm.Interfaces {
		iplist = append(iplist, iface.IpAddressList...)
	}
	return iplist
}

func (vm *VirtualMachine) IsRunning() bool {
	return vm.State == StateRunning
}

func (vm *VirtualMachine) Disks() []*VirtualMachineAttachedVolume {
	disks := []*VirtualMachineAttachedVolume{}
	for _, volume := range vm.Volumes {
		if volume.DeviceType == DeviceTypeDisk {
			disks = append(disks, volume)
		}
	}
	return disks
}

func (vm *VirtualMachine) Cdroms() []*VirtualMachineAttachedVolume {
	cdroms := []*VirtualMachineAttachedVolume{}
	for _, volume := range vm.Volumes {
		if volume.DeviceType == DeviceTypeCdrom {
			cdroms = append(cdroms, volume)
		}
	}
	return cdroms
}

type VirtualMachineAttachedVolume struct {
	Path       string
	DeviceName string
	Type       VolumeType
	Format     VolumeFormat
	DeviceType DeviceType
	DeviceBus  DeviceBus
}

type VirtualMachineAttachedInterface struct {
	NetworkName   string
	Mac           string
	Model         string
	IpAddressList []string
	AccessVlan    uint
}
