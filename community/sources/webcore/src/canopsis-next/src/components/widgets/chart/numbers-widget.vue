<template lang="pug">
  v-layout.py-2(column)
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
      chart-loader(v-if="aggregatedMetricsPending", :has-metrics="hasMetrics")
      numbers-metrics(
        v-if="hasMetrics",
        :metrics="aggregatedMetrics",
        :title="widget.parameters.chart_title",
        :show-trend="widget.parameters.show_trend"
      )
</template>

<script>
import { isRatioMetric } from '@/helpers/metrics';
import { convertFilterToQuery } from '@/helpers/query';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';
import { permissionsWidgetsNumbersInterval } from '@/mixins/permissions/widgets/chart/numbers/interval';
import { permissionsWidgetsNumbersSampling } from '@/mixins/permissions/widgets/chart/numbers/sampling';
import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import NumbersMetrics from './partials/numbers-metrics.vue';

export default {
  inject: ['$system'],
  components: {
    ChartWidgetFilters,
    ChartLoader,
    NumbersMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    entitiesAggregatedMetricsMixin,
    permissionsWidgetsNumbersInterval,
    permissionsWidgetsNumbersSampling,
    permissionsWidgetsNumbersFilters,
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
  computed: {
    hasMetrics() {
      return !!this.aggregatedMetrics.length;
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
        widget_filters: convertFilterToQuery(this.query.filter),
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
