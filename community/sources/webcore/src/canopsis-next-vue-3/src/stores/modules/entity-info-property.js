import { computed } from 'vue';

import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useEntityInfoPropertyStore = createCRUDStore({
  id: 'entityInfoProperty',
  route: API_ROUTES.entityInfosProperties,
  bulkRoute: API_ROUTES.bulkEntityInfosProperties,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  setup({ items }) {
    const itemsWithAlias = computed(() => items.value.filter((item) => !!item.alias));
    const itemsWithoutAlias = computed(() => items.value.filter((item) => !item.alias));

    return {
      itemsWithAlias,
      itemsWithoutAlias,
    };
  },
});
