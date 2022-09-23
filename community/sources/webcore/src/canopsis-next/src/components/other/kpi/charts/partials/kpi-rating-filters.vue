<template lang="pug">
  div.kpi-rating-toolbar
    v-layout.ml-4.my-4(wrap)
      c-quick-date-interval-field.mr-4(v-field="query.interval", :min="minDate")
      c-filter-field.mr-4.kpi-rating-toolbar__filters(v-field="query.filter", :disabled="isUserMetric")
      kpi-rating-criteria-field.mr-4.kpi-rating-toolbar__criteria(
        :value="query.criteria",
        mandatory,
        @input="updateCriteria"
      )
      kpi-rating-metric-field.mr-4.kpi-rating-toolbar__metric(
        v-field="query.metric",
        :criteria="query.criteria"
      )
      c-records-per-page-field(v-field="query.rowsPerPage")
</template>

<script>
import { USER_METRIC_PARAMETERS } from '@/constants';

import { getAvailableMetricByCriteria } from '@/helpers/metrics';

import { formMixin } from '@/mixins/form';

import KpiRatingCriteriaField from './kpi-rating-criteria-field.vue';
import KpiRatingMetricField from './kpi-rating-metric-field.vue';

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
    minDate: {
      type: Number,
      required: false,
    },
  },
  computed: {
    isUserMetric() {
      return this.query.metric === USER_METRIC_PARAMETERS.totalUserActivity;
    },
  },
  methods: {
    updateCriteria(criteria) {
      this.updateModel({
        ...this.query,
        metric: getAvailableMetricByCriteria(this.query.metric, criteria?.label),
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
