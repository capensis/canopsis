<template lang="pug">
  div.kpi-rating-toolbar
    v-layout.ml-4.my-4(wrap)
      c-quick-date-interval-field.mr-4(v-field="query.interval")
      c-filters-field.mr-4.kpi-rating-toolbar__filters(v-field="query.filter")
      kpi-rating-criteria-field.mr-4.kpi-rating-toolbar__criteria(:value="query.criteria", @input="updateCriteria")
      kpi-rating-metric-field.mr-4.kpi-rating-toolbar__metric(v-field="query.metric", :criteria="query.criteria")
      c-records-per-page-field(v-field="query.rowsPerPage")
</template>

<script>
import KpiRatingCriteriaField from './kpi-rating-criteria-field.vue';
import KpiRatingMetricField from './kpi-rating-metric-field.vue';

import { getAvailableMetricByCriteria } from '@/helpers/metrics';

import { formMixin } from '@/mixins/form';

export default {
  components: { KpiRatingMetricField, KpiRatingCriteriaField },
  mixins: [formMixin],
  model: {
    prop: 'query',
    event: 'input',
  },
  props: {
    query: {
      type: Object,
      required: true,
    },
  },
  methods: {
    updateCriteria(criteria) {
      this.updateModel({
        ...this.query,
        metric: getAvailableMetricByCriteria(this.query.metric, criteria),
        criteria,
      });
    },
  },
};
</script>

<style scoped lang="scss">
.kpi-rating-toolbar {
  &__filters, &__criteria {
    max-width: 200px;
  }

  &__metric {
    max-width: 250px;
  }
}
</style>
