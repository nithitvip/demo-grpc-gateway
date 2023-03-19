generate by
```  
protoc -Iproto --go_out=auth/proto --go_opt=paths=source_relative --go-grpc_out=auth/proto --go-grpc_opt=paths=source_relative proto/auth/v1/auth.proto
```