<template>
  <div>
    <c-page-header name="notifications">
      {{ $tc('common.notification', 2) }}
    </c-page-header>
    <v-card class="ma-6">
      <v-card-text>
        <v-tabs
          v-if="tabs.length"
          :value="$route.path"
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
        <v-tabs-items :value="$route.path" touchless>
          <v-tab-item v-for="tab in tabs" :key="tab.key" :value="tab.to">
            <component
              :is="tab.component"
              v-bind="tab.componentProps"
              v-on="tab.componentOn"
            />
          </v-tab-item>
        </v-tabs-items>
      </v-card-text>
    </v-card>
    <c-fab-btn @refresh="fetchAllLists" />
  </div>
</template>

<script>
import { computed, onBeforeMount, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router/composables';

import { SOCKET_ROOMS } from '@/config';
import { NOTIFICATIONS_PAGE_TABS_KEYS, QUICK_RANGES, USER_PERMISSIONS } from '@/constants';

import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';
import { useSocketRoom } from '@/hooks/socket';
import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useEventFilter } from '@/hooks/store/modules/event-filter';
import { useRemdeitionInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useRemdeitionInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

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
    activeId: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const router = useRouter();
    const route = useRoute();
    const { t } = useI18n();

    const {
      hasReadAccess: hasReadInstructionApprovalAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.remediationInstructionApprove);

    const {
      hasReadAccess: hasReadExecuteManualInstructionsAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.business.alarmsList.actions.executeInstruction);

    const {
      hasReadAccess: hasReadEventFilterAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.eventFilter);

    const hasReadAccess = computed(() => (
      hasReadInstructionApprovalAccess.value
      || hasReadExecuteManualInstructionsAccess.value
      || hasReadEventFilterAccess.value
    ));

    /**
     * Generates a tab href by key-value pair for routing
     *
     * @param {string} key - The parameter key
     * @param {string} value - The parameter value
     * @returns {string} The resolved href for the route
     */
    const getTabHrefByKeyValue = (key, value) => {
      const { href } = router.resolve({ params: { [key]: value } }, route);

      return href;
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

    const { fetchRemediationInstructionStatsListWithoutStore } = useRemdeitionInstructionStats();
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
          with_flags: true,
        },
      }),
    });

    /**
     * Computed property that converts the instructions rate query interval to timestamp
     *
     * @returns {Object} The converted timestamp interval
     */
    const instructionsRateQueryInterval = computed(() => (
      convertMetricIntervalToTimestamp({ interval: instructionsToRateOptions.value.interval })
    ));

    const { fetchEventFiltersListWithoutStore } = useEventFilter();
    const {
      data: eventFilterFailures,
      meta: eventFilterFailuresMeta,
      pending: eventFilterFailuresPending,
      options: eventFilterFailuresOptions,
      updateOptions: eventFilterFailuresUpdateOptions,
      fetchList: fetchEventFilterFailuresList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: { page: 1, itemsPerPage: 10 },
      fetchListHandler: ({ params }) => fetchEventFiltersListWithoutStore({
        params: {
          ...params,
          only_unread_failure: true,
        },
      }),
    });

    /**
     * Gets the active ID based on the current tab ID
     *
     * @param {string} tabId - The tab identifier
     * @returns {string|undefined} The active ID if the tab matches, undefined otherwise
     */
    const getActiveIdByTabId = tabId => (tabId === props.tabId ? props.activeId : undefined);

    /**
     * Computed property that defines the tabs configuration with their properties
     *
     * @returns {Array} Array of tab objects with their configuration
     */
    const tabs = computed(() => {
      const result = [];

      if (hasReadInstructionApprovalAccess.value) {
        result.push({
          key: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove,
          label: t('notifications.tabs.instructionsToApprove'),
          to: getTabHrefByKeyValue('tabId', NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove),
          pending: instructionsToApprovePending.value,
          count: instructionsToApproveMeta.value.total_count,
          component: InstructionsToApproveTab,
          componentProps: {
            items: instructionsToApprove.value,
            pending: instructionsToApprovePending.value,
            meta: instructionsToApproveMeta.value,
            options: instructionsToApproveOptions.value,
            activeId: getActiveIdByTabId(NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove),
          },
          componentOn: {
            refresh: fetchInstructionsToApproveList,
            'update:options': instructionsToApproveUpdateOptions,
          },
        });
      }

      if (hasReadExecuteManualInstructionsAccess.value) {
        result.push({
          key: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate,
          label: t('notifications.tabs.instructionsToRate'),
          to: getTabHrefByKeyValue('tabId', NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate),
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
            activeId: getActiveIdByTabId(NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate),
          },
          componentOn: {
            refresh: fetchInstructionsToRateList,
            'update:options': instructionsToRateUpdateOptions,
          },
        });
      }

      if (hasReadEventFilterAccess.value) {
        result.push({
          key: NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures,
          label: t('notifications.tabs.eventFilterFailures'),
          to: getTabHrefByKeyValue('tabId', NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures),
          pending: eventFilterFailuresPending.value,
          count: eventFilterFailuresMeta.value.total_count,
          hasAccess: hasReadEventFilterAccess.value,
          component: EventFilterFailuresTab,
          componentProps: {
            items: eventFilterFailures.value,
            pending: eventFilterFailuresPending.value,
            meta: eventFilterFailuresMeta.value,
            options: eventFilterFailuresOptions.value,
            activeId: getActiveIdByTabId(NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures),
          },
          componentOn: {
            refresh: fetchEventFilterFailuresList,
            'update:options': eventFilterFailuresUpdateOptions,
          },
        });
      }

      return result.map(tab => ({ ...tab, active: tab.key === props.tabId }));
    });

    /**
     * Fetches all notification lists (event filter failures, instructions to approve, and instructions to rate)
     */
    const fetchAllLists = () => {
      fetchEventFilterFailuresList();
      fetchInstructionsToApproveList();
      fetchInstructionsToRateList();
    };

    useSocketRoom(SOCKET_ROOMS.notifications, fetchAllLists);

    onBeforeMount(() => {
      if (!props.tabId) {
        router.replace({ query: { tabId: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove } });
      }
    });

    onMounted(fetchAllLists);

    return {
      tabs,

      hasReadAccess,
      hasReadInstructionApprovalAccess,
      hasReadExecuteManualInstructionsAccess,
      hasReadEventFilterAccess,

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
