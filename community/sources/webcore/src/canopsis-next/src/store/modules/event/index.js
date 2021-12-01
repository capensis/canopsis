import { isArray } from 'lodash';

import { API_ROUTES } from '@/config';

import i18n from '@/i18n';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    async create({ dispatch, rootGetters }, { data }) {
      try {
        const currentUser = rootGetters['auth/currentUser'];
        const prepareEventByUserFields = event => ({
          ...event,

          user_id: currentUser._id,
          author: currentUser.name,
        });

        const preparedData = isArray(data)
          ? data.map(prepareEventByUserFields)
          : prepareEventByUserFields(data);

        await request.post(API_ROUTES.event, preparedData);

        await dispatch('popups/success', { text: i18n.t('success.default') }, { root: true });
      } catch (e) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });
        console.warn(e);
      }
    },
  },
};
