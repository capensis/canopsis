import { computed, ref } from 'vue';
import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createBasicCRUDActions } from '@/stores/plugins/crud';

export const useCorporatePatternStore = defineStore('corporatePattern', () => {
  const items = ref([]);
  const meta = ref({});
  const pending = ref(false);
  const fetchingParams = ref();

  const hasItems = computed(() => items.value.length > 0);

  const fetchList = async ({ params } = {}) => {
    try {
      const preparedParams = {
        ...params,
        corporate: true,
      };

      pending.value = true;
      fetchingParams.value = preparedParams;

      const response = await request.get(API_ROUTES.pattern.list, { params: preparedParams });

      items.value = response?.data ?? [];
      meta.value = response?.meta ?? {};

      return response;
    } catch (err) {
      console.error(err);

      throw err;
    } finally {
      pending.value = false;
    }
  };

  const fetchListWithPreviousParams = () => fetchList({ params: fetchingParams.value });

  return {
    items,
    meta,
    pending,
    fetchingParams,
    hasItems,
    fetchList,
    fetchListWithPreviousParams,
    ...createBasicCRUDActions({ route: API_ROUTES.pattern.list }),
  };
});
