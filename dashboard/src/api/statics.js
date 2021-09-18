import request from '@/utils/request'

export default {
  getEventInfo(query) {
    return request({
      url: '/api/v1/event/info',
      method: 'get',
      params: { ...query }
    })
  }
}
