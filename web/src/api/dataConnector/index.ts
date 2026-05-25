import { http } from '@/utils/http/axios';

function formatUnexpectedResponse(response, requestUrl = '') {
  const text = typeof response === 'string' ? response.trim() : JSON.stringify(response || '');
  const suffix = requestUrl ? `，实际请求地址：${requestUrl}` : '';
  if (!text) {
    return `接口无返回内容${suffix}`;
  }
  if (text.includes('<!DOCTYPE') || text.includes('<html')) {
    return `接口返回了前端页面，请检查请求前缀 /admin 或后端路由是否已注册${suffix}`;
  }
  if (text.includes('Not Found')) {
    return `后端接口未找到，请确认 /dataConnector/kafkaTopics 路由已注册${suffix}`;
  }
  return `接口返回格式异常：${text.slice(0, 120)}${suffix}`;
}

function getDevAdminProxyUrl(path: string) {
  if (!import.meta.env.DEV) {
    return '';
  }
  const proxyConfig = import.meta.env.VITE_PROXY as unknown;
  const proxyList =
    typeof proxyConfig === 'string'
      ? JSON.parse(proxyConfig || '[]')
      : Array.isArray(proxyConfig)
        ? proxyConfig
        : [];
  const adminProxy = proxyList.find((item) => Array.isArray(item) && item[0] === '/admin');
  if (!adminProxy) {
    return '';
  }
  const target = String(adminProxy[1] || '').replace(/\/$/, '');
  if (!target) {
    return '';
  }
  return `${target}${path.replace(/^\/admin/, '')}`;
}

function requestKafkaTopics(url: string, params) {
  return http
    .request(
      {
        url,
        method: 'GET',
        params,
      },
      {
        isTransformResponse: false,
        isReturnNativeResponse: true,
        ignoreCancelToken: true,
        joinPrefix: false,
      }
    )
    .then((response) => {
      const data = response?.data;
      const requestUrl = response?.request?.responseURL || response?.config?.url || '';
      if (!data || typeof data !== 'object') {
        throw new Error(formatUnexpectedResponse(data, requestUrl));
      }
      if (data.code === 0) {
        return data.data || { topics: [] };
      }
      throw new Error(data.message || '读取 Kafka Topic 列表失败');
    });
}

// 获取数据源列表
export function List(params) {
  return http.request({
    url: '/dataConnector/list',
    method: 'get',
    params,
  });
}

// 删除/批量删除数据源
export function Delete(params) {
  return http.request({
    url: '/dataConnector/delete',
    method: 'POST',
    params,
  });
}

// 添加/编辑数据源
export function Edit(params) {
  return http.request({
    url: '/dataConnector/edit',
    method: 'POST',
    params,
  });
}

// 修改数据源状态
export function Status(params) {
  return http.request({
    url: '/dataConnector/status',
    method: 'POST',
    params,
  });
}

// 获取数据源指定详情
export function View(params) {
  return http.request({
    url: '/dataConnector/view',
    method: 'GET',
    params,
  });
}

// 测试数据源配置
export function Test(params) {
  return http.request({
    url: '/dataConnector/test',
    method: 'POST',
    params,
  });
}

// 获取Kafka Topic列表
export function KafkaTopics(params) {
  const path = '/admin/dataConnector/kafkaTopics';
  return requestKafkaTopics(path, params).catch((error) => {
    if (!String(error?.message || '').includes('接口返回了前端页面')) {
      throw error;
    }
    const fallbackUrl = getDevAdminProxyUrl(path);
    if (!fallbackUrl) {
      throw error;
    }
    return requestKafkaTopics(fallbackUrl, params);
  });
}
