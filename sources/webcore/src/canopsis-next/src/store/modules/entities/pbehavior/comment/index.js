import { API_ROUTES } from '@/config';
import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    create(context, { data, pbehaviorId }) {
      return request.post(API_ROUTES.pbehavior.comment.create, {}, {
        params: {
          ...data,

          pbehavior_id: pbehaviorId,
        },
      });
    },

    update(context, { data, pbehaviorId, commentId }) {
      return request.put(API_ROUTES.pbehavior.comment.update, {}, {
        params: {
          ...data,

          _id: commentId,
          pbehavior_id: pbehaviorId,
        },
      });
    },

    remove(context, { id, pbehaviorId }) {
      return request.delete(API_ROUTES.pbehavior.comment.delete, {
        params: {
          _id: id,
          pbehavior_id: pbehaviorId,
        },
      });
    },
  },
};
