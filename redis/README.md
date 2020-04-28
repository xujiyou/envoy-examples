# Reids 流量代理

Envoy 配置：

```yaml
admin:
  access_log_path: /dev/stdout
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 8081

node:
  cluster: hello-service
  id: node1

static_resources:
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 82
    filter_chains:
    - filters:
      - name: envoy.filters.network.redis_proxy
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.redis_proxy.v2.RedisProxy
          stat_prefix: redis-stat
          settings:
            op_timeout: 3s
          prefix_routes:
            routes:
              - prefix: ""
                cluster: redis
  clusters:
  - name: redis
    connect_timeout: 0.25s
    type: strict_dns
    lb_policy: round_robin
    hosts:
    - socket_address:
        address: 127.0.0.1
        port_value: 6379
    health_checks:
      timeout: 1s
      interval: 10s
      unhealthy_threshold: 9
      healthy_threshold: 0
      custom_health_check:
        name: envoy.health_checkers.redis
        typed_config:
          "@type": type.googleapis.com/envoy.config.health_checker.redis.v2.Redis
          key: foo
```

注意

- `envoy.filters.network.redis_proxy` 过滤器必须是最后一个过滤器
- `settings.op_timeout` 是必须要配置
- 这里加入了对 Redis 的健康检查，如何设置了 `key` 则使用 ` EXISTS <key>` 进行检查，如果结果为 0 则成功。如果没设置这个 `key` ，则使用 `PING` 命令进行检查。

测试访问：

```bash
$ redis-cli -h 127.0.0.1 -p 82
```

查看健康检查统计结果：
```
cluster.redis.health_check.attempt: 2
cluster.redis.health_check.degraded: 0
cluster.redis.health_check.failure: 0
cluster.redis.health_check.healthy: 1
cluster.redis.health_check.network_failure: 0
cluster.redis.health_check.passive_failure: 0
cluster.redis.health_check.success: 2
cluster.redis.health_check.verify_cluster: 0
```

