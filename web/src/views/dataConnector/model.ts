import { ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { renderOptionTag, MemberSumma } from '@/utils';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  public id = 0; // 连接ID
  public name = ''; // 连接名称
  public code = ''; // 连接编码
  public direction = 'source'; // 连接方向
  public connectorType = 'http'; // 连接类型
  public config: Record<string, any> | null = null; // 连接配置
  public status = 1; // 状态
  public remark = ''; // 备注
  public createdBy = 0; // 创建者
  public createdBySumma?: null | MemberSumma = null; // 创建者摘要信息
  public updatedBy = 0; // 更新者
  public updatedBySumma?: null | MemberSumma = null; // 更新者摘要信息
  public deletedBy = 0; // 删除者
  public createdAt = ''; // 创建时间
  public updatedAt = ''; // 修改时间
  public deletedAt = ''; // 删除时间

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

// 表单验证规则
export const rules = {
  name: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入连接名称',
  },
  code: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入连接编码',
  },
  direction: {
    required: true,
    trigger: ['blur', 'change'],
    type: 'string',
    message: '请选择连接方向',
  },
  connectorType: {
    required: true,
    trigger: ['blur', 'change'],
    type: 'string',
    message: '请选择连接类型',
  },
  status: {
    required: true,
    trigger: ['blur', 'change'],
    type: 'number',
    message: '请选择状态',
  },
};

// 表格搜索表单
export const schemas = ref<FormSchema[]>([
  {
    field: 'keyword',
    component: 'NInput',
    label: '关键词',
    componentProps: {
      placeholder: '搜索名称/编码',
    },
  },
  {
    field: 'direction',
    component: 'NSelect',
    label: '用途',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择用途',
      options: dict.getOption('DataConnectorDirectionOptions'),
    },
  },
  {
    field: 'connectorType',
    component: 'NSelect',
    label: '类型',
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

// 表格列
export const columns = [
  {
    title: '用途',
    key: 'direction',
    align: 'left',
    width: 100,
    render(row: State) {
      return renderOptionTag('DataConnectorDirectionOptions', row.direction);
    },
  },
  {
    title: '数据源名称',
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
    title: '类型',
    key: 'connectorType',
    align: 'left',
    width: 140,
    render(row: State) {
      return renderOptionTag('DataConnectorTypeOptions', row.connectorType);
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

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions([
    'DataConnectorDirectionOptions',
    'DataConnectorTypeOptions',
    'DataConnectorSourceTypeOptions',
    'DataConnectorSinkTypeOptions',
    'DataKafkaAuthModeOptions',
    'DataS3CredentialModeOptions',
    'DataLogLevelOptions',
    'sys_normal_disable',
  ]);
}
