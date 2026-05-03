import request from '@/utils/request'

export function getFavoriteList(params) {
  return request({
    url: '/favorite/list',
    method: 'get',
    params
  })
}

export function addFavorite(data) {
  return request({
    url: '/favorite/add',
    method: 'post',
    data
  })
}

export function deleteFavorite(id) {
  return request({
    url: `/favorite/delete/${id}`,
    method: 'delete'
  })
}
