import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('broadcastMessage');

export const entitiesBroadcastMessageMixin = {
  computed: {
    ...mapGetters({
      broadcastMessages: 'items',
      broadcastMessagesPending: 'pending',
      broadcastMessagesMeta: 'meta',
      activeMessages: 'activeMessages',
    }),
  },
  methods: {
    ...mapActions({
      fetchBroadcastMessagesList: 'fetchList',
      fetchActiveBroadcastMessagesList: 'fetchActiveList',
      createBroadcastMessage: 'create',
      updateBroadcastMessage: 'update',
      removeBroadcastMessage: 'remove',
    }),
  },
};
