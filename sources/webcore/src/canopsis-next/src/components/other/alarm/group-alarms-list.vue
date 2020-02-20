<template lang="pug">
  div
    alarms-table(
      :widget="widget",
      :alarms="alarms",
      :totalItems="totalItems",
      :isEditingMode="isEditingMode",
      :hasColumns="hasGroupColumns",
      :columns="groupColumns",
      ref="alarmsTable"
    )
</template>

<script>
import RecordsPerPage from '@/components/tables/records-per-page.vue';

import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPaginationMixin from '@/mixins/widget/pagination';
import widgetRecordsPerPageMixin from '@/mixins/widget/records-per-page';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  components: { RecordsPerPage },
  mixins: [
    widgetColumnsMixin,
    widgetPaginationMixin,
    widgetRecordsPerPageMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    alarms: {
      type: Array,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    totalItems() {
      return this.alarms.length;
    },
  },
};
</script>
