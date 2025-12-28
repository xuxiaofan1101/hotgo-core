--
-- Indexes for dumped tables
--

--
-- Indexes for table `hg_addon_hgexample_table`
--
ALTER TABLE `hg_addon_hgexample_table`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hg_addon_hgexample_tenant_order`
--
ALTER TABLE `hg_addon_hgexample_tenant_order`
  ADD PRIMARY KEY (`id`),
  ADD KEY `order_sn` (`order_sn`),
  ADD KEY `member_id` (`user_id`),
  ADD KEY `merchant_id` (`merchant_id`),
  ADD KEY `agent_id` (`tenant_id`);

--
-- Indexes for table `hg_admin_cash`
--
ALTER TABLE `hg_admin_cash`
  ADD PRIMARY KEY (`id`),
  ADD KEY `admin_id` (`member_id`);

--
-- Indexes for table `hg_admin_credits_log`
--
ALTER TABLE `hg_admin_credits_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `member_id` (`member_id`);

--
-- Indexes for table `hg_admin_dept`
--
ALTER TABLE `hg_admin_dept`
  ADD PRIMARY KEY (`id`),
  ADD KEY `pid` (`pid`);

--
-- Indexes for table `hg_admin_member`
--
ALTER TABLE `hg_admin_member`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `invite_code` (`invite_code`),
  ADD KEY `dept_id` (`dept_id`),
  ADD KEY `pid` (`pid`);

--
-- Indexes for table `hg_admin_member_post`
--
ALTER TABLE `hg_admin_member_post`
  ADD PRIMARY KEY (`member_id`,`post_id`);

--
-- Indexes for table `hg_admin_member_role`
--
ALTER TABLE `hg_admin_member_role`
  ADD PRIMARY KEY (`member_id`,`role_id`);

--
-- Indexes for table `hg_admin_menu`
--
ALTER TABLE `hg_admin_menu`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD KEY `pid` (`pid`),
  ADD KEY `status` (`status`),
  ADD KEY `type` (`type`);

--
-- Indexes for table `hg_admin_notice`
--
ALTER TABLE `hg_admin_notice`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hg_admin_notice_read`
--
ALTER TABLE `hg_admin_notice_read`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `notice_id` (`notice_id`,`member_id`);

--
-- Indexes for table `hg_admin_oauth`
--
ALTER TABLE `hg_admin_oauth`
  ADD PRIMARY KEY (`id`),
  ADD KEY `oauth_client` (`oauth_client`,`oauth_openid`),
  ADD KEY `member_id` (`member_id`);

--
-- Indexes for table `hg_admin_order`
--
ALTER TABLE `hg_admin_order`
  ADD PRIMARY KEY (`id`),
  ADD KEY `order_sn` (`order_sn`),
  ADD KEY `member_id` (`member_id`);

--
-- Indexes for table `hg_admin_post`
--
ALTER TABLE `hg_admin_post`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hg_admin_role`
--
ALTER TABLE `hg_admin_role`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hg_admin_role_casbin`
--
ALTER TABLE `hg_admin_role_casbin`
  ADD PRIMARY KEY (`id`) USING BTREE;

--
-- Indexes for table `hg_admin_role_menu`
--
ALTER TABLE `hg_admin_role_menu`
  ADD PRIMARY KEY (`role_id`,`menu_id`);

--
-- Indexes for table `hg_pay_log`
--
ALTER TABLE `hg_pay_log`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `order_sn` (`order_sn`),
  ADD KEY `member_id` (`member_id`);

--
-- Indexes for table `hg_pay_refund`
--
ALTER TABLE `hg_pay_refund`
  ADD PRIMARY KEY (`id`),
  ADD KEY `order_sn` (`order_sn`);

--
-- Indexes for table `hg_sys_addons_config`
--
ALTER TABLE `hg_sys_addons_config`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `addon_name_2` (`addon_name`,`key`),
  ADD KEY `addon_name` (`addon_name`),
  ADD KEY `addon_name_3` (`addon_name`,`group`);

--
-- Indexes for table `hg_sys_addons_install`
--
ALTER TABLE `hg_sys_addons_install`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD UNIQUE KEY `name` (`name`);

--
-- Indexes for table `hg_sys_attachment`
--
ALTER TABLE `hg_sys_attachment`
  ADD PRIMARY KEY (`id`),
  ADD KEY `md5` (`md5`);

--
-- Indexes for table `hg_sys_blacklist`
--
ALTER TABLE `hg_sys_blacklist`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD UNIQUE KEY `name` (`ip`);

--
-- Indexes for table `hg_sys_config`
--
ALTER TABLE `hg_sys_config`
  ADD PRIMARY KEY (`id`),
  ADD KEY `group` (`group`),
  ADD KEY `key` (`key`);

--
-- Indexes for table `hg_sys_cron`
--
ALTER TABLE `hg_sys_cron`
  ADD PRIMARY KEY (`id`) USING BTREE;

--
-- Indexes for table `hg_sys_cron_group`
--
ALTER TABLE `hg_sys_cron_group`
  ADD PRIMARY KEY (`id`) USING BTREE;

--
-- Indexes for table `hg_sys_dict_data`
--
ALTER TABLE `hg_sys_dict_data`
  ADD PRIMARY KEY (`id`),
  ADD KEY `dict_data_idx` (`type`);

--
-- Indexes for table `hg_sys_dict_type`
--
ALTER TABLE `hg_sys_dict_type`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD UNIQUE KEY `dict_type` (`type`);

--
-- Indexes for table `hg_sys_ems_log`
--
ALTER TABLE `hg_sys_ems_log`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD KEY `email` (`email`);

--
-- Indexes for table `hg_sys_gen_codes`
--
ALTER TABLE `hg_sys_gen_codes`
  ADD PRIMARY KEY (`id`) USING BTREE;

--
-- Indexes for table `hg_sys_gen_curd_demo`
--
ALTER TABLE `hg_sys_gen_curd_demo`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hg_sys_gen_tree_demo`
--
ALTER TABLE `hg_sys_gen_tree_demo`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hg_sys_log`
--
ALTER TABLE `hg_sys_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `error_code` (`error_code`),
  ADD KEY `req_id` (`req_id`),
  ADD KEY `member_id` (`member_id`);

--
-- Indexes for table `hg_sys_login_log`
--
ALTER TABLE `hg_sys_login_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `member_id` (`member_id`),
  ADD KEY `req_id` (`req_id`);

--
-- Indexes for table `hg_sys_provinces`
--
ALTER TABLE `hg_sys_provinces`
  ADD PRIMARY KEY (`id`),
  ADD KEY `pid` (`pid`);

--
-- Indexes for table `hg_sys_serve_license`
--
ALTER TABLE `hg_sys_serve_license`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `appid` (`appid`);

--
-- Indexes for table `hg_sys_serve_log`
--
ALTER TABLE `hg_sys_serve_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `member_id` (`level_format`),
  ADD KEY `traceid` (`trace_id`);

--
-- Indexes for table `hg_sys_sms_log`
--
ALTER TABLE `hg_sys_sms_log`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD KEY `mobile` (`mobile`);

--
-- Indexes for table `hg_test_category`
--
ALTER TABLE `hg_test_category`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `hg_addon_hgexample_table`
--
ALTER TABLE `hg_addon_hgexample_table`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',AUTO_INCREMENT=7;
--
-- AUTO_INCREMENT for table `hg_addon_hgexample_tenant_order`
--
ALTER TABLE `hg_addon_hgexample_tenant_order`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_admin_cash`
--
ALTER TABLE `hg_admin_cash`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_admin_credits_log`
--
ALTER TABLE `hg_admin_credits_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '变动ID',AUTO_INCREMENT=8;
--
-- AUTO_INCREMENT for table `hg_admin_dept`
--
ALTER TABLE `hg_admin_dept`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '部门ID',AUTO_INCREMENT=113;
--
-- AUTO_INCREMENT for table `hg_admin_member`
--
ALTER TABLE `hg_admin_member`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '管理员ID',AUTO_INCREMENT=14;
--
-- AUTO_INCREMENT for table `hg_admin_menu`
--
ALTER TABLE `hg_admin_menu`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',AUTO_INCREMENT=2431;
--
-- AUTO_INCREMENT for table `hg_admin_notice`
--
ALTER TABLE `hg_admin_notice`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '公告ID',AUTO_INCREMENT=33;
--
-- AUTO_INCREMENT for table `hg_admin_notice_read`
--
ALTER TABLE `hg_admin_notice_read`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '记录ID',AUTO_INCREMENT=9;
--
-- AUTO_INCREMENT for table `hg_admin_oauth`
--
ALTER TABLE `hg_admin_oauth`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键';
--
-- AUTO_INCREMENT for table `hg_admin_order`
--
ALTER TABLE `hg_admin_order`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_admin_post`
--
ALTER TABLE `hg_admin_post`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '岗位ID',AUTO_INCREMENT=7;
--
-- AUTO_INCREMENT for table `hg_admin_role`
--
ALTER TABLE `hg_admin_role`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',AUTO_INCREMENT=211;
--
-- AUTO_INCREMENT for table `hg_admin_role_casbin`
--
ALTER TABLE `hg_admin_role_casbin`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;
--
-- AUTO_INCREMENT for table `hg_pay_log`
--
ALTER TABLE `hg_pay_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_pay_refund`
--
ALTER TABLE `hg_pay_refund`
  MODIFY `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID';
--
-- AUTO_INCREMENT for table `hg_sys_addons_config`
--
ALTER TABLE `hg_sys_addons_config`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_sys_addons_install`
--
ALTER TABLE `hg_sys_addons_install`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_sys_attachment`
--
ALTER TABLE `hg_sys_attachment`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',AUTO_INCREMENT=9;
--
-- AUTO_INCREMENT for table `hg_sys_blacklist`
--
ALTER TABLE `hg_sys_blacklist`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '黑名单ID',AUTO_INCREMENT=8;
--
-- AUTO_INCREMENT for table `hg_sys_config`
--
ALTER TABLE `hg_sys_config`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',AUTO_INCREMENT=129;
--
-- AUTO_INCREMENT for table `hg_sys_cron`
--
ALTER TABLE `hg_sys_cron`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '任务ID',AUTO_INCREMENT=11;
--
-- AUTO_INCREMENT for table `hg_sys_cron_group`
--
ALTER TABLE `hg_sys_cron_group`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '任务分组ID',AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `hg_sys_dict_data`
--
ALTER TABLE `hg_sys_dict_data`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',AUTO_INCREMENT=171;
--
-- AUTO_INCREMENT for table `hg_sys_dict_type`
--
ALTER TABLE `hg_sys_dict_type`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',AUTO_INCREMENT=45;
--
-- AUTO_INCREMENT for table `hg_sys_ems_log`
--
ALTER TABLE `hg_sys_ems_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',AUTO_INCREMENT=5;
--
-- AUTO_INCREMENT for table `hg_sys_gen_codes`
--
ALTER TABLE `hg_sys_gen_codes`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '生成ID',AUTO_INCREMENT=12;
--
-- AUTO_INCREMENT for table `hg_sys_gen_curd_demo`
--
ALTER TABLE `hg_sys_gen_curd_demo`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',AUTO_INCREMENT=16;
--
-- AUTO_INCREMENT for table `hg_sys_gen_tree_demo`
--
ALTER TABLE `hg_sys_gen_tree_demo`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',AUTO_INCREMENT=41;
--
-- AUTO_INCREMENT for table `hg_sys_log`
--
ALTER TABLE `hg_sys_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID';
--
-- AUTO_INCREMENT for table `hg_sys_login_log`
--
ALTER TABLE `hg_sys_login_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID';
--
-- AUTO_INCREMENT for table `hg_sys_serve_license`
--
ALTER TABLE `hg_sys_serve_license`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '许可ID',AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `hg_sys_serve_log`
--
ALTER TABLE `hg_sys_serve_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID';
--
-- AUTO_INCREMENT for table `hg_sys_sms_log`
--
ALTER TABLE `hg_sys_sms_log`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT for table `hg_test_category`
--
ALTER TABLE `hg_test_category`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分类ID',AUTO_INCREMENT=5;


ALTER TABLE `hg_addon_hgexample_table`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_addon_hgexample_tenant_order`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_cash`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_credits_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_dept`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_member`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_member_post`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_member_role`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_menu`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_notice`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_notice_read`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_oauth`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_order`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_post`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_role`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_role_casbin`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_admin_role_menu`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_pay_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_pay_refund`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_addons_config`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_addons_install`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_attachment`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_blacklist`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_config`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_cron`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_cron_group`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_dict_data`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_dict_type`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_ems_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_gen_codes`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_gen_curd_demo`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_gen_tree_demo`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_login_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_provinces`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_serve_license`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_serve_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_sys_sms_log`
    AUTO_INCREMENT = 10000;
ALTER TABLE `hg_test_category`
    AUTO_INCREMENT = 10000;
