ALTER TABLE `hg_admin_menu`
    MODIFY tree VARCHAR(255) NOT NULL DEFAULT '' COMMENT '关系树';

-- 表的结构 `hg_kafka_clusters`
CREATE TABLE IF NOT EXISTS `hg_kafka_clusters` (
                                                   `id` SERIAL PRIMARY KEY COMMENT 'ID',
                                                   `cluster_name` VARCHAR(100) NOT NULL COMMENT '集群名',
    `bootstrap_servers` VARCHAR(1000) NOT NULL COMMENT 'Servers',
    `description` TEXT COMMENT '集群描述',
    `created_by` BIGINT(20) DEFAULT '0' COMMENT '创建者',
    `updated_by` BIGINT(20) DEFAULT '0' COMMENT '更新者',
    `created_at` DATETIME DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL COMMENT '更新时间',
    UNIQUE (cluster_name)
    ) COMMENT='clusters';

-- 表的结构 `hg_kafka_users`
CREATE TABLE IF NOT EXISTS `hg_kafka_users` (
                                                `id` SERIAL PRIMARY KEY COMMENT 'ID',
                                                `cluster_name` VARCHAR(100) DEFAULT NULL COMMENT '集群名',
    `username` VARCHAR(255) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `role` VARCHAR(50) NOT NULL COMMENT '角色',
    `created_by` BIGINT(20) DEFAULT '0' COMMENT '创建者',
    `updated_by` BIGINT(20) DEFAULT '0' COMMENT '更新者',
    `created_at` DATETIME DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL COMMENT '更新时间',
    unique (username)
    )COMMENT='users';

-- 表的结构 `hg_kafka_topics`
CREATE TABLE IF NOT EXISTS `hg_kafka_topics` (
                                                 `id` SERIAL PRIMARY KEY COMMENT 'ID',
                                                 `cluster_name` VARCHAR(100) DEFAULT NULL COMMENT '集群名',
    `topic_name` VARCHAR(255) NOT NULL COMMENT 'topic',
    `partitions` INT NOT NULL COMMENT 'partitions',
    `replication_factor` INT NOT NULL COMMENT 'replication',
    `message_format` VARCHAR(50) COMMENT '消息格式',
    `description` TEXT COMMENT '描述',
    `created_by` BIGINT(20) DEFAULT '0' COMMENT '创建者',
    `updated_by` BIGINT(20) DEFAULT '0' COMMENT '更新者',
    `created_at` DATETIME DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL COMMENT '更新时间',
    unique (topic_name)
    )COMMENT='topics';

-- 表的结构 `hg_kafka_consumers`
CREATE TABLE IF NOT EXISTS `hg_kafka_consumers` (
                                                    `id` SERIAL PRIMARY KEY COMMENT 'ID',
                                                    `cluster_name` VARCHAR(100) DEFAULT NULL COMMENT '集群名',
    `consumer_group` VARCHAR(255) NOT NULL COMMENT '消费组',
    `topic_id` INT NOT NULL COMMENT 'topic ID',
    `topic_name` INT NOT NULL COMMENT 'topic',
    `offset_strategy` VARCHAR(50) COMMENT '消费策略',
    `lag` BIGINT UNSIGNED COMMENT '消费者处理消息的延迟量',
    `current_offset` BIGINT COMMENT '当前消费者的处理位移',
    `created_by` BIGINT(20) DEFAULT '0' COMMENT '创建者',
    `updated_by` BIGINT(20) DEFAULT '0' COMMENT '更新者',
    `created_at` DATETIME DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL COMMENT '更新时间'
    )COMMENT='consumers';

-- 表的结构 `hg_kafka_acls`
CREATE TABLE IF NOT EXISTS `hg_kafka_acls` (
                                               `id` SERIAL PRIMARY KEY COMMENT 'ID',
                                               `cluster_name` VARCHAR(100) DEFAULT NULL COMMENT '集群名',
    `kafka_user` VARCHAR(255) NOT NULL COMMENT '用户名',
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `resource_name` VARCHAR(50) NOT NULL COMMENT '资源名称',
    `operation` VARCHAR(50) NOT NULL COMMENT '操作类型',
    `permission_type` VARCHAR(50) NOT NULL COMMENT '权限类型',
    `host` VARCHAR(50) NOT NULL COMMENT '主机',
    `created_by` BIGINT(20) DEFAULT '0' COMMENT '创建者',
    `updated_by` BIGINT(20) DEFAULT '0' COMMENT '更新者',
    `created_at` DATETIME DEFAULT NULL COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT NULL COMMENT '更新时间',
    INDEX idx_kafka_acl (`kafka_user`, `resource_type`, `resource_name`, `operation`, `permission_type`, `host`)
    )COMMENT='acls';
