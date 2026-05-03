import request from '@/utils/request'

export function createOrder(data) {
  return request({
    url: '/order/create',
    method: 'post',
    data
  })
}

export function getOrderList(params) {
  return request({
    url: '/order/list',
    method: 'get',
    params
  })
}

export function getOrderDetail(id) {
  return request({
    url: `/order/detail/${id}`,
    method: 'get'
  })
}

export function cancelOrder(id) {
  return request({
    url: `/order/cancel/${id}`,
    method: 'put'
  })
}

export function payOrder(id) {
  return request({
    url: `/order/pay/${id}`,
    method: 'put'
  })
}

export function receiveOrder(id) {
  return request({
    url: `/order/receive/${id}`,
    method: 'put'
  })
}
