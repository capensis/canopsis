<template>
  <v-layout
    class="kpi-sli-toolbar ml-4 my-4"
    wrap
  >
    <c-quick-date-interval-field
      class="mr-4"
      v-field="query.interval"
      :min="minDate"
      :quick-ranges="quickRanges"
    />
    <c-sampling-field
      class="mr-4 kpi-sli-toolbar__sampling"
      v-field="query.sampling"
    />
    <kpi-sli-show-mode-field
      class="mr-4 kpi-sli-toolbar__show-mode"
      v-field="query.type"
    />
    <c-filter-field
      class="kpi-sli-toolbar__filters"
      v-field="query.filter"
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
