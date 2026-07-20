import { defineStore } from 'pinia'

import { API_ROUTES } from '@/config'

import request from '@/services/request'

export const usePbehaviorCommentStore = defineStore('pbehaviorComment', () => {
  const create = ({ data, pbehaviorId } = {}) =>
    request.post(API_ROUTES.pbehavior.pbehaviorComments, {
      ...data,
      pbehavior: pbehaviorId,
    })

  const update = ({ data, pbehaviorId, commentId } = {}) =>
    request.put(API_ROUTES.pbehavior.pbehaviorComments, {
      ...data,
      _id: commentId,
      pbehavior: pbehaviorId,
    })

  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.pbehavior.pbehaviorComments}/${id}`)

  return {
    create,
    update,
    remove,
  }
})
