lint:
	staticcheck ./*go

buid:
	go build

local-dev: # a messy-but-hey-it-works way to run it on my raspian pi for testing
	rm -f goruuvitag ; go build && sudo setcap 'cap_net_raw,cap_net_admin=eip' goruuvitag && ./goruuvitag

test:
	go test
