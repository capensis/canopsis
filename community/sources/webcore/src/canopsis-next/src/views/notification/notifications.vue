<template>
  <div>
    <c-page-header />
    <v-card class="ma-6">
      <v-tabs
        :value="$route.fullPath"
        slider-color="primary"
        centered
      >
        <v-tab v-for="tab in tabs" :key="tab.key" :to="tab.to">
          {{ tab.label }}
          <c-circle-badge
            :outlined="!tab.active"
            :pending="tab.pending"
            class="ml-2"
            color="primary"
          >
            {{ tab.count }}
          </c-circle-badge>
        </v-tab>
      </v-tabs>
      <v-tabs-items :value="$route.fullPath">
        <v-tab-item v-for="tab in tabs" :key="tab.key">
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
import { onMounted, onBeforeUnmount, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router/composables';

import { SOCKET_ROOMS } from '@/config';

import { useI18n } from '@/hooks/i18n';
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
    const router = useRouter();
    const route = useRoute();
    const { t } = useI18n();
    const instance = useComponentInstance();

    const getTabHrefByKeyValue = (key, value) => { // TODO: move this function to router.js helper file
      const { href } = router.resolve({ query: { [key]: value } }, route);

      return href.replace('#', '');
    };

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

    const TABS_KEYS = {
      instructionsToApprove: 'instructionsToApprove',
      instructionsToRate: 'instructionsToRate',
      eventFilterFailures: 'eventFilterFailures',
    };

    const tabs = computed(() => [
      {
        key: TABS_KEYS.instructionsToApprove,
        label: t('notifications.tabs.instructionsToApprove'),
        to: getTabHrefByKeyValue('tabId', TABS_KEYS.instructionsToApprove),
        pending: instructionsToApprovePending.value,
        count: instructionsToApproveMeta.value.total_count,
      },
      {
        key: TABS_KEYS.instructionsToRate,
        label: t('notifications.tabs.instructionsToRate'),
        to: getTabHrefByKeyValue('tabId', TABS_KEYS.instructionsToRate),
        pending: instructionsToRatePending.value,
        count: instructionsToRateMeta.value.total_count,
      },
      {
        key: TABS_KEYS.eventFilterFailures,
        label: t('notifications.tabs.eventFilterFailures'),
        to: getTabHrefByKeyValue('tabId', TABS_KEYS.eventFilterFailures),
        pending: eventFilterFailuresPending.value,
        count: eventFilterFailuresMeta.value.total_count,
      },
    ].map(tab => ({ ...tab, active: tab.key === route.query.tabId })));

    const fetchAllLists = () => {
      fetchEventFilterFailuresList();
      fetchInstructionsToApproveList();
      fetchInstructionsToRateList();
    };

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
      tabs,

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
