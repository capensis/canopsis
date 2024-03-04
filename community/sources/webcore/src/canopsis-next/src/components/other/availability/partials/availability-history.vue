<template>
  <v-layout class="gap-5" column>
    <availability-history-filters
      :sampling="query.sampling"
      :show-type.sync="showType"
      :min-interval-date="minDate"
      @update:sampling="$emit('update:sampling', $event)"
    />
    <availability-line-chart
      :availabilities="availabilities"
      :sampling="query.sampling"
      :display-parameter="displayParameter"
      :show-type="showType"
      v-on="$listeners"
    />
  </v-layout>
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
      required: true,
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
  },
  setup(props) {
    const showType = ref(props.defaultShowType);

    return {
      showType,
    };
  },
};
</script>
