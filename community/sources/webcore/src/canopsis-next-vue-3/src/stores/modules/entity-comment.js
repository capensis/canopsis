import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useEntityCommentStore = defineStore('entityComment', () => {
  const create = ({ data } = {}) => request.post(API_ROUTES.entityComments, data);
  const update = ({ id, data } = {}) => request.put(`${API_ROUTES.entityComments}/${id}`, data);
  const fetchListWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityComments, { params });

  return {
    create,
    update,
    fetchListWithoutStore,
  };
});
