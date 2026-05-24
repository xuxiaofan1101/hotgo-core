import { http } from '@/utils/http/axios';

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
