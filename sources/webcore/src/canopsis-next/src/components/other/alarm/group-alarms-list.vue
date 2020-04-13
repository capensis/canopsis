<template lang="pug">
  div
    alarms-list-table(
      :widget="widget",
      :alarms="displayedAlarms",
      :totalItems="alarmsMeta.total",
      :pagination.sync="vDataTablePagination",
      :isEditingMode="isEditingMode",
      :hasColumns="hasGroupColumns",
      :columns="groupColumns",
      ref="alarmsTable"
    )
    v-layout.white(align-center)
      v-flex(xs10)
        pagination(
          :page="query.page",
          :limit="query.limit",
          :total="alarmsMeta.total",
          @input="updateQueryPage"
        )
      v-spacer
      v-flex(xs2)
        records-per-page(:value="query.limit", @input="updateRecordsPerPage")
</template>

<script>
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPaginationMixin from '@/mixins/widget/pagination';
import widgetGroupFetchQueryMixin from '@/mixins/widget/group-fetch-query';
import widgetExpandPanelAlarm from '@/mixins/widget/expand-panel/alarm/expand-panel';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  components: {
    RecordsPerPage,
  },
  inject: ['$periodicRefresh'],
  mixins: [
    widgetGroupFetchQueryMixin,
    widgetColumnsMixin,
    widgetPaginationMixin,
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
