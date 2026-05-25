<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" :title="pageTitle" />
    </div>
    <n-card :bordered="false" class="proCard data-clean-page-card">
      <n-spin :show="loading" description="请稍候...">
        <n-form
          ref="formRef"
          :model="formValue"
          :rules="rules"
          :label-placement="settingStore.isMobile ? 'top' : 'left'"
          :label-width="120"
          class="data-clean-page-form"
        >
          <n-steps v-model:current="currentStep" class="data-clean-steps" size="small">
            <n-step title="输入配置" />
            <n-step title="清洗规则" />
            <n-step title="输出分发" />
          </n-steps>

          <div v-show="currentStep === 1" class="data-clean-pane">
            <n-grid
              cols="1 s:1 m:2 l:2 xl:2 2xl:2"
              responsive="screen"
              :x-gap="16"
              class="data-clean-control-grid"
            >
              <n-gi>
                <n-form-item label="任务名称" path="name">
                  <n-input v-model:value="formValue.name" placeholder="请输入任务名称" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="任务编码" path="code">
                  <n-input v-model:value="formValue.code" placeholder="请输入任务编码" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="输入数据源" path="sourceId">
                  <n-select
                    v-model:value="formValue.sourceId"
                    :options="sourceOptions"
                    filterable
                    clearable
                    placeholder="请选择输入数据源"
                    @update:value="handleSourceChange"
                  />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="标准模板">
                  <n-select
                    v-model:value="formValue.templateId"
                    :options="templateOptions"
                    filterable
                    clearable
                    placeholder="可选，用于日志标准化"
                  />
                </n-form-item>
              </n-gi>
            </n-grid>

            <div class="data-clean-section">
              <div class="data-clean-section-title">读取配置</div>
              <n-grid
                cols="1 s:1 m:2 l:2 xl:2 2xl:2"
                responsive="screen"
                :x-gap="16"
                class="data-clean-control-grid"
              >
                <template v-if="selectedSourceType === 'kafka'">
                  <n-gi>
                    <n-form-item label="Topic" required>
                      <n-select
                        v-model:value="sourceConfigForm.topic"
                        :options="kafkaTopicOptions"
                        :loading="kafkaTopicLoading"
                        filterable
                        tag
                        clearable
                        placeholder="请选择或输入 Topic"
                        @focus="loadKafkaTopics"
                        @update:value="handleKafkaTopicChange"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="Consumer Group">
                      <n-input
                        v-model:value="sourceConfigForm.groupId"
                        @update:value="handleKafkaGroupIdChange"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="起始位置">
                      <n-select
                        v-model:value="sourceConfigForm.startMode"
                        :options="kafkaStartModeOptions"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi v-if="sourceConfigForm.startMode === 'offset'">
                    <n-form-item label="起始 Offset">
                      <n-input-number
                        v-model:value="sourceConfigForm.startOffset"
                        class="data-clean-full"
                        :show-button="false"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="采样模式">
                      <n-select
                        v-model:value="sourceConfigForm.sampleMode"
                        :options="sampleModeOptions"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="采样条数">
                      <n-input-number
                        v-model:value="sourceConfigForm.sampleLimit"
                        class="data-clean-full"
                        :min="1"
                        :max="500"
                        :show-button="false"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="采样超时(秒)">
                      <n-input-number
                        v-model:value="sourceConfigForm.sampleTimeoutSeconds"
                        class="data-clean-full"
                        :min="1"
                        :max="60"
                        :show-button="false"
                      />
                    </n-form-item>
                  </n-gi>
                </template>

                <template v-else-if="selectedSourceType === 's3'">
                  <n-gi>
                    <n-form-item label="Bucket" required>
                      <n-input v-model:value="sourceConfigForm.bucket" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="Prefix">
                      <n-input v-model:value="sourceConfigForm.prefix" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="最大对象数">
                      <n-input-number
                        v-model:value="sourceConfigForm.maxObjects"
                        class="data-clean-full"
                        :min="1"
                        :max="500"
                        :show-button="false"
                      />
                    </n-form-item>
                  </n-gi>
                </template>

                <template v-else-if="selectedSourceType === 'log'">
                  <n-gi>
                    <n-form-item label="日志路径" required>
                      <n-input
                        v-model:value="sourceConfigForm.path"
                        placeholder="/var/log/nginx/access.log"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="读取上限">
                      <n-input-number
                        v-model:value="sourceConfigForm.logMaxBytes"
                        class="data-clean-full"
                        :min="1024"
                        :show-button="false"
                      />
                    </n-form-item>
                  </n-gi>
                </template>

                <template v-else-if="selectedSourceType === 'http'">
                  <n-gi>
                    <n-form-item label="接入 Token">
                      <n-input
                        v-model:value="sourceConfigForm.token"
                        placeholder="可选，覆盖数据源默认 Token"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="限流/秒">
                      <n-input-number
                        v-model:value="sourceConfigForm.rateLimit"
                        class="data-clean-full"
                        :min="0"
                        :show-button="false"
                      />
                    </n-form-item>
                  </n-gi>
                </template>

                <template v-else>
                  <n-gi>
                    <n-form-item label="来源说明">
                      <n-input
                        v-model:value="sourceConfigForm.description"
                        placeholder="后台手工导入或测试样例"
                      />
                    </n-form-item>
                  </n-gi>
                </template>

                <n-gi>
                  <n-form-item label="数据格式">
                    <n-select
                      v-model:value="sourceConfigForm.payloadFormat"
                      :options="payloadFormatOptions"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi span="2" class="data-clean-sample-grid-item">
                  <n-form-item label="样例数据">
                    <div class="data-clean-full">
                      <n-input
                        v-model:value="sourceConfigForm.samplePayload"
                        type="textarea"
                        placeholder="粘贴 JSON 对象、JSON 数组或 JSON Lines，用于立即提取字段"
                        :autosize="{ minRows: 5, maxRows: 10 }"
                      />
                      <div class="data-clean-sample-actions">
                        <n-button size="small" :loading="sampleLoading" @click="handleSample">
                          采样字段
                        </n-button>
                      </div>
                    </div>
                  </n-form-item>
                </n-gi>
              </n-grid>
            </div>
          </div>

          <div v-show="currentStep === 2" class="data-clean-pane">
            <n-grid
              cols="1 s:1 m:2 l:2 xl:2 2xl:2"
              responsive="screen"
              :x-gap="16"
              class="data-clean-control-grid"
            >
              <n-gi>
                <n-form-item label="失败策略">
                  <n-select
                    v-model:value="formValue.cleanErrorPolicy"
                    :options="cleanErrorPolicyOptions"
                  />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="未知字段">
                  <n-select
                    v-model:value="formValue.unknownFieldPolicy"
                    :options="unknownFieldPolicyOptions"
                  />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="解析深度">
                  <n-input-number
                    v-model:value="cleanRuleForm.parseDepth"
                    class="data-clean-full"
                    :min="0"
                    :max="16"
                    :show-button="false"
                  />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="字段过滤">
                  <n-switch v-model:value="cleanRuleForm.filterEnabled" />
                </n-form-item>
              </n-gi>
            </n-grid>

            <div v-if="cleanRuleForm.filterEnabled" class="data-clean-section">
              <div class="data-clean-toolbar">
                <div class="data-clean-section-title">过滤条件</div>
                <n-space>
                  <n-button size="small" @click="addFilterCondition(0)">
                    <template #icon>
                      <n-icon><PlusOutlined /></n-icon>
                    </template>
                    条件
                  </n-button>
                  <n-button size="small" @click="addFilterGroup">
                    <template #icon>
                      <n-icon><PlusOutlined /></n-icon>
                    </template>
                    条件组
                  </n-button>
                </n-space>
              </div>
              <div
                v-for="(group, groupIndex) in filterGroups"
                :key="group.key"
                class="data-clean-filter-group"
              >
                <div class="data-clean-filter-head">
                  <n-radio-group v-model:value="group.logic" size="small">
                    <n-radio-button value="and">AND</n-radio-button>
                    <n-radio-button value="or">OR</n-radio-button>
                  </n-radio-group>
                  <n-space>
                    <n-button size="tiny" @click="addFilterCondition(groupIndex)">
                      <template #icon>
                        <n-icon><PlusOutlined /></n-icon>
                      </template>
                      条件
                    </n-button>
                    <n-button
                      size="tiny"
                      type="error"
                      :disabled="filterGroups.length === 1"
                      @click="removeFilterGroup(groupIndex)"
                    >
                      <template #icon>
                        <n-icon><DeleteOutlined /></n-icon>
                      </template>
                    </n-button>
                  </n-space>
                </div>
                <div
                  v-for="(condition, conditionIndex) in group.conditions"
                  :key="condition.key"
                  class="data-clean-condition"
                >
                  <n-select
                    v-model:value="condition.field"
                    :options="fieldPathOptions"
                    filterable
                    clearable
                    placeholder="字段"
                  />
                  <n-select
                    v-model:value="condition.operator"
                    :options="filterOperatorOptions"
                    placeholder="条件"
                  />
                  <n-input
                    v-if="!['empty', 'not_empty'].includes(condition.operator)"
                    v-model:value="condition.value"
                    placeholder="值"
                  />
                  <n-input
                    v-if="condition.operator === 'between'"
                    v-model:value="condition.valueEnd"
                    placeholder="结束值"
                  />
                  <n-button
                    size="small"
                    type="error"
                    :disabled="group.conditions.length === 1"
                    @click="removeFilterCondition(groupIndex, conditionIndex)"
                  >
                    <template #icon>
                      <n-icon><DeleteOutlined /></n-icon>
                    </template>
                  </n-button>
                </div>
              </div>
            </div>

            <div class="data-clean-section data-clean-field-section">
              <div class="data-clean-toolbar">
                <div class="data-clean-section-title">字段清洗</div>
                <n-space>
                  <n-button size="small" @click="applySelectedTemplate">导入模板</n-button>
                  <n-button size="small" :loading="sampleLoading" @click="handleSample">
                    采样字段
                  </n-button>
                  <n-button size="small" :loading="fieldLoading" @click="loadCollectedFields">
                    刷新采集字段
                  </n-button>
                  <n-button size="small" type="primary" @click="addFieldRow()">
                    <template #icon>
                      <n-icon><PlusOutlined /></n-icon>
                    </template>
                    添加字段
                  </n-button>
                </n-space>
              </div>
              <n-empty v-if="fieldRows.length === 0" description="暂无字段，请采样或手动添加" />
              <div v-else class="data-clean-field-list">
                <div class="data-clean-field-table">
                  <div
                    class="data-clean-field-header"
                    :style="{ gridTemplateColumns: fieldColumnTemplate }"
                  >
                    <div
                      v-for="(column, columnIndex) in fieldColumns"
                      :key="column.key"
                      class="data-clean-field-th"
                    >
                      <span>{{ column.title }}</span>
                      <span
                        v-if="column.resizable"
                        class="data-clean-field-resizer"
                        @mousedown.prevent="startResizeFieldColumn(columnIndex, $event)"
                      />
                    </div>
                  </div>
                  <template v-for="(field, index) in fieldRows" :key="field.key">
                    <div
                      class="data-clean-field-row"
                      :style="{ gridTemplateColumns: fieldColumnTemplate }"
                    >
                      <div class="data-clean-field-original">
                        <n-input
                          v-model:value="field.fieldPath"
                          placeholder="message.level"
                          :bordered="false"
                        />
                      </div>
                      <n-select v-model:value="field.fieldType" :options="fieldTypeOptions" />
                      <n-switch v-model:value="field.enabled" />
                      <n-input v-model:value="field.fieldName" />
                      <n-select
                        v-model:value="field.action"
                        :options="cleanActionOptions"
                        clearable
                        placeholder="清洗功能"
                      />
                      <div
                        class="data-clean-field-params"
                        :class="{
                          'data-clean-field-params-double': field.action === 'regex_replace',
                        }"
                      >
                        <n-input
                          v-if="field.action === 'rename'"
                          v-model:value="field.targetField"
                          placeholder="目标字段"
                        />
                        <n-input
                          v-else-if="field.action === 'default' || field.action === 'set_field'"
                          v-model:value="field.defaultValue"
                          placeholder="字段值"
                        />
                        <template v-else-if="field.action === 'regex_replace'">
                          <n-input v-model:value="field.pattern" placeholder="正则" />
                          <n-input v-model:value="field.replacement" placeholder="替换为" />
                        </template>
                      </div>
                      <n-button size="small" type="error" @click="removeFieldRow(index)">
                        <template #icon>
                          <n-icon><DeleteOutlined /></n-icon>
                        </template>
                      </n-button>
                    </div>
                    <div v-if="field.sampleValues.length > 0" class="data-clean-field-samples">
                      <span>样例：</span>
                      <n-tag
                        v-for="sample in field.sampleValues"
                        :key="sample"
                        size="small"
                        class="data-clean-sample-tag"
                      >
                        {{ sample }}
                      </n-tag>
                    </div>
                  </template>
                </div>
              </div>
            </div>
          </div>

          <div v-show="currentStep === 3" class="data-clean-pane">
            <n-grid
              cols="1 s:1 m:2 l:2 xl:2 2xl:2"
              responsive="screen"
              :x-gap="16"
              class="data-clean-control-grid"
            >
              <n-gi>
                <n-form-item label="分发模式">
                  <n-select v-model:value="sinkRuleForm.mode" :options="sinkModeOptions" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="失败策略">
                  <n-select
                    v-model:value="sinkRuleForm.failurePolicy"
                    :options="sinkFailurePolicyOptions"
                  />
                </n-form-item>
              </n-gi>
            </n-grid>

            <div class="data-clean-section">
              <div class="data-clean-toolbar">
                <div class="data-clean-section-title">输出目标</div>
                <n-button size="small" type="primary" @click="addDispatchTarget()">
                  <template #icon>
                    <n-icon><PlusOutlined /></n-icon>
                  </template>
                  添加目标
                </n-button>
              </div>
              <div
                v-for="(target, index) in dispatchTargets"
                :key="target.key"
                class="data-clean-target"
              >
                <n-grid
                  cols="1 s:1 m:2 l:4 xl:4 2xl:4"
                  responsive="screen"
                  :x-gap="12"
                  class="data-clean-target-grid"
                >
                  <n-gi>
                    <n-form-item label="启用">
                      <n-switch v-model:value="target.enabled" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="目标类型">
                      <n-select
                        v-model:value="target.type"
                        :options="sinkTargetTypeOptions"
                        @update:value="handleTargetTypeChange(target)"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi v-if="target.type === 'sink'" span="2" class="data-clean-wide-grid-item">
                    <n-form-item label="输出数据源">
                      <n-select
                        v-model:value="target.sinkCode"
                        :options="sinkOptions"
                        filterable
                        clearable
                        placeholder="请选择输出数据源"
                        @update:value="handleTargetSinkChange(target)"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="删除">
                      <n-button
                        size="small"
                        type="error"
                        :disabled="dispatchTargets.length === 1"
                        @click="removeDispatchTarget(index)"
                      >
                        <template #icon>
                          <n-icon><DeleteOutlined /></n-icon>
                        </template>
                      </n-button>
                    </n-form-item>
                  </n-gi>
                  <n-gi span="4" v-if="target.type === 'sink'">
                    <n-grid
                      cols="1 s:1 m:2 l:2 xl:2 2xl:2"
                      responsive="screen"
                      :x-gap="12"
                      class="data-clean-target-config-grid"
                    >
                      <template v-if="selectedDispatchSinkType(target) === 'kafka'">
                        <n-gi>
                          <n-form-item label="输出 Topic" required>
                            <n-input
                              v-model:value="target.config.topic"
                              placeholder="如 clean-events"
                            />
                          </n-form-item>
                        </n-gi>
                        <n-gi>
                          <n-form-item label="Key 模板">
                            <n-input
                              v-model:value="target.config.keyTemplate"
                              placeholder="{{taskCode}}"
                            />
                          </n-form-item>
                        </n-gi>
                      </template>
                      <template v-else-if="selectedDispatchSinkType(target) === 's3'">
                        <n-gi>
                          <n-form-item label="Bucket" required>
                            <n-input
                              v-model:value="target.config.bucket"
                              placeholder="cleaned-events"
                            />
                          </n-form-item>
                        </n-gi>
                        <n-gi>
                          <n-form-item label="Object Key">
                            <n-input
                              v-model:value="target.config.keyTemplate"
                              placeholder="{{eventId}}.json"
                            />
                          </n-form-item>
                        </n-gi>
                      </template>
                      <template
                        v-else-if="
                          selectedDispatchSinkType(target) === 'elasticsearch' ||
                          selectedDispatchSinkType(target) === 'opensearch'
                        "
                      >
                        <n-gi>
                          <n-form-item label="输出 Index" required>
                            <n-input v-model:value="target.config.index" />
                          </n-form-item>
                        </n-gi>
                      </template>
                      <template v-else-if="selectedDispatchSinkType(target) === 'splunk'">
                        <n-gi>
                          <n-form-item label="输出 Index">
                            <n-input v-model:value="target.config.index" />
                          </n-form-item>
                        </n-gi>
                        <n-gi>
                          <n-form-item label="Source">
                            <n-input v-model:value="target.config.source" />
                          </n-form-item>
                        </n-gi>
                        <n-gi>
                          <n-form-item label="SourceType">
                            <n-input v-model:value="target.config.sourcetype" />
                          </n-form-item>
                        </n-gi>
                      </template>
                      <template v-else-if="selectedDispatchSinkType(target) === 'http'">
                        <n-gi>
                          <n-form-item label="请求路径">
                            <n-input v-model:value="target.config.path" />
                          </n-form-item>
                        </n-gi>
                        <n-gi>
                          <n-form-item label="超时秒数">
                            <n-input-number
                              v-model:value="target.config.timeoutSeconds"
                              class="data-clean-full"
                              :min="1"
                              :show-button="false"
                            />
                          </n-form-item>
                        </n-gi>
                      </template>
                      <template v-else>
                        <n-gi span="2" class="data-clean-wide-grid-item">
                          <n-form-item label="消息模板">
                            <n-input
                              v-model:value="target.config.messageTemplate"
                              type="textarea"
                              placeholder="策略告警：{{message}}"
                              :autosize="{ minRows: 2, maxRows: 4 }"
                            />
                          </n-form-item>
                        </n-gi>
                      </template>
                    </n-grid>
                  </n-gi>
                </n-grid>
              </div>
            </div>

            <n-form-item label="备注" path="remark" class="data-clean-wide-form-item">
              <n-input
                v-model:value="formValue.remark"
                type="textarea"
                placeholder="备注"
                :autosize="{ minRows: 2, maxRows: 4 }"
              />
            </n-form-item>
          </div>
        </n-form>
        <div class="data-clean-page-actions">
          <n-space>
            <n-button v-if="currentStep > 1" @click="currentStep--">上一步</n-button>
            <n-button v-if="currentStep < 3" type="primary" @click="handleNextStep"
              >下一步</n-button
            >
            <n-button :loading="testBtnLoading" @click="handleTest">测试配置</n-button>
            <n-button @click="closeForm">返回</n-button>
            <n-button type="info" :loading="formBtnLoading" @click="confirmForm">保存</n-button>
          </n-space>
        </div>
      </n-spin>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, nextTick, onBeforeUnmount, reactive, ref, watch } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { PlusOutlined, DeleteOutlined } from '@vicons/antd';
  import { useMessage } from 'naive-ui';
  import { useDictStore } from '@/store/modules/dict';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { useTabsViewStore } from '@/store/modules/tabsView';
  import { KafkaTopics, List as ConnectorList } from '@/api/dataConnector';
  import { List as TemplateList } from '@/api/dataFieldTemplate';
  import { Edit, FieldList, Sample, Test, View } from '@/api/dataClean';
  import { State, newState, rules } from './model';

  type FieldSource = 'agent' | 'sample' | 'template' | 'manual';

  interface SelectOption {
    label: string;
    value: string | number;
    disabled?: boolean;
    connectorType?: string;
    code?: string;
  }

  interface ConnectorOption extends SelectOption {
    connectorType: string;
    code: string;
    config?: Record<string, any>;
  }

  interface TemplateItem {
    id: number;
    name: string;
    code: string;
    eventType: string;
    fields: any;
  }

  interface FieldRow {
    key: number;
    fieldPath: string;
    fieldName: string;
    fieldType: string;
    enabled: boolean;
    action: string | null;
    targetField: string;
    defaultValue: string;
    pattern: string;
    replacement: string;
    required: boolean;
    description: string;
    sampleValues: string[];
    source: FieldSource;
  }

  interface FilterCondition {
    key: number;
    field: string;
    operator: string;
    value: string;
    valueEnd: string;
  }

  interface FilterGroup {
    key: number;
    logic: 'and' | 'or';
    conditions: FilterCondition[];
  }

  interface DispatchTarget {
    key: number;
    enabled: boolean;
    type: string;
    sinkCode: string;
    config: Record<string, any>;
  }

  const message = useMessage();
  const dict = useDictStore();
  const route = useRoute();
  const router = useRouter();
  const settingStore = useProjectSettingStore();
  const tabsViewStore = useTabsViewStore();

  const loading = ref(false);
  const sampleLoading = ref(false);
  const fieldLoading = ref(false);
  const testBtnLoading = ref(false);
  const formBtnLoading = ref(false);
  const currentStep = ref(1);
  const formRef = ref<any>({});
  const formValue = ref<State>(newState(null));
  const sourceOptions = ref<ConnectorOption[]>([]);
  const sinkOptions = ref<ConnectorOption[]>([]);
  const kafkaTopicLoading = ref(false);
  const kafkaTopicLoadingSourceId = ref<number | null>(null);
  const kafkaTopicLoadedSourceId = ref<number | null>(null);
  const kafkaTopicOptions = ref<SelectOption[]>([]);
  const kafkaGroupIdManual = ref(false);
  const templates = ref<TemplateItem[]>([]);
  const fieldRows = ref<FieldRow[]>([]);
  const filterGroups = ref<FilterGroup[]>([]);
  const dispatchTargets = ref<DispatchTarget[]>([]);
  const fieldColumnWidths = ref([360, 130, 76, 300, 170, 430, 54]);

  const fieldColumns = [
    { key: 'fieldPath', title: '原始字段', resizable: true },
    { key: 'fieldType', title: '字段类型', resizable: true },
    { key: 'enabled', title: '启用', resizable: true },
    { key: 'fieldName', title: '字段名称', resizable: true },
    { key: 'action', title: '清洗功能', resizable: true },
    { key: 'params', title: '参数', resizable: true },
    { key: 'operation', title: '', resizable: false },
  ];
  const fieldColumnMinWidths = [180, 96, 64, 180, 112, 220, 42];

  let resizingFieldColumn: {
    index: number;
    startWidth: number;
    startX: number;
  } | null = null;

  let seed = 1;

  const pageTitle = computed(() => {
    return formValue.value.id > 0 ? `编辑数据清洗 #${formValue.value.id}` : '新增数据清洗';
  });
  const fieldColumnTemplate = computed(() =>
    fieldColumnWidths.value.map((width) => `${width}px`).join(' ')
  );
  const selectedSource = computed(() => {
    return sourceOptions.value.find(
      (item) => Number(item.value) === Number(formValue.value.sourceId)
    );
  });
  const selectedSourceType = computed(() => selectedSource.value?.connectorType || 'manual');
  const templateOptions = computed<SelectOption[]>(() => {
    return templates.value.map((item) => ({
      label: `${item.name}（${item.code}）`,
      value: item.id,
    }));
  });
  const selectedTemplate = computed(() => {
    return templates.value.find((item) => Number(item.id) === Number(formValue.value.templateId));
  });
  const selectedDispatchSink = (target: DispatchTarget) => {
    return sinkOptions.value.find((item) => item.code === target.sinkCode);
  };
  const selectedDispatchSinkType = (target: DispatchTarget) => {
    return selectedDispatchSink(target)?.connectorType || '';
  };
  const fieldPathOptions = computed<SelectOption[]>(() => {
    return fieldRows.value
      .filter((item) => item.fieldPath)
      .map((item) => ({
        label: item.fieldName ? `${item.fieldName}（${item.fieldPath}）` : item.fieldPath,
        value: item.fieldPath,
      }));
  });
  const fieldTypeOptions = computed(() => dict.getOptionUnRef('DataFieldTypeOptions'));
  const cleanActionOptions = computed(() => dict.getOptionUnRef('DataCleanActionOptions'));
  const filterOperatorOptions = computed(() =>
    dict.getOptionUnRef('DataCleanFilterOperatorOptions')
  );
  const sinkTargetTypeOptions = computed(() => dict.getOptionUnRef('DataSinkTargetTypeOptions'));
  const cleanErrorPolicyOptions = computed(() =>
    dict.getOptionUnRef('DataCleanErrorPolicyOptions')
  );
  const unknownFieldPolicyOptions = computed(() =>
    dict.getOptionUnRef('DataUnknownFieldPolicyOptions')
  );

  const kafkaStartModeOptions = computed(() => dict.getOptionUnRef('DataKafkaStartModeOptions'));
  const sampleModeOptions = computed(() => dict.getOptionUnRef('DataSampleModeOptions'));
  const payloadFormatOptions = computed(() => dict.getOptionUnRef('DataPayloadFormatOptions'));
  const sinkModeOptions = computed(() => dict.getOptionUnRef('DataSinkModeOptions'));
  const sinkFailurePolicyOptions = computed(() =>
    dict.getOptionUnRef('DataSinkFailurePolicyOptions')
  );

  const sourceConfigForm = reactive({
    bucket: '',
    description: '后台手工触发',
    groupId: '',
    logMaxBytes: 10485760,
    maxObjects: 100,
    maxBytes: 1048576,
    maxWaitSeconds: 3,
    minBytes: 1,
    path: '',
    payloadFormat: 'json',
    prefix: '',
    rateLimit: 0,
    sampleLimit: 50,
    sampleMode: 'latest',
    samplePayload: '',
    sampleTimeoutSeconds: 10,
    startMode: 'latest',
    startOffset: 0,
    token: '',
    topic: '',
  });

  const cleanRuleForm = reactive({
    filterEnabled: false,
    parseDepth: 8,
  });

  const sinkRuleForm = reactive({
    failurePolicy: 'continue',
    mode: 'parallel',
  });

  function normalizeObject(value: any): Record<string, any> {
    if (!value) return {};
    if (typeof value === 'string') {
      try {
        return JSON.parse(value);
      } catch {
        return {};
      }
    }
    if (typeof value === 'object') {
      return value;
    }
    return {};
  }

  function normalizeArray(value: any): any[] {
    if (!value) return [];
    if (Array.isArray(value)) return value;
    if (typeof value === 'string') {
      try {
        const parsed = JSON.parse(value);
        return Array.isArray(parsed) ? parsed : [];
      } catch {
        return [];
      }
    }
    return [];
  }

  function normalizeSampleValues(value: any): string[] {
    if (!value) return [];
    if (Array.isArray(value)) {
      return value
        .map((item) => String(item))
        .filter(Boolean)
        .slice(0, 5);
    }
    return [String(value)].filter(Boolean);
  }

  function fieldNameFromPath(path: string) {
    const normalized = String(path || '').replace(/\[\]$/g, '');
    const parts = normalized.split('.');
    return parts[parts.length - 1] || normalized;
  }

  function newFieldRow(raw: Record<string, any> = {}, source: FieldSource = 'manual'): FieldRow {
    const fieldPath = raw.fieldPath ?? raw.path ?? raw.name ?? '';
    return {
      key: seed++,
      fieldPath,
      fieldName: raw.fieldName ?? raw.title ?? fieldNameFromPath(fieldPath),
      fieldType: raw.fieldType ?? raw.type ?? 'string',
      enabled: raw.enabled ?? true,
      action: raw.action ?? null,
      targetField: raw.targetField ?? raw.to ?? '',
      defaultValue: raw.defaultValue ?? raw.value ?? '',
      pattern: raw.pattern ?? '',
      replacement: raw.replacement ?? '',
      required: Boolean(raw.required),
      description: raw.description ?? '',
      sampleValues: normalizeSampleValues(raw.sampleValues ?? raw.sampleValue),
      source,
    };
  }

  function addFieldRow(raw: Record<string, any> = {}, source: FieldSource = 'manual') {
    fieldRows.value.push(newFieldRow(raw, source));
  }

  function removeFieldRow(index: number) {
    fieldRows.value.splice(index, 1);
  }

  function startResizeFieldColumn(index: number, event: MouseEvent) {
    if (index >= fieldColumnWidths.value.length - 1) {
      return;
    }
    resizingFieldColumn = {
      index,
      startWidth: fieldColumnWidths.value[index],
      startX: event.clientX,
    };
    window.addEventListener('mousemove', handleResizeFieldColumn);
    window.addEventListener('mouseup', stopResizeFieldColumn);
  }

  function handleResizeFieldColumn(event: MouseEvent) {
    if (!resizingFieldColumn) {
      return;
    }
    const { index, startWidth, startX } = resizingFieldColumn;
    const nextWidths = [...fieldColumnWidths.value];
    nextWidths[index] = Math.max(fieldColumnMinWidths[index], startWidth + event.clientX - startX);
    fieldColumnWidths.value = nextWidths;
  }

  function stopResizeFieldColumn() {
    resizingFieldColumn = null;
    window.removeEventListener('mousemove', handleResizeFieldColumn);
    window.removeEventListener('mouseup', stopResizeFieldColumn);
  }

  function mergeFieldRows(rows: Record<string, any>[], source: FieldSource) {
    rows.forEach((raw) => {
      const next = newFieldRow(raw, source);
      if (!next.fieldPath) return;
      const existing = fieldRows.value.find((item) => item.fieldPath === next.fieldPath);
      if (!existing) {
        fieldRows.value.push(next);
        return;
      }
      if (!existing.fieldType || existing.fieldType === 'unknown') {
        existing.fieldType = next.fieldType;
      }
      if (!existing.fieldName) {
        existing.fieldName = next.fieldName;
      }
      next.sampleValues.forEach((sample) => {
        if (!existing.sampleValues.includes(sample) && existing.sampleValues.length < 5) {
          existing.sampleValues.push(sample);
        }
      });
    });
  }

  function newFilterCondition(raw: Record<string, any> = {}): FilterCondition {
    return {
      key: seed++,
      field: raw.field ?? raw.fieldPath ?? '',
      operator: raw.operator ?? 'eq',
      value: raw.value ?? '',
      valueEnd: raw.valueEnd ?? '',
    };
  }

  function newFilterGroup(raw: Record<string, any> = {}): FilterGroup {
    const conditions = normalizeArray(raw.conditions);
    return {
      key: seed++,
      logic: raw.logic === 'or' ? 'or' : 'and',
      conditions:
        conditions.length > 0 ? conditions.map(newFilterCondition) : [newFilterCondition()],
    };
  }

  function resetFilterGroups(filter: Record<string, any> = {}) {
    cleanRuleForm.filterEnabled = Boolean(filter.enabled);
    const groups = normalizeArray(filter.groups);
    filterGroups.value = groups.length > 0 ? groups.map(newFilterGroup) : [newFilterGroup()];
  }

  function addFilterGroup() {
    filterGroups.value.push(newFilterGroup());
  }

  function removeFilterGroup(index: number) {
    filterGroups.value.splice(index, 1);
  }

  function addFilterCondition(groupIndex: number) {
    if (!filterGroups.value[groupIndex]) {
      filterGroups.value.push(newFilterGroup());
    }
    filterGroups.value[groupIndex].conditions.push(newFilterCondition());
  }

  function removeFilterCondition(groupIndex: number, conditionIndex: number) {
    filterGroups.value[groupIndex].conditions.splice(conditionIndex, 1);
  }

  function newDispatchTarget(raw: Record<string, any> = {}): DispatchTarget {
    const target = {
      key: seed++,
      enabled: raw.enabled !== false,
      type: raw.type ?? 'strategy_engine',
      sinkCode: raw.sinkCode ?? '',
      config: normalizeObject(raw.config),
    };
    if (target.type === 'sink') {
      target.config = buildDefaultDispatchTargetConfig(selectedDispatchSink(target), target.config);
    }
    return {
      ...target,
    };
  }

  function addDispatchTarget(raw: Record<string, any> = {}) {
    dispatchTargets.value.push(newDispatchTarget(raw));
  }

  function removeDispatchTarget(index: number) {
    dispatchTargets.value.splice(index, 1);
  }

  function handleTargetTypeChange(target: DispatchTarget) {
    if (target.type === 'strategy_engine') {
      target.sinkCode = '';
      target.config = {};
      return;
    }
    target.config = buildDefaultDispatchTargetConfig(selectedDispatchSink(target));
  }

  function handleTargetSinkChange(target: DispatchTarget) {
    target.config = buildDefaultDispatchTargetConfig(selectedDispatchSink(target), target.config);
  }

  function normalizeTaskSlug() {
    return (
      String(formValue.value.code || formValue.value.name || selectedSourceType.value || 'clean')
        .trim()
        .toLowerCase()
        .replace(/[^a-z0-9_-]+/g, '-')
        .replace(/^-+|-+$/g, '') || 'clean'
    );
  }

  function buildDefaultDispatchTargetConfig(
    sink?: ConnectorOption,
    existing: Record<string, any> = {}
  ) {
    const sinkConfig = sink?.config ?? {};
    const slug = normalizeTaskSlug();
    let defaults: Record<string, any> = {};
    switch (sink?.connectorType) {
      case 'kafka':
        defaults = {
          keyTemplate: sinkConfig.keyTemplate ?? '{{taskCode}}',
          topic: existing.topic ?? `${slug}-events`,
        };
        break;
      case 's3':
        defaults = {
          bucket: sinkConfig.bucket ?? sourceConfigForm.bucket ?? 'cleaned-events',
          keyTemplate: `${slug}/{{eventId}}.json`,
        };
        break;
      case 'elasticsearch':
      case 'opensearch':
        defaults = { index: `${slug}-events` };
        break;
      case 'splunk':
        defaults = {
          index: slug,
          source: formValue.value.code || slug,
          sourcetype: sinkConfig.sourcetype ?? '_json',
        };
        break;
      case 'http':
        defaults = {
          path: sinkConfig.path ?? '',
          timeoutSeconds: Number(sinkConfig.timeoutSeconds) || 5,
        };
        break;
      default:
        defaults = { messageTemplate: sinkConfig.messageTemplate ?? '策略告警：{{message}}' };
        break;
    }
    return { ...defaults, ...existing };
  }

  function randomKafkaGroupSuffix() {
    return Math.random().toString(36).slice(2, 6).padEnd(4, '0');
  }

  function generateKafkaGroupId(topic: string) {
    const topicName = String(topic || '')
      .trim()
      .replace(/\s+/g, '_');
    return topicName ? `dataclean_${topicName}_${randomKafkaGroupSuffix()}` : '';
  }

  function refreshKafkaGroupIdByTopic(force = false) {
    if (selectedSourceType.value !== 'kafka') {
      return;
    }
    if (!force && kafkaGroupIdManual.value) {
      return;
    }
    const groupId = generateKafkaGroupId(sourceConfigForm.topic);
    if (groupId) {
      sourceConfigForm.groupId = groupId;
    }
  }

  function resetSourceConfig(rawConfig: Record<string, any> = {}) {
    kafkaGroupIdManual.value = Boolean(rawConfig.groupId);
    Object.assign(sourceConfigForm, {
      bucket: rawConfig.bucket ?? '',
      description: rawConfig.description ?? '后台手工触发',
      groupId: rawConfig.groupId ?? '',
      logMaxBytes: Number(rawConfig.maxBytes || rawConfig.logMaxBytes) || 10485760,
      maxObjects: Number(rawConfig.maxObjects) || 100,
      maxBytes: Number(rawConfig.maxBytes) || 1048576,
      maxWaitSeconds: Number(rawConfig.maxWaitSeconds) || 3,
      minBytes: Number(rawConfig.minBytes) || 1,
      path: rawConfig.path ?? '',
      payloadFormat: rawConfig.payloadFormat ?? 'json',
      prefix: rawConfig.prefix ?? '',
      rateLimit: Number(rawConfig.rateLimit) || 0,
      sampleLimit: Number(rawConfig.sampleLimit) || 50,
      sampleMode: rawConfig.sampleMode ?? 'latest',
      samplePayload: rawConfig.samplePayload ?? '',
      sampleTimeoutSeconds: Number(rawConfig.sampleTimeoutSeconds) || 10,
      startMode: rawConfig.startMode ?? 'latest',
      startOffset: Number(rawConfig.startOffset) || 0,
      token: rawConfig.token ?? '',
      topic: rawConfig.topic ?? '',
    });
    setKafkaTopicOptions([]);
    kafkaTopicLoadedSourceId.value = null;
    refreshKafkaGroupIdByTopic();
  }

  function setKafkaTopicOptions(topics: string[]) {
    const values = Array.from(
      new Set(
        topics
          .concat(sourceConfigForm.topic ? [sourceConfigForm.topic] : [])
          .map((topic) => String(topic || '').trim())
          .filter(Boolean)
      )
    );
    kafkaTopicOptions.value = values.map((topic) => ({
      label: topic,
      value: topic,
    }));
  }

  function loadKafkaTopics() {
    if (selectedSourceType.value !== 'kafka') {
      setKafkaTopicOptions([]);
      kafkaTopicLoadedSourceId.value = null;
      return;
    }
    const sourceId = Number(formValue.value.sourceId);
    if (!sourceId) {
      setKafkaTopicOptions([]);
      kafkaTopicLoadedSourceId.value = null;
      return;
    }
    if (kafkaTopicLoadingSourceId.value === sourceId) {
      return;
    }
    if (kafkaTopicLoadedSourceId.value === sourceId) {
      return;
    }

    kafkaTopicLoading.value = true;
    kafkaTopicLoadingSourceId.value = sourceId;
    KafkaTopics({ id: sourceId })
      .then((res) => {
        if (Number(formValue.value.sourceId) !== sourceId) {
          return;
        }
        setKafkaTopicOptions(res?.topics || []);
        kafkaTopicLoadedSourceId.value = sourceId;
      })
      .catch((error) => {
        if (Number(formValue.value.sourceId) === sourceId) {
          if (!error) {
            return;
          }
          setKafkaTopicOptions([]);
          const errorMessage =
            error?.message || error?.response?.data?.message || error?.data?.message || '';
          message.warning(
            errorMessage
              ? `读取 Kafka Topic 列表失败：${errorMessage}`
              : '读取 Kafka Topic 列表失败，可手动输入 Topic'
          );
        }
      })
      .finally(() => {
        if (kafkaTopicLoadingSourceId.value === sourceId) {
          kafkaTopicLoading.value = false;
          kafkaTopicLoadingSourceId.value = null;
        }
      });
  }

  function resetCleanConfig(rawConfig: Record<string, any> = {}) {
    cleanRuleForm.parseDepth = Number(rawConfig.parseDepth) || 8;
    resetFilterGroups(normalizeObject(rawConfig.filter));
    fieldRows.value = normalizeArray(rawConfig.fields).map((item) => newFieldRow(item, 'manual'));
  }

  function resetSinkConfig(rawConfig: Record<string, any> = {}) {
    sinkRuleForm.failurePolicy = rawConfig.failurePolicy ?? 'continue';
    sinkRuleForm.mode = rawConfig.mode ?? 'parallel';
    const targets = normalizeArray(rawConfig.targets);
    dispatchTargets.value =
      targets.length > 0 ? targets.map((item) => newDispatchTarget(item)) : [newDispatchTarget()];
  }

  function buildSourceConfig() {
    const common = {
      payloadFormat: sourceConfigForm.payloadFormat,
      samplePayload: sourceConfigForm.samplePayload,
    };
    if (selectedSourceType.value === 'kafka') {
      return {
        ...common,
        groupId: sourceConfigForm.groupId,
        maxBytes: sourceConfigForm.maxBytes,
        maxWaitSeconds: sourceConfigForm.maxWaitSeconds,
        minBytes: sourceConfigForm.minBytes,
        sampleLimit: sourceConfigForm.sampleLimit,
        sampleMode: sourceConfigForm.sampleMode,
        sampleTimeoutSeconds: sourceConfigForm.sampleTimeoutSeconds,
        startMode: sourceConfigForm.startMode,
        startOffset: sourceConfigForm.startOffset,
        topic: sourceConfigForm.topic,
      };
    }
    if (selectedSourceType.value === 's3') {
      return {
        ...common,
        bucket: sourceConfigForm.bucket,
        maxObjects: sourceConfigForm.maxObjects,
        prefix: sourceConfigForm.prefix,
      };
    }
    if (selectedSourceType.value === 'log') {
      return {
        ...common,
        maxBytes: sourceConfigForm.logMaxBytes,
        path: sourceConfigForm.path,
      };
    }
    if (selectedSourceType.value === 'http') {
      return {
        ...common,
        rateLimit: sourceConfigForm.rateLimit,
        token: sourceConfigForm.token,
      };
    }
    return {
      ...common,
      description: sourceConfigForm.description,
    };
  }

  function buildFilterConfig() {
    return {
      enabled: cleanRuleForm.filterEnabled,
      groups: filterGroups.value.map((group) => ({
        logic: group.logic,
        conditions: group.conditions
          .filter((condition) => condition.field && condition.operator)
          .map((condition) => ({
            field: condition.field,
            operator: condition.operator,
            value: condition.value,
            valueEnd: condition.valueEnd,
          })),
      })),
      logic: 'and',
    };
  }

  function fieldRowToStep(field: FieldRow) {
    const step: Record<string, any> = {
      enabled: field.enabled,
      field: field.fieldPath,
      type: field.action,
    };
    if (field.action === 'rename') {
      step.to = field.targetField;
    }
    if (field.action === 'default' || field.action === 'set_field') {
      step.value = field.defaultValue;
    }
    if (field.action === 'regex_replace') {
      step.pattern = field.pattern;
      step.replacement = field.replacement;
    }
    return step;
  }

  function buildCleanConfig() {
    const fields = fieldRows.value
      .filter((field) => field.fieldPath)
      .map((field) => ({
        action: field.action,
        defaultValue: field.defaultValue,
        description: field.description,
        enabled: field.enabled,
        fieldName: field.fieldName,
        fieldPath: field.fieldPath,
        fieldType: field.fieldType,
        pattern: field.pattern,
        replacement: field.replacement,
        required: field.required,
        source: field.source,
        targetField: field.targetField,
      }));
    return {
      fields,
      filter: buildFilterConfig(),
      parseDepth: cleanRuleForm.parseDepth,
      steps: fieldRows.value
        .filter((field) => field.enabled && field.action)
        .map((field) => fieldRowToStep(field)),
    };
  }

  function buildSinkConfig() {
    const targets = dispatchTargets.value
      .filter((target) => target.enabled)
      .map((target) => {
        const row: Record<string, any> = {
          enabled: target.enabled,
          type: target.type,
        };
        if (target.type === 'sink') {
          row.sinkCode = target.sinkCode;
          row.config = target.config;
        }
        return row;
      });
    return {
      failurePolicy: sinkRuleForm.failurePolicy,
      mode: sinkRuleForm.mode,
      targets: targets.length > 0 ? targets : [{ enabled: true, type: 'strategy_engine' }],
    };
  }

  function hasCleanRules() {
    if (cleanRuleForm.filterEnabled) return true;
    return fieldRows.value.some((field) => field.enabled && Boolean(field.action));
  }

  function validateSourceConfig() {
    if (!formValue.value.sourceId) {
      message.error('请选择输入数据源');
      return false;
    }
    if (selectedSourceType.value === 'kafka' && !sourceConfigForm.topic) {
      message.error('请填写 Kafka Topic');
      return false;
    }
    if (selectedSourceType.value === 's3' && !sourceConfigForm.bucket) {
      message.error('请填写 S3 Bucket');
      return false;
    }
    if (selectedSourceType.value === 'log' && !sourceConfigForm.path) {
      message.error('请填写日志路径');
      return false;
    }
    return true;
  }

  function validateSinkConfig() {
    const invalid = dispatchTargets.value.find(
      (target) => target.enabled && target.type === 'sink' && !target.sinkCode
    );
    if (invalid) {
      message.error('请选择输出数据源');
      return false;
    }
    const kafkaWithoutTopic = dispatchTargets.value.find(
      (target) =>
        target.enabled &&
        target.type === 'sink' &&
        selectedDispatchSinkType(target) === 'kafka' &&
        !String(target.config.topic || '').trim()
    );
    if (kafkaWithoutTopic) {
      message.error('请填写 Kafka 输出 Topic');
      return false;
    }
    return true;
  }

  function handleNextStep() {
    if (currentStep.value === 1 && !validateSourceConfig()) {
      return;
    }
    currentStep.value++;
  }

  function handleSourceChange() {
    resetSourceConfig();
    loadKafkaTopics();
  }

  function handleKafkaTopicChange() {
    refreshKafkaGroupIdByTopic();
  }

  function handleKafkaGroupIdChange() {
    kafkaGroupIdManual.value = true;
  }

  function applySelectedTemplate() {
    if (!selectedTemplate.value) {
      message.error('请选择标准模板');
      return;
    }
    mergeFieldRows(normalizeArray(selectedTemplate.value.fields), 'template');
    message.success('模板字段已导入');
  }

  function handleSample(e?: Event) {
    e?.preventDefault();
    if (!validateSourceConfig()) return;
    sampleLoading.value = true;
    Sample({
      cleanConfig: buildCleanConfig(),
      id: formValue.value.id,
      limit: sourceConfigForm.sampleLimit,
      payload: sourceConfigForm.samplePayload,
      sourceConfig: buildSourceConfig(),
      sourceId: formValue.value.sourceId,
      timeoutSeconds: sourceConfigForm.sampleTimeoutSeconds,
    })
      .then((res) => {
        mergeFieldRows(res?.list || [], 'sample');
        message.success(`采样完成，发现 ${res?.total || 0} 个字段`);
      })
      .finally(() => {
        sampleLoading.value = false;
      });
  }

  function loadCollectedFields() {
    if (formValue.value.id < 1) {
      message.error('请先保存任务，再刷新 Agent 采集字段');
      return;
    }
    fieldLoading.value = true;
    FieldList({ taskId: formValue.value.id })
      .then((res) => {
        mergeFieldRows(res?.list || res || [], 'agent');
        message.success(`已加载 ${res?.total || res?.length || 0} 个采集字段`);
      })
      .finally(() => {
        fieldLoading.value = false;
      });
  }

  function handleTest(e: Event) {
    e.preventDefault();
    if (!validateSourceConfig() || !validateSinkConfig()) {
      return;
    }
    testBtnLoading.value = true;
    Test({
      cleanConfig: buildCleanConfig(),
      sinkConfig: buildSinkConfig(),
      sourceConfig: buildSourceConfig(),
      sourceId: formValue.value.sourceId,
    })
      .then((res) => {
        if (res?.success) {
          message.success(res.message || '测试通过');
          return;
        }
        message.error(res?.message || '测试失败');
      })
      .finally(() => {
        testBtnLoading.value = false;
      });
  }

  function confirmForm(e: Event) {
    e.preventDefault();
    formRef.value.validate((errors) => {
      if (errors) {
        message.error('请填写完整信息');
        return;
      }
      if (!validateSourceConfig() || !validateSinkConfig()) {
        return;
      }
      formBtnLoading.value = true;
      Edit({
        ...formValue.value,
        cleanConfig: buildCleanConfig(),
        cleanEnabled: hasCleanRules() ? 1 : 0,
        sinkConfig: buildSinkConfig(),
        sourceId: Number(formValue.value.sourceId) || 0,
        sourceConfig: buildSourceConfig(),
        templateId: Number(formValue.value.templateId) || 0,
      })
        .then((_res) => {
          message.success('操作成功');
          closeForm();
        })
        .finally(() => {
          formBtnLoading.value = false;
        });
    });
  }

  async function closeForm() {
    const currentFullPath = route.fullPath;
    await router.push({ name: 'dataClean' });
    await nextTick();
    closeTabByFullPath(currentFullPath);
  }

  function resetForm() {
    currentStep.value = 1;
    formValue.value = newState(null);
    formValue.value.cleanErrorPolicy = 'skip';
    formValue.value.unknownFieldPolicy = 'selected_only';
    formValue.value.status = 2;
    resetSourceConfig();
    resetCleanConfig();
    resetSinkConfig();
  }

  function loadOptions() {
    return Promise.all([
      ConnectorList({ direction: 'source', page: 1, perPage: 200, status: 1 }).then((res) => {
        sourceOptions.value = (res?.list || []).map((item) => ({
          code: item.code,
          connectorType: item.connectorType,
          label: `${item.name}（${item.code} / ${item.connectorType}）`,
          value: item.id,
        }));
      }),
      ConnectorList({ direction: 'sink', page: 1, perPage: 200, status: 1 }).then((res) => {
        sinkOptions.value = (res?.list || []).map((item) => ({
          code: item.code,
          config: item.config || {},
          connectorType: item.connectorType,
          label: `${item.name}（${item.code} / ${item.connectorType}）`,
          value: item.code,
        }));
      }),
      TemplateList({ page: 1, perPage: 200, status: 1 }).then((res) => {
        templates.value = res?.list || [];
      }),
    ]);
  }

  function closeTabByFullPath(fullPath: string) {
    if (tabsViewStore.tabsList.length <= 1) {
      return;
    }
    const current = tabsViewStore.tabsList.find((item) => item.fullPath === fullPath);
    if (current) {
      tabsViewStore.closeCurrentTab(current);
    }
  }

  function loadPage() {
    const id = Number(route.query.id) || 0;
    loading.value = true;
    loadOptions()
      .then(() => {
        if (id < 1) {
          resetForm();
          return;
        }
        currentStep.value = 1;
        return View({ id }).then((res) => {
          formValue.value = newState(res);
          if (!formValue.value.sourceId) {
            formValue.value.sourceId = null;
          }
          if (!formValue.value.templateId) {
            formValue.value.templateId = null;
          }
          resetSourceConfig(normalizeObject(res?.sourceConfig));
          resetCleanConfig(normalizeObject(res?.cleanConfig));
          resetSinkConfig(normalizeObject(res?.sinkConfig));
          loadKafkaTopics();
          if (fieldRows.value.length === 0) {
            loadCollectedFields();
          }
        });
      })
      .finally(() => {
        loading.value = false;
      });
  }

  watch(
    () => route.query.id,
    () => loadPage(),
    { immediate: true }
  );

  onBeforeUnmount(() => {
    stopResizeFieldColumn();
  });
</script>

<style lang="less" scoped>
  .data-clean-page-card {
    --data-clean-content-width: 980px;
    --data-clean-field-section-width: 1560px;
    --data-clean-form-label-width: 120px;
    --data-clean-steps-width: calc(var(--data-clean-content-width) + 32px);
    min-height: calc(100vh - 180px);
  }

  .data-clean-page-form {
    padding: 16px 0 24px;
  }

  .data-clean-page-actions {
    background: var(--card-color);
    bottom: 0;
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
    position: sticky;
    z-index: 2;
  }

  .data-clean-steps {
    box-sizing: border-box;
    margin: 0 auto 24px;
    max-width: 100%;
    width: var(--data-clean-steps-width);
  }

  :deep(.data-clean-steps > .n-step:last-child) {
    flex: 0 0 auto;
  }

  .data-clean-pane {
    min-height: calc(100vh - 360px);
  }

  .data-clean-section {
    border: 1px solid rgba(128, 128, 128, 0.22);
    border-radius: 6px;
    box-sizing: border-box;
    margin: 16px auto 0;
    max-width: var(--data-clean-content-width);
    padding: 16px;
    width: min(100%, var(--data-clean-content-width));
  }

  .data-clean-section-title {
    color: var(--text-color-1);
    font-weight: 600;
  }

  .data-clean-toolbar,
  .data-clean-filter-head {
    align-items: center;
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
  }

  .data-clean-filter-group,
  .data-clean-field,
  .data-clean-target {
    border: 1px solid rgba(128, 128, 128, 0.22);
    border-radius: 6px;
    margin-bottom: 12px;
    padding: 12px;
  }

  .data-clean-condition {
    align-items: center;
    column-gap: 8px;
    display: grid;
    grid-template-columns:
      minmax(180px, 1.2fr) minmax(140px, 0.8fr) minmax(160px, 1fr) minmax(120px, 0.8fr)
      42px;
    margin-bottom: 8px;
  }

  .data-clean-control-grid {
    box-sizing: border-box;
    margin-left: auto;
    margin-right: auto;
    max-width: var(--data-clean-content-width);
    width: min(100%, var(--data-clean-content-width));
  }

  .data-clean-target-grid {
    box-sizing: border-box;
    margin-left: auto;
    margin-right: auto;
    max-width: var(--data-clean-content-width);
    width: min(100%, var(--data-clean-content-width));
  }

  .data-clean-target-config-grid {
    box-sizing: border-box;
    margin-left: auto;
    margin-right: auto;
    max-width: var(--data-clean-content-width);
    width: min(100%, var(--data-clean-content-width));
  }

  .data-clean-section > .data-clean-section-title {
    margin-left: auto;
    margin-right: auto;
    max-width: 100%;
    width: 100%;
  }

  .data-clean-section > .data-clean-toolbar {
    box-sizing: border-box;
    margin-left: auto;
    margin-right: auto;
    max-width: 100%;
    width: 100%;
  }

  .data-clean-field-section {
    max-width: var(--data-clean-field-section-width);
    width: min(100%, var(--data-clean-field-section-width));
  }

  .data-clean-wide-form-item {
    box-sizing: border-box;
    margin: 16px auto 0;
    max-width: var(--data-clean-content-width);
    width: min(100%, var(--data-clean-content-width));
  }

  :deep(.data-clean-control-grid > *),
  :deep(.data-clean-target-grid > *),
  :deep(.data-clean-target-config-grid > *) {
    min-width: 0;
  }

  :deep(.data-clean-control-grid .n-form-item),
  :deep(.data-clean-target-grid .n-form-item),
  :deep(.data-clean-target-config-grid .n-form-item) {
    min-width: 0;
  }

  :deep(.data-clean-control-grid .n-form-item-blank),
  :deep(.data-clean-target-config-grid .n-form-item-blank) {
    flex: 0 1 340px;
    max-width: min(340px, calc(100% - var(--data-clean-form-label-width)));
    min-width: 0;
    width: 340px;
  }

  :deep(.data-clean-sample-grid-item .n-form-item-blank) {
    flex: 0 1 820px;
    max-width: min(820px, calc(100% - var(--data-clean-form-label-width)));
    min-width: 0;
    width: 820px;
  }

  :deep(.data-clean-wide-grid-item .n-form-item-blank),
  :deep(.data-clean-wide-form-item .n-form-item-blank) {
    flex: 0 1 720px;
    max-width: min(720px, calc(100% - var(--data-clean-form-label-width)));
    min-width: 0;
    width: 720px;
  }

  :deep(.data-clean-control-grid .n-form-item-blank > *),
  :deep(.data-clean-target-config-grid .n-form-item-blank > *),
  :deep(.data-clean-sample-grid-item .n-form-item-blank > *),
  :deep(.data-clean-wide-grid-item .n-form-item-blank > *),
  :deep(.data-clean-wide-form-item .n-form-item-blank > *) {
    max-width: 100%;
    min-width: 0;
  }

  .data-clean-full {
    width: 100%;
  }

  .data-clean-field-list {
    max-height: calc(100vh - 360px);
    min-height: 360px;
    overflow: auto;
  }

  .data-clean-field-table {
    border: 1px solid rgba(128, 128, 128, 0.22);
    border-radius: 6px;
    overflow: hidden;
    min-width: 100%;
    width: max-content;
  }

  .data-clean-field-header,
  .data-clean-field-row {
    align-items: center;
    display: grid;
  }

  .data-clean-field-header {
    background: rgba(128, 128, 128, 0.08);
    color: var(--text-color-2);
    font-weight: 600;
    min-height: 42px;
    position: sticky;
    top: 0;
    z-index: 1;
  }

  .data-clean-field-th,
  .data-clean-field-row > * {
    min-width: 0;
    padding: 8px 10px;
  }

  .data-clean-field-th {
    align-items: center;
    display: flex;
    height: 100%;
    justify-content: space-between;
    position: relative;
  }

  .data-clean-field-resizer {
    align-self: stretch;
    cursor: col-resize;
    margin-right: -10px;
    position: relative;
    width: 10px;
  }

  .data-clean-field-resizer::after {
    background: rgba(128, 128, 128, 0.28);
    bottom: 8px;
    content: '';
    position: absolute;
    right: 4px;
    top: 8px;
    width: 1px;
  }

  .data-clean-field-resizer:hover::after {
    background: var(--primary-color);
  }

  .data-clean-field-row {
    border-top: 1px solid rgba(128, 128, 128, 0.16);
    min-height: 58px;
  }

  .data-clean-field-row:hover {
    background: rgba(128, 128, 128, 0.04);
  }

  .data-clean-field-row > .n-switch {
    justify-self: start;
  }

  .data-clean-field-original {
    font-weight: 600;
  }

  :deep(.data-clean-field-original .n-input__input-el) {
    font-weight: 600;
  }

  .data-clean-field-params {
    column-gap: 8px;
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    min-height: 34px;
  }

  .data-clean-field-params-double {
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  }

  .data-clean-field-samples {
    align-items: center;
    border-top: 1px solid rgba(128, 128, 128, 0.12);
    color: var(--text-color-2);
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    padding: 8px 10px 8px 18px;
  }

  .data-clean-sample-tag {
    max-width: 320px;
  }

  .data-clean-sample-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 8px;
  }

  @media (max-width: 768px) {
    .data-clean-control-grid,
    .data-clean-target-grid,
    .data-clean-target-config-grid,
    .data-clean-section,
    .data-clean-wide-form-item {
      max-width: none;
      width: 100%;
    }

    .data-clean-steps {
      max-width: none;
    }

    :deep(.data-clean-control-grid .n-form-item-blank),
    :deep(.data-clean-target-config-grid .n-form-item-blank),
    :deep(.data-clean-sample-grid-item .n-form-item-blank),
    :deep(.data-clean-wide-grid-item .n-form-item-blank),
    :deep(.data-clean-wide-form-item .n-form-item-blank) {
      flex: 1 1 auto;
      max-width: none;
      width: 100%;
    }

    .data-clean-condition {
      grid-template-columns: 1fr;
      row-gap: 8px;
    }

    .data-clean-toolbar,
    .data-clean-filter-head {
      align-items: flex-start;
      flex-direction: column;
      row-gap: 8px;
    }
  }
</style>
