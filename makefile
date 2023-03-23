
#  TODO

docker:
	echo "not implemented"
# TODO include environment variables from host

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative apis/grpc/proto/IPBlockerService.proto