import { createEntityModule } from '@/store/plugins/entities';

import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  route: API_ROUTES.users,
  entityType: ENTITIES_TYPES.user,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
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

    /**
     * Fetch users list with previous params
     *
     * @param {Function} dispatch
     * @param {Object} state
     * @returns {*}
     */
    fetchListWithPreviousParams({ dispatch, state }) {
      return dispatch('fetchList', {
        params: state.fetchingParams,
      });
    },

    /**
     * Fetch users list without store
     *
     * @param {VuexActionContext} context
     * @param {Object} [params]
     * @returns {*}
     */
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.users, params);
    },
  },
});
