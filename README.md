# libvirt Monitor (go)

Connects to the libvirt socket and provides a simple web dashboard on port 8080

## Local build

```bash
go mod init libvirt-monitor
go get -u
go mod tidy
git status
go build

# run
./libvirt-monitor
```

## Podman build

```bash
sudo podman login quay.io
sudo podman build . -t libvirt-monitor
sudo podman tag libvirt-monitor quay.io/kmendez/libvirt-monitor:latest
sudo podman push quay.io/kmendez/libvirt-monitor:latest

# Test it
sudo podman run -u0 --privileged -v /var/run/libvirt/virtqemud-sock:/var/run/libvirt/virtqemud-sock:Z -p 8080:8080 quay.io/kmendez/libvirt-monitor:latest
```

## Connect to libvirt via go locally

```bash
podman run -it --rm --privileged -v /var/run/libvirt/virtqemud-sock:/var/run/libvirt/virtqemud-sock:Z fedora:latest /bin/bash
dnf install -y libvirt-devel go
mkdir test
cd test
go mod init libvirt-monitor
go get -u
go mod tidy
go build
./go-libvirt
```

## Use container within truenase scale
```bash
docker run -d --name libvirt-monitor -u0 --privileged -v /run/truenas_libvirt/libvirt-sock-ro:/var/run/libvirt/virtqemud-sock -p 8888:8080  quay.io/kmendez/libvirt-monitor:latest
```

### virsh commands

```bash
virsh -c "qemu+unix:///system?socket=/run/truenas_libvirt/libvirt-sock-ro" dommemstat 6_coreos
virsh -c "qemu+unix:///system?socket=/run/truenas_libvirt/libvirt-sock-ro" domstats | more
```