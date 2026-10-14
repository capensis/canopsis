export const entitiesFilterMixin = {
  computed: {
    filters: () => [],
    filtersPending: () => false,
    filtersMeta: () => ({}),
  },
  methods: {
    fetchFiltersListWithoutStore: () => Promise.resolve({ data: [] }),
    fetchFiltersListWithPreviousParams: () => Promise.resolve({ data: [] }),
    fetchFiltersList: () => Promise.resolve({ data: [] }),
    removeFilter: () => Promise.resolve(),
    createFilter: () => Promise.resolve(),
    updateFilter: () => Promise.resolve(),
  },
};
