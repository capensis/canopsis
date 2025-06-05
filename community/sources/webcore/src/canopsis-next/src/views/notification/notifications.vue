<template>
  <div>
    <c-page-header />
    <v-card class="ma-6">
      <v-tabs
        v-model="activeTab"
        slider-color="primary"
        centered
      >
        <v-tab>
          {{ $t('notifications.tabs.instructionsToApprove') }}
          <c-circle-badge
            :outlined="!instructionsToApproveMeta.total_count"
            :pending="instructionsToApprovePending"
            class="ml-2"
            color="primary"
          >
            {{ instructionsToApproveMeta.total_count }}
          </c-circle-badge>
        </v-tab>
        <v-tab>
          {{ $t('notifications.tabs.instructionsToRate') }}
          <c-circle-badge
            :outlined="!instructionsToRateMeta.total_count"
            :pending="instructionsToRatePending"
            class="ml-2"
            color="primary"
          >
            {{ instructionsToRateMeta.total_count }}
          </c-circle-badge>
        </v-tab>
        <v-tab>
          {{ $t('notifications.tabs.eventFilterFailures') }}
          <c-circle-badge
            :outlined="!eventFilterFailuresMeta.total_count"
            :pending="eventFilterFailuresPending"
            class="ml-2"
            color="primary"
          >
            {{ eventFilterFailuresMeta.total_count }}
          </c-circle-badge>
        </v-tab>
      </v-tabs>
      <v-tabs-items v-model="activeTab">
        <v-tab-item>
          <instructions-to-approve-tab
            :items="instructionsToApprove"
            :pending="instructionsToApprovePending"
            :meta="instructionsToApproveMeta"
            :options="instructionsToApproveOptions"
            @update:options="instructionsToApproveUpdateOptions"
            @refresh="fetchInstructionsToApproveList"
          />
        </v-tab-item>
        <v-tab-item>
          <instructions-to-rate-tab
            :items="instructionsToRate"
            :pending="instructionsToRatePending"
            :meta="instructionsToRateMeta"
            :options="instructionsToRateOptions"
            @update:options="instructionsToRateUpdateOptions"
            @refresh="fetchInstructionsToRateList"
          />
        </v-tab-item>
        <v-tab-item>
          <event-filter-failures-tab
            :items="eventFilterFailures"
            :pending="eventFilterFailuresPending"
            :meta="eventFilterFailuresMeta"
            :options="eventFilterFailuresOptions"
            @update:options="eventFilterFailuresUpdateOptions"
          />
        </v-tab-item>
      </v-tabs-items>
    </v-card>
    <c-fab-btn @refresh="fetchAllLists" />
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { useComponentInstance } from '@/hooks/vue';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useEventFilterStore } from '@/hooks/store/modules/event-filter';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useRemdeitionInstructionStatsStore } from '@/hooks/store/modules/remediation-instruction-stats';

import EventFilterFailuresTab from '@/components/other/notification/event-filter-failures-tab.vue';
import InstructionsToRateTab from '@/components/other/notification/instructions-to-rate-tab.vue';
import InstructionsToApproveTab from '@/components/other/notification/instructions-to-approve-tab.vue';

export default {
  components: {
    EventFilterFailuresTab,
    InstructionsToRateTab,
    InstructionsToApproveTab,
  },
  setup() {
    const instance = useComponentInstance();

    const { fetchRemediationInstructionsListWithoutStore } = useRemdeitionInstruction();
    const {
      data: instructionsToApprove,
      meta: instructionsToApproveMeta,
      pending: instructionsToApprovePending,
      options: instructionsToApproveOptions,
      updateOptions: instructionsToApproveUpdateOptions,
      fetchList: fetchInstructionsToApproveList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: { page: 1, itemsPerPage: 10 },
      fetchListHandler: ({ params }) => fetchRemediationInstructionsListWithoutStore({
        params: {
          ...params,
          only_to_approve: true,
        },
      }),
    });

    const { fetchEventFilterRulesListWithoutStore } = useEventFilterStore();
    const {
      data: eventFilterFailures,
      meta: eventFilterFailuresMeta,
      pending: eventFilterFailuresPending,
      options: eventFilterFailuresOptions,
      updateOptions: eventFilterFailuresUpdateOptions,
      fetchList: fetchEventFilterFailuresList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: { page: 1, itemsPerPage: 10 },
      fetchListHandler: ({ params }) => fetchEventFilterRulesListWithoutStore({
        params: {
          ...params,
          only_unread_failure: true,
        },
      }),
    });

    const { fetchRemediationInstructionStatsListWithoutStore } = useRemdeitionInstructionStatsStore();
    const {
      data: instructionsToRate,
      meta: instructionsToRateMeta,
      pending: instructionsToRatePending,
      options: instructionsToRateOptions,
      updateOptions: instructionsToRateUpdateOptions,
      fetchList: fetchInstructionsToRateList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: { page: 1, itemsPerPage: 10 },
      fetchListHandler: ({ params }) => fetchRemediationInstructionStatsListWithoutStore({
        params: {
          ...params,
          only_to_rate: true,
        },
      }),
    });

    const fetchAllLists = () => {
      fetchEventFilterFailuresList();
      fetchInstructionsToApproveList();
      fetchInstructionsToRateList();
    };

    // Reactive data
    const activeTab = ref(0);

    const handleNotificationUpdate = () => {};

    // Lifecycle
    onMounted(() => {
      // Initialize data
      fetchAllLists();

      // Join notifications websocket room
      if (instance?.$socket) {
        instance.$socket.join(SOCKET_ROOMS.notifications).addListener(handleNotificationUpdate);
      }
    });

    onBeforeUnmount(() => {
      // Leave notifications websocket room
      if (instance?.$socket) {
        instance.$socket.leave(SOCKET_ROOMS.notifications).removeListener(handleNotificationUpdate);
      }
    });

    return {
      activeTab,

      eventFilterFailures,
      eventFilterFailuresMeta,
      eventFilterFailuresPending,
      eventFilterFailuresOptions,
      eventFilterFailuresUpdateOptions,
      fetchEventFilterFailuresList,

      instructionsToApprove,
      instructionsToApproveMeta,
      instructionsToApprovePending,
      instructionsToApproveOptions,
      instructionsToApproveUpdateOptions,
      fetchInstructionsToApproveList,

      instructionsToRate,
      instructionsToRateMeta,
      instructionsToRatePending,
      instructionsToRateOptions,
      instructionsToRateUpdateOptions,
      fetchInstructionsToRateList,

      fetchAllLists,
    };
  },
};
</script>
