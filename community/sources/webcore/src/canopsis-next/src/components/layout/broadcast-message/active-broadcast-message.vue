<template lang="pug">
  div(v-if="activeMessages.length")
    broadcast-message(
      v-for="activeMessage in activeMessages",
      :key="activeMessage._id",
      :message="activeMessage.message",
      :color="activeMessage.color"
    )
      template(v-if="activeMessage.maintenance", #actions="")
        v-btn.my-0.ml-0.mr-2(
          color="white",
          outline,
          round,
          small,
          @click="showEditBroadcastMessageModal(activeMessage)"
        )
          v-icon(small) edit
        v-btn.my-0.ml-0.mr-2(color="white", outline, round, small, @click="logout")
          v-icon(small) logout
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import { SOCKET_ROOMS } from '@/config';
import { MODALS } from '@/constants';

import { maintenanceActionsMixin } from '@/mixins/maintenance/maintenance-actions';

import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';

const { mapActions } = createNamespacedHelpers('broadcastMessage');

export default {
  components: { BroadcastMessage },
  mixins: [maintenanceActionsMixin],
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
      fetchBroadcastMessagesListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchActiveBroadcastMessagesListWithoutStore: 'fetchActiveListWithoutStore',
      updateBroadcastMessage: 'update',
    }),

    setActiveMessages(activeMessages) {
      this.activeMessages = activeMessages;
    },

    async fetchList() {
      const data = await this.fetchActiveBroadcastMessagesListWithoutStore();

      this.setActiveMessages(data);
    },

    showEditBroadcastMessageModal(broadcastMessage) {
      this.$modals.show({
        name: MODALS.createMaintenance,
        config: {
          title: '',
          maintenance: pick(broadcastMessage, ['message', 'color']),
          action: async (data) => {
            await this.updateBroadcastMessage({
              id: broadcastMessage._id,
              data: { ...broadcastMessage, ...data },
            });

            this.fetchBroadcastMessagesListWithPreviousParams();
          },
        },
      });
    },
  },
};
</script>
