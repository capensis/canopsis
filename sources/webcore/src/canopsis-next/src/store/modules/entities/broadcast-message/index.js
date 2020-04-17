import i18n from '@/i18n';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import { broadcastMessageSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  FETCH_ACTIVE_LIST: 'FETCH_ACTIVE_LIST',
  FETCH_ACTIVE_LIST_COMPLETED: 'FETCH_ACTIVE_LIST_COMPLETED',
  FETCH_ACTIVE_LIST_FAILED: 'FETCH_ACTIVE_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.broadcastMessage.list,
  entityType: ENTITIES_TYPES.broadcastMessage,
  withFetchingParams: true,
}, {
  state: {
    activeMessagesIds: [],
    activeMessagesPending: false,
  },
  getters: {
    activeMessages: (state, getters, rootState, rootGetters) =>
      rootGetters['entities/getList'](ENTITIES_TYPES.broadcastMessage, state.activeMessagesIds),

    activeMessagesPending: state => state.activeMessagesPending,
  },
  mutations: {
    [types.FETCH_ACTIVE_LIST](state) {
      state.activeMessagesPending = true;
    },
    [types.FETCH_ACTIVE_LIST_COMPLETED](state, { allIds }) {
      state.activeMessagesIds = allIds;
    },
    [types.FETCH_ACTIVE_LIST_FAILED](state) {
      state.activeMessagesPending = false;
    },
  },
  actions: {
    async fetchActiveList({ commit, dispatch }) {
      try {
        commit(types.FETCH_ACTIVE_LIST);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.broadcastMessage.activeList,
          schema: [broadcastMessageSchema],
        }, { root: true });

        commit(types.FETCH_ACTIVE_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        commit(types.FETCH_ACTIVE_LIST_FAILED);
      }
    },
  },
});
