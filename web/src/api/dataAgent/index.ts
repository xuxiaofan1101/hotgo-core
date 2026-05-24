import { http } from '@/utils/http/axios';

export function List(params) {
  return http.request({
    url: '/dataAgent/list',
    method: 'get',
    params,
  });
}

export function Approve(params) {
  return http.request({
    url: '/dataAgent/approve',
    method: 'POST',
    params,
  });
}

export function Dispatch(params) {
  return http.request({
    url: '/dataAgent/dispatch',
    method: 'POST',
    params,
  });
}

export function Reject(params) {
  return http.request({
    url: '/dataAgent/reject',
    method: 'POST',
    params,
  });
}
