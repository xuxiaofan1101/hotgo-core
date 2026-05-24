import { h, ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { NTag, NTooltip } from 'naive-ui';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { renderOptionTag } from '@/utils';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  public id = 0;
  public agentId = '';
  public name = '';
  public hostname = '';
  public version = '';
  public agentIp: string[] = [];
  public registerStatus = 'pending';
  public dispatchStatus = 'disabled';
  public dispatchAllowed = false;
  public onlineStatus = 'offline';
  public online = false;
  public lastSeenAt = '';
  public approvedBy = 0;
  public approvedAt = '';
  public disabledAt = '';
  public remark = '';
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

export const schemas = ref<FormSchema[]>([
  {
    field: 'keyword',
    component: 'NInput',
    label: '关键词',
    componentProps: {
      placeholder: '搜索Agent ID/名称/主机名',
    },
  },
  {
    field: 'registerStatus',
    component: 'NSelect',
    label: '注册状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择注册状态',
      options: dict.getOption('DataAgentRegisterStatusOptions'),
    },
  },
  {
    field: 'dispatchStatus',
    component: 'NSelect',
    label: '调度状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择调度状态',
      options: dict.getOption('DataAgentDispatchStatusOptions'),
    },
  },
  {
    field: 'onlineStatus',
    component: 'NSelect',
    label: '在线状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择在线状态',
      options: dict.getOption('DataAgentOnlineStatusOptions'),
    },
  },
  {
    field: 'createdAt',
    component: 'NDatePicker',
    label: '注册时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
]);

function renderAgentIp(row: State) {
  const text = row.agentIp && row.agentIp.length > 0 ? row.agentIp.join(', ') : '-';
  return h(NTooltip, null, {
    trigger: () => h('span', { class: 'agent-ip-text' }, text),
    default: () => text,
  });
}

export const columns = [
  {
    title: 'Agent ID',
    key: 'agentId',
    align: 'left',
    width: 210,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '节点名称',
    key: 'name',
    align: 'left',
    width: 140,
    render(row: State) {
      return row.name || '-';
    },
  },
  {
    title: '主机名',
    key: 'hostname',
    align: 'left',
    width: 150,
    render(row: State) {
      return row.hostname || '-';
    },
  },
  {
    title: 'Agent IP',
    key: 'agentIp',
    align: 'left',
    width: 180,
    render: renderAgentIp,
  },
  {
    title: '在线',
    key: 'onlineStatus',
    align: 'left',
    width: 90,
    render(row: State) {
      return renderOptionTag('DataAgentOnlineStatusOptions', row.onlineStatus);
    },
  },
  {
    title: '注册状态',
    key: 'registerStatus',
    align: 'left',
    width: 110,
    render(row: State) {
      return renderOptionTag('DataAgentRegisterStatusOptions', row.registerStatus);
    },
  },
  {
    title: '调度状态',
    key: 'dispatchStatus',
    align: 'left',
    width: 110,
    render(row: State) {
      return row.dispatchAllowed
        ? h(NTag, { type: 'success', size: 'small' }, { default: () => '可调度' })
        : renderOptionTag('DataAgentDispatchStatusOptions', row.dispatchStatus);
    },
  },
  {
    title: '版本',
    key: 'version',
    align: 'left',
    width: 110,
    render(row: State) {
      return row.version || '-';
    },
  },
  {
    title: '最近心跳',
    key: 'lastSeenAt',
    align: 'left',
    width: 180,
    render(row: State) {
      return row.lastSeenAt || '-';
    },
  },
  {
    title: '注册时间',
    key: 'createdAt',
    align: 'left',
    width: 180,
  },
];

export function loadOptions() {
  dict.loadOptions([
    'DataAgentRegisterStatusOptions',
    'DataAgentDispatchStatusOptions',
    'DataAgentOnlineStatusOptions',
  ]);
}
