import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('filter');

export const entitiesFilterMixin = {
  computed: {
    ...mapGetters({
      filters: 'items',
      getFilterById: 'getItemById',
      filtersPending: 'pending',
      filtersMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchFiltersListWithoutStore: 'fetchListWithoutStore',
      fetchFiltersListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchFiltersList: 'fetchList',
      removeFilter: 'remove',
      createFilter: 'create',
      updateFilter: 'update',
    }),
  },
};
