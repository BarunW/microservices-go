protogo:
	protoc -I protos/ protos/currency.proto --go_out=:.

protogrpc:
	protoc -I protos/ protos/currency.proto --go-grpc_out=:.

proto: 
	protoc -I protos/ protos/currency.proto --go_out=:. && protoc -I protos/ protos/currency.proto --go-grpc_out=:.
