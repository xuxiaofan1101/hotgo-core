-- phpMyAdmin SQL Dump
-- version 4.4.15.10
-- https://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: 2024-08-27 19:04:42
-- 服务器版本： 5.7.37-log
-- PHP Version: 5.6.40

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `hotgo`
--

-- --------------------------------------------------------

--
-- 表的结构 `hg_addon_hgexample_table`
--

CREATE TABLE IF NOT EXISTS `hg_addon_hgexample_table` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `pid` bigint(20) DEFAULT '0' COMMENT '上级ID',
  `level` int(11) DEFAULT '1' COMMENT '树等级',
  `tree` varchar(512) DEFAULT NULL COMMENT '关系树',
  `category_id` bigint(20) DEFAULT NULL COMMENT '分类ID',
  `flag` json DEFAULT NULL COMMENT '标签',
  `title` varchar(255) NOT NULL COMMENT '标题',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `content` longtext COMMENT '内容',
  `image` varchar(255) DEFAULT NULL COMMENT '单图',
  `images` json DEFAULT NULL COMMENT '多图',
  `attachfile` varchar(255) DEFAULT NULL COMMENT '附件',
  `attachfiles` json DEFAULT NULL COMMENT '多附件',
  `map` json DEFAULT NULL COMMENT '动态键值对',
  `star` decimal(5,1) DEFAULT '0.0' COMMENT '推荐星',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `views` bigint(20) DEFAULT NULL COMMENT '浏览次数',
  `activity_at` date DEFAULT NULL COMMENT '活动时间',
  `start_at` datetime DEFAULT NULL COMMENT '开启时间',
  `end_at` datetime DEFAULT NULL COMMENT '结束时间',
  `switch` tinyint(1) DEFAULT NULL COMMENT '开关',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `avatar` varchar(255) DEFAULT '' COMMENT '头像',
  `sex` tinyint(1) DEFAULT NULL COMMENT '性别',
  `qq` varchar(20) DEFAULT '' COMMENT 'qq',
  `email` varchar(60) DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT '' COMMENT '手机号码',
  `hobby` json DEFAULT NULL COMMENT '爱好',
  `channel` int(11) DEFAULT '1' COMMENT '渠道',
  `city_id` bigint(20) DEFAULT '0' COMMENT '所在城市',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='插件_案例_表格';


--
-- 表的结构 `hg_addon_hgexample_tenant_order`
--

CREATE TABLE IF NOT EXISTS `hg_addon_hgexample_tenant_order` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `tenant_id` bigint(20) DEFAULT NULL COMMENT '租户ID',
  `merchant_id` bigint(20) NOT NULL COMMENT '商户ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `product_name` varchar(255) DEFAULT NULL COMMENT '购买产品',
  `order_sn` varchar(64) DEFAULT NULL COMMENT '订单号',
  `money` decimal(10,2) NOT NULL COMMENT '充值金额',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(4) DEFAULT '1' COMMENT '订单状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='多租户_充值订单';


--
-- 表的结构 `hg_admin_cash`
--

CREATE TABLE IF NOT EXISTS `hg_admin_cash` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `member_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `money` decimal(10,2) NOT NULL COMMENT '提现金额',
  `fee` decimal(10,2) NOT NULL COMMENT '手续费',
  `last_money` decimal(10,2) NOT NULL COMMENT '最终到账金额',
  `ip` varchar(128) NOT NULL COMMENT '申请人IP',
  `status` bigint(20) NOT NULL COMMENT '状态码',
  `msg` varchar(128) DEFAULT NULL COMMENT '处理结果',
  `handle_at` datetime DEFAULT NULL COMMENT '处理时间',
  `created_at` datetime NOT NULL COMMENT '申请时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_提现记录表';


--
-- 表的结构 `hg_admin_credits_log`
--

CREATE TABLE IF NOT EXISTS `hg_admin_credits_log` (
  `id` bigint(20) NOT NULL COMMENT '变动ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '管理员ID',
  `app_id` varchar(64) DEFAULT NULL COMMENT '应用id',
  `addons_name` varchar(100) NOT NULL DEFAULT '' COMMENT '插件名称',
  `credit_type` varchar(32) NOT NULL DEFAULT '' COMMENT '变动类型',
  `credit_group` varchar(32) DEFAULT NULL COMMENT '变动组别',
  `before_num` decimal(10,2) DEFAULT '0.00' COMMENT '变动前',
  `num` decimal(10,2) DEFAULT '0.00' COMMENT '变动数据',
  `after_num` decimal(10,2) DEFAULT '0.00' COMMENT '变动后',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `ip` varchar(20) DEFAULT NULL COMMENT '操作人IP',
  `map_id` bigint(20) DEFAULT '0' COMMENT '关联ID',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_资产变动表';


--
-- 表的结构 `hg_admin_dept`
--

CREATE TABLE IF NOT EXISTS `hg_admin_dept` (
  `id` bigint(20) NOT NULL COMMENT '部门ID',
  `pid` bigint(20) DEFAULT '0' COMMENT '父部门ID',
  `name` varchar(32) DEFAULT NULL COMMENT '部门名称',
  `code` varchar(255) DEFAULT NULL COMMENT '部门编码',
  `type` varchar(10) DEFAULT NULL COMMENT '部门类型',
  `leader` varchar(32) DEFAULT NULL COMMENT '负责人',
  `phone` varchar(11) DEFAULT NULL COMMENT '联系电话',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `level` int(11) NOT NULL COMMENT '关系树等级',
  `tree` varchar(512) DEFAULT NULL COMMENT '关系树',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '部门状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_部门';


--
-- 表的结构 `hg_admin_member`
--

CREATE TABLE IF NOT EXISTS `hg_admin_member` (
  `id` bigint(20) NOT NULL COMMENT '管理员ID',
  `dept_id` bigint(20) DEFAULT '0' COMMENT '部门ID',
  `role_id` bigint(20) DEFAULT '10' COMMENT '角色ID',
  `real_name` varchar(32) DEFAULT '' COMMENT '真实姓名',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '帐号',
  `password_hash` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(16) NOT NULL COMMENT '密码盐',
  `password_reset_token` varchar(150) DEFAULT '' COMMENT '密码重置令牌',
  `integral` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '积分',
  `balance` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '余额',
  `avatar` char(150) DEFAULT '' COMMENT '头像',
  `sex` tinyint(1) DEFAULT '1' COMMENT '性别',
  `qq` varchar(20) DEFAULT '' COMMENT 'qq',
  `email` varchar(60) DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(20) DEFAULT '' COMMENT '手机号码',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `city_id` bigint(20) DEFAULT '0' COMMENT '城市编码',
  `address` varchar(100) DEFAULT '' COMMENT '联系地址',
  `pid` bigint(20) NOT NULL COMMENT '上级管理员ID',
  `level` int(11) DEFAULT '1' COMMENT '关系树等级',
  `tree` varchar(512) NOT NULL COMMENT '关系树',
  `invite_code` varchar(12) DEFAULT NULL COMMENT '邀请码',
  `cash` json DEFAULT NULL COMMENT '提现配置',
  `last_active_at` datetime DEFAULT NULL COMMENT '最后活跃时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户表';


--
-- 表的结构 `hg_admin_member_post`
--

CREATE TABLE IF NOT EXISTS `hg_admin_member_post` (
  `member_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `post_id` bigint(20) NOT NULL COMMENT '岗位ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户岗位关联';


--
-- 表的结构 `hg_admin_member_role`
--

CREATE TABLE IF NOT EXISTS `hg_admin_member_role` (
  `member_id` bigint(20) NOT NULL COMMENT '管理员ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_用户角色关联';


--
-- 表的结构 `hg_admin_menu`
--

CREATE TABLE IF NOT EXISTS `hg_admin_menu` (
  `id` bigint(20) NOT NULL COMMENT '菜单ID',
  `pid` bigint(20) DEFAULT '0' COMMENT '父菜单ID',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '关系树等级',
  `tree` varchar(255) DEFAULT NULL COMMENT '关系树',
  `title` varchar(64) NOT NULL COMMENT '菜单名称',
  `name` varchar(128) NOT NULL COMMENT '名称编码',
  `path` varchar(200) DEFAULT NULL COMMENT '路由地址',
  `icon` varchar(128) DEFAULT NULL COMMENT '菜单图标',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '菜单类型（1目录 2菜单 3按钮）',
  `redirect` varchar(255) DEFAULT NULL COMMENT '重定向地址',
  `permissions` varchar(512) DEFAULT NULL COMMENT '菜单包含权限集合',
  `permission_name` varchar(64) DEFAULT NULL COMMENT '权限名称',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `always_show` tinyint(1) DEFAULT '0' COMMENT '取消自动计算根路由模式',
  `active_menu` varchar(255) DEFAULT NULL COMMENT '高亮菜单编码',
  `is_root` tinyint(1) DEFAULT '0' COMMENT '是否跟路由',
  `is_frame` tinyint(1) DEFAULT '1' COMMENT '是否内嵌',
  `frame_src` varchar(512) DEFAULT NULL COMMENT '内联外部地址',
  `keep_alive` tinyint(1) DEFAULT '0' COMMENT '缓存该路由',
  `hidden` tinyint(1) DEFAULT '0' COMMENT '是否隐藏',
  `affix` tinyint(1) DEFAULT '0' COMMENT '是否固定',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '菜单状态',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间'
) ENGINE=InnoDB AUTO_INCREMENT=2431 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_菜单权限';


--
-- 表的结构 `hg_admin_notice`
--

CREATE TABLE IF NOT EXISTS `hg_admin_notice` (
  `id` bigint(20) NOT NULL COMMENT '公告ID',
  `title` varchar(64) NOT NULL COMMENT '公告标题',
  `type` bigint(20) NOT NULL COMMENT '公告类型',
  `tag` int(11) DEFAULT NULL COMMENT '标签',
  `content` longtext NOT NULL COMMENT '公告内容',
  `receiver` json DEFAULT NULL COMMENT '接收者',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '公告状态',
  `created_by` bigint(20) DEFAULT NULL COMMENT '发送人',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '修改人',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_通知公告';


--
-- 表的结构 `hg_admin_notice_read`
--

CREATE TABLE IF NOT EXISTS `hg_admin_notice_read` (
  `id` bigint(20) NOT NULL COMMENT '记录ID',
  `notice_id` bigint(20) NOT NULL COMMENT '公告ID',
  `member_id` bigint(20) NOT NULL COMMENT '会员ID',
  `clicks` int(11) DEFAULT '1' COMMENT '已读次数',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `created_at` datetime DEFAULT NULL COMMENT '阅读时间'
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_公告已读记录';


--
-- 表的结构 `hg_admin_oauth`
--

CREATE TABLE IF NOT EXISTS `hg_admin_oauth` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `member_id` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `unionid` varchar(64) DEFAULT '' COMMENT '唯一ID',
  `oauth_client` varchar(32) DEFAULT NULL COMMENT '授权组别',
  `oauth_openid` varchar(128) DEFAULT NULL COMMENT '授权开放ID',
  `sex` tinyint(1) DEFAULT '1' COMMENT '性别',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称',
  `head_portrait` varchar(512) DEFAULT NULL COMMENT '头像',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `country` varchar(100) DEFAULT '' COMMENT '国家',
  `province` varchar(100) DEFAULT '' COMMENT '省',
  `city` varchar(100) DEFAULT '' COMMENT '市',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(1) DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_第三方登录';

-- --------------------------------------------------------

--
-- 表的结构 `hg_admin_order`
--

CREATE TABLE IF NOT EXISTS `hg_admin_order` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `member_id` bigint(20) DEFAULT '0' COMMENT '管理员id',
  `order_type` varchar(32) NOT NULL COMMENT '订单类型',
  `product_id` bigint(20) DEFAULT NULL COMMENT '产品id',
  `order_sn` varchar(64) DEFAULT '' COMMENT '关联订单号',
  `money` decimal(10,2) NOT NULL COMMENT '充值金额',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `refund_reason` varchar(255) DEFAULT NULL COMMENT '退款原因',
  `reject_refund_reason` varchar(255) DEFAULT NULL COMMENT '拒绝退款原因',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_充值订单';

-- --------------------------------------------------------

--
-- 表的结构 `hg_admin_post`
--

CREATE TABLE IF NOT EXISTS `hg_admin_post` (
  `id` bigint(20) NOT NULL COMMENT '岗位ID',
  `code` varchar(64) NOT NULL COMMENT '岗位编码',
  `name` varchar(50) NOT NULL COMMENT '岗位名称',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `sort` int(11) NOT NULL COMMENT '排序',
  `status` tinyint(1) NOT NULL COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_岗位';


--
-- 表的结构 `hg_admin_role`
--

CREATE TABLE IF NOT EXISTS `hg_admin_role` (
  `id` bigint(20) NOT NULL COMMENT '角色ID',
  `name` varchar(32) NOT NULL COMMENT '角色名称',
  `key` varchar(128) NOT NULL COMMENT '角色权限字符串',
  `data_scope` tinyint(1) DEFAULT '1' COMMENT '数据范围',
  `custom_dept` json DEFAULT NULL COMMENT '自定义部门权限',
  `pid` bigint(20) DEFAULT '0' COMMENT '上级角色ID',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '关系树等级',
  `tree` varchar(512) DEFAULT NULL COMMENT '关系树',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '角色状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=211 DEFAULT CHARSET=utf8mb4 COMMENT='管理员_角色信息';


--
-- 表的结构 `hg_admin_role_casbin`
--

CREATE TABLE IF NOT EXISTS `hg_admin_role_casbin` (
  `id` bigint(20) NOT NULL,
  `p_type` varchar(64) DEFAULT NULL,
  `v0` varchar(256) DEFAULT NULL,
  `v1` varchar(256) DEFAULT NULL,
  `v2` varchar(256) DEFAULT NULL,
  `v3` varchar(256) DEFAULT NULL,
  `v4` varchar(256) DEFAULT NULL,
  `v5` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='管理员_casbin权限表';

-- --------------------------------------------------------

--
-- 表的结构 `hg_admin_role_menu`
--

CREATE TABLE IF NOT EXISTS `hg_admin_role_menu` (
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) NOT NULL COMMENT '菜单ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员_角色菜单关联';


--
-- 表的结构 `hg_pay_log`
--

CREATE TABLE IF NOT EXISTS `hg_pay_log` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `member_id` bigint(20) DEFAULT '0' COMMENT '会员ID',
  `app_id` varchar(50) DEFAULT NULL COMMENT '应用ID',
  `addons_name` varchar(100) DEFAULT '' COMMENT '插件名称',
  `order_sn` varchar(64) DEFAULT '' COMMENT '关联订单号',
  `order_group` varchar(32) DEFAULT '' COMMENT '组别[默认统一支付类型]',
  `openid` varchar(50) DEFAULT '' COMMENT 'openid',
  `mch_id` varchar(20) DEFAULT '' COMMENT '商户支付账户',
  `subject` varchar(255) DEFAULT NULL COMMENT '订单标题',
  `detail` json DEFAULT NULL COMMENT '支付商品详情',
  `auth_code` varchar(50) DEFAULT '' COMMENT '刷卡码',
  `out_trade_no` varchar(128) DEFAULT '' COMMENT '商户订单号',
  `transaction_id` varchar(128) DEFAULT NULL COMMENT '交易号',
  `pay_type` varchar(32) NOT NULL COMMENT '支付类型',
  `pay_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '支付金额',
  `actual_amount` decimal(10,2) DEFAULT NULL COMMENT '实付金额',
  `pay_status` tinyint(4) DEFAULT '0' COMMENT '支付状态',
  `pay_at` datetime DEFAULT NULL COMMENT '支付时间',
  `trade_type` varchar(16) DEFAULT '' COMMENT '交易类型',
  `refund_sn` varchar(128) DEFAULT NULL COMMENT '退款单号',
  `is_refund` tinyint(4) DEFAULT '0' COMMENT '是否退款 ',
  `custom` text COMMENT '自定义参数',
  `create_ip` varchar(128) DEFAULT NULL COMMENT '创建者IP',
  `pay_ip` varchar(128) DEFAULT NULL COMMENT '支付者IP',
  `notify_url` varchar(255) DEFAULT NULL COMMENT '支付通知回调地址',
  `return_url` varchar(255) DEFAULT NULL COMMENT '买家付款成功跳转地址',
  `trace_ids` json DEFAULT NULL COMMENT '链路ID集合',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='支付_支付日志';

-- --------------------------------------------------------

--
-- 表的结构 `hg_pay_refund`
--

CREATE TABLE IF NOT EXISTS `hg_pay_refund` (
  `id` bigint(20) unsigned NOT NULL COMMENT '主键ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '会员ID',
  `app_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '应用ID',
  `order_sn` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '业务订单号',
  `refund_trade_no` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '退款交易号',
  `refund_money` decimal(10,2) DEFAULT NULL COMMENT '退款金额',
  `refund_way` tinyint(4) DEFAULT '1' COMMENT '退款方式',
  `ip` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '申请者IP',
  `reason` varchar(255) DEFAULT NULL COMMENT '申请退款原因',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '退款备注',
  `status` tinyint(4) DEFAULT '1' COMMENT '退款状态',
  `created_at` datetime DEFAULT NULL COMMENT '申请时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付_退款记录';

-- --------------------------------------------------------

--
-- 表的结构 `hg_sys_addons_config`
--

CREATE TABLE IF NOT EXISTS `hg_sys_addons_config` (
  `id` bigint(20) NOT NULL COMMENT '配置ID',
  `addon_name` varchar(128) NOT NULL COMMENT '插件名称',
  `group` varchar(128) NOT NULL COMMENT '分组',
  `name` varchar(100) DEFAULT '' COMMENT '参数名称',
  `type` varchar(32) NOT NULL COMMENT '键值类型:string,int,uint,bool,datetime,date',
  `key` varchar(100) DEFAULT '' COMMENT '参数键名',
  `value` varchar(500) DEFAULT '' COMMENT '参数键值',
  `default_value` varchar(500) NOT NULL COMMENT '默认值',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `tip` varchar(500) DEFAULT NULL COMMENT '变量描述',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否为系统默认',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统_插件配置';


--
-- 表的结构 `hg_sys_addons_install`
--

CREATE TABLE IF NOT EXISTS `hg_sys_addons_install` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '插件名称',
  `version` varchar(128) NOT NULL DEFAULT '' COMMENT '版本号',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统_插件安装记录';


--
-- 表的结构 `hg_sys_attachment`
--

CREATE TABLE IF NOT EXISTS `hg_sys_attachment` (
  `id` bigint(20) NOT NULL COMMENT '文件ID',
  `app_id` varchar(64) NOT NULL COMMENT '应用ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '管理员ID',
  `cate_id` bigint(20) unsigned DEFAULT '0' COMMENT '上传分类',
  `drive` varchar(64) DEFAULT NULL COMMENT '上传驱动',
  `name` varchar(1000) DEFAULT NULL COMMENT '文件原始名',
  `kind` varchar(16) DEFAULT NULL COMMENT '上传类型',
  `mime_type` varchar(128) NOT NULL DEFAULT '' COMMENT '扩展类型',
  `naive_type` varchar(32) DEFAULT NULL COMMENT 'NaiveUI类型',
  `path` varchar(1000) DEFAULT NULL COMMENT '本地路径',
  `file_url` varchar(1000) DEFAULT NULL COMMENT 'url',
  `size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `ext` varchar(50) DEFAULT NULL COMMENT '扩展名',
  `md5` varchar(32) DEFAULT NULL COMMENT 'md5校验码',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='系统_附件管理';

-- --------------------------------------------------------

--
-- 表的结构 `hg_sys_blacklist`
--

CREATE TABLE IF NOT EXISTS `hg_sys_blacklist` (
  `id` bigint(20) NOT NULL COMMENT '黑名单ID',
  `ip` varchar(100) DEFAULT '' COMMENT 'IP地址',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='系统_访问黑名单';


--
-- 表的结构 `hg_sys_config`
--

CREATE TABLE IF NOT EXISTS `hg_sys_config` (
  `id` bigint(20) NOT NULL COMMENT '配置ID',
  `group` varchar(128) NOT NULL COMMENT '配置分组',
  `name` varchar(100) DEFAULT '' COMMENT '参数名称',
  `type` varchar(32) NOT NULL COMMENT '键值类型:string,int,uint,bool,datetime,date',
  `key` varchar(100) DEFAULT '' COMMENT '参数键名',
  `value` longtext COMMENT '参数键值',
  `default_value` varchar(500) NOT NULL COMMENT '默认值',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `tip` varchar(500) DEFAULT NULL COMMENT '变量描述',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否为系统默认',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=129 DEFAULT CHARSET=utf8mb4 COMMENT='系统_配置';


--
-- 表的结构 `hg_sys_cron`
--

CREATE TABLE IF NOT EXISTS `hg_sys_cron` (
  `id` bigint(20) NOT NULL COMMENT '任务ID',
  `group_id` bigint(20) NOT NULL COMMENT '分组ID',
  `title` varchar(128) NOT NULL COMMENT '任务标题',
  `name` varchar(100) DEFAULT NULL COMMENT '任务方法',
  `params` varchar(255) DEFAULT NULL COMMENT '函数参数',
  `pattern` varchar(64) NOT NULL COMMENT '表达式',
  `policy` bigint(20) NOT NULL DEFAULT '1' COMMENT '策略',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT '执行次数',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '任务状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COMMENT='系统_定时任务';


--
-- 表的结构 `hg_sys_cron_group`
--

CREATE TABLE IF NOT EXISTS `hg_sys_cron_group` (
  `id` bigint(20) NOT NULL COMMENT '任务分组ID',
  `pid` bigint(20) NOT NULL COMMENT '父类任务分组ID',
  `name` varchar(100) DEFAULT '' COMMENT '分组名称',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '分组状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='系统_定时任务分组';


--
-- 表的结构 `hg_sys_dict_data`
--

CREATE TABLE IF NOT EXISTS `hg_sys_dict_data` (
  `id` bigint(20) NOT NULL COMMENT '字典数据ID',
  `label` varchar(100) DEFAULT NULL COMMENT '字典标签',
  `value` varchar(100) DEFAULT NULL COMMENT '字典键值',
  `value_type` varchar(255) NOT NULL DEFAULT 'string' COMMENT '键值数据类型：string,int,uint,bool,datetime,date',
  `type` varchar(100) DEFAULT NULL COMMENT '字典类型',
  `list_class` varchar(100) DEFAULT NULL COMMENT '表格回显样式',
  `is_default` tinyint(1) DEFAULT '2' COMMENT '是否为系统默认',
  `sort` int(11) DEFAULT '0' COMMENT '字典排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=171 DEFAULT CHARSET=utf8mb4 COMMENT='系统_字典数据';


--
-- 表的结构 `hg_sys_dict_type`
--

CREATE TABLE IF NOT EXISTS `hg_sys_dict_type` (
  `id` bigint(20) NOT NULL COMMENT '字典类型ID',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '父类字典类型ID',
  `name` varchar(100) DEFAULT '' COMMENT '字典类型名称',
  `type` varchar(100) DEFAULT '' COMMENT '字典类型',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '字典类型状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8mb4 COMMENT='系统_字典类型';


--
-- 表的结构 `hg_sys_ems_log`
--

CREATE TABLE IF NOT EXISTS `hg_sys_ems_log` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `event` varchar(64) NOT NULL COMMENT '事件',
  `email` varchar(512) NOT NULL COMMENT '邮箱地址，多个用;隔开',
  `code` varchar(256) DEFAULT '' COMMENT '验证码',
  `times` bigint(20) NOT NULL COMMENT '验证次数',
  `content` longtext COMMENT '邮件内容',
  `ip` varchar(128) DEFAULT NULL COMMENT 'ip地址',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态(1未验证,2已验证)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='系统_邮件发送记录';


--
-- 表的结构 `hg_sys_gen_codes`
--

CREATE TABLE IF NOT EXISTS `hg_sys_gen_codes` (
  `id` bigint(20) NOT NULL COMMENT '生成ID',
  `gen_type` int(10) unsigned NOT NULL COMMENT '生成类型',
  `gen_template` int(11) DEFAULT '0' COMMENT '生成模板',
  `var_name` varchar(255) NOT NULL COMMENT '实体命名',
  `options` json DEFAULT NULL COMMENT '配置选项',
  `db_name` varchar(128) DEFAULT NULL COMMENT '数据库名称',
  `table_name` varchar(255) NOT NULL COMMENT '主表名称',
  `table_comment` varchar(255) DEFAULT NULL COMMENT '主表注释',
  `dao_name` varchar(255) DEFAULT NULL COMMENT '主表dao模型',
  `master_columns` json DEFAULT NULL COMMENT '主表字段',
  `addon_name` varchar(128) DEFAULT NULL COMMENT '插件名称',
  `status` tinyint(1) DEFAULT '1' COMMENT '生成状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COMMENT='系统_代码生成记录';


--
-- 表的结构 `hg_sys_gen_curd_demo`
--

CREATE TABLE IF NOT EXISTS `hg_sys_gen_curd_demo` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `category_id` bigint(20) DEFAULT '0' COMMENT '分类ID',
  `title` varchar(64) NOT NULL COMMENT '标题',
  `description` varchar(255) DEFAULT '' COMMENT '描述',
  `content` text COMMENT '内容',
  `image` varchar(255) DEFAULT NULL COMMENT '单图',
  `attachfile` varchar(255) DEFAULT NULL COMMENT '附件',
  `city_id` bigint(20) DEFAULT '0' COMMENT '所在城市',
  `switch` int(11) DEFAULT '1' COMMENT '显示开关',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `deleted_by` bigint(20) DEFAULT '0' COMMENT '删除者',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COMMENT='系统_生成curd演示';


--
-- 表的结构 `hg_sys_gen_tree_demo`
--

CREATE TABLE IF NOT EXISTS `hg_sys_gen_tree_demo` (
  `id` bigint(20) NOT NULL COMMENT 'ID',
  `pid` bigint(20) DEFAULT NULL COMMENT '上级ID',
  `level` int(11) DEFAULT '1' COMMENT '关系树级别',
  `tree` varchar(512) DEFAULT NULL COMMENT '关系树',
  `category_id` bigint(20) DEFAULT '0' COMMENT '分类ID',
  `title` varchar(64) NOT NULL COMMENT '标题',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_by` bigint(20) DEFAULT '0' COMMENT '创建者',
  `updated_by` bigint(20) DEFAULT '0' COMMENT '更新者',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COMMENT='系统_生成树表演示';


--
-- 表的结构 `hg_sys_log`
--

CREATE TABLE IF NOT EXISTS `hg_sys_log` (
  `id` bigint(20) NOT NULL COMMENT '日志ID',
  `req_id` varchar(50) DEFAULT NULL COMMENT '对外ID',
  `app_id` varchar(50) DEFAULT '' COMMENT '应用ID',
  `merchant_id` bigint(20) unsigned DEFAULT '0' COMMENT '商户ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `method` varchar(20) DEFAULT NULL COMMENT '提交类型',
  `module` varchar(50) DEFAULT NULL COMMENT '访问模块',
  `url` varchar(1000) DEFAULT NULL COMMENT '提交url',
  `get_data` json DEFAULT NULL COMMENT 'get数据',
  `post_data` json DEFAULT NULL COMMENT 'post数据',
  `header_data` json DEFAULT NULL COMMENT 'header数据',
  `ip` varchar(128) DEFAULT NULL COMMENT 'IP地址',
  `province_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '省编码',
  `city_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '市编码',
  `error_code` int(11) DEFAULT '0' COMMENT '报错code',
  `error_msg` longtext COMMENT '对外错误提示',
  `error_data` json DEFAULT NULL COMMENT '报错日志',
  `user_agent` varchar(512) DEFAULT NULL COMMENT 'UA信息',
  `take_up_time` bigint(20) DEFAULT '0' COMMENT '请求耗时',
  `timestamp` bigint(20) DEFAULT '0' COMMENT '响应时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_全局日志';

-- --------------------------------------------------------

--
-- 表的结构 `hg_sys_login_log`
--

CREATE TABLE IF NOT EXISTS `hg_sys_login_log` (
  `id` bigint(20) NOT NULL COMMENT '日志ID',
  `req_id` varchar(50) DEFAULT NULL COMMENT '请求ID',
  `member_id` bigint(20) DEFAULT '0' COMMENT '用户ID',
  `username` varchar(64) DEFAULT NULL COMMENT '用户名',
  `response` json DEFAULT NULL COMMENT '响应数据',
  `login_at` datetime DEFAULT NULL COMMENT '登录时间',
  `login_ip` varchar(128) DEFAULT NULL COMMENT '登录IP',
  `province_id` bigint(20) DEFAULT NULL COMMENT '省编码',
  `city_id` bigint(20) DEFAULT NULL COMMENT '市编码',
  `user_agent` varchar(512) DEFAULT NULL COMMENT 'UA信息',
  `err_msg` varchar(1000) DEFAULT NULL COMMENT '错误提示',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_登录日志';

-- --------------------------------------------------------

--
-- 表的结构 `hg_sys_provinces`
--

CREATE TABLE IF NOT EXISTS `hg_sys_provinces` (
  `id` bigint(20) NOT NULL COMMENT '省市区ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '栏目名称',
  `pinyin` varchar(100) DEFAULT '' COMMENT '拼音',
  `lng` varchar(20) DEFAULT '' COMMENT '经度',
  `lat` varchar(20) DEFAULT '' COMMENT '纬度',
  `pid` bigint(20) NOT NULL DEFAULT '0' COMMENT '父栏目',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '关系树等级',
  `tree` varchar(200) DEFAULT NULL COMMENT '关系树',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='系统_省市区编码';


--
-- 表的结构 `hg_sys_serve_license`
--

CREATE TABLE IF NOT EXISTS `hg_sys_serve_license` (
  `id` bigint(20) NOT NULL COMMENT '许可ID',
  `group` varchar(50) NOT NULL COMMENT '分组',
  `name` varchar(128) NOT NULL COMMENT '许可名称',
  `appid` varchar(64) NOT NULL COMMENT '应用ID',
  `secret_key` varchar(255) DEFAULT NULL COMMENT '应用秘钥',
  `remote_addr` varchar(64) DEFAULT NULL COMMENT '最后连接地址',
  `online_limit` int(11) DEFAULT '1' COMMENT '在线限制',
  `login_times` bigint(20) DEFAULT NULL COMMENT '登录次数',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_active_at` datetime DEFAULT NULL COMMENT '最后心跳',
  `routes` json DEFAULT NULL COMMENT '路由表，空使用默认分组路由',
  `allowed_ips` varchar(512) DEFAULT NULL COMMENT 'IP白名单',
  `end_at` datetime NOT NULL COMMENT '授权有效期',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='系统_服务许可证';


--
-- 表的结构 `hg_sys_serve_log`
--

CREATE TABLE IF NOT EXISTS `hg_sys_serve_log` (
  `id` bigint(20) NOT NULL COMMENT '日志ID',
  `trace_id` varchar(50) DEFAULT NULL COMMENT '链路ID',
  `level_format` varchar(32) DEFAULT NULL COMMENT '日志级别',
  `content` text COMMENT '日志内容',
  `stack` json DEFAULT NULL COMMENT '打印堆栈',
  `line` varchar(255) NOT NULL COMMENT '调用行',
  `trigger_ns` bigint(20) DEFAULT NULL COMMENT '触发时间(ns)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统_服务日志';

-- --------------------------------------------------------

--
-- 表的结构 `hg_sys_sms_log`
--

CREATE TABLE IF NOT EXISTS `hg_sys_sms_log` (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `event` varchar(64) NOT NULL COMMENT '事件',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `code` varchar(256) DEFAULT '' COMMENT '验证码或短信内容',
  `times` bigint(20) NOT NULL COMMENT '验证次数',
  `ip` varchar(128) DEFAULT NULL COMMENT 'ip地址',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态(1未验证,2已验证)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统_短信发送记录';


--
-- 表的结构 `hg_test_category`
--

CREATE TABLE IF NOT EXISTS `hg_test_category` (
  `id` bigint(20) NOT NULL COMMENT '分类ID',
  `name` varchar(255) NOT NULL COMMENT '分类名称',
  `short_name` varchar(128) DEFAULT NULL COMMENT '简称',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `sort` int(11) NOT NULL COMMENT '排序',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='测试分类';
