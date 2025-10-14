<template>
  <div v-if="activeMessages.length">
    <broadcast-message
      v-for="activeMessage in activeMessages"
      :key="activeMessage._id"
      :message="activeMessage.message"
      :color="activeMessage.color"
    >
      <template
        v-if="isLoggedIn"
        #actions=""
      >
        <template v-if="activeMessage.maintenance">
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
        <v-btn
          v-else
          class="my-0 ml-0 mr-2"
          color="white"
          outlined
          rounded
          small
          @click="showConfirmationMarkAsRead(activeMessage)"
        >
          <v-icon small>
            visibility_off
          </v-icon>
        </v-btn>
      </template>
    </broadcast-message>
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue';
import { pick } from 'lodash';

import { SOCKET_ROOMS } from '@/config';
import { MODALS } from '@/constants';

import { useAuth } from '@/hooks/auth';
import { useInfo } from '@/hooks/store/modules/info';
import { useBroadcastMessages } from '@/hooks/store/modules/broadcast-message';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSocketRoom } from '@/hooks/socket';

import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';

export default {
  components: { BroadcastMessage },
  setup() {
    const { t } = useI18n();
    const modals = useModals();

    const { isLoggedIn } = useAuth();
    const { maintenance, updateMaintenanceMode, fetchAppInfo } = useInfo();

    const {
      fetchBroadcastMessagesListWithPreviousParams,
      fetchActiveBroadcastMessagesListWithoutStore,
      updateBroadcastMessage,
      markBroadcastMessageAsRead,
    } = useBroadcastMessages();

    const activeMessages = ref([]);

    /**
     * Sets the active broadcast messages.
     *
     * @param {Array} [messages=[]] - Array of active broadcast message objects
     */
    const setActiveMessages = (messages = []) => activeMessages.value = messages;

    /**
     * Fetches the list of active broadcast messages and updates the local state.
     */
    const fetchList = async () => {
      const data = await fetchActiveBroadcastMessagesListWithoutStore();

      setActiveMessages(data);
    };

    /**
     * Disables maintenance mode by updating the maintenance mode setting
     * and refreshing the application info.
     */
    const disableMaintenanceMode = async () => {
      await updateMaintenanceMode({
        data: { enabled: false },
      });

      return fetchAppInfo();
    };

    /**
     * Shows a modal to edit a broadcast message.
     * Allows editing the message content and color for maintenance mode messages.
     *
     * @param {Object} [broadcastMessage={}] - The broadcast message object to edit
     * @param {string} broadcastMessage._id - The ID of the broadcast message
     * @param {string} broadcastMessage.message - The message text
     * @param {string} broadcastMessage.color - The message color
     */
    const showEditBroadcastMessageModal = (broadcastMessage = {}) => modals.show({
      name: MODALS.createMaintenance,
      config: {
        title: t('modals.createMaintenance.edit.title'),
        warningText: t('maintenance.maintenanceModeIsOn'),
        maintenance: pick(broadcastMessage, ['message', 'color']),
        action: async (data) => {
          await updateBroadcastMessage({
            id: broadcastMessage._id,
            data: { ...broadcastMessage, ...data },
          });

          return fetchBroadcastMessagesListWithPreviousParams();
        },
      },
    });

    /**
     * Shows a confirmation modal for leaving maintenance mode.
     * Upon confirmation, disables maintenance mode and refreshes the broadcast messages list.
     */
    const showConfirmationLeaveMaintenanceMode = () => modals.show({
      name: MODALS.confirmation,
      config: {
        title: t('modals.confirmationLeaveMaintenance.title'),
        text: t('modals.confirmationLeaveMaintenance.text'),
        action: async () => {
          await disableMaintenanceMode();

          return fetchList();
        },
      },
    });

    /**
     * Shows a confirmation modal for marking a broadcast message as read.
     * Upon confirmation, disables maintenance mode and refreshes the broadcast messages list.
     */
    const showConfirmationMarkAsRead = (broadcaseMessage = {}) => modals.show({
      name: MODALS.confirmation,
      config: {
        title: t('modals.confirmationMarkAsRead.title'),
        text: t('modals.confirmationMarkAsRead.text'),
        action: async () => {
          await markBroadcastMessageAsRead({ id: broadcaseMessage._id });

          return fetchList();
        },
      },
    });

    useSocketRoom(SOCKET_ROOMS.broadcastMessages, setActiveMessages);

    watch(maintenance, fetchList);

    onMounted(fetchList);

    return {
      activeMessages,
      isLoggedIn,
      showEditBroadcastMessageModal,
      showConfirmationLeaveMaintenanceMode,
      showConfirmationMarkAsRead,
    };
  },
};
</script>
