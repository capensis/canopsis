<template>
  <v-layout
    class="py-2"
    column
  >
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreference.filters"
      :widget-filters="widget.filters"
      :locked-value="lockedFilter"
      :filters="mainFilter"
      :interval="query.interval"
      :min-interval-date="minAvailableDate"
      :sampling="query.sampling"
      :show-filter="hasAccessToUserFilter"
      :show-interval="hasAccessToInterval"
      :show-sampling="hasAccessToSampling"
      :filter-disabled="!hasAccessToListFilters"
      :filter-addable="hasAccessToAddFilter"
      :filter-editable="hasAccessToEditFilter"
      class="mx-3"
      @update:filters="updateSelectedFilter"
      @update:sampling="updateSampling"
      @update:interval="updateInterval"
    />
    <v-layout
      class="pa-3"
      column
    >
      <chart-loader
        v-if="vectorMetricsPending"
        :has-data="hasMetrics"
      />
      <bar-chart-metrics
        v-if="hasMetrics"
        :chart-id="widget._id"
        :metrics="preparedMetrics"
        :title="widget.parameters.chart_title"
        :sampling="query.sampling"
        :stacked="widget.parameters.stacked"
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

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertFilterToQuery } from '@/helpers/entities/shared/query';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';
import { entitiesVectorMetricsMixin } from '@/mixins/entities/vector-metrics';
import { permissionsWidgetsBarChartInterval } from '@/mixins/permissions/widgets/chart/bar/interval';
import { permissionsWidgetsBarChartSampling } from '@/mixins/permissions/widgets/chart/bar/sampling';
import { permissionsWidgetsBarChartFilters } from '@/mixins/permissions/widgets/chart/bar/filters';

import KpiWidgetFilters from '../partials/kpi-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import BarChartMetrics from './partials/bar-chart-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$system'],
  components: {
    KpiWidgetFilters,
    ChartLoader,
    BarChartMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    widgetChartMetricsMap,
    entitiesVectorMetricsMixin,
    permissionsWidgetsBarChartInterval,
    permissionsWidgetsBarChartSampling,
    permissionsWidgetsBarChartFilters,
    widgetChartExportMixinCreator({
      createExport: 'createKpiAlarmExport',
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
  computed: {
    hasMetrics() {
      return !!this.vectorMetrics.length;
    },

    minAvailableDate() {
      const { min_date: minDate } = this.vectorMetricsMeta;

      return minDate
        ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
        : null;
    },

    preparedMetrics() {
      return convertMetricsToTimezone(this.vectorMetrics, this.$system.timezone).map((metric) => {
        const parameters = this.widgetMetricsMap[metric.title] ?? {};

        return {
          ...metric,

          color: parameters.color,
          label: parameters.label,
        };
      });
    },
  },
  watch: {
    minAvailableDate() {
      const { from } = this.getIntervalQuery();

      if (from < this.minAvailableDate) {
        this.query = {
          ...this.query,
          interval: {
            from: this.minAvailableDate,
            to: this.query.interval.to,
          },
        };
      }
    },
  },
  methods: {
    ...mapMetricsActions({
      createKpiAlarmExport: 'createKpiAlarmExport',
      fetchMetricExport: 'fetchMetricExport',
    }),

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'sampling', 'with_history']),
        widget_filters: convertFilterToQuery(this.query.filter),
      };
    },

    async fetchList() {
      await this.fetchVectorMetricsList({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });

      this.setWidgetMetricsMap();
    },
  },
};
</script>
