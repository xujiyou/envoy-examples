
# envoy-examples

学习 envoy 过程中的示例代码和文档。

构建环境为：https://www.getenvoy.io/
常用的运行命令为：
```shell script
$ sudo getenvoy run standard:1.14.1 -- --config-path ./envoy-config.yaml
```

- [stats](./stats) ： envoy 如何将统计信息传输到 Prometheus。 
- [metrics_service](./metrics_service)： envoy 如何自定义传输统计信息。
- [rls](./rls)： RLS 即 Rate limit service，Envoy 访问速率限制。
- [udp](./udp)： envoy 如何代理 UDP 流量。
- [jaeger](./jaeger)： envoy 集成 jaeger
- [buffer](./buffer)： 限制请求大小
- [cors](./cors)： 允许跨域访问
- [authz](./authz)： 外部认证
- [fault](./fault)： 故障注入
- [gzip](./gzip)： envoy 对响应数据进行 gzip 压缩
- [ip-tagging](./ip-tagging) 给客户端的访问的目标 IP 打上标签，然后用于统计。
- [jwt](./jwt)： envoy 实现 jwt 认证
- [lua](./lua) lua 脚本过滤器的使用 
- [http_inspector](./http_inspector) Listen 过滤器，实现对 HTTP 协议版本的统计
