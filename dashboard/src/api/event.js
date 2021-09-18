import request from '@/utils/request'

export default {
  getList(query, current = 1, size = 20) {
    return request({
      url: '/api/v1/event',
      method: 'get',
      params: { ...query, page: current, size: size }
    })
  },
  getEventDatabase(query) {
    return request({
      url: '/api/v1/event/db',
      method: 'get',
      params: { ...query }
    })
  }
}
