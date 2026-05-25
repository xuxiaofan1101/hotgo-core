import { ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { renderOptionTag, MemberSumma } from '@/utils';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  public id = 0;
  public sourceId: number | null = null;
  public sourceName = '';
  public sourceCode = '';
  public sourceType = '';
  public templateId: number | null = null;
  public templateName = '';
  public name = '';
  public code = '';
  public eventType = 'custom.event';
  public sourceConfig: Record<string, any> | null = null;
  public cleanEnabled = 0;
  public cleanConfig: Record<string, any> | null = null;
  public cleanConfigVersion = 1;
  public fieldSchemaVersion = 1;
  public unknownFieldPolicy = 'selected_only';
  public cleanErrorPolicy = 'skip';
  public sinkConfig: Record<string, any> | null = null;
  public status = 2;
  public remark = '';
  public createdBy = 0;
  public createdBySumma?: null | MemberSumma = null;
  public updatedBy = 0;
  public updatedBySumma?: null | MemberSumma = null;
  public createdAt = '';
  public updatedAt = '';

  constructor(state?: Partial<State>) {
    if (state) {
      Object.assign(this, state);
    }
  }
}

export function newState(state: State | Record<string, any> | null): State {
  if (state !== null) {
    if (state instanceof State) {
      return cloneDeep(state);
    }
    return new State(state);
  }
  return new State();
}

export const rules = {
  name: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入清洗任务名称',
  },
  code: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入清洗任务编码',
  },
  sourceId: {
    required: true,
    trigger: ['blur', 'change'],
    type: 'number',
    message: '请选择输入数据源',
  },
};

export const schemas = ref<FormSchema[]>([
  {
    field: 'keyword',
    component: 'NInput',
    label: '关键词',
    componentProps: {
      placeholder: '搜索任务/编码/数据源',
    },
  },
  {
    field: 'sourceType',
    component: 'NSelect',
    label: '数据源类型',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择类型',
      options: dict.getOption('DataConnectorTypeOptions'),
    },
  },
  {
    field: 'status',
    component: 'NSelect',
    label: '状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择状态',
      options: dict.getOption('sys_normal_disable'),
    },
  },
  {
    field: 'createdAt',
    component: 'NDatePicker',
    label: '创建时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
]);

export const columns = [
  {
    title: '任务名称',
    key: 'name',
    align: 'left',
    width: 180,
  },
  {
    title: '编码',
    key: 'code',
    align: 'left',
    width: 160,
  },
  {
    title: '数据源',
    key: 'sourceCode',
    align: 'left',
    width: 180,
    render(row: State) {
      return row.sourceName ? `${row.sourceName}（${row.sourceCode}）` : row.sourceCode;
    },
  },
  {
    title: '类型',
    key: 'sourceType',
    align: 'left',
    width: 130,
    render(row: State) {
      return renderOptionTag('DataConnectorTypeOptions', row.sourceType);
    },
  },
  {
    title: '清洗',
    key: 'cleanEnabled',
    align: 'left',
    width: 90,
    render(row: State) {
      return row.cleanEnabled === 1 ? '已配置' : '直通';
    },
  },
  {
    title: '状态',
    key: 'status',
    align: 'left',
    width: 100,
    render(row: State) {
      return renderOptionTag('sys_normal_disable', row.status);
    },
  },
  {
    title: '创建时间',
    key: 'createdAt',
    align: 'left',
    width: 180,
  },
];

export function loadOptions() {
  dict.loadOptions([
    'DataConnectorTypeOptions',
    'DataCleanErrorPolicyOptions',
    'DataUnknownFieldPolicyOptions',
    'DataFieldTypeOptions',
    'DataCleanActionOptions',
    'DataCleanFilterOperatorOptions',
    'DataSinkTargetTypeOptions',
    'DataKafkaStartModeOptions',
    'DataSampleModeOptions',
    'DataPayloadFormatOptions',
    'DataSinkModeOptions',
    'DataSinkFailurePolicyOptions',
    'sys_normal_disable',
  ]);
}
