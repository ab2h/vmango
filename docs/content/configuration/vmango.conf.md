+++
weight = 15
title = "vmango.conf"
date = "2017-02-05T17:59:37+03:00"
toc = true
+++

Main configuration file in HCL format: https://github.com/hashicorp/hcl

All vm configuration templates are golang text templates: https://golang.org/pkg/text/template/

## Web server options

**listen** - Web server listen address

**session_secret** - Secret key for session cookie encryption.

**static_cache** - Static files cache duration, e.g: "1d", "10m", "60s". Used mainly for development.

**trusted_proxies** - List of trusted ip addresses for X-Forwarded-For or X-Real-IP headers processing.

**ssl_key** - Path to private key file

**ssl_cert** - Path to SSL certificate

## Hypervisor

{{% notice info %}}
If you use SSH connection url, make sure it present in known_hosts file and remote user has permissions to access libvirt socket.
{{% /notice %}}

{{% notice tip %}}
Libvirt socket path can be changed via ?socket=/path/to/libvirt-sock url option.
{{% /notice %}}

**hypervisor** - Hypervisor definition, may be specified multiple times.

**hypervisor.url** - Libvirt connection URL.

**hypervisor.image_storage_pool** - Libvirt storage pool name for VM images.

**hypervisor.root_storage_pool** - Libvirt storage pool name for root disks.

**hypervisor.ignored_vms** - List of ignored virtual machines names.

**hypervisor.network** - Libvirt network name or bridge name if network_script specified.

**hypervisor.network_script** - Path to executable file to integrate vmango with network not managed by libvirt (see [network]({{% ref "hypervisor.md#scripted-network" %}}) chapter).

**hypervisor.vm_template** - Path to go template file (relative to vmango.conf) with libvirt domain XML. Used to create a new machine.

Execution context:

* Machine    [VirtualMachine](https://github.com/subuk/vmango/blob/master/src/vmango/domain/vm.go#L57)
* Image  - [Image](https://github.com/subuk/vmango/blob/master/src/vmango/domain/image.go#L16)
* Plan       [Plan](https://github.com/subuk/vmango/blob/master/src/vmango/domain/plan.go#L3)
* VolumePath string
* Network    string

**hypervisor.volume_template** - Path to go template file (relative to vmango.conf) with libvirt volume XML. It is used to create a new root volume for a new machine.

Execution context:

* Machine    [VirtualMachine](https://github.com/subuk/vmango/blob/master/src/vmango/domain/vm.go#L57)
* Image  - [Image](https://github.com/subuk/vmango/blob/master/src/vmango/domain/image.go#L16)
* Plan       [Plan](https://github.com/subuk/vmango/blob/master/src/vmango/domain/plan.go#L3)

Example:

    hypervisor "LOCAL1" {
        url = "qemu:///system"
        image_storage_pool = "vmango-images"
        root_storage_pool = "default"
        network = "vmango"
        vm_template = "vm.xml.in"
        volume_template = "volume.xml.in"
    }


## Authentication users

More about users see in [Users]({{% ref "auth_users.md" %}}) section.

**user** - user definition, may be used multiple times.

**user.password** - password in bcrypt hashed form.

Example:

    user "admin" {
        password = "$2a$10$..."
    }

## Plans

Plans are limiting availaible hardware resources for machines.

**plan** - plan definition, may be used multiple times.

**plan.memory** - memory limit in megabytes.

**plan.cpus** - cpu count.

**disk_size** - disk size in gigabytes.

Example:

    plan "small" {
        memory = 512
        cpus = 1
        disk_size = 5
    }

## SSH Keys

List of availaible ssh keys.

**ssh_key** - key definition, may be used multiple times.

**ssh_key.public** - full public key in ssh format.

Example:

    ssh_key "test" {
        public = "ssh-rsa AAAAB3NzaC1y..."
    }
