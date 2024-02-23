<template>
  <v-layout
    class="kpi-sli-toolbar ml-4 my-4"
    wrap
  >
    <c-quick-date-interval-field
      v-field="query.interval"
      :min="minDate"
      :quick-ranges="quickRanges"
      class="mr-4"
    />
    <c-sampling-field
      v-field="query.sampling"
      class="mr-4 kpi-sli-toolbar__sampling"
    />
    <kpi-sli-show-mode-field
      v-field="query.type"
      class="mr-4 kpi-sli-toolbar__show-mode"
    />
    <c-filter-field
      v-field="query.filter"
      class="kpi-sli-toolbar__filters"
    />
  </v-layout>
</template>

<script>
import { METRICS_QUICK_RANGES } from '@/constants';

import KpiSliShowModeField from '../form/fields/kpi-sli-show-mode-field.vue';

export default {
  components: { KpiSliShowModeField },
  model: {
    prop: 'query',
    event: 'input',
  },
  props: {
    query: {
      type: Object,
      required: true,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },
  computed: {
    quickRanges() {
      return Object.values(METRICS_QUICK_RANGES);
    },
  },
};
</script>

<style lang="scss">
.kpi-sli-toolbar {
  &__sampling {
    max-width: 200px;
  }

  &__show-mode {
    max-width: 150px;
  }

  &__filters {
    max-width: 200px;
  }
}
</style>
