<template>
  <div class="kpi-rating-toolbar">
    <v-layout
      class="ml-4 my-4"
      wrap
    >
      <c-quick-date-interval-field
        class="mr-4"
        v-field="query.interval"
        :min="minDate"
        :quick-ranges="quickRanges"
      />
      <c-filter-field
        class="mr-4 kpi-rating-toolbar__filters"
        v-field="query.filter"
        :disabled="isUserMetric"
      />
      <kpi-rating-criteria-field
        class="mr-4 kpi-rating-toolbar__criteria"
        :value="query.criteria"
        mandatory
        @input="updateCriteria"
      />
      <kpi-rating-metric-field
        class="mr-4 kpi-rating-toolbar__metric"
        v-field="query.metric"
        :type="criteriaType"
        hide-details
      />
      <c-items-per-page-field
        class="mt-4"
        v-field="query.itemsPerPage"
      />
    </v-layout>
  </div>
</template>

<script>
import { KPI_RATING_SETTINGS_TYPES, METRICS_QUICK_RANGES, USER_METRIC_PARAMETERS } from '@/constants';

import { isUserCriteria } from '@/helpers/entities/metric/form';
import { getAvailableMetricByCriteria } from '@/helpers/entities/metric/list';

import { formMixin } from '@/mixins/form';

import KpiRatingCriteriaField from '../form/fields/kpi-rating-criteria-field.vue';
import KpiRatingMetricField from '../form/fields/kpi-rating-metric-field.vue';

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
    quickRanges() {
      return Object.values(METRICS_QUICK_RANGES);
    },

    isUserMetric() {
      return this.query.metric === USER_METRIC_PARAMETERS.totalUserActivity;
    },

    criteriaType() {
      return isUserCriteria(this.query.criteria?.label)
        ? KPI_RATING_SETTINGS_TYPES.user
        : KPI_RATING_SETTINGS_TYPES.entity;
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
