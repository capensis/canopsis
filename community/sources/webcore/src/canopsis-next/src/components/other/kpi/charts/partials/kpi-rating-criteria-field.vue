<template lang="pug">
  v-select(
    v-field="value",
    :items="criteriaTypes",
    :label="$t('kpiMetrics.criteria')",
    hide-details
  )
</template>

<script>
import { KPI_RATING_CRITERIA } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: true,
    },
  },
  computed: {
    criteriaTypes() {
      return Object.values(KPI_RATING_CRITERIA).map(value => ({
        value,
        text: this.$t(`kpiMetrics.criteriaTypes.${value}`),
      }));
    },
  },
};
</script>
