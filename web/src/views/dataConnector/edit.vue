<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      title="数据源配置"
      :style="{ width: dialogWidth }"
    >
      <n-scrollbar class="pr-5 data-connector-scrollbar">
        <n-spin :show="loading" description="请稍候...">
          <n-form
            ref="formRef"
            :model="formValue"
            :rules="rules"
            :label-placement="settingStore.isMobile ? 'top' : 'left'"
            :label-width="120"
            class="py-4"
          >
            <n-grid cols="1" responsive="screen">
              <n-gi>
                <n-form-item label="用途" path="direction">
                  <n-radio-group
                    v-model:value="formValue.direction"
                    name="direction"
                    :disabled="formValue.id > 0"
                    @update:value="handleDirectionChange"
                  >
                    <n-radio-button
                      v-for="item in directionOptions"
                      :key="item.value"
                      :value="item.value"
                    >
                      {{ item.label }}
                    </n-radio-button>
                  </n-radio-group>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="名称" path="name">
                  <n-input placeholder="请输入名称" v-model:value="formValue.name" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="编码" path="code">
                  <n-input placeholder="请输入编码" v-model:value="formValue.code" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="类型" path="connectorType">
                  <n-select
                    v-model:value="formValue.connectorType"
                    :options="formTypeOptions"
                    placeholder="请选择类型"
                    @update:value="handleTypeChange"
                  />
                </n-form-item>
              </n-gi>

              <template v-if="isSource && formValue.connectorType === 'http'">
                <n-gi>
                  <n-form-item label="接入 Token">
                    <n-input v-model:value="configForm.token" placeholder="可选，留空表示不校验" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="说明">
                    <n-input
                      v-model:value="configForm.description"
                      placeholder="POST /api/v1/data-integration/ingest"
                    />
                  </n-form-item>
                </n-gi>
              </template>

              <template v-else-if="formValue.connectorType === 'kafka'">
                <n-gi>
                  <n-form-item label="Brokers" required>
                    <n-input
                      v-model:value="configForm.brokers"
                      placeholder="127.0.0.1:9092,127.0.0.1:9093"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="认证模式">
                    <n-select
                      v-model:value="configForm.authMode"
                      :options="kafkaAuthOptions"
                      placeholder="请选择认证模式"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi v-if="kafkaAuthNeedsPassword">
                  <n-form-item label="用户名">
                    <n-input v-model:value="configForm.username" />
                  </n-form-item>
                </n-gi>
                <n-gi v-if="kafkaAuthNeedsPassword">
                  <n-form-item label="密码">
                    <n-input
                      v-model:value="configForm.password"
                      type="password"
                      show-password-on="click"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi v-if="kafkaAuthNeedsToken">
                  <n-form-item label="Token">
                    <n-input
                      v-model:value="configForm.token"
                      type="password"
                      show-password-on="click"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="TLS">
                    <n-switch v-model:value="configForm.tls" />
                  </n-form-item>
                </n-gi>
                <template v-if="configForm.tls">
                  <n-gi>
                    <n-form-item label="跳过证书校验">
                      <n-switch v-model:value="configForm.insecureSkipVerify" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="CA 证书">
                      <n-input
                        v-model:value="configForm.caCert"
                        type="textarea"
                        :autosize="{ minRows: 2, maxRows: 5 }"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="mTLS 客户端证书">
                      <n-input
                        v-model:value="configForm.clientCert"
                        type="textarea"
                        :autosize="{ minRows: 2, maxRows: 5 }"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="mTLS 客户端私钥">
                      <n-input
                        v-model:value="configForm.clientKey"
                        type="textarea"
                        :autosize="{ minRows: 2, maxRows: 5 }"
                      />
                    </n-form-item>
                  </n-gi>
                </template>
                <n-gi v-if="!isSource">
                  <n-form-item label="超时秒数">
                    <n-input-number
                      v-model:value="configForm.timeoutSeconds"
                      class="data-connector-full"
                      :min="1"
                      :show-button="false"
                    />
                  </n-form-item>
                </n-gi>
              </template>

              <template v-else-if="formValue.connectorType === 's3'">
                <n-gi>
                  <n-form-item label="Endpoint">
                    <n-input
                      v-model:value="configForm.endpoint"
                      placeholder="可选，S3 兼容服务再填写"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="Region" required>
                    <n-input v-model:value="configForm.region" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="凭证方式">
                    <n-select
                      v-model:value="configForm.credentialMode"
                      :options="credentialModeOptions"
                    />
                  </n-form-item>
                </n-gi>
                <template v-if="configForm.credentialMode === 'static'">
                  <n-gi>
                    <n-form-item label="AccessKey">
                      <n-input v-model:value="configForm.accessKey" />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="SecretKey">
                      <n-input
                        v-model:value="configForm.secretKey"
                        type="password"
                        show-password-on="click"
                      />
                    </n-form-item>
                  </n-gi>
                  <n-gi>
                    <n-form-item label="SessionToken">
                      <n-input
                        v-model:value="configForm.sessionToken"
                        type="password"
                        show-password-on="click"
                        placeholder="临时凭证可选"
                      />
                    </n-form-item>
                  </n-gi>
                </template>
              </template>

              <template v-else-if="isSource && formValue.connectorType === 'log'">
                <n-gi>
                  <n-form-item label="日志目录" required>
                    <n-input v-model:value="configForm.basePath" />
                  </n-form-item>
                </n-gi>
              </template>

              <template v-else-if="!isSource && formValue.connectorType === 'http'">
                <n-gi>
                  <n-form-item label="Base URL" required>
                    <n-input v-model:value="configForm.baseUrl" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="Authorization">
                    <n-input v-model:value="configForm.authHeader" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="超时秒数">
                    <n-input-number
                      v-model:value="configForm.timeoutSeconds"
                      class="data-connector-full"
                      :min="1"
                      :show-button="false"
                    />
                  </n-form-item>
                </n-gi>
              </template>

              <template
                v-else-if="
                  !isSource &&
                  (formValue.connectorType === 'feishu' || formValue.connectorType === 'lark')
                "
              >
                <n-gi>
                  <n-form-item label="Webhook" required>
                    <n-input v-model:value="configForm.webhook" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="签名密钥">
                    <n-input
                      v-model:value="configForm.secret"
                      type="password"
                      show-password-on="click"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="消息模板">
                    <n-input
                      v-model:value="configForm.messageTemplate"
                      type="textarea"
                      :autosize="{ minRows: 2, maxRows: 4 }"
                    />
                  </n-form-item>
                </n-gi>
              </template>

              <template v-else-if="!isSource && formValue.connectorType === 'log'">
                <n-gi>
                  <n-form-item label="日志等级">
                    <n-select v-model:value="configForm.level" :options="logLevelOptions" />
                  </n-form-item>
                </n-gi>
              </template>

              <template
                v-else-if="
                  !isSource &&
                  (formValue.connectorType === 'elasticsearch' ||
                    formValue.connectorType === 'opensearch')
                "
              >
                <n-gi>
                  <n-form-item label="Base URL" required>
                    <n-input v-model:value="configForm.baseUrl" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="Authorization">
                    <n-input v-model:value="configForm.authHeader" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="超时秒数">
                    <n-input-number
                      v-model:value="configForm.timeoutSeconds"
                      class="data-connector-full"
                      :min="1"
                      :show-button="false"
                    />
                  </n-form-item>
                </n-gi>
              </template>

              <template v-else-if="!isSource && formValue.connectorType === 'splunk'">
                <n-gi>
                  <n-form-item label="HEC URL" required>
                    <n-input v-model:value="configForm.hecUrl" />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="HEC Token" required>
                    <n-input
                      v-model:value="configForm.token"
                      type="password"
                      show-password-on="click"
                    />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="超时秒数">
                    <n-input-number
                      v-model:value="configForm.timeoutSeconds"
                      class="data-connector-full"
                      :min="1"
                      :show-button="false"
                    />
                  </n-form-item>
                </n-gi>
              </template>

              <template v-else>
                <n-gi>
                  <n-form-item label="说明">
                    <n-input v-model:value="configForm.description" />
                  </n-form-item>
                </n-gi>
              </template>

              <n-gi v-if="!isSource">
                <n-form-item label="测试消息">
                  <n-input v-model:value="testMessage" />
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
              <n-gi>
                <n-form-item label="备注" path="remark">
                  <n-input
                    type="textarea"
                    placeholder="备注"
                    v-model:value="formValue.remark"
                    :autosize="{ minRows: 2, maxRows: 4 }"
                  />
                </n-form-item>
              </n-gi>
            </n-grid>
          </n-form>
        </n-spin>
      </n-scrollbar>
      <template #action>
        <n-space>
          <n-button :loading="testBtnLoading" @click="handleTest"> 测试配置 </n-button>
          <n-button @click="closeForm"> 取消 </n-button>
          <n-button type="info" :loading="formBtnLoading" @click="confirmForm"> 保存 </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, reactive, ref } from 'vue';
  import { Edit, Test, View } from '@/api/dataConnector';
  import { State, newState, rules } from './model';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { useMessage } from 'naive-ui';
  import { adaModalWidth } from '@/utils/hotgo';

  type ConnectorDirection = 'source' | 'sink';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const settingStore = useProjectSettingStore();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref<State>(newState(null));
  const formRef = ref<any>({});
  const formBtnLoading = ref(false);
  const testBtnLoading = ref(false);
  const testMessage = ref('HotGo 策略引擎测试告警');
  const dialogWidth = computed(() => {
    return adaModalWidth(920);
  });

  const directionOptions = [
    { label: '输入', value: 'source' },
    { label: '输出', value: 'sink' },
  ];

  const sourceTypeOptions = [
    { label: 'HTTP', value: 'http' },
    { label: 'Kafka', value: 'kafka' },
    { label: 'S3', value: 's3' },
    { label: '日志', value: 'log' },
    { label: '手动', value: 'manual' },
  ];

  const sinkTypeOptions = [
    { label: 'HTTP Webhook', value: 'http' },
    { label: 'Kafka', value: 'kafka' },
    { label: 'Feishu/Lark 飞书群机器人', value: 'feishu' },
    { label: 'Lark', value: 'lark' },
    { label: 'S3', value: 's3' },
    { label: 'Elasticsearch', value: 'elasticsearch' },
    { label: 'OpenSearch', value: 'opensearch' },
    { label: 'Splunk HEC', value: 'splunk' },
    { label: '日志', value: 'log' },
  ];

  const kafkaAuthOptions = [
    { label: '无认证', value: 'none' },
    { label: 'SASL/PLAIN', value: 'plain' },
    { label: 'SCRAM-SHA-256', value: 'scram-sha-256' },
    { label: 'SCRAM-SHA-512', value: 'scram-sha-512' },
    { label: 'OAuth Bearer', value: 'oauthbearer' },
  ];

  const credentialModeOptions = [
    { label: '云角色', value: 'role' },
    { label: '静态 AK/SK', value: 'static' },
    { label: '匿名访问', value: 'anonymous' },
  ];

  const logLevelOptions = [
    { label: 'debug', value: 'debug' },
    { label: 'info', value: 'info' },
    { label: 'warning', value: 'warning' },
    { label: 'error', value: 'error' },
  ];

  const configForm = reactive({
    accessKey: '',
    authMode: 'none',
    authHeader: '',
    basePath: './logs',
    baseUrl: 'https://example.com',
    brokers: '127.0.0.1:9092',
    caCert: '',
    clientCert: '',
    clientKey: '',
    description: '',
    endpoint: '',
    hecUrl: 'https://splunk.example.com:8088/services/collector/event',
    insecureSkipVerify: false,
    level: 'warning',
    messageTemplate: '策略告警：{{message}}',
    region: 'us-east-1',
    secret: '',
    secretKey: '',
    credentialMode: 'role',
    sessionToken: '',
    timeoutSeconds: 5,
    tls: false,
    token: '',
    username: '',
    password: '',
    webhook: 'https://open.feishu.cn/open-apis/bot/v2/hook/xxxx',
  });

  const isSource = computed(() => formValue.value.direction === 'source');
  const formTypeOptions = computed(() => {
    return isSource.value ? sourceTypeOptions : sinkTypeOptions;
  });
  const kafkaAuthNeedsPassword = computed(() => {
    return (
      formValue.value.connectorType === 'kafka' &&
      !['none', 'oauthbearer'].includes(configForm.authMode)
    );
  });
  const kafkaAuthNeedsToken = computed(() => {
    return formValue.value.connectorType === 'kafka' && configForm.authMode === 'oauthbearer';
  });

  function splitList(value: string) {
    return value
      .split(',')
      .map((item) => item.trim())
      .filter(Boolean);
  }

  function normalizeConfig(config: unknown): Record<string, any> {
    if (!config) {
      return {};
    }
    if (typeof config === 'string') {
      try {
        return JSON.parse(config);
      } catch {
        return {};
      }
    }
    if (typeof config === 'object') {
      return config as Record<string, any>;
    }
    return {};
  }

  function kafkaConfigValue(config: Record<string, any>, key: string, fallback = '') {
    return config[key] ?? config.sasl?.[key] ?? fallback;
  }

  function normalizeKafkaAuthMode(config: Record<string, any>) {
    const mode = String(
      kafkaConfigValue(config, 'authMode', kafkaConfigValue(config, 'mechanism', 'none'))
    )
      .toLowerCase()
      .replace(/_/g, '-');
    if (mode === 'sasl/plain' || mode === 'sasl-plain') return 'plain';
    if (mode === 'scram-sha256') return 'scram-sha-256';
    if (mode === 'scram-sha512') return 'scram-sha-512';
    if (mode === 'oauth' || mode === 'oauth-bearer' || mode === 'sasl/oauthbearer') {
      return 'oauthbearer';
    }
    return mode || 'none';
  }

  function buildKafkaConfig(includeTimeout = false) {
    const config: Record<string, any> = {
      authMode: configForm.authMode,
      brokers: splitList(configForm.brokers),
      tls: configForm.tls,
    };
    if (includeTimeout) {
      config.timeoutSeconds = configForm.timeoutSeconds;
    }
    if (configForm.authMode !== 'none') {
      if (configForm.authMode === 'oauthbearer') {
        config.token = configForm.token;
      } else {
        config.username = configForm.username;
        config.password = configForm.password;
      }
    }
    if (configForm.tls) {
      config.insecureSkipVerify = configForm.insecureSkipVerify;
      config.caCert = configForm.caCert;
      config.clientCert = configForm.clientCert;
      config.clientKey = configForm.clientKey;
    }
    return config;
  }

  function resetConfig(direction: ConnectorDirection, type: string, rawConfig: unknown = {}) {
    const config = normalizeConfig(rawConfig);
    Object.assign(configForm, {
      accessKey: config.accessKey ?? '',
      authMode: normalizeKafkaAuthMode(config),
      authHeader: config.headers?.Authorization ?? '',
      basePath: config.basePath ?? './logs',
      baseUrl: config.baseUrl ?? 'https://example.com',
      brokers: Array.isArray(config.brokers) ? config.brokers.join(',') : '127.0.0.1:9092',
      caCert: config.caCert ?? '',
      clientCert: config.clientCert ?? '',
      clientKey: config.clientKey ?? '',
      description: config.description ?? '',
      endpoint: config.endpoint ?? '',
      hecUrl:
        config.hecUrl ?? config.url ?? 'https://splunk.example.com:8088/services/collector/event',
      insecureSkipVerify: Boolean(config.insecureSkipVerify),
      level: config.level ?? 'warning',
      messageTemplate:
        config.messageTemplate ?? (direction === 'sink' ? '策略告警：{{message}}' : ''),
      password: kafkaConfigValue(config, 'password'),
      region: config.region ?? 'us-east-1',
      secret: config.secret ?? '',
      secretKey: config.secretKey ?? '',
      credentialMode:
        config.credentialMode ?? (config.accessKey || config.secretKey ? 'static' : 'role'),
      sessionToken: config.sessionToken ?? '',
      timeoutSeconds: Number(config.timeoutSeconds) || 5,
      tls: Boolean(config.tls),
      token: config.token ?? '',
      username: kafkaConfigValue(config, 'username'),
      webhook:
        config.webhook ??
        (type === 'lark'
          ? 'https://open.larksuite.com/open-apis/bot/v2/hook/xxxx'
          : 'https://open.feishu.cn/open-apis/bot/v2/hook/xxxx'),
    });
  }

  function buildConfig() {
    if (formValue.value.direction === 'source') {
      if (formValue.value.connectorType === 'http') {
        return { description: configForm.description, token: configForm.token };
      }
      if (formValue.value.connectorType === 'kafka') {
        return buildKafkaConfig();
      }
      if (formValue.value.connectorType === 's3') {
        const config: Record<string, any> = {
          credentialMode: configForm.credentialMode,
          endpoint: configForm.endpoint,
          region: configForm.region,
        };
        if (configForm.credentialMode === 'static') {
          config.accessKey = configForm.accessKey;
          config.secretKey = configForm.secretKey;
          config.sessionToken = configForm.sessionToken;
        }
        return config;
      }
      if (formValue.value.connectorType === 'log') {
        return { basePath: configForm.basePath };
      }
      return { description: configForm.description || '后台手工触发' };
    }

    if (formValue.value.connectorType === 'http') {
      return {
        baseUrl: configForm.baseUrl,
        headers: { Authorization: configForm.authHeader },
        timeoutSeconds: configForm.timeoutSeconds,
      };
    }
    if (formValue.value.connectorType === 'kafka') {
      return buildKafkaConfig(true);
    }
    if (formValue.value.connectorType === 'feishu' || formValue.value.connectorType === 'lark') {
      return {
        messageTemplate: configForm.messageTemplate,
        secret: configForm.secret,
        webhook: configForm.webhook,
      };
    }
    if (formValue.value.connectorType === 's3') {
      const config: Record<string, any> = {
        credentialMode: configForm.credentialMode,
        endpoint: configForm.endpoint,
        region: configForm.region,
      };
      if (configForm.credentialMode === 'static') {
        config.accessKey = configForm.accessKey;
        config.secretKey = configForm.secretKey;
        config.sessionToken = configForm.sessionToken;
      }
      return config;
    }
    if (
      formValue.value.connectorType === 'elasticsearch' ||
      formValue.value.connectorType === 'opensearch'
    ) {
      return {
        baseUrl: configForm.baseUrl,
        headers: { Authorization: configForm.authHeader },
        timeoutSeconds: configForm.timeoutSeconds,
      };
    }
    if (formValue.value.connectorType === 'splunk') {
      return {
        hecUrl: configForm.hecUrl,
        token: configForm.token,
        timeoutSeconds: configForm.timeoutSeconds,
      };
    }
    return { level: configForm.level };
  }

  function validateConfig() {
    if (formValue.value.connectorType === 'kafka' && splitList(configForm.brokers).length === 0) {
      message.error('请填写 Kafka Brokers');
      return false;
    }
    if (formValue.value.connectorType === 'kafka' && kafkaAuthNeedsPassword.value) {
      if (!configForm.username || !configForm.password) {
        message.error('请填写 Kafka 用户名和密码');
        return false;
      }
    }
    if (
      formValue.value.connectorType === 'kafka' &&
      kafkaAuthNeedsToken.value &&
      !configForm.token
    ) {
      message.error('请填写 Kafka Token');
      return false;
    }
    if (formValue.value.connectorType === 's3' && !configForm.region) {
      message.error('请填写 S3 Region');
      return false;
    }
    if (formValue.value.connectorType === 's3' && configForm.credentialMode === 'static') {
      if (!configForm.accessKey || !configForm.secretKey) {
        message.error('请填写 S3 AccessKey 和 SecretKey');
        return false;
      }
    }
    if (isSource.value && formValue.value.connectorType === 'log' && !configForm.basePath) {
      message.error('请填写日志目录');
      return false;
    }
    if (
      !isSource.value &&
      ['http', 'elasticsearch', 'opensearch'].includes(formValue.value.connectorType)
    ) {
      if (!configForm.baseUrl) {
        message.error('请填写 Base URL');
        return false;
      }
    }
    if (
      !isSource.value &&
      ['feishu', 'lark'].includes(formValue.value.connectorType) &&
      !configForm.webhook
    ) {
      message.error('请填写 Webhook');
      return false;
    }
    if (!isSource.value && formValue.value.connectorType === 'splunk') {
      if (!configForm.hecUrl || !configForm.token) {
        message.error('请填写 Splunk HEC URL 和 HEC Token');
        return false;
      }
    }
    return true;
  }

  function handleDirectionChange(value: ConnectorDirection) {
    formValue.value.direction = value;
    formValue.value.connectorType = value === 'source' ? 'http' : 'feishu';
    resetConfig(value, formValue.value.connectorType);
  }

  function handleTypeChange(value: string) {
    formValue.value.connectorType = value;
    resetConfig(formValue.value.direction as ConnectorDirection, value);
  }

  function resetForm() {
    formValue.value = newState(null);
    formValue.value.direction = 'source';
    formValue.value.connectorType = 'http';
    formValue.value.status = 1;
    resetConfig('source', 'http');
    testMessage.value = 'HotGo 策略引擎测试告警';
  }

  function submitEdit(config: Record<string, any>) {
    formBtnLoading.value = true;
    Edit({ ...formValue.value, config })
      .then((_res) => {
        message.success('操作成功');
        closeForm();
        emit('reloadTable');
      })
      .finally(() => {
        formBtnLoading.value = false;
      });
  }

  function confirmForm(e) {
    e.preventDefault();
    formRef.value.validate((errors) => {
      if (errors) {
        message.error('请填写完整信息');
        return;
      }
      if (!validateConfig()) {
        return;
      }
      submitEdit(buildConfig());
    });
  }

  function handleTest(e) {
    e.preventDefault();
    if (!validateConfig()) {
      return;
    }
    testBtnLoading.value = true;
    Test({
      config: buildConfig(),
      connectorType: formValue.value.connectorType,
      direction: formValue.value.direction,
      message: testMessage.value,
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

  function closeForm() {
    showModal.value = false;
    loading.value = false;
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
        resetConfig(
          formValue.value.direction as ConnectorDirection,
          formValue.value.connectorType,
          formValue.value.config
        );
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
  .data-connector-scrollbar {
    max-height: 87vh;
  }

  .data-connector-full {
    width: 100%;
  }
</style>
