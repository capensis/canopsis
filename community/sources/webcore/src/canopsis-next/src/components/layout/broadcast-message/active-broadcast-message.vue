<template>
  <div v-if="filteredActiveMessagesByRoute.length">
    <broadcast-message
      v-for="activeMessage in filteredActiveMessagesByRoute"
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
import { pick } from 'lodash';
import { computed, ref, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router/composables';

import { SOCKET_ROOMS } from '@/config';
import { MODALS, ROUTES_NAMES_TO_BROADCAST_MESSAGES } from '@/constants';

import { isBroadcastMessageViewMatchingRoute } from '@/helpers/entities/broadcast-message/list';

import { useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSocketRoom } from '@/hooks/socket';
import { useInfo } from '@/hooks/store/modules/info';
import { useBroadcastMessages } from '@/hooks/store/modules/broadcast-message';
import { useViewGroup } from '@/hooks/store/modules/view';

import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';

export default {
  components: { BroadcastMessage },
  setup() {
    const route = useRoute();

    const activeMessages = ref([]);

    const { t } = useI18n();
    const modals = useModals();
    const { isLoggedIn } = useAuth();
    const { maintenance, updateMaintenanceMode, fetchAppInfo } = useInfo();
    const {
      updateBroadcastMessage,
      fetchBroadcastMessagesListWithPreviousParams,
      fetchActiveBroadcastMessagesListWithoutStore,
    } = useBroadcastMessages();
    const { getViewById } = useViewGroup();

    /**
     * Filters active broadcast messages based on current route and view context
     * Shows messages that are configured for:
     * - The current route's broadcast message view type
     * - Specific view IDs that match the current route parameter
     * - View group IDs when available
     */
    const filteredActiveMessagesByRoute = computed(() => {
      const routeView = ROUTES_NAMES_TO_BROADCAST_MESSAGES[route.name];
      const { id: routeId } = route.params;
      const currentView = getViewById.value(routeId);

      return activeMessages.value.filter(({ views: messageViews }) => (messageViews || []).some(
        messageView => isBroadcastMessageViewMatchingRoute(messageView, routeView, routeId, currentView),
      ));
    });

    /**
     * Sets the active broadcast messages
     *
     * @param {Array} [messages=[]] - Array of broadcast messages to set as active
     */
    const setActiveMessages = (messages = []) => activeMessages.value = messages;

    /**
     * Fetches the list of active broadcast messages without storing in Vuex
     */
    const fetchList = async () => {
      const data = await fetchActiveBroadcastMessagesListWithoutStore();

      setActiveMessages(data);
    };

    /**
     * Disables maintenance mode by updating the maintenance mode status and refreshing app info
     */
    const disableMaintenanceMode = async () => {
      await updateMaintenanceMode({
        data: { enabled: false },
      });

      await fetchAppInfo();
    };

    /**
     * Shows a modal to edit the broadcast message associated with maintenance mode
     *
     * @param {Object} broadcastMessage - The broadcast message object to edit
     * @param {string} broadcastMessage._id - Unique identifier of the broadcast message
     * @param {string} broadcastMessage.message - Message content
     * @param {string} broadcastMessage.color - Message color
     */
    const showEditBroadcastMessageModal = (broadcastMessage) => {
      modals.show({
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

            fetchBroadcastMessagesListWithPreviousParams();
          },
        },
      });
    };

    /**
     * Shows a confirmation modal to leave maintenance mode
     */
    const showConfirmationLeaveMaintenanceMode = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.confirmationLeaveMaintenance.title'),
          text: t('modals.confirmationLeaveMaintenance.text'),
          action: async () => {
            await disableMaintenanceMode();
            await fetchList();
          },
        },
      });
    };

    useSocketRoom({
      room: SOCKET_ROOMS.broadcastMessages,
      data: null,
      needAuth: false,
      listener: setActiveMessages,
    });

    watch(maintenance, () => fetchList());

    onMounted(fetchList);

    return {
      filteredActiveMessagesByRoute,
      isLoggedIn,
      showEditBroadcastMessageModal,
      showConfirmationLeaveMaintenanceMode,
    };
  },
};
</script>
