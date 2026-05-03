import request from '@/utils/request'

export function getProductList(params) {
  return request({
    url: '/product/list',
    method: 'get',
    params
  })
}

export function getProductDetail(id) {
  return request({
    url: `/product/detail/${id}`,
    method: 'get'
  })
}

export function getRecommendProducts(limit = 10) {
  return request({
    url: '/product/recommend',
    method: 'get',
    params: { limit }
  })
}

export function getCategoryList() {
  return request({
    url: '/category/list',
    method: 'get'
  })
}

export function getBannerList(type = 1) {
  return request({
    url: '/banner/list',
    method: 'get',
    params: { type }
  })
}
