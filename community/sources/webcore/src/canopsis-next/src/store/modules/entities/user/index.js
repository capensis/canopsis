import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createCRUDModule({
  route: API_ROUTES.users,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
}, {
  actions: {
    /**
     * Action for user removing
     *
     * @param {ActionContext} context
     * @param {string} id
     * @returns {AxiosPromise<any>}
     */
    remove(context, { id }) {
      return request.delete(`${API_ROUTES.users}/${encodeURIComponent(id)}`);
    },

    /**
     * Method for update current user
     *
     * @param {ActionContext} context
     * @param {string} id
     * @param {User} data
     * @returns {AxiosPromise<any>}
     */
    updateCurrentUser(context, { data }) {
      return request.put(API_ROUTES.currentUser, data);
    },
  },
});
