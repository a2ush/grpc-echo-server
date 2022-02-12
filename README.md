## Demo with grcp_cli

List "Service".
```
$ grpc_cli ls localhost:50051                  
grpc.reflection.v1alpha.ServerReflection
rpc.Echo
rpc.Greating
```

Detail
```
$ grpc_cli ls localhost:50051 -l
filename: reflection/grpc_reflection_v1alpha/reflection.proto
package: grpc.reflection.v1alpha;
service ServerReflection {
  rpc ServerReflectionInfo(stream grpc.reflection.v1alpha.ServerReflectionRequest) returns (stream grpc.reflection.v1alpha.ServerReflectionResponse) {}
}

filename: grpc-echo-server.proto
package: rpc;
service Echo {
  rpc Reply(rpc.ClientRequest) returns (rpc.ServerResponse) {}
}

filename: grpc-echo-server.proto
package: rpc;
service Greating {
  rpc Reply(rpc.GreatingClientRequest) returns (rpc.ServerGreatingResponse) {}
}
```
You can know "rpc.Echo" and "rpc.Greating" service have "Reply()" method.

Confirm type.
```
$ grpc_cli type localhost:50051 rpc.ClientRequest                 
message ClientRequest {
  string message = 1 [json_name = "message"];
}

$ grpc_cli type localhost:50051 rpc.GreatingClientRequest         
message GreatingClientRequest {
  .rpc.Format.Greeting client_greeting = 1 [json_name = "clientGreeting"];
}
```

Test
```
$ grpc_cli call localhost:50051 rpc.Echo.Reply "message: 'hello'"
connecting to localhost:50051
name: "EchoBot"
message: "hello by EchoBot"
create_time {
  seconds: 1644665166
  nanos: 960671000
}
Rpc succeeded with OK status

$ grpc_cli call localhost:50051 rpc.Greating.Reply "client_greeting: 1"
connecting to localhost:50051
name: "GreetingBot"
format {
  echo: "How are you?"
}
create_time {
  seconds: 1644665114
  nanos: 613346000
}
Rpc succeeded with OK status
```

## How to generate a golang script from a proto file

.proto file
```
option go_package = "./;rpc";
```

```
ex1)
$ protoc \
  -Irpc \
  --go_out=plugins=grpc:./rpc \
  rpc/*.proto

$ tree ./
./
└── proto
    ├── grpc-echo-server.pb.go
    └── grpc-echo-server.proto
    
ex2)
$ protoc \ 
  -Irpc \
  --go_out=plugins=grpc:./api \
  rpc/*.proto

$ tree ./
./
├── api
│   └── grpc-echo-server.pb.go
└── rpc
    └── grpc-echo-server.proto
```
You should specify the directory where proto file is after `-I` option.

## How to implement

Use `RegisterEchoServer()` and `RegisterGreatingServer`. <br>
At this time, you must pass `server` and `Struct` to args.
```
type EchoServerHandler struct{}

type GreetingServerHandler struct {
	name string
}

func NewServerHandler() *GreetingServerHandler {
	return &GreetingServerHandler{
		name: "GreetingBot",
	}
}

func main() {
...
	rpc.RegisterEchoServer(
		server,
		&EchoServerHandler{},
	)

	rpc.RegisterGreatingServer(
		server,
		NewServerHandler(),
	)
```