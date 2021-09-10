<template lang="pug">
  div(v-if="activeMessages.length")
    broadcast-message(
      v-for="activeMessage in activeMessages",
      :key="activeMessage._id",
      :message="activeMessage.message",
      :color="activeMessage.color"
    )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ACTIVE_BROADCAST_MESSAGE_FETCHING_INTERVAL } from '@/config';

import { broadcastMessageSchema } from '@/store/schemas';

import { createPollingMixin } from '@/mixins/polling';
import { registrableMixin } from '@/mixins/registrable';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('broadcastMessage');

export default {
  components: { BroadcastMessage },
  mixins: [
    registrableMixin([broadcastMessageSchema], 'activeMessages'),
    createPollingMixin({
      method: 'fetchActiveBroadcastMessagesList',
      delay: ACTIVE_BROADCAST_MESSAGE_FETCHING_INTERVAL,
      startOnMount: true,
    }),
  ],
  computed: {
    ...mapGetters(['activeMessages']),
  },
  mounted() {
    this.fetchActiveBroadcastMessagesList();
  },
  methods: {
    ...mapActions({
      fetchActiveBroadcastMessagesList: 'fetchActiveList',
    }),
  },
};
</script>
