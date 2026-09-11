<template>
  <div v-if="filteredActiveMessagesByRoute.length">
    <broadcast-message
      v-for="activeMessage in filteredActiveMessagesByRoute"
      :key="activeMessage._id"
      :message="activeMessage.message"
      :color="activeMessage.color"
    >
      <template v-if="isLoggedIn" #actions>
        <template v-if="activeMessage.maintenance">
          <v-btn
            class="mr-2"
            color="white"
            variant="outlined"
            rounded="pill"
            size="small"
            @click="showEditBroadcastMessageModal(activeMessage)"
          >
            <c-icon small>
              edit
            </c-icon>
          </v-btn>
          <v-btn
            class="my-0 ml-0 mr-2"
            color="white"
            variant="outlined"
            rounded="pill"
            size="small"
            @click="showConfirmationLeaveMaintenanceMode"
          >
            <c-icon small>
              logout
            </c-icon>
          </v-btn>
        </template>
        <v-btn
          v-if="!activeMessage.maintenance"
          class="my-0 ml-0 mr-2"
          color="white"
          variant="outlined"
          rounded="pill"
          size="small"
          @click="showConfirmationMarkAsRead(activeMessage)"
        >
          <c-icon small>
            close
          </c-icon>
        </v-btn>
      </template>
    </broadcast-message>
  </div>
</template>

<script>
import { pick } from 'lodash';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { SOCKET_ROOMS } from '@/config';
import { MODALS, ROUTES_NAMES_TO_BROADCAST_MESSAGES } from '@/constants';

import { isBroadcastMessageViewMatchingRoute } from '@/helpers/entities/broadcast-message/list';

import { useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useSocketRoom } from '@/hooks/socket';
import { useBroadcastMessages } from '@/hooks/store/modules/broadcast-message';
import { useInfo } from '@/hooks/store/modules/info';
import { useView } from '@/hooks/store/modules/view';

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
      markBroadcastMessageAsRead,
      fetchBroadcastMessagesListWithPreviousParams,
      fetchActiveBroadcastMessagesListWithoutStore,
    } = useBroadcastMessages();
    const { getViewById } = useView();

    const filteredActiveMessagesByRoute = computed(() => {
      const routeView = ROUTES_NAMES_TO_BROADCAST_MESSAGES[route.name];
      const { id: routeId } = route.params;
      const currentView = getViewById(routeId);

      const hasMatchingView = message => (message.views || []).some(
        messageView => isBroadcastMessageViewMatchingRoute(messageView, routeView, routeId, currentView),
      );

      return activeMessages.value.filter(message => message.maintenance || hasMatchingView(message));
    });

    const setActiveMessages = (messages = []) => (activeMessages.value = messages);

    const fetchList = async () => {
      const data = await fetchActiveBroadcastMessagesListWithoutStore();

      setActiveMessages(data);
    };

    const disableMaintenanceMode = async () => {
      await updateMaintenanceMode({
        data: { enabled: false },
      });

      return fetchAppInfo();
    };

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

    const showConfirmationLeaveMaintenanceMode = () => {
      modals.show({
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
    };

    const showConfirmationMarkAsRead = (broadcastMessage = {}) => modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('modals.confirmationMarkAsRead.title'),
          text: t('modals.confirmationMarkAsRead.text'),
          action: async () => {
            await markBroadcastMessageAsRead({ id: broadcastMessage._id });

            return fetchList();
          },
        },
      });

    useSocketRoom({
      room: SOCKET_ROOMS.broadcastMessages,
      payload: null,
      authNeeded: false,
      listener: setActiveMessages,
    });

    watch(maintenance, () => fetchList());
    watch(isLoggedIn, fetchList);

    onMounted(fetchList);

    return {
      filteredActiveMessagesByRoute,
      isLoggedIn,
      showEditBroadcastMessageModal,
      showConfirmationLeaveMaintenanceMode,
      showConfirmationMarkAsRead,
    };
  },
};
</script>
