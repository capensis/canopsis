import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

import i18n from '@/i18n';

import { createCRUDModule, DEFAULT_ENTITY_MODULE_TYPES } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.entityInfosProperties,
  withWithoutStore: true,
}, {
  getters: {
    itemsWithAlias: state => state.items.filter(item => !!item.alias),
    itemsWithoutAlias: state => state.items.filter(item => !item.alias),
  },
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        await useRequestCancelling(async (source) => {
          commit(DEFAULT_ENTITY_MODULE_TYPES.FETCH_LIST, { params });

          const { data, meta } = await request.get(
            API_ROUTES.entityInfosProperties,
            { params, cancelToken: source.token },
          );

          commit(DEFAULT_ENTITY_MODULE_TYPES.FETCH_LIST_COMPLETED, {
            data,
            meta,
          });
        }, 'entity-infos-properties-list');
      } catch (err) {
        console.error(err);

        await dispatch('popups/error', { text: i18n.t('errors.default') }, { root: true });

        commit(DEFAULT_ENTITY_MODULE_TYPES.FETCH_LIST_FAILED, {});
      }
    },
  },
});
