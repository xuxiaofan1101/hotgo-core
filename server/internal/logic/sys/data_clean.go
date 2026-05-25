package sys

import (
	"context"
	"sort"
	"strings"
	"time"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/hgorm"
	"hotgo/internal/library/hgorm/handler"
	"hotgo/internal/model/do"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sSysDataClean struct{}

func NewSysDataClean() *sSysDataClean {
	return &sSysDataClean{}
}

func init() {
	service.RegisterSysDataClean(NewSysDataClean())
}

func (s *sSysDataClean) Model(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.DataCleanTask.Ctx(ctx), option...)
}

func (s *sSysDataClean) FieldModel(ctx context.Context, option ...*handler.Option) *gdb.Model {
	return handler.Model(dao.DataField.Ctx(ctx), option...)
}

func (s *sSysDataClean) List(ctx context.Context, in *sysin.DataCleanTaskListInp) (list []*sysin.DataCleanTaskListModel, totalCount int, err error) {
	taskCols := dao.DataCleanTask.Columns()
	connCols := dao.DataConnector.Columns()
	tplCols := dao.DataFieldTemplate.Columns()
	mod := dao.DataCleanTask.Ctx(ctx).As("t").
		LeftJoin(dao.DataConnector.Table()+" c", "c."+connCols.Id+"=t."+taskCols.SourceId).
		LeftJoin(dao.DataFieldTemplate.Table()+" tpl", "tpl."+tplCols.Id+"=t."+taskCols.TemplateId)

	mod = mod.Fields("t.*, c." + connCols.Name + " AS source_name, c." + connCols.Code + " AS source_code, c." + connCols.ConnectorType + " AS source_type, tpl." + tplCols.Name + " AS template_name")

	if in.Id > 0 {
		mod = mod.Where("t."+taskCols.Id, in.Id)
	}
	if in.Keyword != "" {
		keyword := "%" + strings.TrimSpace(in.Keyword) + "%"
		mod = mod.Where("(t."+taskCols.Name+" LIKE ? OR t."+taskCols.Code+" LIKE ? OR c."+connCols.Name+" LIKE ? OR c."+connCols.Code+" LIKE ?)", keyword, keyword, keyword, keyword)
	}
	if in.SourceId > 0 {
		mod = mod.Where("t."+taskCols.SourceId, in.SourceId)
	}
	if in.SourceType != "" {
		mod = mod.Where("c."+connCols.ConnectorType, in.SourceType)
	}
	if in.TemplateId > 0 {
		mod = mod.Where("t."+taskCols.TemplateId, in.TemplateId)
	}
	if in.EventType != "" {
		mod = mod.Where("t."+taskCols.EventType, in.EventType)
	}
	if in.Topic != "" {
		mod = mod.Where("JSON_UNQUOTE(JSON_EXTRACT(t."+taskCols.SourceConfig+", '$.topic')) LIKE ?", "%"+in.Topic+"%")
	}
	if in.GroupId != "" {
		mod = mod.Where("JSON_UNQUOTE(JSON_EXTRACT(t."+taskCols.SourceConfig+", '$.groupId')) LIKE ?", "%"+in.GroupId+"%")
	}
	if in.Bucket != "" {
		mod = mod.Where("JSON_UNQUOTE(JSON_EXTRACT(t."+taskCols.SourceConfig+", '$.bucket')) LIKE ?", "%"+in.Bucket+"%")
	}
	if in.Prefix != "" {
		mod = mod.Where("JSON_UNQUOTE(JSON_EXTRACT(t."+taskCols.SourceConfig+", '$.prefix')) LIKE ?", "%"+in.Prefix+"%")
	}
	if in.Path != "" {
		mod = mod.Where("JSON_UNQUOTE(JSON_EXTRACT(t."+taskCols.SourceConfig+", '$.path')) LIKE ?", "%"+in.Path+"%")
	}
	if in.Status > 0 {
		mod = mod.Where("t."+taskCols.Status, in.Status)
	}
	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween("t."+taskCols.CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	mod = mod.Page(in.Page, in.PerPage).OrderDesc("t." + taskCols.Id)
	if err = mod.ScanAndCount(&list, &totalCount, false); err != nil {
		return nil, 0, gerror.Wrap(err, "获取数据清洗任务列表失败，请稍后重试！")
	}
	return
}

func (s *sSysDataClean) Edit(ctx context.Context, in *sysin.DataCleanTaskEditInp) (err error) {
	source, err := s.getSourceConnector(ctx, in.SourceId)
	if err != nil {
		return err
	}
	if err = validateDataCleanSourceConfig(source.ConnectorType, in.SourceConfig); err != nil {
		return err
	}
	if err = validateDataCleanConfig(in.CleanConfig); err != nil {
		return err
	}
	in.SinkConfig = normalizeDataSinkConfig(in.SinkConfig)
	if err = validateDataSinkConfig(in.SinkConfig); err != nil {
		return err
	}
	if err = hgorm.IsUnique(ctx, &dao.DataCleanTask, g.Map{dao.DataCleanTask.Columns().Code: in.Code}, "清洗任务编码已存在", in.Id); err != nil {
		return err
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		if in.Id > 0 {
			current, err := s.getTask(ctx, in.Id)
			if err != nil {
				return err
			}
			in.CleanConfigVersion = current.CleanConfigVersion + 1
			if in.CleanConfigVersion <= 0 {
				in.CleanConfigVersion = 1
			}
			in.UpdatedBy = contexts.GetUserId(ctx)
			if _, err = s.Model(ctx).
				Fields(sysin.DataCleanTaskUpdateFields{}).
				WherePri(in.Id).Data(in).Update(); err != nil {
				return gerror.Wrap(err, "修改数据清洗任务失败，请稍后重试！")
			}
			return nil
		}

		in.CreatedBy = contexts.GetUserId(ctx)
		in.CleanConfigVersion = 1
		in.FieldSchemaVersion = 1
		if _, err = s.Model(ctx, &handler.Option{FilterAuth: false}).
			Fields(sysin.DataCleanTaskInsertFields{}).
			Data(in).OmitEmptyData().Insert(); err != nil {
			return gerror.Wrap(err, "新增数据清洗任务失败，请稍后重试！")
		}
		return nil
	})
}

func (s *sSysDataClean) Delete(ctx context.Context, in *sysin.DataCleanTaskDeleteInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(do.DataCleanTask{
		DeletedBy: contexts.GetUserId(ctx),
		DeletedAt: gtime.Now(),
	}).Update(); err != nil {
		return gerror.Wrap(err, "删除数据清洗任务失败，请稍后重试！")
	}
	return nil
}

func (s *sSysDataClean) View(ctx context.Context, in *sysin.DataCleanTaskViewInp) (res *sysin.DataCleanTaskViewModel, err error) {
	taskCols := dao.DataCleanTask.Columns()
	connCols := dao.DataConnector.Columns()
	tplCols := dao.DataFieldTemplate.Columns()
	err = dao.DataCleanTask.Ctx(ctx).As("t").
		LeftJoin(dao.DataConnector.Table()+" c", "c."+connCols.Id+"=t."+taskCols.SourceId).
		LeftJoin(dao.DataFieldTemplate.Table()+" tpl", "tpl."+tplCols.Id+"=t."+taskCols.TemplateId).
		Fields("t.*, c."+connCols.Name+" AS source_name, c."+connCols.Code+" AS source_code, c."+connCols.ConnectorType+" AS source_type, tpl."+tplCols.Name+" AS template_name").
		Where("t."+taskCols.Id, in.Id).
		Scan(&res)
	if err != nil {
		return nil, gerror.Wrap(err, "获取数据清洗任务信息失败，请稍后重试！")
	}
	return
}

func (s *sSysDataClean) Status(ctx context.Context, in *sysin.DataCleanTaskStatusInp) (err error) {
	if _, err = s.Model(ctx).WherePri(in.Id).Data(do.DataCleanTask{
		Status:    in.Status,
		UpdatedBy: contexts.GetUserId(ctx),
	}).Update(); err != nil {
		return gerror.Wrap(err, "更新数据清洗任务状态失败，请稍后重试！")
	}
	return nil
}

func (s *sSysDataClean) Test(ctx context.Context, in *sysin.DataCleanTaskTestInp) (res *sysin.DataCleanTaskTestModel, err error) {
	source, err := s.getSourceConnector(ctx, in.SourceId)
	if err != nil {
		return nil, err
	}
	if err = validateDataCleanSourceConfig(source.ConnectorType, in.SourceConfig); err != nil {
		return &sysin.DataCleanTaskTestModel{Success: false, Message: err.Error()}, nil
	}
	if err = validateDataCleanConfig(in.CleanConfig); err != nil {
		return &sysin.DataCleanTaskTestModel{Success: false, Message: err.Error()}, nil
	}
	if err = validateDataSinkConfig(in.SinkConfig); err != nil {
		return &sysin.DataCleanTaskTestModel{Success: false, Message: err.Error()}, nil
	}
	return &sysin.DataCleanTaskTestModel{Success: true, Message: "清洗配置校验通过"}, nil
}

func (s *sSysDataClean) Sample(ctx context.Context, in *sysin.DataCleanTaskSampleInp) (res *sysin.DataCleanTaskSampleModel, err error) {
	source, err := s.getSourceConnector(ctx, in.SourceId)
	if err != nil {
		return nil, err
	}
	payload := strings.TrimSpace(in.Payload)
	if payload == "" && in.SourceConfig != nil {
		payload = strings.TrimSpace(in.SourceConfig.Get("samplePayload").String())
	}
	var payloads []map[string]interface{}
	if payload != "" {
		payloads, err = parseDataCleanPayloads(payload)
		if err != nil {
			return nil, err
		}
		payloads = limitDataCleanSamplePayloads(payloads, in.Limit)
	} else {
		config := mergeDataCleanSourceConfig(source.Config, in.SourceConfig)
		payloads, err = sampleDataCleanSourcePayloads(ctx, source.ConnectorType, config, in.Limit, time.Duration(in.TimeoutSeconds)*time.Second)
		if err != nil {
			return nil, err
		}
	}
	parseDepth := 0
	if in.CleanConfig != nil {
		parseDepth = in.CleanConfig.Get("parseDepth").Int()
	}
	profiles := profileDataFields(payloads, parseDepth)
	if in.Id > 0 && len(profiles) > 0 {
		if err = s.upsertDataFields(ctx, in.Id, in.SourceId, profiles); err != nil {
			return nil, err
		}
	}
	list := make([]*sysin.DataFieldListModel, 0, len(profiles))
	for _, profile := range profiles {
		list = append(list, dataFieldProfileToModel(in.Id, in.SourceId, profile))
	}
	return &sysin.DataCleanTaskSampleModel{List: list, Total: len(list)}, nil
}

func sampleDataCleanSourcePayloads(ctx context.Context, sourceType string, config map[string]interface{}, limit int, timeout time.Duration) ([]map[string]interface{}, error) {
	switch sourceType {
	case consts.DataConnectorTypeKafka:
		return sampleDataCleanKafkaPayloads(ctx, config, limit, timeout)
	default:
		return nil, gerror.New("当前数据源类型暂不支持主动采样，请粘贴样例数据")
	}
}

func sampleDataCleanKafkaPayloads(ctx context.Context, config map[string]interface{}, limit int, timeout time.Duration) ([]map[string]interface{}, error) {
	brokers := dataConnectorKafkaStringSlice(config, "brokers")
	topic := dataConnectorKafkaString(config, "topic")
	if len(brokers) == 0 || topic == "" {
		return nil, gerror.New("Kafka 采样需要配置 Brokers 和 Topic")
	}
	kafkaConfig, err := buildDataConnectorKafkaConfig(config)
	if err != nil {
		return nil, gerror.Wrap(err, "Kafka配置无效")
	}
	kafkaConfig.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, kafkaConfig)
	if err != nil {
		return nil, gerror.Wrap(err, "连接Kafka失败")
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return nil, gerror.Wrap(err, "读取Kafka分区失败")
	}
	sort.Slice(partitions, func(i, j int) bool { return partitions[i] < partitions[j] })

	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	payloads := make([]map[string]interface{}, 0, limit)
	for _, partition := range partitions {
		if len(payloads) >= limit {
			break
		}
		offset := dataCleanKafkaSampleStartOffset(config, partition)
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
		if err != nil {
			return nil, gerror.Wrap(err, "消费Kafka分区失败")
		}
		for len(payloads) < limit {
			select {
			case <-ctx.Done():
				partitionConsumer.Close()
				return payloads, ctx.Err()
			case <-deadline.C:
				partitionConsumer.Close()
				return payloads, nil
			case err := <-partitionConsumer.Errors():
				partitionConsumer.Close()
				if err == nil {
					return payloads, nil
				}
				return payloads, err
			case msg := <-partitionConsumer.Messages():
				if msg == nil {
					partitionConsumer.Close()
					goto nextPartition
				}
				parsed, err := parseDataCleanPayloads(string(msg.Value))
				if err != nil {
					partitionConsumer.Close()
					return nil, err
				}
				payloads = appendDataCleanSamplePayloads(payloads, parsed, limit)
			}
		}
		partitionConsumer.Close()
	nextPartition:
	}
	return payloads, nil
}

func dataCleanKafkaSampleStartOffset(config map[string]interface{}, partition int32) int64 {
	mode := strings.ToLower(strings.TrimSpace(dataConnectorKafkaString(config, "sampleMode")))
	if mode == "" {
		mode = "latest"
	}
	switch mode {
	case "latest":
		return sarama.OffsetNewest
	case "specific", "offset":
		offsets := gconv.Map(config["partitionOffsets"])
		if value, ok := offsets[gconv.String(partition)]; ok {
			return gconv.Int64(value)
		}
		return gconv.Int64(config["startOffset"])
	default:
		return sarama.OffsetOldest
	}
}

func (s *sSysDataClean) FieldList(ctx context.Context, in *sysin.DataFieldListInp) (list []*sysin.DataFieldListModel, total int, err error) {
	if err = s.FieldModel(ctx).
		Where(dao.DataField.Columns().TaskId, in.TaskId).
		Where(dao.DataField.Columns().Status, consts.StatusEnabled).
		OrderAsc(dao.DataField.Columns().FieldPath).
		Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, "获取采集字段列表失败，请稍后重试！")
	}
	for _, item := range list {
		item.SampleValues = sampleValueToStrings(item.SampleValue)
	}
	return list, len(list), nil
}

func (s *sSysDataClean) getTask(ctx context.Context, id int64) (task entity.DataCleanTask, err error) {
	err = s.Model(ctx).WherePri(id).Scan(&task)
	if err != nil {
		return task, err
	}
	if task.Id <= 0 {
		return task, gerror.New("数据清洗任务不存在")
	}
	return task, nil
}

func (s *sSysDataClean) getSourceConnector(ctx context.Context, sourceId int64) (connector entity.DataConnector, err error) {
	cols := dao.DataConnector.Columns()
	err = dao.DataConnector.Ctx(ctx).
		Where(cols.Id, sourceId).
		Where(cols.Direction, consts.DataConnectorDirectionSource).
		Where(cols.Status, consts.StatusEnabled).
		Scan(&connector)
	if err != nil {
		return connector, err
	}
	if connector.Id <= 0 {
		return connector, gerror.New("输入数据源不存在或已禁用")
	}
	return connector, nil
}

func validateDataCleanSourceConfig(sourceType string, config *gjson.Json) error {
	switch sourceType {
	case consts.DataConnectorTypeHTTP, consts.DataConnectorTypeManual:
		return nil
	case consts.DataConnectorTypeKafka:
		if dataCleanConfigString(config, "topic") == "" {
			return gerror.New("Kafka 数据清洗需要配置 Topic")
		}
	case consts.DataConnectorTypeS3:
		if dataCleanConfigString(config, "bucket") == "" {
			return gerror.New("S3 数据清洗需要配置 Bucket")
		}
	case consts.DataConnectorTypeLog:
		if dataCleanConfigString(config, "path") == "" {
			return gerror.New("日志数据清洗需要配置日志路径")
		}
	default:
		return gerror.New("输入数据源类型无效")
	}
	return nil
}

func (s *sSysDataClean) upsertDataFields(ctx context.Context, taskId int64, connectorId int64, profiles []dataFieldProfile) error {
	now := gtime.Now()
	newFieldCount := 0
	for _, profile := range profiles {
		var existing entity.DataField
		cols := dao.DataField.Columns()
		if err := dao.DataField.Ctx(ctx).
			Where(cols.TaskId, taskId).
			Where(cols.FieldPath, profile.FieldPath).
			Scan(&existing); err != nil {
			if !isSqlNoRows(err) {
				return err
			}
		}
		if existing.Id > 0 {
			if _, err := dao.DataField.Ctx(ctx).WherePri(existing.Id).Data(do.DataField{
				FieldType:   profile.FieldType,
				FieldTypes:  gjson.New(profile.TypeStats),
				SampleValue: gjson.New(profile.SampleValues),
				LastSeenAt:  now,
				UpdatedAt:   now,
			}).Update(); err != nil {
				return err
			}
			continue
		}
		newFieldCount++
		if _, err := dao.DataField.Ctx(ctx).Data(do.DataField{
			TaskId:      taskId,
			ConnectorId: connectorId,
			FieldPath:   profile.FieldPath,
			FieldName:   defaultDataFieldName(profile.FieldPath),
			FieldType:   profile.FieldType,
			FieldTypes:  gjson.New(profile.TypeStats),
			SampleValue: gjson.New(profile.SampleValues),
			Status:      consts.StatusEnabled,
			FirstSeenAt: now,
			LastSeenAt:  now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}).Insert(); err != nil {
			return err
		}
	}
	if newFieldCount > 0 {
		task, err := s.getTask(ctx, taskId)
		if err != nil {
			return err
		}
		_, err = s.Model(ctx).WherePri(taskId).Data(do.DataCleanTask{
			FieldSchemaVersion: task.FieldSchemaVersion + 1,
			UpdatedBy:          contexts.GetUserId(ctx),
		}).Update()
		return err
	}
	return nil
}

func dataFieldProfileToModel(taskId int64, connectorId int64, profile dataFieldProfile) *sysin.DataFieldListModel {
	return &sysin.DataFieldListModel{
		TaskId:       taskId,
		ConnectorId:  connectorId,
		FieldPath:    profile.FieldPath,
		FieldName:    defaultDataFieldName(profile.FieldPath),
		FieldType:    profile.FieldType,
		FieldTypes:   gjson.New(profile.TypeStats),
		SampleValue:  gjson.New(profile.SampleValues),
		SampleValues: profile.SampleValues,
		Status:       consts.StatusEnabled,
	}
}

func defaultDataFieldName(path string) string {
	path = strings.TrimSuffix(path, "[]")
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return path
	}
	return strings.TrimSuffix(parts[len(parts)-1], "[]")
}
