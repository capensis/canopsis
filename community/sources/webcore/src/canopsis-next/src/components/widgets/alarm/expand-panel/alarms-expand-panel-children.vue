<template lang="pug">
  alarms-list-table-with-pagination(
    v-on="$listeners",
    :parent-alarm="alarm",
    :widget="widget",
    :alarms="alarms",
    :meta="meta",
    :query="query",
    :columns="columns",
    :editing="editing",
    :loading="pending"
  )
</template>

<script>
import { DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS, ALARM_ENTITY_FIELDS } from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';

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
    editing: {
      type: Boolean,
      default: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarms() {
      return this.children?.data ?? [];
    },

    meta() {
      return this.children?.meta ?? {};
    },

    columns() {
      const { widgetGroupColumns = [] } = this.widget.parameters;
      const columns = widgetGroupColumns.length
        ? widgetGroupColumns
        : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS);

      return columns.map(({ value, label, ...column }) => ({
        ...column,
        value,
        text: label,
        sortable: value !== ALARM_ENTITY_FIELDS.extraDetails,
      }));
    },
  },
};
</script>
