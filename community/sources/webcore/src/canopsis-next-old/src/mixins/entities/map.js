import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('map');

export const entitiesMapMixin = {
  computed: {
    ...mapGetters({
      maps: 'items',
      mapsPending: 'pending',
      mapsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchMapsList: 'fetchList',
      fetchMapsListWithoutStore: 'fetchListWithoutStore',
      fetchMapWithoutStore: 'fetchItemWithoutStore',
      fetchMapStateWithoutStore: 'fetchItemStateWithoutStore',
      createMap: 'create',
      updateMap: 'update',
      removeMap: 'remove',
      bulkRemoveMaps: 'bulkRemove',
    }),
  },
};
