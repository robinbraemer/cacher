gen:
	protoc -I cacher/ cacher/cacher.proto --go_out=plugins=grpc:cacher