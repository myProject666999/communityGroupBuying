import request from '@/utils/request'

export function getCartList() {
  return request({
    url: '/cart/list',
    method: 'get'
  })
}

export function addCart(data) {
  return request({
    url: '/cart/add',
    method: 'post',
    data
  })
}

export function updateCart(data) {
  return request({
    url: '/cart/update',
    method: 'put',
    data
  })
}

export function deleteCart(id) {
  return request({
    url: `/cart/delete/${id}`,
    method: 'delete'
  })
}

export function clearCart() {
  return request({
    url: '/cart/clear',
    method: 'post'
  })
}
