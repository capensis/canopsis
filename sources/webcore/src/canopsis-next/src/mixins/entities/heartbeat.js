import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('heartbeat');

export default {
  computed: {
    ...mapGetters({
      heartbeatsPending: 'pending',
      heartbeats: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchHeartbeatsList: 'fetchList',
      refreshHeartbeatsList: 'fetchListWithPreviousParams',
      createHeartbeat: 'create',
      removeHeartbeat: 'remove',
    }),
  },
};
