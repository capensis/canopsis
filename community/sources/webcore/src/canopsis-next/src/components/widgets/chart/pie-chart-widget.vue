<template lang="pug">
  v-layout(column)
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
      template(v-if="aggregatedMetricsPending")
        v-fade-transition(v-if="aggregatedMetrics.length", key="progress", mode="out-in")
          v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
        v-layout.pa-4(v-else, justify-center)
          v-progress-circular(color="primary", indeterminate)
      pie-chart-metrics(
        v-if="aggregatedMetrics.length",
        :metrics="aggregatedMetrics",
        :colors-by-metrics="colorsByMetrics",
        :title="widget.parameters.chart_title",
        :show-mode="widget.parameters.show_mode"
      )
</template>

<script>
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { permissionsWidgetsPieChartInterval } from '@/mixins/permissions/widgets/chart/pie/interval';
import { permissionsWidgetsPieChartSampling } from '@/mixins/permissions/widgets/chart/pie/sampling';
import { permissionsWidgetsPieChartFilters } from '@/mixins/permissions/widgets/chart/pie/filters';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';
import PieChartMetrics from '@/components/widgets/chart/partials/pie-chart-metrics.vue';

export default {
  inject: ['$system'],
  components: {
    PieChartMetrics,
    ChartWidgetFilters,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    permissionsWidgetsPieChartInterval,
    permissionsWidgetsPieChartSampling,
    permissionsWidgetsPieChartFilters,
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
  computed: {
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
