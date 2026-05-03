import request from '@/utils/request'

export function addComment(data) {
  return request({
    url: '/comment/add',
    method: 'post',
    data
  })
}

export function getCommentList(params) {
  return request({
    url: '/comment/list',
    method: 'get',
    params
  })
}
