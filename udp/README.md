## 代理 UDP

官方教程：https://www.envoyproxy.io/docs/envoy/latest/configuration/listeners/udp_filters/udp_proxy

目前 udp 代理还处于 alpha 阶段。

先启动一个 UDP 服务端：

```bash
$ nc -l -u -p 1235
```

- -l ：表示 listener ，监听
- -u：UDP
- -p：指定端口

编写 envoy 的配置文件：

```yaml
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address:
      protocol: TCP
      address: 127.0.0.1
      port_value: 9901

static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address:
        protocol: UDP
        address: 127.0.0.1
        port_value: 1234
    reuse_port: true
    access_log:
      name: envoy.file_access_log
      typed_config:
        "@type": type.googleapis.com/envoy.config.accesslog.v2.FileAccessLog
        path: /dev/stdout
    listener_filters:
      name: envoy.filters.udp_listener.udp_proxy
      typed_config:
        '@type': type.googleapis.com/envoy.config.filter.udp.udp_proxy.v2alpha.UdpProxyConfig
        stat_prefix: service
        cluster: service_udp
  clusters:
  - name: service_udp
    connect_timeout: 0.25s
    type: STATIC
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: service_udp
      endpoints:
      - lb_endpoints:
        - endpoint:
            address:
              socket_address:
                address: 127.0.0.1
                port_value: 1235
```

在另外一个命令行启动：

```bash
$ sudo getenvoy run standard:1.14.1 -- --config-path ./envoy-config.yaml
```

在第三个命令行测试：

```bash
$ nc -u 127.0.0.1 1234
hello
world
```

输入一些信息之后，发现第一个命令行窗口也打印了相同的信息，说明代理成功！！！