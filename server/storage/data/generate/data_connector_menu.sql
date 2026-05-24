-- 数据集成 - 数据源菜单权限
-- 可重复执行：以 hg_admin_menu.name 做幂等更新，并为 superadmin 写入菜单授权。

SET @now := NOW();

START TRANSACTION;

-- 一级目录：数据集成
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, 0, '数据集成', 'dataIntegration', '/dataIntegration', 'DatabaseOutlined', 1,
    '/dataIntegration/dataConnector', '', '', 'LAYOUT', 1, '', 0, 0, '', 0, 0,
    0, 1, '', 220, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `path` = VALUES(`path`),
    `icon` = VALUES(`icon`),
    `type` = VALUES(`type`),
    `redirect` = VALUES(`redirect`),
    `permissions` = VALUES(`permissions`),
    `permission_name` = VALUES(`permission_name`),
    `component` = VALUES(`component`),
    `always_show` = VALUES(`always_show`),
    `active_menu` = VALUES(`active_menu`),
    `is_root` = VALUES(`is_root`),
    `is_frame` = VALUES(`is_frame`),
    `frame_src` = VALUES(`frame_src`),
    `keep_alive` = VALUES(`keep_alive`),
    `hidden` = VALUES(`hidden`),
    `affix` = VALUES(`affix`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `remark` = VALUES(`remark`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

SET @dataIntegrationId := LAST_INSERT_ID();

-- 二级菜单：数据源
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, @dataIntegrationId, '数据源', 'dataConnector', 'dataConnector', '', 2, '',
    '/dataConnector/list', '', '/dataConnector/index', 1, '', 0, 0, '', 0, 0,
    0, 2, CONCAT('tr_', @dataIntegrationId, ' '), 10, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `path` = VALUES(`path`),
    `icon` = VALUES(`icon`),
    `type` = VALUES(`type`),
    `redirect` = VALUES(`redirect`),
    `permissions` = VALUES(`permissions`),
    `permission_name` = VALUES(`permission_name`),
    `component` = VALUES(`component`),
    `always_show` = VALUES(`always_show`),
    `active_menu` = VALUES(`active_menu`),
    `is_root` = VALUES(`is_root`),
    `is_frame` = VALUES(`is_frame`),
    `frame_src` = VALUES(`frame_src`),
    `keep_alive` = VALUES(`keep_alive`),
    `hidden` = VALUES(`hidden`),
    `affix` = VALUES(`affix`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `remark` = VALUES(`remark`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

SET @dataConnectorId := LAST_INSERT_ID();

-- 按钮权限：详情
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, @dataConnectorId, '数据源详情', 'dataConnectorView', '', '', 3, '', '/dataConnector/view', '',
    '', 1, '', 0, 0, '', 0, 1, 0, 3,
    CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 10, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `permissions` = VALUES(`permissions`),
    `type` = VALUES(`type`),
    `hidden` = VALUES(`hidden`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

-- 按钮权限：新增/编辑
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, @dataConnectorId, '新增/编辑数据源', 'dataConnectorEdit', '', '', 3, '', '/dataConnector/edit', '',
    '', 1, '', 0, 0, '', 0, 1, 0, 3,
    CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 20, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `permissions` = VALUES(`permissions`),
    `type` = VALUES(`type`),
    `hidden` = VALUES(`hidden`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

-- 按钮权限：删除
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, @dataConnectorId, '删除数据源', 'dataConnectorDelete', '', '', 3, '', '/dataConnector/delete', '',
    '', 1, '', 0, 0, '', 0, 1, 0, 3,
    CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 30, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `permissions` = VALUES(`permissions`),
    `type` = VALUES(`type`),
    `hidden` = VALUES(`hidden`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

-- 按钮权限：状态更新
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, @dataConnectorId, '更新数据源状态', 'dataConnectorStatus', '', '', 3, '', '/dataConnector/status', '',
    '', 1, '', 0, 0, '', 0, 1, 0, 3,
    CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 40, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `permissions` = VALUES(`permissions`),
    `type` = VALUES(`type`),
    `hidden` = VALUES(`hidden`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

-- 按钮权限：测试配置
INSERT INTO `hg_admin_menu` (
    `id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`,
    `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`,
    `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`
) VALUES (
    NULL, @dataConnectorId, '测试数据源配置', 'dataConnectorTest', '', '', 3, '', '/dataConnector/test', '',
    '', 1, '', 0, 0, '', 0, 1, 0, 3,
    CONCAT('tr_', @dataIntegrationId, ' tr_', @dataConnectorId, ' '), 50, '', 1, @now, @now
) ON DUPLICATE KEY UPDATE
    `id` = LAST_INSERT_ID(`id`),
    `pid` = VALUES(`pid`),
    `title` = VALUES(`title`),
    `permissions` = VALUES(`permissions`),
    `type` = VALUES(`type`),
    `hidden` = VALUES(`hidden`),
    `level` = VALUES(`level`),
    `tree` = VALUES(`tree`),
    `sort` = VALUES(`sort`),
    `status` = VALUES(`status`),
    `updated_at` = @now;

-- 默认授权给超级管理员角色
INSERT IGNORE INTO `hg_admin_role_menu` (`role_id`, `menu_id`)
SELECT 1, `id`
FROM `hg_admin_menu`
WHERE `name` IN (
    'dataIntegration',
    'dataConnector',
    'dataConnectorView',
    'dataConnectorEdit',
    'dataConnectorDelete',
    'dataConnectorStatus',
    'dataConnectorTest'
);

COMMIT;
