package sys

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	"github.com/gogf/gf/v2/frame/g"
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

func (s *sSysDataAgent) HandleAgentEnvelope(ctx context.Context, envelope sysin.DataAgentEnvelope) (res sysin.DataAgentEnvelope, err error) {
	switch envelope.Type {
	case sysin.DataAgentMessageHello:
		var payload sysin.DataAgentHelloPayload
		if err = json.Unmarshal(envelope.Payload, &payload); err != nil {
			return res, gerror.Wrap(err, "解析Agent hello失败")
		}
		dispatchAllowed, err := s.handleAgentHello(ctx, &payload)
		if err != nil {
			return res, err
		}
		return dataAgentResponseEnvelope(envelope.RequestId, sysin.DataAgentMessageHelloAck, map[string]any{
			"accepted":        true,
			"dispatchAllowed": dispatchAllowed,
			"capabilities":    normalizeDataAgentStrings(payload.Capabilities),
			"protocolVersion": sysin.DataAgentProtocolVersion,
		}), nil

	case sysin.DataAgentMessageHeartbeat:
		var payload sysin.DataAgentHeartbeatPayload
		if err = json.Unmarshal(envelope.Payload, &payload); err != nil {
			return res, gerror.Wrap(err, "解析Agent heartbeat失败")
		}
		if err = s.handleAgentHeartbeat(ctx, &payload); err != nil {
			return res, err
		}
		return dataAgentResponseEnvelope(envelope.RequestId, sysin.DataAgentMessageHeartbeatAck, nil), nil

	case sysin.DataAgentMessageFieldsReport:
		var payload sysin.DataAgentFieldsReportPayload
		if err = json.Unmarshal(envelope.Payload, &payload); err != nil {
			return res, gerror.Wrap(err, "解析Agent字段上报失败")
		}
		if err = s.handleAgentFieldsReport(ctx, &payload); err != nil {
			return res, err
		}
		return dataAgentResponseEnvelope(envelope.RequestId, sysin.DataAgentMessageFieldsAck, nil), nil

	case sysin.DataAgentMessageStatsReport:
		var payload sysin.DataAgentStatsReportPayload
		if err = json.Unmarshal(envelope.Payload, &payload); err != nil {
			return res, gerror.Wrap(err, "解析Agent统计上报失败")
		}
		if err = s.handleAgentStatsReport(ctx, &payload); err != nil {
			return res, err
		}
		return dataAgentResponseEnvelope(envelope.RequestId, sysin.DataAgentMessageAck, nil), nil

	case sysin.DataAgentMessageError:
		return dataAgentResponseEnvelope(envelope.RequestId, sysin.DataAgentMessageAck, nil), nil
	default:
		return res, gerror.Newf("不支持的Agent消息类型: %s", envelope.Type)
	}
}

func (s *sSysDataAgent) handleAgentHello(ctx context.Context, payload *sysin.DataAgentHelloPayload) (dispatchAllowed bool, err error) {
	payload.AgentId = strings.TrimSpace(payload.AgentId)
	if payload.AgentId == "" {
		return false, gerror.New("Agent hello缺少agentId")
	}
	payload.Hostname = strings.TrimSpace(payload.Hostname)
	payload.Version = strings.TrimSpace(payload.Version)
	payload.AgentIp = normalizeDataAgentStrings(payload.AgentIp)
	payload.Capabilities = normalizeDataAgentStrings(payload.Capabilities)

	if err = s.upsertAgentHello(ctx, *payload); err != nil {
		return false, err
	}
	dispatchAllowed, err = s.loadAgentDispatchAllowed(ctx, payload.AgentId)
	if err != nil {
		return false, err
	}
	g.Log().Infof(ctx, "Agent hello已登记 agentId:%s hostname:%s dispatchAllowed:%v", payload.AgentId, payload.Hostname, dispatchAllowed)
	return dispatchAllowed, nil
}

func (s *sSysDataAgent) handleAgentHeartbeat(ctx context.Context, payload *sysin.DataAgentHeartbeatPayload) (err error) {
	payload.AgentId = strings.TrimSpace(payload.AgentId)
	if payload.AgentId == "" {
		return gerror.New("Agent heartbeat缺少agentId")
	}
	payload.AgentIp = normalizeDataAgentStrings(payload.AgentIp)
	return s.touchAgentHeartbeat(ctx, *payload)
}

func (s *sSysDataAgent) upsertAgentHello(ctx context.Context, payload sysin.DataAgentHelloPayload) error {
	now := gtime.Now()
	cols := dao.DataAgent.Columns()
	var current entity.DataAgent
	if err := dao.DataAgent.Ctx(ctx).Where(cols.AgentId, payload.AgentId).Scan(&current); err != nil {
		if !isSqlNoRows(err) {
			return gerror.Wrap(err, "查询Agent节点失败")
		}
	}
	data := do.DataAgent{
		Hostname:   payload.Hostname,
		Version:    payload.Version,
		AgentIp:    gjson.New(payload.AgentIp),
		LastSeenAt: now,
	}
	if current.Id > 0 {
		if _, err := dao.DataAgent.Ctx(ctx).WherePri(current.Id).Data(data).Update(); err != nil {
			return gerror.Wrap(err, "更新Agent节点失败")
		}
		return nil
	}

	name := payload.Hostname
	if name == "" {
		name = payload.AgentId
	}
	data.AgentId = payload.AgentId
	data.Name = name
	data.RegisterStatus = consts.DataAgentRegisterStatusPending
	data.DispatchStatus = consts.DataAgentDispatchStatusDisabled
	if _, err := dao.DataAgent.Ctx(ctx).Data(data).Insert(); err != nil {
		return gerror.Wrap(err, "创建Agent节点失败")
	}
	return nil
}

func (s *sSysDataAgent) touchAgentHeartbeat(ctx context.Context, payload sysin.DataAgentHeartbeatPayload) error {
	now := gtime.Now()
	cols := dao.DataAgent.Columns()
	var current entity.DataAgent
	if err := dao.DataAgent.Ctx(ctx).Where(cols.AgentId, payload.AgentId).Scan(&current); err != nil {
		if !isSqlNoRows(err) {
			return gerror.Wrap(err, "查询Agent心跳节点失败")
		}
	}
	if current.Id > 0 {
		if _, err := dao.DataAgent.Ctx(ctx).WherePri(current.Id).Data(do.DataAgent{
			AgentIp:    gjson.New(payload.AgentIp),
			LastSeenAt: now,
		}).Update(); err != nil {
			return gerror.Wrap(err, "更新Agent心跳失败")
		}
		return nil
	}
	if _, err := dao.DataAgent.Ctx(ctx).Data(do.DataAgent{
		AgentId:        payload.AgentId,
		Name:           payload.AgentId,
		AgentIp:        gjson.New(payload.AgentIp),
		RegisterStatus: consts.DataAgentRegisterStatusPending,
		DispatchStatus: consts.DataAgentDispatchStatusDisabled,
		LastSeenAt:     now,
	}).Insert(); err != nil {
		return gerror.Wrap(err, "创建Agent心跳节点失败")
	}
	return nil
}

func (s *sSysDataAgent) loadAgentDispatchAllowed(ctx context.Context, agentId string) (bool, error) {
	cols := dao.DataAgent.Columns()
	count, err := dao.DataAgent.Ctx(ctx).
		Where(cols.AgentId, strings.TrimSpace(agentId)).
		Where(cols.RegisterStatus, consts.DataAgentRegisterStatusApproved).
		Where(cols.DispatchStatus, consts.DataAgentDispatchStatusEnabled).
		Count()
	if err != nil {
		return false, gerror.Wrap(err, "查询Agent调度状态失败")
	}
	return count > 0, nil
}

func (s *sSysDataAgent) handleAgentFieldsReport(ctx context.Context, payload *sysin.DataAgentFieldsReportPayload) error {
	payload.AgentId = strings.TrimSpace(payload.AgentId)
	if payload.AgentId == "" {
		return gerror.New("Agent字段上报缺少agentId")
	}
	if payload.TaskId <= 0 {
		return gerror.New("Agent字段上报缺少taskId")
	}
	task, err := s.loadAgentTask(ctx, payload.TaskId)
	if err != nil {
		return err
	}
	profiles := dataAgentFieldSamplesToProfiles(payload.Fields)
	if len(profiles) == 0 {
		return nil
	}
	return NewSysDataClean().upsertDataFields(ctx, task.Id, task.SourceId, profiles)
}

func (s *sSysDataAgent) handleAgentStatsReport(ctx context.Context, payload *sysin.DataAgentStatsReportPayload) error {
	payload.AgentId = strings.TrimSpace(payload.AgentId)
	if payload.AgentId == "" {
		return gerror.New("Agent统计上报缺少agentId")
	}
	windowStart := payload.WindowStart
	if windowStart.IsZero() {
		windowStart = time.Now().Truncate(time.Minute)
	}
	windowEnd := payload.WindowEnd
	if windowEnd.IsZero() || windowEnd.Before(windowStart) {
		windowEnd = windowStart.Add(time.Minute)
	}
	counters := dataAgentStatsCounters(payload.Metrics)
	for taskId, counter := range counters {
		task, err := s.loadAgentTask(ctx, taskId)
		if err != nil {
			return err
		}
		if err = s.upsertAgentCleanStat(ctx, task, windowStart, windowEnd, counter); err != nil {
			return err
		}
	}
	return nil
}

func (s *sSysDataAgent) loadAgentTask(ctx context.Context, taskId int64) (entity.DataCleanTask, error) {
	var task entity.DataCleanTask
	if err := dao.DataCleanTask.Ctx(ctx).WherePri(taskId).Scan(&task); err != nil {
		if isSqlNoRows(err) {
			return task, gerror.Newf("数据清洗任务不存在: %d", taskId)
		}
		return task, gerror.Wrap(err, "查询数据清洗任务失败")
	}
	if task.Id <= 0 {
		return task, gerror.Newf("数据清洗任务不存在: %d", taskId)
	}
	return task, nil
}

type dataAgentStatCounter struct {
	TotalCount       int64
	SuccessCount     int64
	FailedCount      int64
	CleanFailedCount int64
	DroppedCount     int64
}

func dataAgentStatsCounters(metrics []sysin.DataAgentStatsMetric) map[int64]dataAgentStatCounter {
	counters := make(map[int64]dataAgentStatCounter)
	for _, metric := range metrics {
		if metric.TaskId <= 0 || metric.Value <= 0 {
			continue
		}
		counter := counters[metric.TaskId]
		recognized := true
		switch strings.TrimSpace(metric.Name) {
		case "records.read":
			counter.TotalCount += metric.Value
		case "records.output":
			counter.SuccessCount += metric.Value
		case "records.failed", "records.output_failed", "records.error":
			counter.FailedCount += metric.Value
		case "records.clean_failed":
			counter.CleanFailedCount += metric.Value
		case "records.dropped":
			counter.DroppedCount += metric.Value
		default:
			recognized = false
		}
		if recognized {
			counters[metric.TaskId] = counter
		}
	}
	return counters
}

func (s *sSysDataAgent) upsertAgentCleanStat(ctx context.Context, task entity.DataCleanTask, windowStart time.Time, windowEnd time.Time, counter dataAgentStatCounter) error {
	if counter.TotalCount == 0 && counter.SuccessCount == 0 && counter.FailedCount == 0 && counter.CleanFailedCount == 0 && counter.DroppedCount == 0 {
		return nil
	}
	cols := dao.DataCleanStat.Columns()
	eventType := strings.TrimSpace(task.EventType)
	if eventType == "" {
		eventType = "custom.event"
	}
	var current entity.DataCleanStat
	if err := dao.DataCleanStat.Ctx(ctx).
		Where(cols.TaskId, task.Id).
		Where(cols.ConnectorId, task.SourceId).
		Where(cols.EventType, eventType).
		Where(cols.WindowStart, gtime.New(windowStart)).
		Scan(&current); err != nil {
		if !isSqlNoRows(err) {
			return gerror.Wrap(err, "查询Agent清洗统计失败")
		}
	}
	data := do.DataCleanStat{
		TotalCount:       current.TotalCount + counter.TotalCount,
		SuccessCount:     current.SuccessCount + counter.SuccessCount,
		FailedCount:      current.FailedCount + counter.FailedCount,
		CleanFailedCount: current.CleanFailedCount + counter.CleanFailedCount,
		DroppedCount:     current.DroppedCount + counter.DroppedCount,
	}
	if current.Id > 0 {
		if _, err := dao.DataCleanStat.Ctx(ctx).WherePri(current.Id).Data(data).Update(); err != nil {
			return gerror.Wrap(err, "更新Agent清洗统计失败")
		}
		return nil
	}
	data.TaskId = task.Id
	data.ConnectorId = task.SourceId
	data.EventType = eventType
	data.WindowStart = gtime.New(windowStart)
	data.WindowEnd = gtime.New(windowEnd)
	if _, err := dao.DataCleanStat.Ctx(ctx).Data(data).Insert(); err != nil {
		return gerror.Wrap(err, "创建Agent清洗统计失败")
	}
	return nil
}

func dataAgentFieldSamplesToProfiles(fields []sysin.DataAgentFieldSample) []dataFieldProfile {
	profiles := make(map[string]*dataFieldProfile)
	for _, field := range fields {
		path := strings.TrimSpace(field.Path)
		if path == "" || len(path) > 512 {
			continue
		}
		fieldType := strings.TrimSpace(field.Type)
		if fieldType == "" {
			fieldType = "unknown"
		}
		count := field.Count
		if count <= 0 {
			count = 1
		}
		profile := profiles[path]
		if profile == nil {
			profile = &dataFieldProfile{
				FieldPath: path,
				TypeStats: map[string]int64{},
			}
			profiles[path] = profile
		}
		profile.TypeStats[fieldType] += count
		profile.Count += count
		profile.NullCount += field.NullCount
		sample := maskDataFieldSample(path, field.Sample)
		if sample != "" && !stringSliceContains(profile.SampleValues, sample) && len(profile.SampleValues) < 5 {
			profile.SampleValues = append(profile.SampleValues, sample)
		}
	}
	list := make([]dataFieldProfile, 0, len(profiles))
	for _, profile := range profiles {
		profile.FieldType = mergeDataFieldType(profile.TypeStats)
		list = append(list, *profile)
	}
	return list
}

func normalizeDataAgentStrings(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func dataAgentResponseEnvelope(requestId string, messageType string, payload any) sysin.DataAgentEnvelope {
	raw, _ := json.Marshal(payload)
	return sysin.DataAgentEnvelope{
		Type:      messageType,
		RequestId: requestId,
		SentAt:    time.Now(),
		Payload:   raw,
	}
}

func isSqlNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
