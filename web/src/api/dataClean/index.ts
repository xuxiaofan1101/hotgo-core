import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/dataClean/list',
    method: 'get',
    params,
  });
}

export function Delete(params) {
  return http.request({
    url: '/dataClean/delete',
    method: 'POST',
    params,
  });
}

export function Edit(params) {
  return http.request({
    url: '/dataClean/edit',
    method: 'POST',
    params,
  });
}

export function Status(params) {
  return http.request({
    url: '/dataClean/status',
    method: 'POST',
    params,
  });
}

export function View(params) {
  return http.request({
    url: '/dataClean/view',
    method: 'GET',
    params,
  });
}

export function Test(params) {
  return http.request({
    url: '/dataClean/test',
    method: 'POST',
    params,
  });
}

export function Sample(params) {
  return http.request({
    url: '/dataClean/sample',
    method: 'POST',
    params,
  });
}

export function FieldList(params) {
  return http.request({
    url: '/dataClean/fieldList',
    method: 'GET',
    params,
  });
}
