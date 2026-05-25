# Kafka 3.5 本地认证环境

本目录按认证方式拆分 Kafka 3.5 本地启动配置，用于调试规则引擎的数据源、数据集成和数据输出。镜像使用 `bitnamilegacy/kafka:3.5.2-debian-12-r30`，因为 `bitnami/kafka:3.5.2` 这类短 tag 当前 Docker Hub 已经无法拉取。

| 目录                 | 外部地址          | 认证方式                               |
| -------------------- | ----------------- | -------------------------------------- |
| `plaintext`          | `localhost:19092` | 无认证                                 |
| `sasl-plain`         | `localhost:19093` | SASL/PLAIN                             |
| `sasl-scram-sha-256` | `localhost:19094` | SASL/SCRAM-SHA-256                     |
| `sasl-scram-sha-512` | `localhost:19095` | SASL/SCRAM-SHA-512                     |
| `oauthbearer`        | `localhost:19096` | SASL/OAUTHBEARER，本地 unsecured token |
| `tls`                | `localhost:19097` | TLS 加密，不要求客户端证书             |
| `mtls`               | `localhost:19098` | mTLS，要求客户端证书                   |
| `kafka-ui`           | `localhost:18080` | Kafka UI，默认连接 sasl-plain 环境     |

默认账号：

- SASL 用户：`hotgo`
- SASL 密码：`hotgo-secret`
- Broker 内部用户：`admin`
- Broker 内部密码：`admin-secret`
- TLS 证书密码：`changeit`

启动示例：

```bash
cd docker/kafka/plaintext
docker compose up -d
```

TLS/mTLS 首次启动前先生成证书：

```bash
cd docker/kafka/certs
bash generate-certs.sh
```

然后启动：

```bash
cd docker/kafka/tls
docker compose up -d
```

Kafka UI 默认连接 `sasl-plain` 里的 Kafka 内部监听地址，并使用管理账号 `admin` / `admin-secret`。先启动 sasl-plain Kafka，再启动 UI：

```bash
cd docker/kafka/sasl-plain
docker compose up -d

cd ../kafka-ui
docker compose up -d
```

访问 `http://localhost:18080`。如果要连接其他认证方式的 Kafka，需要同步调整 `kafka-ui/docker-compose.yaml` 里的外部网络、`KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS` 和认证属性。

本地发送测试数据：

```bash
cd agent
go run ./cmd/kafka-producer -topic hotgo-test-events -count 3
```

SCRAM 说明：`sasl-scram-sha-256` 和 `sasl-scram-sha-512` 使用 ZooKeeper 模式启动，并在 broker 启动前写入 SCRAM 用户凭据，避免 KRaft 本地初始化 SCRAM 用户不完整导致认证失败。

OAuth 说明：`oauthbearer` 使用 Kafka 自带的 unsecured OAuthBearer callback handler，只适合本地开发验证，不用于生产。
