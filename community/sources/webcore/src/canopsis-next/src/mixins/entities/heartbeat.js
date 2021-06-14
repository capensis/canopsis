import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('heartbeat');

export default {
  computed: {
    ...mapGetters({
      heartbeatsPending: 'pending',
      heartbeatsMeta: 'meta',
      heartbeats: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchHeartbeatsList: 'fetchList',
      refreshHeartbeatsList: 'fetchListWithPreviousParams',
      createHeartbeat: 'create',
      updateHeartbeat: 'update',
      removeHeartbeat: 'remove',
    }),
  },
};
