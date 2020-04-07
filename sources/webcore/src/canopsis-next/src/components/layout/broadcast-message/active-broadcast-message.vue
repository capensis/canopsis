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

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

const { mapActions } = createNamespacedHelpers('broadcastMessage');

export default {
  components: { BroadcastMessage },
  data() {
    return {
      activeMessages: [],
      timeout: null,
    };
  },
  mounted() {
    this.startPolling();
  },
  beforeDestroy() {
    this.stopPolling();
  },
  methods: {
    ...mapActions({
      fetchActiveBroadcastMessageWithoutStore: 'fetchActiveWithoutStore',
    }),

    async startPolling() {
      this.activeMessages = await this.fetchActiveBroadcastMessageWithoutStore();
      this.timeout = setTimeout(this.startPolling, ACTIVE_BROADCAST_MESSAGE_FETCHING_INTERVAL);
    },

    stopPolling() {
      clearTimeout(this.timeout);
    },
  },
};
</script>
