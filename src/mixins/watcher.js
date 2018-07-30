import { createNamespacedHelpers } from 'vuex';
import watchedEntitiesMixin from '@/mixins/watched-entities';

const { mapActions, mapGetters } = createNamespacedHelpers('watcher');

export default {
  mixins: [
    watchedEntitiesMixin,
  ],
  methods: {
    ...mapActions({
      fetchWatchersList: 'fetchList',
    }),
  },
  computed: {
    ...mapGetters({
      watchers: 'items',
      getWatcher: 'item',
    }),
  },
};
