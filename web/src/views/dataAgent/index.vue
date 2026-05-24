<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="节点管理" />
    </div>
    <n-card :bordered="false" class="proCard">
      <BasicForm
        ref="searchFormRef"
        @register="register"
        @submit="reloadTable"
        @reset="reloadTable"
        @keyup.enter="reloadTable"
      />
      <BasicTable
        ref="actionRef"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        :actionColumn="actionColumn"
        :scroll-x="scrollX"
        :resizeHeightOffset="-10000"
      >
        <template #tableTitle>
          <n-button
            type="primary"
            :loading="refreshing"
            @click="reloadTable"
            class="min-left-space"
          >
            <template #icon>
              <n-icon>
                <ReloadOutlined />
              </n-icon>
            </template>
            刷新
          </n-button>
        </template>
      </BasicTable>
    </n-card>

    <n-modal v-model:show="detailVisible" preset="card" title="节点详情" class="agent-detail-modal">
      <n-descriptions v-if="current" bordered label-placement="left" :column="2">
        <n-descriptions-item label="Agent ID" :span="2">
          {{ current.agentId }}
        </n-descriptions-item>
        <n-descriptions-item label="节点名称">
          {{ emptyText(current.name) }}
        </n-descriptions-item>
        <n-descriptions-item label="主机名">
          {{ emptyText(current.hostname) }}
        </n-descriptions-item>
        <n-descriptions-item label="Agent IP" :span="2">
          {{ formatAgentIp(current.agentIp) }}
        </n-descriptions-item>
        <n-descriptions-item label="在线状态">
          <n-tag :type="current.online ? 'success' : 'default'" size="small">
            {{ dict.getLabel('DataAgentOnlineStatusOptions', current.onlineStatus) }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="注册状态">
          <n-tag size="small">
            {{ dict.getLabel('DataAgentRegisterStatusOptions', current.registerStatus) }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="调度状态">
          <n-tag :type="current.dispatchAllowed ? 'success' : 'warning'" size="small">
            {{
              current.dispatchAllowed
                ? '可调度'
                : dict.getLabel('DataAgentDispatchStatusOptions', current.dispatchStatus)
            }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="版本">
          {{ emptyText(current.version) }}
        </n-descriptions-item>
        <n-descriptions-item label="最近心跳">
          {{ emptyText(current.lastSeenAt) }}
        </n-descriptions-item>
        <n-descriptions-item label="注册时间">
          {{ emptyText(current.createdAt) }}
        </n-descriptions-item>
        <n-descriptions-item label="批准时间">
          {{ emptyText(current.approvedAt) }}
        </n-descriptions-item>
        <n-descriptions-item label="禁用时间">
          {{ emptyText(current.disabledAt) }}
        </n-descriptions-item>
        <n-descriptions-item label="备注" :span="2">
          {{ emptyText(current.remark) }}
        </n-descriptions-item>
      </n-descriptions>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { useDictStore } from '@/store/modules/dict';
  import { Approve, Dispatch, List, Reject } from '@/api/dataAgent';
  import { ReloadOutlined } from '@vicons/antd';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { columns, loadOptions, schemas, State } from './model';

  const dict = useDictStore();
  const dialog = useDialog();
  const message = useMessage();
  const actionRef = ref();
  const searchFormRef = ref<any>({});
  const detailVisible = ref(false);
  const refreshing = ref(false);
  const current = ref<State | null>(null);

  const actionColumn = reactive({
    width: 300,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: State) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '详情',
            onClick: showDetail.bind(null, record),
            auth: ['/dataAgent/list'],
          },
          {
            label: '批准',
            onClick: handleApprove.bind(null, record),
            ifShow: () => record.registerStatus !== 'approved',
            auth: ['/dataAgent/approve'],
          },
          {
            label: '开启调度',
            onClick: handleDispatch.bind(null, record, 'enabled'),
            ifShow: () =>
              record.registerStatus === 'approved' && record.dispatchStatus !== 'enabled',
            auth: ['/dataAgent/dispatch'],
          },
          {
            label: '禁用调度',
            onClick: handleDispatch.bind(null, record, 'disabled'),
            ifShow: () =>
              record.registerStatus === 'approved' && record.dispatchStatus === 'enabled',
            auth: ['/dataAgent/dispatch'],
          },
          {
            label: '拒绝',
            onClick: handleReject.bind(null, record, 'rejected'),
            ifShow: () => record.registerStatus === 'pending',
            auth: ['/dataAgent/reject'],
          },
          {
            label: '吊销',
            onClick: handleReject.bind(null, record, 'revoked'),
            ifShow: () => record.registerStatus !== 'revoked',
            auth: ['/dataAgent/reject'],
          },
        ],
      });
    },
  });

  const scrollX = computed(() => adaTableScrollX(columns, actionColumn.width));

  const [register, {}] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 90,
    schemas,
  });

  const loadDataTable = async (res) => {
    return await List({ ...searchFormRef.value?.formModel, ...res });
  };

  function reloadTable() {
    refreshing.value = true;
    actionRef.value?.reload();
    setTimeout(() => {
      refreshing.value = false;
    }, 300);
  }

  function showDetail(record: State) {
    current.value = record;
    detailVisible.value = true;
  }

  function emptyText(value?: string | number | null) {
    return value === undefined || value === null || value === '' ? '-' : value;
  }

  function formatAgentIp(items?: string[]) {
    return items && items.length > 0 ? items.join(', ') : '-';
  }

  function handleApprove(record: State) {
    dialog.warning({
      title: '批准节点',
      content: `批准 ${record.agentId} 后，该节点可接收平台调度任务。`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        await Approve({ agentId: record.agentId });
        message.success('已批准节点');
        reloadTable();
      },
    });
  }

  function handleDispatch(record: State, dispatchStatus: 'enabled' | 'disabled') {
    dialog.warning({
      title: dispatchStatus === 'enabled' ? '开启调度' : '禁用调度',
      content:
        dispatchStatus === 'enabled'
          ? `开启 ${record.agentId} 调度后，符合条件的任务可以下发到该节点。`
          : `禁用 ${record.agentId} 调度后，新的任务不会再下发到该节点。`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        await Dispatch({ agentId: record.agentId, dispatchStatus });
        message.success(dispatchStatus === 'enabled' ? '已开启调度' : '已禁用调度');
        reloadTable();
      },
    });
  }

  function handleReject(record: State, registerStatus: 'rejected' | 'revoked') {
    dialog.warning({
      title: registerStatus === 'rejected' ? '拒绝节点' : '吊销节点',
      content:
        registerStatus === 'rejected'
          ? `拒绝 ${record.agentId} 后，该节点不能被调度。`
          : `吊销 ${record.agentId} 后，该节点不能再接收调度任务。`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        await Reject({ agentId: record.agentId, registerStatus });
        message.success(registerStatus === 'rejected' ? '已拒绝节点' : '已吊销节点');
        reloadTable();
      },
    });
  }

  onMounted(() => {
    loadOptions();
  });
</script>

<style scoped>
  :deep(.agent-ip-text) {
    display: inline-block;
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    vertical-align: middle;
    white-space: nowrap;
  }

  .agent-detail-modal {
    width: 760px;
  }
</style>
