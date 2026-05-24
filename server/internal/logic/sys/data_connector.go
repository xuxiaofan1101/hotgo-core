// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"strings"

	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/library/hgorm/hook"
	"hotgo/internal/model/do"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sSysDataConnector struct{}

func NewSysDataConnector() *sSysDataConnector {
	return &sSysDataConnector{}
}

func init() {
	service.RegisterSysDataConnector(NewSysDataConnector())
}

// Model 数据源ORM模型
func (s *sSysDataConnector) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.DataConnector.Ctx(ctx), option...)
}

// List 获取数据源列表
func (s *sSysDataConnector) List(ctx context.Context, in *sysin.DataConnectorListInp) (list []*sysin.DataConnectorListModel, totalCount int, err error) {
	mod := s.Model(ctx)

	// 字段过滤
	mod = mod.FieldsPrefix(dao.DataConnector.Table(), sysin.DataConnectorListModel{})

	// 查询连接ID
	if in.Id > 0 {
		mod = mod.Where(dao.DataConnector.Columns().Id, in.Id)
	}

	// 查询关键词
	if in.Keyword != "" {
		cols := dao.DataConnector.Columns()
		keyword := "%" + in.Keyword + "%"
		mod = mod.Where("(`"+cols.Name+"` LIKE ? OR `"+cols.Code+"` LIKE ?)", keyword, keyword)
	}

	// 查询连接名称
	if in.Name != "" {
		mod = mod.WhereLike(dao.DataConnector.Columns().Name, "%"+in.Name+"%")
	}

	// 查询连接编码
	if in.Code != "" {
		mod = mod.WhereLike(dao.DataConnector.Columns().Code, "%"+in.Code+"%")
	}

	// 查询连接方向
	if in.Direction != "" {
		mod = mod.Where(dao.DataConnector.Columns().Direction, in.Direction)
	}

	// 查询连接类型
	if in.ConnectorType != "" {
		mod = mod.Where(dao.DataConnector.Columns().ConnectorType, in.ConnectorType)
	}

	// 查询状态
	if in.Status > 0 {
		mod = mod.Where(dao.DataConnector.Columns().Status, in.Status)
	}

	// 查询创建者
	if in.CreatedBy != "" {
		ids, err := service.AdminMember().GetIdsByKeyword(ctx, in.CreatedBy)
		if err != nil {
			return nil, 0, err
		}
		mod = mod.WhereIn(dao.DataConnector.Columns().CreatedBy, ids)
	}

	// 查询创建时间
	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(dao.DataConnector.Columns().CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	// 分页
	mod = mod.Page(in.Page, in.PerPage)

	// 排序
	mod = mod.OrderDesc(dao.DataConnector.Columns().Id)

	// 操作人摘要信息
	mod = mod.Hook(hook.MemberSummary)

	// 查询数据
	if err = mod.ScanAndCount(&list, &totalCount, false); err != nil {
		err = gerror.Wrap(err, "获取数据源列表失败，请稍后重试！")
		return
	}
	return
}

// Edit 修改/新增数据源
func (s *sSysDataConnector) Edit(ctx context.Context, in *sysin.DataConnectorEditInp) (err error) {
	// 验证连接编码唯一
	if err = hgorm.IsUnique(ctx, &dao.DataConnector, g.Map{dao.DataConnector.Columns().Code: in.Code}, "连接编码已存在", in.Id); err != nil {
		return
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		// 修改
		if in.Id > 0 {
			in.UpdatedBy = contexts.GetUserId(ctx)
			if _, err = s.Model(ctx).
				Fields(sysin.DataConnectorUpdateFields{}).
				WherePri(in.Id).Data(in).Update(); err != nil {
				err = gerror.Wrap(err, "修改数据源失败，请稍后重试！")
			}
			return
		}

		// 新增
		in.CreatedBy = contexts.GetUserId(ctx)
		if _, err = s.Model(ctx, &handler.Option{FilterAuth: false}).
			Fields(sysin.DataConnectorInsertFields{}).
			Data(in).OmitEmptyData().Insert(); err != nil {
			err = gerror.Wrap(err, "新增数据源失败，请稍后重试！")
		}
		return
	})
}

// Delete 删除数据源
func (s *sSysDataConnector) Delete(ctx context.Context, in *sysin.DataConnectorDeleteInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(do.DataConnector{
		DeletedBy: contexts.GetUserId(ctx),
		DeletedAt: gtime.Now(),
	}).Update(); err != nil {
		err = gerror.Wrap(err, "删除数据源失败，请稍后重试！")
		return
	}
	return
}

// View 获取数据源指定信息
func (s *sSysDataConnector) View(ctx context.Context, in *sysin.DataConnectorViewInp) (res *sysin.DataConnectorViewModel, err error) {
	if err = s.Model(ctx).WherePri(in.Id).Hook(hook.MemberSummary).Scan(&res); err != nil {
		err = gerror.Wrap(err, "获取数据源信息，请稍后重试！")
		return
	}
	return
}

// Status 更新数据源状态
func (s *sSysDataConnector) Status(ctx context.Context, in *sysin.DataConnectorStatusInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(do.DataConnector{
		Status:    in.Status,
		UpdatedBy: contexts.GetUserId(ctx),
	}).Update(); err != nil {
		err = gerror.Wrap(err, "更新数据源状态失败，请稍后重试！")
		return
	}
	return
}

// Test 测试数据源配置
func (s *sSysDataConnector) Test(ctx context.Context, in *sysin.DataConnectorTestInp) (res *sysin.DataConnectorTestModel, err error) {
	if err = validateConnectorConfig(in.Direction, in.ConnectorType, in.Config); err != nil {
		return
	}

	res = &sysin.DataConnectorTestModel{
		Success: true,
		Message: "配置格式校验通过，真实连通性测试将在接入具体驱动后执行",
	}
	return
}

func validateConnectorConfig(direction, connectorType string, config *gjson.Json) error {
	if connectorType == "kafka" {
		if !hasConnectorConfigListValue(config, "brokers") {
			return gerror.New("请填写 Kafka Brokers")
		}
		authMode := connectorConfigString(config, "authMode")
		if authMode != "" && authMode != "none" {
			if authMode == "oauthbearer" {
				if connectorConfigString(config, "token") == "" {
					return gerror.New("请填写 Kafka Token")
				}
			} else if connectorConfigString(config, "username") == "" || connectorConfigString(config, "password") == "" {
				return gerror.New("请填写 Kafka 用户名和密码")
			}
		}
		return nil
	}

	if connectorType == "s3" {
		if connectorConfigString(config, "region") == "" {
			return gerror.New("请填写 S3 Region")
		}
		if connectorConfigString(config, "credentialMode") == "static" &&
			(connectorConfigString(config, "accessKey") == "" || connectorConfigString(config, "secretKey") == "") {
			return gerror.New("请填写 S3 AccessKey 和 SecretKey")
		}
		return nil
	}

	if direction == "source" && connectorType == "log" && connectorConfigString(config, "basePath") == "" {
		return gerror.New("请填写日志目录")
	}

	if direction == "sink" {
		if (connectorType == "http" || connectorType == "elasticsearch" || connectorType == "opensearch") &&
			connectorConfigString(config, "baseUrl") == "" {
			return gerror.New("请填写 Base URL")
		}
		if (connectorType == "feishu" || connectorType == "lark") && connectorConfigString(config, "webhook") == "" {
			return gerror.New("请填写 Webhook")
		}
		if connectorType == "splunk" &&
			(connectorConfigString(config, "hecUrl") == "" || connectorConfigString(config, "token") == "") {
			return gerror.New("请填写 Splunk HEC URL 和 HEC Token")
		}
	}

	return nil
}

func connectorConfigString(config *gjson.Json, path string) string {
	if config == nil {
		return ""
	}
	return strings.TrimSpace(config.Get(path).String())
}

func hasConnectorConfigListValue(config *gjson.Json, path string) bool {
	if config == nil {
		return false
	}
	for _, value := range gconv.Strings(config.Get(path).Interface()) {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return connectorConfigString(config, path) != ""
}
