# Connect to libvirt via go

podman run -it --rm --privileged -v /var/run/libvirt/virtqemud-sock:/var/run/libvirt/virtqemud-sock:Z fedora:latest /bin/bash
dnf install -y libvirt-devel go
mkdir test
cd test
go mod init go-libvirt
cat > main.go
go get github.com/libvirt/libvirt-go
go build
./go-libvirt
