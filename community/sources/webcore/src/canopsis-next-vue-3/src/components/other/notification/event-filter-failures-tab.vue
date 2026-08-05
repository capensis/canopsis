<template>
  <event-filters-list
    ref="eventFiltersListEl"
    :event-filters="items"
    :pending="pending"
    :total-items="meta.total_count"
    :options="options"
    :updatable="hasUpdateAnyEventFilterAccess"
    :removable="hasDeleteAnyEventFilterAccess"
    :duplicable="hasCreateAnyEventFilterAccess"
    @update:options="$emit('update:options', $event)"
    @duplicate="showDuplicateRuleModal"
    @remove="showDeleteRuleModal"
    @edit="showEditRuleModal"
    @refresh="refresh"
  />
</template>

<script>
import { nextTick, ref, toRef } from 'vue';

import { EVENT_FILTER_EXPAND_PANEL_TABS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';

import { useEventFilterActions } from '@/components/other/event-filter/hooks/event-filters';
import EventFiltersList from '@/components/other/event-filter/event-filters-list.vue';

import { useNotificationActiveId } from './hooks/notifications';

export default {
  components: { EventFiltersList },

  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    meta: {
      type: Object,
      default: () => ({}),
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    activeId: {
      type: String,
      required: false,
    },
  },

  emits: ['update:options', 'refresh'],

  setup(props, { emit }) {
    const eventFiltersListEl = ref(null);
    const {
      hasCreateAccess: hasCreateAnyEventFilterAccess,
      hasUpdateAccess: hasUpdateAnyEventFilterAccess,
      hasDeleteAccess: hasDeleteAnyEventFilterAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.eventFilter);

    const refresh = () => emit('refresh');
    const {
      showDuplicateRuleModal,
      showEditRuleModal,
      showDeleteRuleModal,
    } = useEventFilterActions(refresh);

    const showRuleExpandPanelErrorTab = (rule = {}) => {
      nextTick(() => {
        eventFiltersListEl.value?.$refs?.advancedDataTable?.toggleExpand?.(rule);

        nextTick(() => {
          if (
            eventFiltersListEl.value?.$refs?.expandPanel?.activeTab &&
            eventFiltersListEl.value.$refs.expandPanel.activeTab !== EVENT_FILTER_EXPAND_PANEL_TABS.errors
          ) {
            eventFiltersListEl.value.$refs.expandPanel.activeTab = EVENT_FILTER_EXPAND_PANEL_TABS.errors;
          }
        });
      });
    };

    useNotificationActiveId({
      activeId: toRef(props, 'activeId'),
      items: toRef(props, 'items'),
      action: showRuleExpandPanelErrorTab,
    });

    return {
      eventFiltersListEl,
      hasCreateAnyEventFilterAccess,
      hasUpdateAnyEventFilterAccess,
      hasDeleteAnyEventFilterAccess,
      showDuplicateRuleModal,
      showEditRuleModal,
      showDeleteRuleModal,
      refresh,
    };
  },
};
</script>
