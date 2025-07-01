<template>
  <div>
    <c-page-header />
    <v-card class="ma-6">
      <v-tabs
        :value="$route.fullPath"
        slider-color="primary"
        centered
      >
        <v-tab
          v-for="tab in tabs"
          :key="tab.key"
          :to="tab.to"
          exact
          ripple
        >
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
      <v-tabs-items :value="$route.fullPath" touchless>
        <v-tab-item v-for="tab in tabs" :key="tab.key" :value="tab.to">
          <component
            :is="tab.component"
            v-bind="tab.componentProps"
            v-on="tab.componentOn"
          />
        </v-tab-item>
      </v-tabs-items>
    </v-card>
    <c-fab-btn @refresh="fetchAllLists" />
  </div>
</template>

<script>
import { onBeforeMount, onMounted, onBeforeUnmount, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router/composables';

import { SOCKET_ROOMS } from '@/config';
import { QUICK_RANGES } from '@/constants';

import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';
import { useComponentInstance } from '@/hooks/vue';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useEventFilterStore } from '@/hooks/store/modules/event-filter';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useRemdeitionInstructionStatsStore } from '@/hooks/store/modules/remediation-instruction-stats';

import InstructionsToApproveTab from '@/components/other/notification/instructions-to-approve-tab.vue';
import InstructionsToRateTab from '@/components/other/notification/instructions-to-rate-tab.vue';
import EventFilterFailuresTab from '@/components/other/notification/event-filter-failures-tab.vue';

export default {
  components: {
    InstructionsToApproveTab,
    InstructionsToRateTab,
    EventFilterFailuresTab,
  },
  props: {
    tabId: {
      type: String,
      default: '',
    },
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

    const { fetchRemediationInstructionStatsListWithoutStore } = useRemdeitionInstructionStatsStore();
    const {
      data: instructionsToRate,
      meta: instructionsToRateMeta,
      pending: instructionsToRatePending,
      options: instructionsToRateOptions,
      updateOptions: instructionsToRateUpdateOptions,
      fetchList: fetchInstructionsToRateList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: 10,
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
      fetchListHandler: ({ params }) => fetchRemediationInstructionStatsListWithoutStore({
        params: {
          ...params,
          only_to_rate: true,
        },
      }),
    });

    const instructionsRateQueryInterval = computed(() => (
      convertMetricIntervalToTimestamp({ interval: instructionsToRateOptions.value.interval })
    ));

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

    const TABS_KEYS = {
      instructionsToApprove: 'instructions-to-approve',
      instructionsToRate: 'instructions-to-rate',
      eventFilterFailures: 'event-filter-failures',
    };

    const tabs = computed(() => [
      {
        key: TABS_KEYS.instructionsToApprove,
        label: t('notifications.tabs.instructionsToApprove'),
        to: getTabHrefByKeyValue('tabId', TABS_KEYS.instructionsToApprove),
        pending: instructionsToApprovePending.value,
        count: instructionsToApproveMeta.value.total_count,
        component: InstructionsToApproveTab,
        componentProps: {
          items: instructionsToApprove.value,
          pending: instructionsToApprovePending.value,
          meta: instructionsToApproveMeta.value,
          options: instructionsToApproveOptions.value,
        },
        componentOn: {
          refresh: fetchInstructionsToApproveList,
          'update:options': instructionsToApproveUpdateOptions,
        },
      },
      {
        key: TABS_KEYS.instructionsToRate,
        label: t('notifications.tabs.instructionsToRate'),
        to: getTabHrefByKeyValue('tabId', TABS_KEYS.instructionsToRate),
        pending: instructionsToRatePending.value,
        count: instructionsToRateMeta.value.total_count,
        component: InstructionsToRateTab,
        componentProps: {
          items: instructionsToRate.value,
          pending: instructionsToRatePending.value,
          meta: instructionsToRateMeta.value,
          options: instructionsToRateOptions.value,
          interval: instructionsRateQueryInterval.value,
          accumulatedBefore: instructionsToRateMeta.value.accumulated_before,
        },
        componentOn: {
          refresh: fetchInstructionsToRateList,
          'update:options': instructionsToRateUpdateOptions,
        },
      },
      {
        key: TABS_KEYS.eventFilterFailures,
        label: t('notifications.tabs.eventFilterFailures'),
        to: getTabHrefByKeyValue('tabId', TABS_KEYS.eventFilterFailures),
        pending: eventFilterFailuresPending.value,
        count: eventFilterFailuresMeta.value.total_count,
        component: EventFilterFailuresTab,
        componentProps: {
          items: eventFilterFailures.value,
          pending: eventFilterFailuresPending.value,
          meta: eventFilterFailuresMeta.value,
          options: eventFilterFailuresOptions.value,
        },
        componentOn: {
          refresh: fetchEventFilterFailuresList,
          'update:options': eventFilterFailuresUpdateOptions,
        },
      },
    ].map(tab => ({ ...tab, active: tab.key === route.query.tabId })));

    const fetchAllLists = () => {
      fetchEventFilterFailuresList();
      fetchInstructionsToApproveList();
      fetchInstructionsToRateList();
    };

    const handleNotificationUpdate = () => {};

    onBeforeMount(() => {
      if (!route.query.tabId) {
        router.replace({ query: { tabId: TABS_KEYS.instructionsToApprove } });
      }
    });

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

      eventFilterFailures,
      eventFilterFailuresMeta,
      eventFilterFailuresPending,
      eventFilterFailuresOptions,
      eventFilterFailuresUpdateOptions,
      fetchEventFilterFailuresList,

      fetchAllLists,
    };
  },
};
</script>
