install gateway plugin
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest google.golang.org/protobuf/cmd/protoc-gen-go@latest google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
generate gateway by
```  
protoc -Iproto --go_out=gateway/proto --go_opt=paths=source_relative --go-grpc_out=gateway/proto --go-grpc_opt=paths=source_relative --grpc-gateway_out=gateway/proto --grpc-gateway_opt=paths=source_relative proto/greet/v1/greet_gw.proto
```

generate open api
```
protoc -Iproto --openapiv2_out gateway/openapiv2 --openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true proto/greet/v1/greet_gw.proto
```