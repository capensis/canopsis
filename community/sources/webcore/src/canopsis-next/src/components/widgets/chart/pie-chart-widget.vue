<template lang="pug">
  v-layout.py-2(column)
    chart-widget-filters.px-3(
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
      pie-chart-metrics(
        v-if="hasMetrics",
        :metrics="aggregatedMetrics",
        :colors-by-metrics="colorsByMetrics",
        :title="widget.parameters.chart_title",
        :show-mode="widget.parameters.show_mode"
      )
</template>

<script>
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';
import { permissionsWidgetsPieChartInterval } from '@/mixins/permissions/widgets/chart/pie/interval';
import { permissionsWidgetsPieChartSampling } from '@/mixins/permissions/widgets/chart/pie/sampling';
import { permissionsWidgetsPieChartFilters } from '@/mixins/permissions/widgets/chart/pie/filters';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import PieChartMetrics from './partials/pie-chart-metrics.vue';

export default {
  inject: ['$system'],
  components: {
    ChartWidgetFilters,
    ChartLoader,
    PieChartMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    entitiesAggregatedMetricsMixin,
    permissionsWidgetsPieChartInterval,
    permissionsWidgetsPieChartSampling,
    permissionsWidgetsPieChartFilters,
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

    colorsByMetrics() {
      return this.widget.parameters.metrics.reduce((acc, { color, metric }) => {
        if (color) {
          acc[metric] = color;
        }

        return acc;
      }, {});
    },
  },
  methods: {
    getQuery() {
      return {
        ...this.getIntervalQuery(),

        parameters: this.widget.parameters.metrics.map(({ metric }) => ({
          metric,
          aggregate_func: this.widget.parameters.aggregate_func,
        })),
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    fetchList() {
      this.fetchAggregatedMetricsList({ widgetId: this.widget._id, params: this.getQuery() });
    },
  },
};
</script>
