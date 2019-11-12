import i18n from '@/i18n';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import { groupSchema } from '@/store/schemas';

import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.viewGroup,
  entityType: ENTITIES_TYPES.group,
}, {
  actions: {
    async fetchList({ commit, dispatch }) {
      try {
        commit(types.FETCH_LIST);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.view,
          schema: [groupSchema],
          dataPreparer: d => Object.keys(d.groups).map(key => ({ _id: key, ...d.groups[key] })),
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        commit(types.FETCH_LIST_FAILED);

        await dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') }, { root: true });
      }
    },
  },
});
