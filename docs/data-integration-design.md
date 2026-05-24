# 数据集成设计说明

本文记录数据集成功能当前阶段的表结构边界、字段采集策略和清洗规则约定，供后续生成代码和开发业务逻辑时参考。

## 设计边界

- 数据集成放在主模块，不做插件。
- HotGo 当前项目不是多租户项目，数据集成表不包含 `tenant_id`。
- SQL 单独放在 `server/storage/data/data.sql`，不合入 `hotgo.sql`，也不设置 `AUTO_INCREMENT = 10000`。
- 不手写或修改 GoFrame 生成目录，包括 `internal/dao/internal`、`internal/model/entity`、`internal/model/do`。
- 不保存全量原始数据。Agent 可以在本地解析原始数据，但 server 只沉淀字段集合、配置和分钟级清洗统计。

## 表结构职责

当前第一版保留 5 张表。

### `hg_data_connector`

连接配置表。输入源和输出目标都放在同一张表，通过 `direction` 区分：

- `source`：输入连接，例如 Kafka、S3、HTTP、日志、手动输入。
- `sink`：输出连接，例如 HTTP Webhook、Kafka、飞书、S3、Elasticsearch、OpenSearch、Splunk、日志。

字段边界：

- `config` 是连接级配置，例如 broker、endpoint、认证信息、AK/SK、token、webhook 等。
- 敏感字段后续业务开发时必须加密存储，编辑回显时做脱敏。

### `hg_data_clean_task`

清洗任务表。任务运行后默认采集字段；只有 `clean_enabled = 1` 时才执行字段清洗规则。

字段边界：

- `source_id` 指向 `hg_data_connector.id`，业务层必须校验对应连接为 `direction = source`。
- `template_id` 指向 `hg_data_field_template.id`，用于标准化/归一化，可为空。
- `source_config` 是任务级输入配置，例如 Kafka topic、consumer group、S3 bucket/prefix、HTTP path、解析方式等。
- `clean_config` 是字段清洗规则，例如字段选择、删除、重命名、类型转换、脱敏、模板映射等。
- `sink_config` 是任务级输出配置，例如输出到哪些 sink 连接、串行/并行、失败策略、目标 topic/index/path、消息模板等。
- `clean_config_version` 用于 Agent 判断清洗规则是否需要重新拉取。
- `field_schema_version` 用于 Agent 判断 server 端字段集合是否有变化。
- `unknown_field_policy` 约定为：
  - `selected_only`：默认策略。新字段不会自动进入输出，只按已选择字段清洗。
  - `strict`：发现未知字段时阻塞或失败，等待用户确认。
  - `passthrough`：未知字段原样透传。
  - `drop`：未知字段丢弃。
- `clean_error_policy` 约定为：
  - `skip`：记录失败并跳过当前数据。
  - `keep_raw`：清洗失败时保留原值继续输出。
  - `stop_task`：停止任务。

### `hg_data_field`

字段集合表，是清洗配置页面字段列表的主要来源。

一条记录表示：

```text
某个清洗任务 + 某个字段路径
```

唯一键为：

```text
task_id + field_path
```

它不是按数据条数增长，而是按字段路径数量增长。例如每天 1 亿条日志，如果最终只发现 300 个字段路径，那么该任务下大约就是 300 条字段记录。

字段边界：

- `field_path` 是字段路径，例如 `user.id`、`items[].price`。
- `field_name` 是展示名，默认可取字段路径最后一段，也允许用户后续修改。
- `field_type` 是当前主类型或最近识别类型。
- `field_types` 保存出现过的类型集合，例如 `["string","number"]`。
- `sample_value` 只保存单字段样例值，不保存整条原始数据；样例必须截断和脱敏。
- `first_seen_at` / `last_seen_at` 用于展示字段首次/最近发现时间。

### `hg_data_field_template`

标准字段模板表，用于日志归一化处理。它不是 Agent 自动发现字段的来源，而是标准目标结构。

典型用途：

- 定义 `nginx.access`、`auth.login.failed`、`security.alert` 等标准事件字段。
- 清洗页面中把 `hg_data_field` 的实际字段映射到 `hg_data_field_template.fields` 的标准字段。
- `examples` 存标准样例，不存业务原始数据。

标准模板不应要求用户直接手写 JSON。`fields` 是后端存储格式，后台应提供表单化维护能力。

建议提供 4 种填入方式：

1. 系统内置模板。初始化数据或后续种子数据中预置常见模板，例如 `nginx.access`、`auth.login`、`security.alert`、`app.log`。用户可以复制内置模板后再调整。
2. 从已发现字段生成。用户在某个清洗任务的 `hg_data_field` 字段列表中勾选字段，点击“生成标准模板”，系统创建初始模板后让用户修改字段名、类型、必填项和说明。
3. 手工新增标准字段。页面上使用表格表单维护标准字段，不暴露 JSON 编辑器作为主要入口。字段包括字段编码、字段名称、类型、是否必填、默认值和说明。
4. 粘贴标准样例生成。用户粘贴一段标准 JSON 样例，前端或后端解析字段路径，生成初始模板字段，再由用户调整。

模板 `fields` 的存储结构可以按如下格式：

```json
[
  {
    "name": "src_ip",
    "title": "源IP",
    "type": "string",
    "required": true,
    "default": "",
    "description": "请求来源IP"
  },
  {
    "name": "http_status",
    "title": "HTTP状态码",
    "type": "number",
    "required": false,
    "default": 0,
    "description": "响应状态码"
  }
]
```

模板只定义“清洗后应该长什么样”。实际字段到标准字段的映射关系必须存在 `hg_data_clean_task.clean_config` 中，而不是存在模板表中。

例如：

```text
hg_data_field.field_path = client_ip
hg_data_field_template.fields[].name = src_ip
clean_config 负责记录 client_ip -> src_ip
```

### `hg_data_clean_stat`

分钟级清洗统计表。它不是逐条数据日志，而是按任务、输入连接、事件类型和分钟窗口聚合处理结果。

聚合粒度：

```text
task_id + connector_id + event_type + window_start
```

也就是同一个任务、同一个输入连接、同一个事件类型在同一分钟内只写一条统计记录。

字段边界：

- `total_count`：窗口内总处理数。
- `success_count`：成功输出数。
- `failed_count`：通用失败数，例如下游发送失败、系统异常等。
- `clean_failed_count`：清洗失败数，例如必填字段缺失、类型转换失败、规则执行失败等。
- `dropped_count`：按规则丢弃数。
- `last_error`：最近一次错误信息。
- `error_samples`：少量错误样例，不存 raw payload。建议最多保留 10 条，记录错误类型、字段路径和错误说明。

`error_samples` 示例：

```json
[
  {
    "errorType": "missing_required_field",
    "fieldPath": "user.id",
    "message": "必填字段缺失"
  },
  {
    "errorType": "type_convert_failed",
    "fieldPath": "price",
    "message": "无法转换为 number"
  }
]
```

这张表最多保留最近 100 万条记录，避免高吞吐场景下无限增长。该约束由后续定时清理任务或写入后的清理逻辑实现，表结构只提供按 `id`、`created_at` 和 `window_start` 查询清理所需的索引。

如需排查字段来源，看 `hg_data_field`；如需排查配置，看 `hg_data_clean_task.clean_config`；如需看运行趋势，看 `hg_data_clean_stat`。

清洗统计必须由 Agent 端先聚合后上报。Agent 不应逐条数据上报处理结果。

推荐流程：

1. Agent 本地按 `task_id + connector_id + event_type + window_start` 维护分钟级计数器。
2. 每条数据处理完成后只更新本地计数器。
3. 每个窗口结束后 1-5 秒内 flush 一次统计数据到 server。
4. 如果上报失败，Agent 本地保留该窗口统计并重试。
5. Server 收到统计后按唯一键累加合并到 `hg_data_clean_stat`。
6. 多个 Agent 同时上报同一分钟窗口时，server 只做累加，不覆盖已有计数。

Agent 上报示例：

```json
{
  "taskId": 1,
  "connectorId": 10,
  "eventType": "nginx.access",
  "windowStart": "2026-05-24 16:30:00",
  "windowEnd": "2026-05-24 16:30:59",
  "totalCount": 100000,
  "successCount": 98500,
  "failedCount": 300,
  "cleanFailedCount": 900,
  "droppedCount": 300,
  "lastError": "price 类型转换失败",
  "errorSamples": [
    {
      "errorType": "type_convert_failed",
      "fieldPath": "price",
      "message": "无法转换为 number"
    }
  ]
}
```

Server 写入时采用累加语义：

```text
total_count = total_count + incoming.total_count
success_count = success_count + incoming.success_count
failed_count = failed_count + incoming.failed_count
clean_failed_count = clean_failed_count + incoming.clean_failed_count
dropped_count = dropped_count + incoming.dropped_count
last_error = incoming.last_error
error_samples = 合并后最多保留 10 条
```

## 字段采集策略

字段采集默认开启，和数据清洗开关解耦：

```text
字段采集：任务运行就默认做
数据清洗：clean_enabled = 1 才做
```

这样即使当前任务还没开启清洗，也能提前积累字段集合。后续需要配置清洗规则时，页面已经有字段列表。

## Agent 字段发现流程

不要每条数据都请求 server，也不要每条数据都按所有 key 计算 schema hash 后请求 server。这会在高吞吐场景压垮 server，也会被可选字段组合打爆。

推荐流程：

1. Agent 启动时拉取任务配置、清洗规则版本、字段集合版本和已知字段集合。
2. Agent 本地维护 `knownFields` 缓存，缓存维度至少包含 `task_id + field_path` 和字段类型集合。
3. 每条数据在 Agent 本地提取字段路径和字段类型。
4. 如果没有新字段、没有新类型，则不请求 server。
5. 如果发现新字段或新类型，放入本地 pending 队列。
6. pending 队列使用短 debounce 和批量上报，例如 1-3 秒内合并一次，或达到 N 个字段立即上报。
7. Server 批量 upsert 到 `hg_data_field`。
8. Server 成功合并新字段后递增 `hg_data_clean_task.field_schema_version`。
9. Agent 根据返回结果更新本地 `knownFields`。

第一版不做 `hg_data_field_report` 批次表。幂等依赖 `hg_data_field` 的唯一键 `task_id + field_path`。

## 字段路径提取规则

字段发现必须做收敛，否则动态 JSON 会撑爆字段集合。

建议规则：

- 数组下标归一化：

```text
items[0].price -> items[].price
items[1].price -> items[].price
```

- 限制最大展开深度，例如 8 层。
- 限制单条数据最多提取字段数，例如 500 个。
- 限制单个任务最多字段数，例如 5000 或 10000 个，超过后停止新增并告警。
- 限制字段路径长度，例如 512 字符。
- 过滤动态 key：明显像 UUID、纯数字 ID、时间戳、手机号、随机串的 key 不作为字段路径。
- 类型归一化为有限集合，例如 `string`、`number`、`boolean`、`object`、`array`、`null`、`unknown`。
- 字段样例只保存单字段值，必须截断和脱敏，例如 token、password、secret、手机号、身份证、邮箱等。

这些限制可以放在 `hg_data_clean_task.source_config` 中：

```json
{
  "fieldDiscovery": {
    "maxDepth": 8,
    "maxFieldsPerRecord": 500,
    "maxFieldsPerTask": 10000,
    "normalizeArrayIndex": true,
    "ignoreDynamicKeys": true,
    "maxPathLength": 512
  }
}
```

## 清洗场景

### 未开启清洗

```text
clean_enabled = 0
```

Agent 仍然默认采集字段，但不执行 `clean_config`。数据按 `sink_config` 直接输出到下游。

### 开启清洗且规则存在

```text
clean_enabled = 1
clean_config 已发布
```

Agent 按 `clean_config` 做字段选择、删除、重命名、转换、脱敏、模板映射等处理后输出到下游。

### 开启清洗但规则不存在

Agent 仍采集字段并上报到 `hg_data_field`，但数据是否输出由策略决定。默认建议不输出，等用户配置清洗规则后再处理后续数据。

### 发现新字段

默认 `selected_only` 策略下，新字段不会自动进入输出。Agent 上报新字段，server 合并到 `hg_data_field`，清洗页面展示新字段，用户后续按需加入规则。

如果任务配置为 `strict`，发现未知字段可阻塞或失败，等待用户确认。

### 字段缺失或类型变化

- 字段缺失：按清洗规则处理，例如默认值、置空、跳过或失败。
- 字段类型变化：上报新的字段类型集合，若规则能转换则继续输出；转换失败则执行 `clean_error_policy`。

## 后续开发注意点

- 后端查询字段列表时从 `hg_data_field` 按 `task_id` 拉取。
- 标准化配置界面应展示两侧字段：
  - 左侧：实际发现字段 `hg_data_field`
  - 右侧：标准模板字段 `hg_data_field_template.fields`
- 保存清洗规则时递增 `clean_config_version`。
- 字段集合新增成功时递增 `field_schema_version`。
- Agent 不应每条数据请求 server。只在发现新字段或新类型时批量上报。
- Server 不保存全量原始数据，不提供原始 payload 查询能力。
- 连接级敏感配置必须加密存储，回显时脱敏。
