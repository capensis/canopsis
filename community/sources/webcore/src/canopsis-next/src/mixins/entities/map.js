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
      createMap: 'create',
      updateMap: 'update',
      removeMap: 'remove',
      bulkRemoveMaps: 'bulkRemove',
    }),
  },
};
