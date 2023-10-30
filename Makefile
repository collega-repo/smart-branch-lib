generate_proto:
	protoc --proto_path=proto --go_out=commons --go_opt=paths=source_relative proto/error.proto
