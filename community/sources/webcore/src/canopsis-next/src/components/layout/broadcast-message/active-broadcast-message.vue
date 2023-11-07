<template>
  <div v-if="activeMessages.length">
    <broadcast-message
      v-for="activeMessage in activeMessages"
      :key="activeMessage._id"
      :message="activeMessage.message"
      :color="activeMessage.color"
    >
      <template
        v-if="isLoggedIn && activeMessage.maintenance"
        #actions=""
      >
        <v-btn
          class="mr-2"
          color="white"
          outlined
          rounded
          small
          @click="showEditBroadcastMessageModal(activeMessage)"
        >
          <v-icon small>
            edit
          </v-icon>
        </v-btn>
        <v-btn
          class="my-0 ml-0 mr-2"
          color="white"
          outlined
          rounded
          small
          @click="showConfirmationLeaveMaintenanceMode"
        >
          <v-icon small>
            logout
          </v-icon>
        </v-btn>
      </template>
    </broadcast-message>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import { SOCKET_ROOMS } from '@/config';
import { MODALS } from '@/constants';

import { maintenanceActionsMixin } from '@/mixins/maintenance/maintenance-actions';
import { authMixin } from '@/mixins/auth';

import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';

const { mapActions } = createNamespacedHelpers('broadcastMessage');

export default {
  components: { BroadcastMessage },
  mixins: [maintenanceActionsMixin, authMixin],
  data() {
    return {
      activeMessages: [],
    };
  },
  watch: {
    maintenance: 'fetchList',
  },
  mounted() {
    this.fetchList();

    this.$socket
      .join(SOCKET_ROOMS.broadcastMessages, null, false)
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
          title: this.$t('modals.createMaintenance.edit.title'),
          warningText: this.$t('maintenance.maintenanceModeIsOn'),
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

    showConfirmationLeaveMaintenanceMode() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          title: this.$t('modals.confirmationLeaveMaintenance.title'),
          text: this.$t('modals.confirmationLeaveMaintenance.text'),
          action: async () => {
            await this.disableMaintenanceMode();
            await this.fetchList();
          },
        },
      });
    },
  },
};
</script>
