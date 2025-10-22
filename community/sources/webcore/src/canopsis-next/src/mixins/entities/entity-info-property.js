import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('entityInfoProperty');

export const entitiesEntityInfoPropertyMixin = {
  computed: {
    ...mapGetters({
      entityInfoProperties: 'items',
      entityInfoPropertiesWithAlias: 'itemsWithAlias',
      entityInfoPropertiesWithoutAlias: 'itemsWithoutAlias',
      entityInfoPropertyMeta: 'meta',
      entityInfoPropertyPending: 'pending',
    }),
  },

  methods: {
    ...mapActions({
      fetchEntityInfoPropertiesList: 'fetchList',
      createEntityInfoProperty: 'create',
      updateEntityInfoProperty: 'update',
      removeEntityInfoProperty: 'remove',
      fetchEntityInfoPropertiesListWithoutStore: 'fetchListWithoutStore',
    }),

    fetchAllEntityInfoPropertiesList() {
      return this.fetchEntityInfoPropertiesList({ params: { paginate: false } });
    },
  },
};
