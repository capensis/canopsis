<template>
  <alarms-list-table-with-pagination
    :parent-alarm="alarm"
    :widget="widget"
    :alarms="alarms"
    :meta="meta"
    :query="query"
    :columns="widget.parameters.widgetGroupColumns"
    :loading="pending"
    :refresh-alarms-list="refreshAlarmsList"
    :columns-settings="columnsSettings"
    expandable
    v-on="$listeners"
  />
</template>

<script>
import { computed } from 'vue';

import AlarmsListTableWithPagination from '../partials/alarms-list-table-with-pagination.vue';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  components: { AlarmsListTableWithPagination },
  props: {
    children: {
      type: Object,
      required: true,
    },
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    refreshAlarmsList: {
      type: Function,
      default: () => () => {},
    },
    columnsSettings: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const alarms = computed(() => props.children?.data ?? []);
    const meta = computed(() => props.children?.meta ?? {});

    return {
      alarms,
      meta,
    };
  },
};
</script>
