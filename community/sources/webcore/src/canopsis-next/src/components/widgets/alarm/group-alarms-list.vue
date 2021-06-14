<template lang="pug">
  div
    alarms-list-table(
      :widget="widget",
      :alarms="displayedAlarms",
      :total-items="alarmsMeta.total_count",
      :pagination.sync="vDataTablePagination",
      :is-editing-mode="isEditingMode",
      :has-columns="hasGroupColumns",
      :columns="groupColumns",
      :parent-alarm="alarm",
      expandable,
      hideGroups,
      ref="alarmsTable"
    )
    c-table-pagination(
      :total-items="alarmsMeta.total_count",
      :rows-per-page="query.limit",
      :page="query.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetGroupFetchQueryMixin from '@/mixins/widget/group-fetch-query';
import widgetExpandPanelAlarm from '@/mixins/widget/expand-panel/alarm/expand-panel';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  mixins: [
    widgetGroupFetchQueryMixin,
    widgetColumnsMixin,
    widgetExpandPanelAlarm,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
