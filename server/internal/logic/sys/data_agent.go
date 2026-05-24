package sys

import (
	"context"
	"strings"
	"time"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/do"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

const dataAgentOnlineWindow = 2 * time.Minute

type sSysDataAgent struct{}

func NewSysDataAgent() *sSysDataAgent {
	return &sSysDataAgent{}
}

func init() {
	service.RegisterSysDataAgent(NewSysDataAgent())
}

func (s *sSysDataAgent) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.DataAgent.Ctx(ctx), option...)
}

func (s *sSysDataAgent) List(ctx context.Context, in *sysin.DataAgentListInp) (list []*sysin.DataAgentListModel, totalCount int, err error) {
	cols := dao.DataAgent.Columns()
	mod := s.Model(ctx)
	keyword := strings.TrimSpace(in.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		mod = mod.Where("("+cols.AgentId+" LIKE ? OR "+cols.Name+" LIKE ? OR "+cols.Hostname+" LIKE ?)", like, like, like)
	}
	if in.RegisterStatus != "" {
		mod = mod.Where(cols.RegisterStatus, strings.TrimSpace(in.RegisterStatus))
	}
	if in.DispatchStatus != "" {
		mod = mod.Where(cols.DispatchStatus, strings.TrimSpace(in.DispatchStatus))
	}
	threshold := gtime.Now().Add(-dataAgentOnlineWindow)
	switch strings.TrimSpace(in.OnlineStatus) {
	case consts.DataAgentOnlineStatusOnline:
		mod = mod.WhereGTE(cols.LastSeenAt, threshold)
	case consts.DataAgentOnlineStatusOffline:
		mod = mod.Where("("+cols.LastSeenAt+" IS NULL OR "+cols.LastSeenAt+" < ?)", threshold)
	}
	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(cols.CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	var rows []*entity.DataAgent
	if err = mod.Page(in.Page, in.PerPage).
		OrderDesc(cols.LastSeenAt).
		OrderDesc(cols.Id).
		ScanAndCount(&rows, &totalCount, false); err != nil {
		return nil, 0, gerror.Wrap(err, "获取Agent节点列表失败，请稍后重试！")
	}

	now := gtime.Now().Time
	list = make([]*sysin.DataAgentListModel, 0, len(rows))
	for _, row := range rows {
		online, onlineStatus := dataAgentOnlineStatus(gtimeToTime(row.LastSeenAt), now)
		list = append(list, &sysin.DataAgentListModel{
			Id:              row.Id,
			AgentId:         row.AgentId,
			Name:            row.Name,
			Hostname:        row.Hostname,
			Version:         row.Version,
			AgentIp:         dataAgentJsonStrings(row.AgentIp),
			RegisterStatus:  row.RegisterStatus,
			DispatchStatus:  row.DispatchStatus,
			DispatchAllowed: row.RegisterStatus == consts.DataAgentRegisterStatusApproved && row.DispatchStatus == consts.DataAgentDispatchStatusEnabled,
			OnlineStatus:    onlineStatus,
			Online:          online,
			LastSeenAt:      row.LastSeenAt,
			ApprovedBy:      row.ApprovedBy,
			ApprovedAt:      row.ApprovedAt,
			DisabledAt:      row.DisabledAt,
			Remark:          row.Remark,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		})
	}
	return
}

func (s *sSysDataAgent) Approve(ctx context.Context, in *sysin.DataAgentApproveInp) (err error) {
	result, err := s.Model(ctx).
		Where(dao.DataAgent.Columns().AgentId, strings.TrimSpace(in.AgentId)).
		Data(do.DataAgent{
			RegisterStatus: consts.DataAgentRegisterStatusApproved,
			DispatchStatus: consts.DataAgentDispatchStatusEnabled,
			ApprovedBy:     contexts.GetUserId(ctx),
			ApprovedAt:     gtime.Now(),
			Remark:         in.Remark,
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "批准Agent节点失败，请稍后重试！")
	}
	return ensureDataAgentAffected(result, in.AgentId)
}

func (s *sSysDataAgent) Dispatch(ctx context.Context, in *sysin.DataAgentDispatchInp) (err error) {
	data := do.DataAgent{
		DispatchStatus: strings.TrimSpace(in.DispatchStatus),
		Remark:         in.Remark,
	}
	if in.DispatchStatus == consts.DataAgentDispatchStatusDisabled {
		data.DisabledAt = gtime.Now()
	}
	result, err := s.Model(ctx).
		Where(dao.DataAgent.Columns().AgentId, strings.TrimSpace(in.AgentId)).
		Where(dao.DataAgent.Columns().RegisterStatus, consts.DataAgentRegisterStatusApproved).
		Data(data).
		Update()
	if err != nil {
		return gerror.Wrap(err, "更新Agent调度状态失败，请稍后重试！")
	}
	return ensureDataAgentAffected(result, in.AgentId)
}

func (s *sSysDataAgent) Reject(ctx context.Context, in *sysin.DataAgentRejectInp) (err error) {
	result, err := s.Model(ctx).
		Where(dao.DataAgent.Columns().AgentId, strings.TrimSpace(in.AgentId)).
		Data(do.DataAgent{
			RegisterStatus: strings.TrimSpace(in.RegisterStatus),
			DispatchStatus: consts.DataAgentDispatchStatusDisabled,
			DisabledAt:     gtime.Now(),
			Remark:         in.Remark,
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "拒绝或吊销Agent节点失败，请稍后重试！")
	}
	return ensureDataAgentAffected(result, in.AgentId)
}

func dataAgentOnlineStatus(lastSeenAt time.Time, now time.Time) (bool, string) {
	if lastSeenAt.IsZero() || lastSeenAt.Before(now.Add(-dataAgentOnlineWindow)) {
		return false, consts.DataAgentOnlineStatusOffline
	}
	return true, consts.DataAgentOnlineStatusOnline
}

func dataAgentJsonStrings(value *gjson.Json) []string {
	if value == nil {
		return []string{}
	}
	items := gconv.Strings(value.Array())
	if items == nil {
		return []string{}
	}
	return items
}

func gtimeToTime(value *gtime.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.Time
}

func ensureDataAgentAffected(result interface{}, agentId string) error {
	if sqlResult, ok := result.(interface{ RowsAffected() (int64, error) }); ok {
		affected, err := sqlResult.RowsAffected()
		if err == nil && affected == 0 {
			return gerror.Newf("Agent节点不存在或状态不允许操作: %s", agentId)
		}
	}
	return nil
}
