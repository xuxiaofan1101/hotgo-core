-- --------------------------------------------------------
-- 数据集成 MySQL 表结构
-- 主模块表前缀：hg_data_*
-- --------------------------------------------------------

--
-- 表的结构 `hg_data_connector`
--

CREATE TABLE IF NOT EXISTS `hg_data_connector` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '连接ID',
  `name` varchar(128) NOT NULL COMMENT '连接名称',
  `code` varchar(128) NOT NULL COMMENT '连接编码',
  `direction` varchar(16) NOT NULL COMMENT '连接方向：source=输入 sink=输出',
  `connector_type` varchar(32) NOT NULL COMMENT '连接类型：http,kafka,s3,log,manual,feishu,lark,elasticsearch,opensearch,splunk',
  `config` json DEFAULT NULL COMMENT '连接配置',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint(20) DEFAULT '0' COMMENT '删除者',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_data_connector_direction` (`direction`,`status`),
  KEY `idx_data_connector_type` (`connector_type`,`status`),
  KEY `idx_data_connector_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_数据源连接';

-- --------------------------------------------------------

--
-- 表的结构 `hg_data_clean_task`
--

CREATE TABLE IF NOT EXISTS `hg_data_clean_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '清洗任务ID',
  `source_id` bigint(20) NOT NULL COMMENT '输入连接ID',
  `template_id` bigint(20) DEFAULT '0' COMMENT '标准字段模板ID',
  `name` varchar(128) NOT NULL COMMENT '任务名称',
  `code` varchar(128) NOT NULL COMMENT '任务编码',
  `event_type` varchar(128) NOT NULL COMMENT '事件类型',
  `source_config` json DEFAULT NULL COMMENT '输入配置',
  `clean_enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用数据清洗',
  `clean_config` json DEFAULT NULL COMMENT '数据清洗配置',
  `clean_config_version` int(11) NOT NULL DEFAULT '1' COMMENT '清洗配置版本',
  `field_schema_version` int(11) NOT NULL DEFAULT '1' COMMENT '字段集合版本',
  `unknown_field_policy` varchar(32) NOT NULL DEFAULT 'selected_only' COMMENT '未知字段策略：selected_only strict passthrough drop',
  `clean_error_policy` varchar(16) NOT NULL DEFAULT 'skip' COMMENT '清洗失败策略：skip keep_raw stop_task',
  `sink_config` json DEFAULT NULL COMMENT '输出配置',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint(20) DEFAULT '0' COMMENT '删除者',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_data_clean_task_source` (`source_id`,`status`),
  KEY `idx_data_clean_task_template` (`template_id`),
  KEY `idx_data_clean_task_event` (`event_type`,`status`),
  KEY `idx_data_clean_task_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_数据清洗任务';

-- --------------------------------------------------------

--
-- 表的结构 `hg_data_field`
--

CREATE TABLE IF NOT EXISTS `hg_data_field` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字段ID',
  `task_id` bigint(20) NOT NULL COMMENT '清洗任务ID',
  `connector_id` bigint(20) NOT NULL COMMENT '输入连接ID',
  `field_path` varchar(512) NOT NULL COMMENT '字段路径',
  `field_name` varchar(128) DEFAULT '' COMMENT '字段名称',
  `field_type` varchar(64) DEFAULT '' COMMENT '字段类型',
  `field_types` json DEFAULT NULL COMMENT '出现过的字段类型集合',
  `sample_value` json DEFAULT NULL COMMENT '样例值',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `first_seen_at` datetime DEFAULT NULL COMMENT '首次发现时间',
  `last_seen_at` datetime DEFAULT NULL COMMENT '最后发现时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_data_field_task_path` (`task_id`,`field_path`),
  KEY `idx_data_field_task` (`task_id`,`status`),
  KEY `idx_data_field_connector` (`connector_id`),
  KEY `idx_data_field_type` (`field_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_数据字段集合';

-- --------------------------------------------------------

--
-- 表的结构 `hg_data_field_template`
--

CREATE TABLE IF NOT EXISTS `hg_data_field_template` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字段模板ID',
  `name` varchar(128) NOT NULL COMMENT '模板名称',
  `code` varchar(128) NOT NULL COMMENT '模板编码',
  `event_type` varchar(128) NOT NULL COMMENT '事件类型',
  `fields` json DEFAULT NULL COMMENT '标准字段定义',
  `examples` json DEFAULT NULL COMMENT '标准样例数据',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint(20) DEFAULT '0' COMMENT '删除者',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_data_field_template_event` (`event_type`,`status`),
  KEY `idx_data_field_template_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_数据字段模板';

-- --------------------------------------------------------

--
-- 表的结构 `hg_data_clean_stat`
--

CREATE TABLE IF NOT EXISTS `hg_data_clean_stat` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '清洗统计ID',
  `task_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '清洗任务ID',
  `connector_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '输入连接ID',
  `event_type` varchar(128) NOT NULL DEFAULT '' COMMENT '事件类型',
  `window_start` datetime NOT NULL COMMENT '统计窗口开始时间',
  `window_end` datetime NOT NULL COMMENT '统计窗口结束时间',
  `total_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '总处理数',
  `success_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '成功数',
  `failed_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '失败数',
  `clean_failed_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '清洗失败数',
  `dropped_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '丢弃数',
  `last_error` text COMMENT '最近错误',
  `error_samples` json DEFAULT NULL COMMENT '错误样例',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_data_clean_stat_window` (`task_id`,`connector_id`,`event_type`,`window_start`),
  KEY `idx_data_clean_stat_task` (`task_id`,`window_start`),
  KEY `idx_data_clean_stat_connector` (`connector_id`,`window_start`),
  KEY `idx_data_clean_stat_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_数据清洗统计';
