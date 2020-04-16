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

import registrableMixin from '@/mixins/registrable';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('broadcastMessage');

export default {
  components: { BroadcastMessage },
  mixins: [registrableMixin([broadcastMessageSchema], 'activeMessages')],
  data() {
    return {
      timeout: null,
    };
  },
  computed: {
    ...mapGetters(['activeMessages']),
  },
  mounted() {
    this.startPolling();
  },
  beforeDestroy() {
    this.stopPolling();
  },
  methods: {
    ...mapActions({
      fetchActiveBroadcastMessagesList: 'fetchActiveList',
    }),

    async startPolling() {
      await this.fetchActiveBroadcastMessagesList();

      this.timeout = setTimeout(this.startPolling, ACTIVE_BROADCAST_MESSAGE_FETCHING_INTERVAL);
    },

    stopPolling() {
      clearTimeout(this.timeout);
    },
  },
};
</script>
