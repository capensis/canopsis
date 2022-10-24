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

import { SOCKET_ROOMS } from '@/config';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

const { mapActions } = createNamespacedHelpers('broadcastMessage');

export default {
  components: { BroadcastMessage },
  data() {
    return {
      activeMessages: [],
    };
  },
  mounted() {
    this.fetchList();

    this.$socket
      .join(SOCKET_ROOMS.broadcastMessages, false)
      .addListener(this.setActiveMessages);
  },
  beforeDestroy() {
    this.$socket
      .leave(SOCKET_ROOMS.broadcastMessages)
      .removeListener(this.setActiveMessages);
  },
  methods: {
    ...mapActions({
      fetchActiveBroadcastMessagesListWithoutStore: 'fetchActiveListWithoutStore',
    }),

    setActiveMessages(activeMessages) {
      this.activeMessages = activeMessages;
    },

    async fetchList() {
      const data = await this.fetchActiveBroadcastMessagesListWithoutStore();

      this.setActiveMessages(data);
    },
  },
};
</script>
