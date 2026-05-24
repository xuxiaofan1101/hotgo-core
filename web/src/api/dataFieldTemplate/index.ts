import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/dataFieldTemplate/list',
    method: 'get',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/dataFieldTemplate/delete',
    method: 'POST',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/dataFieldTemplate/edit',
    method: 'POST',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/dataFieldTemplate/status',
    method: 'POST',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/dataFieldTemplate/view',
    method: 'GET',
    params,
  });
}
