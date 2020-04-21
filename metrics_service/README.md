# metrics_service

在 [stats](../stats) 中提到了 Envoy 内置的 4 种接收器，这里再来测试一下 `envoy.stat_sinks.metrics_service`。

metrics_service 接收器可以让我们灵活处理 metrics，想怎么玩就怎么玩！！！！！！

修改 envoy 的配置文件，添加如下内容：

```yaml
stats_sinks:
  - name: envoy.stat_sinks.metrics_service
    config:
      grpc_service: 
        envoy_grpc:
          cluster_name: grpc-exporter
static_resources:      
  clusters:
  - name: grpc-exporter
    connect_timeout: 0.25s
    type: strict_dns
    lb_policy: ROUND_ROBIN
    http2_protocol_options: {}
    hosts:
    - socket_address:
        address: 127.0.0.1
        port_value: 10001
```

注意其中的 `http2_protocol_options` 不能漏掉！！！因为 gRPC 需要 HTTP2 协议。

服务端要实现这个 proto ：https://github.com/envoyproxy/data-plane-api/blob/master/envoy/service/metrics/v2/metrics_service.proto

其中的关键代码如下：

```protobuf
service MetricsService {
  rpc StreamMetrics(stream StreamMetricsMessage) returns (StreamMetricsResponse) {}
}
```

可以看到是一个客户端流的 gRPC 方法。

编译后的 go 代码在 `go-control-plane` 项目的 `envoy/service/metrics/v2/metrics_service.pb.go` 中。

其中 `MetricsServiceServer` 如下：

```go
type MetricsServiceServer interface {
	StreamMetrics(MetricsService_StreamMetricsServer) error
}
```



下面来写服务端代码：

创建项目 `metrics_service`，初始化模块：

```bash
$ go mod init metrics_service
```

写 server.go :

```go
package main

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v2"
	"io"
	"log"
)

type MyMetricsServer struct {}

func (myServer *MyMetricsServer) StreamMetrics(server v2.MetricsService_StreamMetricsServer) error {
	for {
		message, err := server.Recv()
		if err == io.EOF {
			_ = server.SendAndClose(&v2.StreamMetricsResponse{})
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		log.Println(message)
	}
}
```

注意这里的 `message`，这就是 envoy 传过来的统计信息，这里可以对这个信息做任意处理，比如传到 Fluentd、Logstash、Kafka、或者 Elasticsearch 都是可以的，我这里就是打印到控制台了。

然后写 main.go :

```go
package main

import (
	"flag"
	"fmt"
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v2"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 10001))
	log.Println("metrics server listen to 10001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	v2.RegisterMetricsServiceServer(grpcServer, &MyMetricsServer{})
	_ = grpcServer.Serve(lis)
}
```

在服务端编译并启动:

```bash
$ go build
$ ./metrics_service
```


在另外一个命令行中启动 envoy，之后再对服务进行访问，会发现在 metrics_service 的命令行中打印了一堆日志，这些日志就是 envoy 的统计信息了！