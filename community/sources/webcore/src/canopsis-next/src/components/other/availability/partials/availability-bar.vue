<template>
  <c-layout gap="5" column>
    <availability-filters
      :interval="query.interval"
      :show-type="showType"
      :min-interval-date="minDate"
      @update:interval="$emit('update:interval', $event)"
      @update:showType="showType = $event"
    />
    <availability-bar-chart :availability="availability" :show-type="showType" />
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityBarChart from './availability-bar-chart.vue';
import AvailabilityFilters from './availability-filters.vue';

export default {
  components: { AvailabilityBarChart, AvailabilityFilters },

  props: {
    availability: {
      type: Object,
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
  },

  emits: ['update:interval'],

  setup(props) {
    const showType = ref(props.defaultShowType);

    return {
      showType,
    };
  },
};
</script>
