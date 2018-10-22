import i18n from '@/i18n';

import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async editUser({ dispatch }, { user } = {}) {
      const editedUser = new FormData();
      editedUser.append('user', JSON.stringify(user));
      try {
        request.post(API_ROUTES.user, editedUser);
      } catch (err) {
        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
};
