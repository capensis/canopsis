import { ref } from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDStore } from '@/stores/plugins/crud';

export const usePbehaviorTypeStore = createCRUDStore({
  id: 'pbehaviorType',
  route: API_ROUTES.pbehavior.types,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    fetchNextPriority() {
      return request.get(API_ROUTES.pbehavior.nextTypesPriority);
    },

    bulkHide({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkTypesHide, data);
    },

    bulkUnhide({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkTypesUnhide, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.pbehavior.bulkTypes, { data });
    },
  },
});

export const usePbehaviorTypeFieldStore = () => {
  const items = ref([]);
  const pending = ref(false);

  const fetchList = async ({ params = { paginate: false, with_hidden: true } } = {}) => {
    try {
      pending.value = true;

      const { data = [] } = await request.get(API_ROUTES.pbehavior.types, { params });

      items.value = data;
    } finally {
      pending.value = false;
    }
  };

  return {
    items,
    pending,
    fetchList,
  };
};
