FROM golang:1.13.8-alpine
RUN mkdir -p /go/src/github.com/hightoxicity/k8s-rook-flexvolume-reconciler
WORKDIR /go/src/github.com/hightoxicity/k8s-rook-flexvolume-reconciler
COPY . ./
RUN ls -al /go/src/github.com/hightoxicity/k8s-rook-flexvolume-reconciler
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s -v -extldflags -static" -a main.go
ENV ROOTFS /EXTRAROOTFS
RUN mkdir -p ${ROOTFS}

FROM scratch
COPY --from=0 /go/src/github.com/hightoxicity/k8s-rook-flexvolume-reconciler/main /k8s-rook-flexvolume-reconciler
CMD ["/k8s-rook-flexvolume-reconciler"]
