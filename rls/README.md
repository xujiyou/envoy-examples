## Rate limit service

用于处理 Envoy 的访问速率限制。

需要配置并运行全局速率限制服务，代码地址：https://github.com/envoyproxy/ratelimit

将这份代码下载到服务器，然后进行编译：

```bash
$ make bootstrap
$ make compile
```

这会在当前目录生成一个 bin 目录，这个 bin 目录里面有三个可执行文件，分别是：`ratelimit`、`ratelimit_client`、`ratelimit_config_check`

修改配置文件 `./examples/ratelimit/config/config.yaml`

```yaml
---
domain: rate_per_ip
descriptors:
  - key: generic_key
    value: default
    rate_limit:
      unit: minute
      requests_per_unit: 1
```

这里为了方便看到效果，设置了一分钟只接受一次请求。

ratelimit_config_check 用于检查配置文件是否正确，比如这样：

```bash
$ ./bin/ratelimit_config_check -config_dir ./examples/ratelimit/config/
```

然后在当前目录启动 docker-compose：

```bash
$ sudo cp bin/ratelimit /usr/local/bin/
$ docker-compose up
```

使用客户端工具测试：

```bash
$ ./bin/ratelimit_client -domain rate_per_ip -descriptors generic_key=default
```

返回如下：

```
dial string: localhost:8081
domain: rate_per_ip
descriptors: entries:<key:"generic_key" value:"default" > 
response: overall_code:OK statuses:<code:OK current_limit:<requests_per_unit:1 unit:SECOND > >
```





修改 envoy 中的配置，我的全部相关配置如下：

```yaml
static_resources:
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 82
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager
          codec_type: auto
          stat_prefix: ingress_http
          access_log:
            name: envoy.file_access_log
            typed_config:
              "@type": type.googleapis.com/envoy.config.accesslog.v2.FileAccessLog
              path: /dev/stdout
          route_config:
            name: local_route
            virtual_hosts:
            - name: service
              domains:
                - "*"
              routes:
                - match:
                    prefix: "/hello"
                  route:
                    cluster: hello-service
              rate_limits:
                - stage: 0
                  actions:
                    - {generic_key: {descriptor_value: "default"}}
          http_filters:
          - name: envoy.filters.http.ratelimit
            typed_config:
              "@type": type.googleapis.com/envoy.config.filter.http.rate_limit.v2.RateLimit
              domain: rate_per_ip
              stage: 0
              rate_limit_service:
                grpc_service:
                  envoy_grpc:
                    cluster_name: rate_limit_cluster
          - name: envoy.filters.http.router
  clusters:
  - name: rate_limit_cluster
    type: strict_dns
    connect_timeout: 0.25s
    lb_policy: round_robin
    http2_protocol_options: {}
    hosts:
    - socket_address:
        address: 127.0.0.1
        port_value: 8081
```

需要注意上面的 `rate_limits` 、 `envoy.filters.http.ratelimit` 。关于 `http_filters` ，必须以 `- name: envoy.filters.http.router` 结尾。`actions` 和 `domain` 中的配置，要和上面的 ratelimit 服务端的配置一样！

另外还要注意其中的 `typed_config` 和 `@type` 也是有讲究的！

启动 Envoy。

再次对服务进行访问，会发现一分钟内的请求，只有一个是成功的，其他都返回 429 错误，429 的意思就是太多请求的意思。