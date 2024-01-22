<template>
  <alarms-list-table-with-pagination
    v-on="$listeners"
    :parent-alarm="alarm"
    :widget="widget"
    :alarms="alarms"
    :meta="meta"
    :query="query"
    :columns="widget.parameters.widgetGroupColumns"
    :loading="pending"
    :refresh-alarms-list="refreshAlarmsList"
    expandable
  />
</template>

<script>
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
  },
  computed: {
    alarms() {
      return this.children?.data ?? [];
    },

    meta() {
      return this.children?.meta ?? {};
    },
  },
};
</script>
