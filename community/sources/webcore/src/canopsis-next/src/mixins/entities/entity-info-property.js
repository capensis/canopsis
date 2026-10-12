import { mapActions, mapState } from 'pinia';

import { useEntityInfoPropertyStore } from '@/stores/modules/entity-info-property';

export const entitiesEntityInfoPropertyMixin = {
  computed: {
    ...mapState(useEntityInfoPropertyStore, {
      entityInfoProperties: 'items',
      entityInfoPropertiesWithAlias: 'itemsWithAlias',
      entityInfoPropertiesWithoutAlias: 'itemsWithoutAlias',
      entityInfoPropertyMeta: 'meta',
      entityInfoPropertyPending: 'pending',
    }),
  },

  methods: {
    ...mapActions(useEntityInfoPropertyStore, {
      fetchEntityInfoPropertiesList: 'fetchList',
      createEntityInfoProperty: 'create',
      updateEntityInfoProperty: 'update',
      removeEntityInfoProperty: 'remove',
      bulkRemoveEntityInfoProperty: 'bulkRemove',
      fetchEntityInfoPropertiesListWithoutStore: 'fetchListWithoutStore',
    }),

    fetchAllEntityInfoPropertiesList() {
      return this.fetchEntityInfoPropertiesList({ params: { paginate: false } });
    },
  },
};
