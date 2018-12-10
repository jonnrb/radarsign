from golang:1.11.2 as build
add . /src
run cd /src && CGO_ENABLED=0 go get .

from gcr.io/distroless/base
copy --from=build /go/bin/radarsign /radarsign
entrypoint ["/radarsign"]
cmd ["-logtostderr", "-v=1"]
