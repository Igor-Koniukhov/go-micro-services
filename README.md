

- for css in mail-service
```bash
go get github.com/vanng822/go-premailer/premailer
```
```bash
go get github.com/xhit/go-simple-mail/v2

```
- For gRPC install
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27

```
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

[Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)

- For Linux
```bash
apt install -y protobuf-compiler
protoc --version 
```

- For generating go code execute command placed below in the folder with proto file
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
 
```

 - define the protocol (proto file), compile code, writhe the client code, write the server code
