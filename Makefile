SHELL := /bin/bash

default: release

neutron: cmd/neutron/*.go
	go build -o neutron cmd/neutron/*.go

quasar: cmd/quasar/*.go
	go build -o quasar cmd/quasar/*.go

release: neutron quasar examples/nebula@.service
	mkdir -p release/linux-x86_64
	cp neutron release/linux-x86_64
	cp quasar release/linux-x86_64
	cp examples/nebula@.service release/linux-x86_64
	tar -cvf release/linux-x86_64.tar.gz release/linux-x86_64
