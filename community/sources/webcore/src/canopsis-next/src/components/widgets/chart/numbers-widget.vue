<template lang="pug">
  v-layout(column)
    chart-widget-filters.mx-3(
      :widget-id="widget._id",
      :user-filters="userPreference.filters",
      :widget-filters="widget.filters",
      :locked-value="lockedFilter",
      :filters="mainFilter",
      :interval="query.interval",
      :sampling="query.sampling",
      :show-filter="hasAccessToUserFilter",
      :show-interval="hasAccessToInterval",
      :show-sampling="hasAccessToSampling",
      :filter-disabled="!hasAccessToListFilters",
      :filter-addable="hasAccessToAddFilter",
      :filter-editable="hasAccessToEditFilter",
      @update:filters="updateSelectedFilter",
      @update:sampling="updateSampling",
      @update:interval="updateInterval"
    )
    v-layout.pa-3(column)
      template(v-if="aggregatedMetricsPending")
        v-fade-transition(v-if="aggregatedMetrics.length", key="progress", mode="out-in")
          v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
        v-layout.pa-4(v-else, justify-center)
          v-progress-circular(color="primary", indeterminate)
      numbers-metrics(
        v-if="aggregatedMetrics.length",
        :metrics="aggregatedMetrics",
        :title="widget.parameters.chart_title",
        :show-trend="widget.parameters.show_trend"
      )
</template>

<script>
import { isRatioMetric } from '@/helpers/metrics';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { permissionsWidgetsNumbersInterval } from '@/mixins/permissions/widgets/chart/numbers/interval';
import { permissionsWidgetsNumbersSampling } from '@/mixins/permissions/widgets/chart/numbers/sampling';
import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';
import NumbersMetrics from '@/components/widgets/chart/partials/numbers-metrics.vue';

export default {
  inject: ['$system'],
  components: {
    NumbersMetrics,
    ChartWidgetFilters,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    permissionsWidgetsNumbersInterval,
    permissionsWidgetsNumbersSampling,
    permissionsWidgetsNumbersFilters,
    entitiesAggregatedMetricsMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      default: '',
    },
  },
  methods: {
    getQuery() {
      return {
        ...this.getIntervalQuery(),

        parameters: this.widget.parameters.metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({
          metric,
          aggregate_func: isRatioMetric(metric) ? undefined : aggregateFunc,
        })),
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    fetchList() {
      this.fetchAggregatedMetricsList({
        widgetId: this.widget._id,
        trend: this.widget.parameters.show_trend,
        params: this.getQuery(),
      });
    },
  },
};
</script>
