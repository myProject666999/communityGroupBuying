import request from '@/utils/request'

export function getAddressList() {
  return request({
    url: '/address/list',
    method: 'get'
  })
}

export function getAddressDetail(id) {
  return request({
    url: `/address/detail/${id}`,
    method: 'get'
  })
}

export function addAddress(data) {
  return request({
    url: '/address/add',
    method: 'post',
    data
  })
}

export function updateAddress(id, data) {
  return request({
    url: `/address/update/${id}`,
    method: 'put',
    data
  })
}

export function deleteAddress(id) {
  return request({
    url: `/address/delete/${id}`,
    method: 'delete'
  })
}

export function setDefaultAddress(id) {
  return request({
    url: `/address/default/${id}`,
    method: 'put'
  })
}
