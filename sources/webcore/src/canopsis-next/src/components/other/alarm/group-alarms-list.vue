<template lang="pug">
  div
    alarms-list-table(
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
import widgetColumnsMixin from '@/mixins/widget/columns';

/**
 * Group-alarm-list component
 *
 * @module alarm
 *
 */
export default {
  mixins: [
    widgetColumnsMixin,
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
