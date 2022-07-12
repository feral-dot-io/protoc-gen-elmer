.PHONY: all end-to-end readme real-world
all: example end-to-end readme real-world

example:
	protoc --elmer_out=. --elmer-fuzzer_out=. --elmer-twirp_out=. example.proto

end-to-end:
	make -C end-to-end
readme:
	make -C readme
real-world:
	make -C real-world
