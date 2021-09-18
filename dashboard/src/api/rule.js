import request from '@/utils/request'

export default {
  add(data) {
    return request({
      url: `/api/v1/rule`,
      method: 'post',
      data: data
    })
  },
  deleteById(id) {
    return request({
      url: `/api/v1/rule/${id}`,
      method: 'delete'
    })
  },
  updateById(id, data) {
    return request({
      url: `/api/v1/rule/${id}`,
      method: 'put',
      data
    })
  },
  getById(id) {
    return request({
      url: `/api/v1/rule/${id}`,
      method: 'get'
    })
  },
  getList(query, current = 1, size = 20) {
    return request({
      url: '/api/v1/rule',
      method: 'get',
      params: { ...query, page: current, size: size }
    })
  }
}
