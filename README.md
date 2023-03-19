generate for blog service by
```  
protoc -Iproto --go_out=blog/proto --go_opt=paths=source_relative --go-grpc_out=blog/proto --go-grpc_opt=paths=source_relative proto/blog/v1/blog.proto
```

generate for auth service by
```  
protoc -Iproto --go_out=auth/proto --go_opt=paths=source_relative --go-grpc_out=auth/proto --go-grpc_opt=paths=source_relative proto/auth/v1/auth.proto
```

generate gateway by
```  
protoc -Iproto --go_out=gateway/proto --go_opt=paths=source_relative \
--go-grpc_out=gateway/proto --go-grpc_opt=paths=source_relative \
--grpc-gateway_out=gateway/proto --grpc-gateway_opt=paths=source_relative \
 proto/auth/v1/auth.proto
 ```

generate open api
```
protoc -Iproto --openapiv2_out gateway/openapiv2 --openapiv2_opt logtostderr=true proto/blog/v1/blog.proto proto/auth/v1/auth.proto
```