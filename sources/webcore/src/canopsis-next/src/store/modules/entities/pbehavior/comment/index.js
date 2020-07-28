import qs from 'qs';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    create(context, { data, pbehaviorId }) {
      return request.post(API_ROUTES.pbehavior.comment.create, qs.stringify({
        ...data,

        pbehavior_id: pbehaviorId,
      }), {
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
      });
    },

    update(context, { data, pbehaviorId, commentId }) {
      return request.put(API_ROUTES.pbehavior.comment.update, qs.stringify({
        ...data,

        _id: commentId,
        pbehavior_id: pbehaviorId,
      }), {
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
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
