-- --------------------------------------------------------
-- 数据集成 MySQL 表结构
-- 主模块表前缀：hg_data_*
-- --------------------------------------------------------

--
-- 表的结构 `hg_data_connector`
--

CREATE TABLE IF NOT EXISTS `hg_data_connector` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '连接ID',
  `name` varchar(128) NOT NULL COMMENT '连接名称',
  `code` varchar(128) NOT NULL COMMENT '连接编码',
  `direction` varchar(16) NOT NULL COMMENT '连接方向：source=输入 sink=输出',
  `connector_type` varchar(32) NOT NULL COMMENT '连接类型：http,kafka,s3,log,manual,feishu,lark,elasticsearch,opensearch,splunk',
  `config` json DEFAULT NULL COMMENT '连接配置',
  `status` tinyint DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint DEFAULT '0' COMMENT '删除者',
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
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '清洗任务ID',
  `source_id` bigint NOT NULL COMMENT '输入连接ID',
  `template_id` bigint DEFAULT '0' COMMENT '标准字段模板ID',
  `name` varchar(128) NOT NULL COMMENT '任务名称',
  `code` varchar(128) NOT NULL COMMENT '任务编码',
  `event_type` varchar(128) NOT NULL COMMENT '事件类型',
  `source_config` json DEFAULT NULL COMMENT '输入配置',
  `clean_enabled` tinyint NOT NULL DEFAULT '0' COMMENT '是否启用数据清洗',
  `clean_config` json DEFAULT NULL COMMENT '数据清洗配置',
  `clean_config_version` int NOT NULL DEFAULT '1' COMMENT '清洗配置版本',
  `field_schema_version` int NOT NULL DEFAULT '1' COMMENT '字段集合版本',
  `unknown_field_policy` varchar(32) NOT NULL DEFAULT 'selected_only' COMMENT '未知字段策略：selected_only strict passthrough drop',
  `clean_error_policy` varchar(16) NOT NULL DEFAULT 'skip' COMMENT '清洗失败策略：skip keep_raw stop_task',
  `sink_config` json DEFAULT NULL COMMENT '输出配置',
  `status` tinyint DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint DEFAULT '0' COMMENT '删除者',
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
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段ID',
  `task_id` bigint NOT NULL COMMENT '清洗任务ID',
  `connector_id` bigint NOT NULL COMMENT '输入连接ID',
  `field_path` varchar(512) NOT NULL COMMENT '字段路径',
  `field_name` varchar(128) DEFAULT '' COMMENT '字段名称',
  `field_type` varchar(64) DEFAULT '' COMMENT '字段类型',
  `field_types` json DEFAULT NULL COMMENT '出现过的字段类型集合',
  `sample_value` json DEFAULT NULL COMMENT '样例值',
  `status` tinyint DEFAULT '1' COMMENT '状态',
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
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字段模板ID',
  `name` varchar(128) NOT NULL COMMENT '模板名称',
  `code` varchar(128) NOT NULL COMMENT '模板编码',
  `event_type` varchar(128) NOT NULL COMMENT '事件类型',
  `fields` json DEFAULT NULL COMMENT '标准字段定义',
  `examples` json DEFAULT NULL COMMENT '标准样例数据',
  `status` tinyint DEFAULT '1' COMMENT '状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_by` bigint DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint DEFAULT '0' COMMENT '删除者',
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
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '清洗统计ID',
  `task_id` bigint NOT NULL DEFAULT '0' COMMENT '清洗任务ID',
  `connector_id` bigint NOT NULL DEFAULT '0' COMMENT '输入连接ID',
  `event_type` varchar(128) NOT NULL DEFAULT '' COMMENT '事件类型',
  `window_start` datetime NOT NULL COMMENT '统计窗口开始时间',
  `window_end` datetime NOT NULL COMMENT '统计窗口结束时间',
  `total_count` bigint NOT NULL DEFAULT '0' COMMENT '总处理数',
  `success_count` bigint NOT NULL DEFAULT '0' COMMENT '成功数',
  `failed_count` bigint NOT NULL DEFAULT '0' COMMENT '失败数',
  `clean_failed_count` bigint NOT NULL DEFAULT '0' COMMENT '清洗失败数',
  `dropped_count` bigint NOT NULL DEFAULT '0' COMMENT '丢弃数',
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

-- --------------------------------------------------------

--
-- 表的结构 `hg_data_agent`
--

CREATE TABLE IF NOT EXISTS `hg_data_agent` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '节点ID',
  `agent_id` varchar(128) NOT NULL COMMENT 'Agent唯一标识',
  `name` varchar(128) DEFAULT '' COMMENT '节点名称',
  `hostname` varchar(255) DEFAULT '' COMMENT '主机名',
  `version` varchar(64) DEFAULT '' COMMENT 'Agent版本',
  `agent_ip` json DEFAULT NULL COMMENT 'Agent上报IP列表',
  `register_status` varchar(16) NOT NULL DEFAULT 'pending' COMMENT '注册状态：pending approved rejected revoked',
  `dispatch_status` varchar(16) NOT NULL DEFAULT 'disabled' COMMENT '调度状态：enabled disabled',
  `last_seen_at` datetime DEFAULT NULL COMMENT '最近心跳时间',
  `approved_by` bigint DEFAULT '0' COMMENT '批准人',
  `approved_at` datetime DEFAULT NULL COMMENT '批准时间',
  `disabled_at` datetime DEFAULT NULL COMMENT '禁用/拒绝/吊销时间',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_data_agent_agent_id` (`agent_id`),
  KEY `idx_data_agent_register` (`register_status`),
  KEY `idx_data_agent_dispatch` (`dispatch_status`),
  KEY `idx_data_agent_last_seen` (`last_seen_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_数据集成Agent节点';

-- --------------------------------------------------------

--
-- 数据源菜单权限，可独立执行
--

SET @now := NOW();

START TRANSACTION;

-- 一级目录：数据集成
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, 0, '数据集成', 'dataIntegration', '/dataIntegration', 'DatabaseOutlined', 1, '/dataIntegration/dataConnector', '', '', 'LAYOUT', 1, '', 0, 0, '', 0, 0, 0, 1, '', 220, '', 1, @now, @now);

SET @dataIntegrationId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataIntegration' LIMIT 1);

-- 二级菜单：数据源
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataIntegrationId, '数据源', 'dataConnector', 'dataConnector', '', 2, '', '/dataConnector/list', '', '/dataConnector/index', 1, '', 0, 0, '', 0, 0, 0, 2, CONCAT('tr_', @dataIntegrationId, ' '), 10, '', 1, @now, @now);

SET @dataConnectorId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataConnector' LIMIT 1);

-- 数据源按钮权限
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataConnectorId, '数据源详情', 'dataConnectorView', '', '', 3, '', '/dataConnector/view', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 10, '', 1, @now, @now),(NULL, @dataConnectorId, '新增/编辑数据源', 'dataConnectorEdit', '', '', 3, '', '/dataConnector/edit', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 20, '', 1, @now, @now),(NULL, @dataConnectorId, '删除数据源', 'dataConnectorDelete', '', '', 3, '', '/dataConnector/delete', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 30, '', 1, @now, @now),(NULL, @dataConnectorId, '更新数据源状态', 'dataConnectorStatus', '', '', 3, '', '/dataConnector/status', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 40, '', 1, @now, @now),(NULL, @dataConnectorId, '测试数据源配置', 'dataConnectorTest', '', '', 3, '', '/dataConnector/test', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 50, '', 1, @now, @now);

-- 数据源默认授权给超级管理员角色
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`)
SELECT r.`id`, m.`id` FROM `hg_admin_role` r JOIN `hg_admin_menu` m WHERE r.`key` = 'superadmin' AND m.`name` IN ('dataIntegration', 'dataConnector', 'dataConnectorView', 'dataConnectorEdit', 'dataConnectorDelete', 'dataConnectorStatus', 'dataConnectorTest');

COMMIT;

--
-- 数据清洗菜单权限，可独立执行
--

SET @now := NOW();

START TRANSACTION;

-- 一级目录：数据集成
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, 0, '数据集成', 'dataIntegration', '/dataIntegration', 'DatabaseOutlined', 1, '/dataIntegration/dataConnector', '', '', 'LAYOUT', 1, '', 0, 0, '', 0, 0, 0, 1, '', 220, '', 1, @now, @now);

SET @dataIntegrationId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataIntegration' LIMIT 1);

-- 二级菜单：数据清洗
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataIntegrationId, '数据清洗', 'dataClean', 'dataClean', '', 2, '', '/dataClean/list', '', '/dataClean/index', 1, '', 0, 0, '', 0, 0, 0, 2, CONCAT('tr_', @dataIntegrationId, ' '), 20, '', 1, @now, @now);

SET @dataCleanId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataClean' LIMIT 1);

-- 数据清洗按钮权限
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataCleanId, '数据清洗详情', 'dataCleanView', '', '', 3, '', '/dataClean/view', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 10, '', 1, @now, @now),(NULL, @dataCleanId, '新增/编辑数据清洗', 'dataCleanEdit', '', '', 3, '', '/dataClean/edit', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 20, '', 1, @now, @now),(NULL, @dataCleanId, '删除数据清洗', 'dataCleanDelete', '', '', 3, '', '/dataClean/delete', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 30, '', 1, @now, @now),(NULL, @dataCleanId, '更新数据清洗状态', 'dataCleanStatus', '', '', 3, '', '/dataClean/status', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 40, '', 1, @now, @now),(NULL, @dataCleanId, '测试数据清洗配置', 'dataCleanTest', '', '', 3, '', '/dataClean/test', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 50, '', 1, @now, @now),(NULL, @dataCleanId, '采样数据清洗字段', 'dataCleanSample', '', '', 3, '', '/dataClean/sample', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 60, '', 1, @now, @now),(NULL, @dataCleanId, '获取采集字段', 'dataCleanFieldList', '', '', 3, '', '/dataClean/fieldList', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataCleanId, ' '), 70, '', 1, @now, @now);

-- 数据清洗默认授权给超级管理员角色
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`)
SELECT r.`id`, m.`id` FROM `hg_admin_role` r JOIN `hg_admin_menu` m WHERE r.`key` = 'superadmin' AND m.`name` IN ('dataIntegration', 'dataClean', 'dataCleanView', 'dataCleanEdit', 'dataCleanDelete', 'dataCleanStatus', 'dataCleanTest', 'dataCleanSample', 'dataCleanFieldList');

COMMIT;

--
-- 字段模板菜单权限，可独立执行
--

SET @now := NOW();

START TRANSACTION;

-- 一级目录：数据集成
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, 0, '数据集成', 'dataIntegration', '/dataIntegration', 'DatabaseOutlined', 1, '/dataIntegration/dataConnector', '', '', 'LAYOUT', 1, '', 0, 0, '', 0, 0, 0, 1, '', 220, '', 1, @now, @now);

SET @dataIntegrationId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataIntegration' LIMIT 1);

-- 二级菜单：字段模板
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataIntegrationId, '字段模板', 'dataFieldTemplate', 'dataFieldTemplate', '', 2, '', '/dataFieldTemplate/list', '', '/dataFieldTemplate/index', 1, '', 0, 0, '', 0, 0, 0, 2, CONCAT('tr_', @dataIntegrationId, ' '), 30, '', 1, @now, @now);

SET @dataFieldTemplateId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataFieldTemplate' LIMIT 1);

-- 字段模板按钮权限
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataFieldTemplateId, '字段模板详情', 'dataFieldTemplateView', '', '', 3, '', '/dataFieldTemplate/view', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataFieldTemplateId, ' '), 10, '', 1, @now, @now),(NULL, @dataFieldTemplateId, '新增/编辑字段模板', 'dataFieldTemplateEdit', '', '', 3, '', '/dataFieldTemplate/edit', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataFieldTemplateId, ' '), 20, '', 1, @now, @now),(NULL, @dataFieldTemplateId, '删除字段模板', 'dataFieldTemplateDelete', '', '', 3, '', '/dataFieldTemplate/delete', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataFieldTemplateId, ' '), 30, '', 1, @now, @now),(NULL, @dataFieldTemplateId, '更新字段模板状态', 'dataFieldTemplateStatus', '', '', 3, '', '/dataFieldTemplate/status', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataFieldTemplateId, ' '), 40, '', 1, @now, @now);

-- 字段模板默认授权给超级管理员角色
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`)
SELECT r.`id`, m.`id` FROM `hg_admin_role` r JOIN `hg_admin_menu` m WHERE r.`key` = 'superadmin' AND m.`name` IN ('dataIntegration', 'dataFieldTemplate', 'dataFieldTemplateView', 'dataFieldTemplateEdit', 'dataFieldTemplateDelete', 'dataFieldTemplateStatus');

COMMIT;

--
-- 节点管理菜单权限，可独立执行
--

SET @now := NOW();

START TRANSACTION;

-- 一级目录：数据集成
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, 0, '数据集成', 'dataIntegration', '/dataIntegration', 'DatabaseOutlined', 1, '/dataIntegration/dataConnector', '', '', 'LAYOUT', 1, '', 0, 0, '', 0, 0, 0, 1, '', 220, '', 1, @now, @now);

SET @dataIntegrationId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataIntegration' LIMIT 1);

-- 二级菜单：节点管理
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataIntegrationId, '节点管理', 'dataAgent', 'dataAgent', '', 2, '', '/dataAgent/list', '', '/dataAgent/index', 1, '', 0, 0, '', 0, 0, 0, 2, CONCAT('tr_', @dataIntegrationId, ' '), 40, '', 1, @now, @now);

SET @dataAgentId := (SELECT `id` FROM `hg_admin_menu` WHERE `name` = 'dataAgent' LIMIT 1);

-- 节点管理按钮权限
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(NULL, @dataAgentId, '查看节点', 'dataAgentList', '', '', 3, '', '/dataAgent/list', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataAgentId, ' '), 10, '', 1, @now, @now),(NULL, @dataAgentId, '批准节点', 'dataAgentApprove', '', '', 3, '', '/dataAgent/approve', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataAgentId, ' '), 20, '', 1, @now, @now),(NULL, @dataAgentId, '设置节点调度', 'dataAgentDispatch', '', '', 3, '', '/dataAgent/dispatch', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataAgentId, ' '), 30, '', 1, @now, @now),(NULL, @dataAgentId, '拒绝/吊销节点', 'dataAgentReject', '', '', 3, '', '/dataAgent/reject', '', '', 1, '', 0, 0, '', 0, 1, 0, 3, CONCAT('tr_', @dataIntegrationId, ' tr_', @dataAgentId, ' '), 40, '', 1, @now, @now);

-- 节点管理默认授权给超级管理员角色
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`)
SELECT r.`id`, m.`id` FROM `hg_admin_role` r JOIN `hg_admin_menu` m WHERE r.`key` = 'superadmin' AND m.`name` IN ('dataIntegration', 'dataAgent', 'dataAgentList', 'dataAgentApprove', 'dataAgentDispatch', 'dataAgentReject');

COMMIT;
