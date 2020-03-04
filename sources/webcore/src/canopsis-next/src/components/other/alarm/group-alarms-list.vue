<template lang="pug">
  div
    alarms-list-table(
      :widget="widget",
      :alarms="alarms",
      :totalItems="totalItems",
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
import entitiesAlarmMixin from '@/mixins/entities/alarm';
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
    widgetColumnsMixin,
    entitiesAlarmMixin,
    widgetPaginationMixin,
    widgetGroupFetchQueryMixin,
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
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    totalItems() {
      return this.alarms.length;
    },
  },
  created() {
    this.$periodicRefresh.subscribe(this.fetchGroupAlarmListData);
  },
  beforeDestroy() {
    this.$periodicRefresh.unsubscribe(this.fetchGroupAlarmListData);
  },
  methods: {
    fetchGroupAlarmListData() {
      const query = this.getQuery();

      if (this.alarm.causes) {
        query.with_causes = true;
      }

      if (this.alarm.consequences) {
        query.with_consequences = true;
      }

      this.fetchAlarmItemWithParams(this.alarm, { ...query, with_steps: true });
    },
  },
};
</script>
