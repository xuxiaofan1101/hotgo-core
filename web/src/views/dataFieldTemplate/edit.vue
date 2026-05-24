<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="formValue.id > 0 ? '编辑字段模板 #' + formValue.id : '新增字段模板'"
      :style="{ width: dialogWidth }"
    >
      <n-scrollbar class="pr-5 data-template-scrollbar">
        <n-spin :show="loading" description="请稍候...">
          <n-form
            ref="formRef"
            :model="formValue"
            :rules="rules"
            :label-placement="settingStore.isMobile ? 'top' : 'left'"
            :label-width="110"
            class="py-4"
          >
            <n-grid cols="1 s:1 m:2 l:2 xl:2 2xl:2" responsive="screen" :x-gap="16">
              <n-gi>
                <n-form-item label="模板名称" path="name">
                  <n-input v-model:value="formValue.name" placeholder="请输入模板名称" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="模板编码" path="code">
                  <n-input v-model:value="formValue.code" placeholder="请输入模板编码" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="事件类型" path="eventType">
                  <n-input v-model:value="formValue.eventType" placeholder="如 nginx.access" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="状态" path="status">
                  <n-switch
                    v-model:value="formValue.status"
                    :checked-value="1"
                    :unchecked-value="2"
                  />
                </n-form-item>
              </n-gi>
            </n-grid>

            <div class="data-template-section">
              <div class="data-template-toolbar">
                <div class="data-template-title">标准字段</div>
                <n-button size="small" type="primary" @click="addTemplateField()">
                  <template #icon>
                    <n-icon><PlusOutlined /></n-icon>
                  </template>
                  添加字段
                </n-button>
              </div>
              <n-empty v-if="templateFields.length === 0" description="暂无标准字段" />
              <div
                v-for="(field, index) in templateFields"
                v-else
                :key="field.key"
                class="data-template-field"
              >
                <n-grid cols="1 s:1 m:2 l:4 xl:4 2xl:4" responsive="screen" :x-gap="12">
                  <n-gi>
                    <n-form-item label="字段路径">
                      <n-input v-model:value="field.fieldPath" placeholder="message.level" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="字段名称">
                      <n-input v-model:value="field.fieldName" placeholder="等级" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="字段类型">
                      <n-select v-model:value="field.fieldType" :options="fieldTypeOptions" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="必填">
                      <n-space align="center">
                        <n-switch v-model:value="field.required" />
                        <n-button size="small" type="error" @click="removeTemplateField(index)">
                          <template #icon>
                            <n-icon><DeleteOutlined /></n-icon>
                          </template>
                        </n-button>
                      </n-space>
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="默认值">
                      <n-input v-model:value="field.defaultValue" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="枚举值">
                      <n-input v-model:value="field.enumText" placeholder="逗号分隔" />
                    </n-form-item>
                  </n-gi>
                  <n-gi span="2">
                    <n-form-item label="说明">
                      <n-input v-model:value="field.description" />
                    </n-form-item>
                  </n-gi>
                </n-grid>
              </div>
            </div>

            <n-form-item label="样例数据">
              <n-input
                v-model:value="exampleText"
                type="textarea"
                placeholder="可选，填写 JSON 数组或对象，作为模板样例"
                :autosize="{ minRows: 4, maxRows: 10 }"
              />
            </n-form-item>
            <n-form-item label="备注" path="remark">
              <n-input
                v-model:value="formValue.remark"
                type="textarea"
                placeholder="备注"
                :autosize="{ minRows: 2, maxRows: 4 }"
              />
            </n-form-item>
          </n-form>
        </n-spin>
      </n-scrollbar>
      <template #action>
        <n-space>
          <n-button @click="closeForm">取消</n-button>
          <n-button type="info" :loading="formBtnLoading" @click="confirmForm">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import { PlusOutlined, DeleteOutlined } from '@vicons/antd';
  import { useMessage } from 'naive-ui';
  import { useDictStore } from '@/store/modules/dict';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { adaModalWidth } from '@/utils/hotgo';
  import { Edit, View } from '@/api/dataFieldTemplate';
  import { State, newState, rules } from './model';

  interface TemplateField {
    key: number;
    fieldPath: string;
    fieldName: string;
    fieldType: string;
    required: boolean;
    defaultValue: string;
    enumText: string;
    description: string;
  }

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const dict = useDictStore();
  const settingStore = useProjectSettingStore();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref<State>(newState(null));
  const formRef = ref<any>({});
  const formBtnLoading = ref(false);
  const templateFields = ref<TemplateField[]>([]);
  const exampleText = ref('');
  let seed = 1;

  const dialogWidth = computed(() => adaModalWidth(980));
  const fieldTypeOptions = computed(() => dict.getOptionUnRef('DataFieldTypeOptions'));

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

  function newTemplateField(raw: Record<string, any> = {}): TemplateField {
    const enums = raw.enums ?? raw.enumValues ?? [];
    return {
      key: seed++,
      fieldPath: raw.fieldPath ?? raw.path ?? raw.name ?? '',
      fieldName: raw.fieldName ?? raw.title ?? '',
      fieldType: raw.fieldType ?? raw.type ?? 'string',
      required: Boolean(raw.required),
      defaultValue: raw.defaultValue ?? '',
      enumText: Array.isArray(enums) ? enums.join(',') : String(enums || ''),
      description: raw.description ?? '',
    };
  }

  function addTemplateField(raw: Record<string, any> = {}) {
    templateFields.value.push(newTemplateField(raw));
  }

  function removeTemplateField(index: number) {
    templateFields.value.splice(index, 1);
  }

  function resetTemplateFields(fields: any) {
    templateFields.value = normalizeArray(fields).map((item) => newTemplateField(item));
    if (templateFields.value.length === 0) {
      addTemplateField();
    }
  }

  function resetExamples(examples: any) {
    if (!examples) {
      exampleText.value = '';
      return;
    }
    if (typeof examples === 'string') {
      exampleText.value = examples;
      return;
    }
    exampleText.value = JSON.stringify(examples, null, 2);
  }

  function buildFields() {
    return templateFields.value
      .filter((field) => field.fieldPath)
      .map((field) => ({
        defaultValue: field.defaultValue,
        description: field.description,
        enums: field.enumText
          .split(',')
          .map((item) => item.trim())
          .filter(Boolean),
        fieldName: field.fieldName,
        fieldPath: field.fieldPath,
        fieldType: field.fieldType,
        required: field.required,
      }));
  }

  function buildExamples() {
    const text = exampleText.value.trim();
    if (!text) {
      return [];
    }
    return JSON.parse(text);
  }

  function confirmForm(e: Event) {
    e.preventDefault();
    formRef.value.validate((errors) => {
      if (errors) {
        message.error('请填写完整信息');
        return;
      }
      const fields = buildFields();
      if (fields.length === 0) {
        message.error('请至少配置一个标准字段');
        return;
      }
      let examples: any;
      try {
        examples = buildExamples();
      } catch {
        message.error('样例数据不是合法 JSON');
        return;
      }
      formBtnLoading.value = true;
      Edit({
        ...formValue.value,
        examples,
        fields,
      })
        .then((_res) => {
          message.success('操作成功');
          closeForm();
          emit('reloadTable');
        })
        .finally(() => {
          formBtnLoading.value = false;
        });
    });
  }

  function closeForm() {
    showModal.value = false;
    loading.value = false;
  }

  function resetForm() {
    formValue.value = newState(null);
    formValue.value.eventType = 'custom.event';
    formValue.value.status = 1;
    resetTemplateFields([]);
    resetExamples([]);
  }

  function openModal(state: State) {
    showModal.value = true;

    if (!state || state.id < 1) {
      resetForm();
      return;
    }

    loading.value = true;
    View({ id: state.id })
      .then((res) => {
        formValue.value = newState(res);
        resetTemplateFields(res?.fields);
        resetExamples(res?.examples);
      })
      .finally(() => {
        loading.value = false;
      });
  }

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped>
  .data-template-scrollbar {
    max-height: 87vh;
  }

  .data-template-section {
    border: 1px solid rgba(128, 128, 128, 0.22);
    border-radius: 6px;
    margin-bottom: 16px;
    padding: 16px;
  }

  .data-template-toolbar {
    align-items: center;
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
  }

  .data-template-title {
    color: var(--text-color-1);
    font-weight: 600;
  }

  .data-template-field {
    border: 1px solid rgba(128, 128, 128, 0.22);
    border-radius: 6px;
    margin-bottom: 12px;
    padding: 12px;
  }

  @media (max-width: 768px) {
    .data-template-toolbar {
      align-items: flex-start;
      flex-direction: column;
      row-gap: 8px;
    }
  }
</style>
