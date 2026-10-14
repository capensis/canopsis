<template>
  <div>
    <c-page-header name="notifications">
      {{ t('common.notification', 2) }}
    </c-page-header>

    <v-card class="ma-6">
      <v-card-text>
        <v-tabs
          v-if="tabs.length"
          :model-value="tabId"
          color="primary"
          align-tabs="center"
        >
          <v-tab
            v-for="tab in tabs"
            :key="tab.key"
            :value="tab.key"
            :to="tab.to"
            exact
          >
            {{ tab.label }}
            <c-circle-badge
              :outlined="!tab.active"
              :pending="tab.pending"
              :class="{ 'notifications-page__tab-badge--outlined': !tab.active }"
              class="notifications-page__tab-badge ml-2"
              color="primary"
            >
              {{ tab.count }}
            </c-circle-badge>
          </v-tab>
        </v-tabs>

        <v-window :model-value="tabId" touchless>
          <v-window-item v-for="tab in tabs" :key="tab.key" :value="tab.key">
            <component :is="tab.component" v-bind="tab.bind" />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <c-fab-btn @refresh="fetchAllLists" />
  </div>
</template>

<script>
import { computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { SOCKET_ROOMS } from '@/config';
import {
  NOTIFICATIONS_PAGE_TABS_KEYS,
  QUICK_RANGES,
  ROUTES_NAMES,
  USER_PERMISSIONS,
} from '@/constants';

import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { useCanPermission, useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useSocketRoom } from '@/hooks/socket';
import { useEventFilter } from '@/hooks/store/modules/event-filter';
import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useRemediationInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

import EventFilterFailuresTab from '@/components/other/notification/event-filter-failures-tab.vue';
import InstructionsToApproveTab from '@/components/other/notification/instructions-to-approve-tab.vue';
import InstructionsToRateTab from '@/components/other/notification/instructions-to-rate-tab.vue';

export default {
  components: {
    EventFilterFailuresTab,
    InstructionsToApproveTab,
    InstructionsToRateTab,
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
      hasAccess: hasInstructionApprovalAccess,
    } = useCanPermission(USER_PERMISSIONS.technical.exploitation.remediationInstructionApprove);
    const {
      hasAccess: hasExecuteManualInstructionsAccess,
    } = useCanPermission(USER_PERMISSIONS.business.alarmsList.actions.executeInstruction);
    const {
      hasReadAccess: hasReadEventFilterAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.eventFilter);

    const { fetchRemediationInstructionsListWithoutStore } = useRemediationInstruction();
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

    const { fetchRemediationInstructionStatsListWithoutStore } = useRemediationInstructionStats();
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

    const instructionsRateQueryInterval = computed(
      () => convertMetricIntervalToTimestamp({ interval: instructionsToRateOptions.value.interval }),
    );

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

    const getTabRoute = tabId => ({
      name: ROUTES_NAMES.notifications,
      params: { tabId },
      query: route.query,
    });
    const getActiveIdByTabId = tabId => (tabId === props.tabId ? props.activeId : undefined);

    const tabs = computed(() => {
      const result = [];

      if (hasInstructionApprovalAccess.value) {
        result.push({
          key: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove,
          label: t('notifications.tabs.instructionsToApprove'),
          to: getTabRoute(NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove),
          pending: instructionsToApprovePending.value,
          count: instructionsToApproveMeta.value.total_count,
          component: InstructionsToApproveTab,
          bind: {
            items: instructionsToApprove.value,
            pending: instructionsToApprovePending.value,
            meta: instructionsToApproveMeta.value,
            options: instructionsToApproveOptions.value,
            activeId: getActiveIdByTabId(NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove),
            onRefresh: fetchInstructionsToApproveList,
            'onUpdate:options': instructionsToApproveUpdateOptions,
          },
        });
      }

      if (hasExecuteManualInstructionsAccess.value) {
        result.push({
          key: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate,
          label: t('notifications.tabs.instructionsToRate'),
          to: getTabRoute(NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate),
          pending: instructionsToRatePending.value,
          count: instructionsToRateMeta.value.total_count,
          component: InstructionsToRateTab,
          bind: {
            items: instructionsToRate.value,
            pending: instructionsToRatePending.value,
            meta: instructionsToRateMeta.value,
            options: instructionsToRateOptions.value,
            interval: instructionsRateQueryInterval.value,
            accumulatedBefore: instructionsToRateMeta.value.accumulated_before,
            activeId: getActiveIdByTabId(NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToRate),
            onRefresh: fetchInstructionsToRateList,
            'onUpdate:options': instructionsToRateUpdateOptions,
          },
        });
      }

      if (hasReadEventFilterAccess.value) {
        result.push({
          key: NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures,
          label: t('notifications.tabs.eventFilterFailures'),
          to: getTabRoute(NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures),
          pending: eventFilterFailuresPending.value,
          count: eventFilterFailuresMeta.value.total_count,
          component: EventFilterFailuresTab,
          bind: {
            items: eventFilterFailures.value,
            pending: eventFilterFailuresPending.value,
            meta: eventFilterFailuresMeta.value,
            options: eventFilterFailuresOptions.value,
            activeId: getActiveIdByTabId(NOTIFICATIONS_PAGE_TABS_KEYS.eventFilterFailures),
            onRefresh: fetchEventFilterFailuresList,
            'onUpdate:options': eventFilterFailuresUpdateOptions,
          },
        });
      }

      return result.map(tab => ({ ...tab, active: tab.key === props.tabId }));
    });

    const fetchAllLists = () => {
      fetchEventFilterFailuresList();
      fetchInstructionsToApproveList();
      fetchInstructionsToRateList();
    };

    watch(tabs, (availableTabs) => {
      if (!availableTabs.length || availableTabs.some(tab => tab.key === props.tabId)) {
        return;
      }

      router.replace(getTabRoute(availableTabs[0].key));
    }, { immediate: true });

    useSocketRoom({ room: SOCKET_ROOMS.notifications, listener: fetchAllLists });
    onMounted(fetchAllLists);

    return {
      t,
      tabs,
      fetchAllLists,
    };
  },
};
</script>

<style lang="scss" scoped>
.notifications-page__tab-badge {
  &--outlined {
    color: rgb(var(--v-theme-primary));
  }
}
</style>
