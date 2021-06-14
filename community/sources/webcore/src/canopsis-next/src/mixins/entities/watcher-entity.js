import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('watcher/entity');

export default {
  computed: {
    ...mapGetters({
      getWatcherEntitiesListByWatcherId: 'getListByWatcherId',
      getWatcherEntitiesPendingByWatcherId: 'getPendingByWatcherId',
      getWatcherEntitiesMetaByWatcherId: 'getMetaByWatcherId',
    }),
    watcherEntities() {
      return this.getWatcherEntitiesListByWatcherId(this.watcher._id);
    },
    watcherEntitiesPending() {
      return this.getWatcherEntitiesPendingByWatcherId(this.watcher._id);
    },
    watcherEntitiesMeta() {
      return this.getWatcherEntitiesMetaByWatcherId(this.watcher._id);
    },
  },
  methods: {
    ...mapActions({
      fetchWatcherEntitiesList: 'fetchList',
    }),
  },
};
