import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('watcher/watchedEntity');

export default {
  computed: {
    ...mapGetters({
      watchedEntities: 'items',
      watchedEntitiesPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchWatchedEntities: 'fetchList',
    }),
  },
};
