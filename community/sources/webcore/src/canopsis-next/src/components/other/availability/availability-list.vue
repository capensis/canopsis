<template>
  <c-advanced-data-table
    :headers="headers"
    :items="availabilities"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #value="{ item }">
      <availability-list-column-value
        :availability="item"
        :display-parameter="displayParameter"
        :show-type="showType"
      />
    </template>
    <template #expand="{ item }">
      <availability-list-expand-panel
        :availability="item"
        :active-alarms-columns="activeAlarmsColumns"
        :resolved-alarms-columns="resolvedAlarmsColumns"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityListColumnValue from '@/components/other/availability/partials/availability-list-column-value.vue';

import AvailabilityListExpandPanel from './partials/availability-list-expand-panel.vue';

export default {
  components: { AvailabilityListColumnValue, AvailabilityListExpandPanel },
  props: {
    availabilities: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    activeAlarmsColumns: {
      type: Array,
      default: () => [],
    },
    resolvedAlarmsColumns: {
      type: Array,
      default: () => [],
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
  },
  computed: {
    isUptimeParameter() {
      return this.displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime;
    },

    headers() {
      const headers = [];

      headers.push({
        text: this.$t(`common.${this.isUptimeParameter ? 'uptime' : 'downtime'}`),
        value: 'value',
        sortable: false,
      });

      headers.push(...this.columns.map(column => ({
        ...column,
        value: `entity.${column.value}`,
      })));

      return headers;
    },
  },
};
</script>
