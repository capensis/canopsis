import { API_ROUTES } from '@/config';
import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    create(context, { data, pbehaviorId }) {
      return request.post(API_ROUTES.planning.pbehaviorComments, {
        ...data,
        pbehavior: pbehaviorId,
      });
    },

    update(context, { data, pbehaviorId, commentId }) {
      return request.put(API_ROUTES.planning.pbehaviorComments, {
        ...data,
        _id: commentId,
        pbehavior: pbehaviorId,
      });
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.planning.pbehaviorComments}/${id}`);
    },
  },
};
