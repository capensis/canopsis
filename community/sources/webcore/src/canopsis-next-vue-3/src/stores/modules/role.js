import { computed, ref } from 'vue';
import { defineStore } from 'pinia';
import { find } from 'lodash';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createBasicCRUDActions } from '@/stores/plugins/crud';

export const useRoleStore = defineStore('role', () => {
  const items = ref([]);
  const meta = ref({});
  const pending = ref(false);
  const fetchingParams = ref();
  const hasItems = computed(() => items.value.length > 0);

  const getItemById = id => find(items.value, { _id: id });

  const fetchList = async ({ params } = {}) => {
    try {
      pending.value = true;
      fetchingParams.value = params;

      const response = await request.get(API_ROUTES.roles.list, { params });

      items.value = response.data ?? [];
      meta.value = response.meta ?? {};

      return response;
    } finally {
      pending.value = false;
    }
  };

  const fetchListWithPreviousParams = () => fetchList({ params: fetchingParams.value });

  const fetchListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.roles.list, { params });

  const fetchTemplatesListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.roles.templates, { params });

  const bulkUpdatePermissions = ({ data } = {}) => request.put(API_ROUTES.roles.bulkPermissions, data);

  return {
    items,
    meta,
    pending,
    fetchingParams,
    hasItems,

    getItemById,
    fetchList,
    fetchListWithPreviousParams,
    fetchListWithoutStore,
    fetchTemplatesListWithoutStore,
    bulkUpdatePermissions,

    ...createBasicCRUDActions({ route: API_ROUTES.roles.list }),
  };
});
