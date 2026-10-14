<template>
  <c-layout gap="5" column>
    <availability-history-filters
      :sampling="query.sampling"
      :show-type="showType"
      :min-interval-date="minDate"
      @update:sampling="$emit('update:sampling', $event)"
      @update:show-type="showType = $event"
    />
    <availability-line-chart
      :availabilities="availabilities"
      :sampling="query.sampling"
      :display-parameter="displayParameter"
      :show-type="showType"
      :downloading="downloading"
      @export:csv="$emit('export:csv')"
      @export:png="$emit('export:png', $event)"
    />
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityHistoryFilters from '@/components/other/availability/partials/availability-history-filters.vue';
import AvailabilityLineChart from '@/components/other/availability/partials/availability-line-chart.vue';

export default {
  components: { AvailabilityLineChart, AvailabilityHistoryFilters },

  props: {
    availabilities: {
      type: Array,
      default: () => [],
    },
    query: {
      type: Object,
      required: true,
    },
    minDate: {
      type: Number,
      required: false,
    },
    defaultShowType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    displayParameter: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:sampling', 'export:csv', 'export:png'],

  setup(props) {
    const showType = ref(props.defaultShowType);

    return {
      showType,
    };
  },
};
</script>
