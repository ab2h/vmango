package compute

import (
	"fmt"
)

type EventVirtualMachineCreated struct {
	vm *VirtualMachine
}

func NewEventVirtualMachineCreated(vm *VirtualMachine) *EventVirtualMachineCreated {
	return &EventVirtualMachineCreated{vm: vm}
}

func (e *EventVirtualMachineCreated) Name() string {
	return "vm_created"
}

func (e *EventVirtualMachineCreated) Plain() map[string]string {
	data := map[string]string{
		"event":              e.Name(),
		"vm_id":              e.vm.Id,
		"vm_cpus":            fmt.Sprintf("%d", e.vm.VCpus),
		"vm_memory_mib":      fmt.Sprintf("%d", e.vm.MemoryMiB()),
		"vm_volume_count":    fmt.Sprintf("%d", len(e.vm.Volumes)),
		"vm_interface_count": fmt.Sprintf("%d", len(e.vm.Interfaces)),
	}
	for idx, volume := range e.vm.Volumes {
		data[fmt.Sprintf("vm_volume_%d_path", idx)] = volume.Path
		data[fmt.Sprintf("vm_volume_%d_format", idx)] = volume.Format.String()
		data[fmt.Sprintf("vm_volume_%d_device", idx)] = volume.Device.String()
		data[fmt.Sprintf("vm_volume_%d_type", idx)] = volume.Type.String()
	}
	for idx, iface := range e.vm.Interfaces {
		data[fmt.Sprintf("vm_interface_%d_mac", idx)] = iface.Mac
		data[fmt.Sprintf("vm_interface_%d_network", idx)] = iface.NetworkName
		data[fmt.Sprintf("vm_interface_%d_type", idx)] = iface.NetworkType.String()
	}
	return data
}
