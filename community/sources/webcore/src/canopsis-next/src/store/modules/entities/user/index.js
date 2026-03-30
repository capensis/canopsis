import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import { userToForm, formToUserRequest } from '@/helpers/entities/user/form';

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
     * @param {ActionContext} { commit }
     * @param {Object} params - The parameters object
     * @param {Object} params.data - The user data to update
     * @returns {Promise<User>}
     */
    async updateCurrentUser({ commit }, { data } = {}) {
      const newCurrentUser = await request.put(API_ROUTES.currentUser, data);

      commit('auth/FETCH_USER_COMPLETED', newCurrentUser, { root: true });

      return newCurrentUser;
    },

    updateCurrentUserTours({ rootGetters, dispatch }, { data } = {}) {
      const userForm = userToForm(rootGetters['auth/currentUser']);

      return dispatch('updateCurrentUser', {
        data: formToUserRequest({
          ...userForm,

          ui_tours: {
            ...userForm.ui_tours,
            ...data,
          },
        }),
      });
    },
  },
});
