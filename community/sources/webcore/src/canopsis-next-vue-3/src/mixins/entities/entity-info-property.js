export const entitiesEntityInfoPropertyMixin = {
  computed: {
    entityInfoProperties: () => [],
    entityInfoPropertiesWithAlias: () => [],
    entityInfoPropertiesWithoutAlias: () => [],
    entityInfoPropertyMeta: () => ({}),
    entityInfoPropertyPending: () => false,
  },
  methods: {
    fetchEntityInfoPropertiesList: () => Promise.resolve({ data: [] }),
    createEntityInfoProperty: () => Promise.resolve(),
    updateEntityInfoProperty: () => Promise.resolve(),
    removeEntityInfoProperty: () => Promise.resolve(),
    bulkRemoveEntityInfoProperty: () => Promise.resolve(),
    fetchEntityInfoPropertiesListWithoutStore: () => Promise.resolve({ data: [] }),
    fetchAllEntityInfoPropertiesList() {
      return this.fetchEntityInfoPropertiesList({ params: { paginate: false } });
    },
  },
};
