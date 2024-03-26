<template>
  <v-layout
    class="py-2"
    column
  >
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreference.filters"
      :widget-filters="widget.filters"
      :locked-value="query.lockedFilter"
      :filters="query.filter"
      :interval="query.interval"
      :sampling="query.sampling"
      :show-filter="hasAccessToUserFilter"
      :show-interval="hasAccessToInterval"
      :show-sampling="hasAccessToSampling"
      :filter-disabled="!hasAccessToListFilters"
      :filter-addable="hasAccessToAddFilter"
      :filter-editable="hasAccessToEditFilter"
      class="px-3"
      @update:filters="updateSelectedFilter"
      @update:sampling="updateSampling"
      @update:interval="updateInterval"
    />
    <v-layout
      class="pa-3"
      column
    >
      <chart-loader
        v-if="aggregatedMetricsPending"
        :has-data="hasMetrics"
      />
      <pie-chart-metrics
        v-if="hasMetrics"
        :chart-id="widget._id"
        :metrics="preparedMetrics"
        :title="widget.parameters.chart_title"
        :show-mode="widget.parameters.show_mode"
        :downloading="downloading"
        @export:png="exportMetricsAsPng"
        @export:csv="exportMetricsAsCsv"
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import { convertFilterToQuery } from '@/helpers/entities/shared/query';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';
import { permissionsWidgetsPieChartInterval } from '@/mixins/permissions/widgets/chart/pie/interval';
import { permissionsWidgetsPieChartSampling } from '@/mixins/permissions/widgets/chart/pie/sampling';
import { permissionsWidgetsPieChartFilters } from '@/mixins/permissions/widgets/chart/pie/filters';

import KpiWidgetFilters from '../partials/kpi-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import PieChartMetrics from './partials/pie-chart-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$system'],
  components: {
    KpiWidgetFilters,
    ChartLoader,
    PieChartMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    widgetChartMetricsMap,
    entitiesAggregatedMetricsMixin,
    permissionsWidgetsPieChartInterval,
    permissionsWidgetsPieChartSampling,
    permissionsWidgetsPieChartFilters,
    widgetChartExportMixinCreator({
      createExport: 'createKpiAlarmAggregateExport',
      fetchExport: 'fetchMetricExport',
    }),
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
  data() {
    return {
      widgetMetricsMap: {},
    };
  },
  computed: {
    hasMetrics() {
      return !!this.aggregatedMetrics.length;
    },

    preparedMetrics() {
      return this.aggregatedMetrics.map((metric) => {
        const parameters = this.widgetMetricsMap[metric.title] ?? {};

        return {
          ...metric,

          color: parameters.color,
          label: parameters.label,
        };
      });
    },
  },
  created() {
    this.setWidgetMetricsMap();
  },
  methods: {
    ...mapMetricsActions({
      createKpiAlarmAggregateExport: 'createKpiAlarmAggregateExport',
      fetchMetricExport: 'fetchMetricExport',
    }),

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'sampling']),
        widget_filters: convertFilterToQuery(this.query.filter),
      };
    },

    async fetchList() {
      await this.fetchAggregatedMetricsList({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });

      this.setWidgetMetricsMap();
    },
  },
};
</script>
