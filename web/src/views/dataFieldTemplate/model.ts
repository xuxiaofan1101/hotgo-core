import { ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { renderOptionTag, MemberSumma } from '@/utils';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  public id = 0;
  public name = '';
  public code = '';
  public eventType = 'custom.event';
  public fields: any[] | null = null;
  public examples: any[] | null = null;
  public fieldCount = 0;
  public status = 1;
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
    message: '请输入模板名称',
  },
  code: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入模板编码',
  },
  eventType: {
    required: true,
    trigger: ['blur', 'input'],
    type: 'string',
    message: '请输入事件类型',
  },
};

export const schemas = ref<FormSchema[]>([
  {
    field: 'keyword',
    component: 'NInput',
    label: '关键词',
    componentProps: {
      placeholder: '搜索名称/编码/事件类型',
    },
  },
  {
    field: 'eventType',
    component: 'NInput',
    label: '事件类型',
    componentProps: {
      placeholder: '如 nginx.access',
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
    title: '模板名称',
    key: 'name',
    align: 'left',
    width: 180,
  },
  {
    title: '模板编码',
    key: 'code',
    align: 'left',
    width: 180,
  },
  {
    title: '事件类型',
    key: 'eventType',
    align: 'left',
    width: 180,
  },
  {
    title: '字段数',
    key: 'fieldCount',
    align: 'left',
    width: 100,
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
  dict.loadOptions(['DataFieldTypeOptions', 'sys_normal_disable']);
}
