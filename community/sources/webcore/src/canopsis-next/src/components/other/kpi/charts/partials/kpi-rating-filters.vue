<template>
  <div class="kpi-rating-toolbar">
    <c-layout class="ml-4 my-4" gap="3" wrap>
      <c-quick-date-interval-field
        :interval="query.interval"
        :min="minDate"
        :quick-ranges="quickRanges"
        @update:interval="updateQueryField('interval', $event)"
      />
      <c-filter-field
        :value="query.filter"
        :disabled="isUserMetric"
        class="kpi-rating-toolbar__filters"
        @update:value="updateQueryField('filter', $event)"
      />
      <kpi-rating-criteria-field
        :value="query.criteria"
        class="kpi-rating-toolbar__criteria"
        mandatory
        @update:value="updateCriteria"
      />
      <kpi-rating-metric-field
        :value="query.metric"
        :type="criteriaType"
        class="kpi-rating-toolbar__metric"
        hide-details
        @update:value="updateQueryField('metric', $event)"
      />
      <c-items-per-page-field
        :model-value="query.itemsPerPage"
        class="mt-4"
        @update:model-value="updateQueryField('itemsPerPage', $event)"
      />
    </c-layout>
  </div>
</template>

<script>
import {
  KPI_RATING_SETTINGS_TYPES,
  METRICS_QUICK_RANGES,
  USER_METRIC_PARAMETERS,
} from '@/constants';

import { isUserCriteria } from '@/helpers/entities/metric/form';
import { getAvailableMetricByCriteria } from '@/helpers/entities/metric/list';

import KpiRatingCriteriaField from '../form/fields/kpi-rating-criteria-field.vue';
import KpiRatingMetricField from '../form/fields/kpi-rating-metric-field.vue';

export default {
  components: {
    KpiRatingCriteriaField,
    KpiRatingMetricField,
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

  emits: ['update:query'],

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
    updateQueryField(field, value) {
      this.$emit('update:query', {
        ...this.query,
        [field]: value,
      });
    },

    updateCriteria(criteria) {
      this.$emit('update:query', {
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
  &__filters,
  &__criteria {
    max-width: 200px;
  }

  &__metric {
    max-width: 250px;
  }
}
</style>
