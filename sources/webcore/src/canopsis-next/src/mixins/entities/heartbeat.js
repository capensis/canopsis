import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';

const { mapActions, mapGetters } = createNamespacedHelpers('heartbeat');

export default {
  mixins: [popupMixin],
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
      createHeartbeats: 'create',
      removeHeartbeats: 'remove',
    }),
  },
};
