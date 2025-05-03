FROM fedora:latest

USER 0

EXPOSE 8080

RUN mkdir -p /opt/app && \
    cd /opt/app && \
    dnf install -y go libvirt-devel && \
    dnf clean all

WORKDIR /opt/app

COPY main.go /opt/app

RUN go mod init libvirt-monitor && \
    go get -u && \
    go mod tidy && \
    go build 

CMD /opt/app/libvirt-monitor
