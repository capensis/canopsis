import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('watcher/entity');

export default {
  computed: {
    ...mapGetters({
      getWatcherEntitiesListByWatcherId: 'getListByWatcherId',
      getWatcherEntitiesPendingByWatcherId: 'getPendingByWatcherId',
    }),
    watcherEntities() {
      return this.getWatcherEntitiesListByWatcherId(this.watcher.entity_id);
    },
    watcherEntitiesPending() {
      return this.getWatcherEntitiesPendingByWatcherId(this.watcher.entity_id);
    },
  },
  methods: {
    ...mapActions({
      fetchWatcherEntitiesList: 'fetchList',
      fetchWatcherEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
