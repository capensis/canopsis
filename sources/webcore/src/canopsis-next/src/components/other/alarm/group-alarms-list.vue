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
import widgetQueryMixin from '@/mixins/widget/query';

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
  mixins: [
    widgetColumnsMixin,
    entitiesAlarmMixin,
    widgetPaginationMixin,
    widgetQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    details: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    alarmId: {
      type: String,
      required: true,
    },
  },
  computed: {
    totalItems() {
      return this.alarms.length;
    },
  },
};
</script>
