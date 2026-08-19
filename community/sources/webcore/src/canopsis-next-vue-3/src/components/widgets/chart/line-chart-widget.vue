<template>
  <c-layout class="py-2" column>
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreferenceFilters"
      :widget-filters="widget.filters"
      :locked-filters="query.lockedFilter"
      :filters="query.filter"
      :interval="query.interval"
      :min-interval-date="minAvailableDate"
      :sampling="query.sampling"
      :show-filter="hasAccessToUserFilter"
      :show-interval="hasAccessToInterval"
      :show-sampling="hasAccessToSampling"
      class="mx-3"
      @update:filters="updateSelectedFilter"
      @update:sampling="updateSampling"
      @update:interval="updateInterval"
    />

    <c-layout class="pa-3" column>
      <chart-loader
        v-if="vectorMetricsPending"
        :has-data="hasMetrics"
      />

      <line-chart-metrics
        v-if="hasMetrics"
        :chart-id="widget._id"
        :metrics="preparedMetrics"
        :title="widget.parameters.chart_title"
        :sampling="query.sampling"
        :downloading="downloading"
        :hide-actions="hideActions"
        @export:png="exportMetricsAsPng"
        @export:csv="exportMetricsAsCsv"
      />
    </c-layout>
  </c-layout>
</template>

<script>
import { pick } from 'lodash';
import { mapActions, mapState } from 'pinia';

import { useMetricsStore } from '@/stores/modules/metrics';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';
import { convertFiltersToQuery } from '@/helpers/entities/shared/query';

import { permissionsWidgetsLineChartFilters } from '@/mixins/permissions/widgets/chart/line/filters';
import { permissionsWidgetsLineChartInterval } from '@/mixins/permissions/widgets/chart/line/interval';
import { permissionsWidgetsLineChartSampling } from '@/mixins/permissions/widgets/chart/line/sampling';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import KpiWidgetFilters from '@/components/widgets/partials/kpi-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import LineChartMetrics from './partials/line-chart-metrics.vue';

export default {
  inject: ['$system'],

  components: {
    ChartLoader,
    KpiWidgetFilters,
    LineChartMetrics,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetPeriodicRefreshMixin,
    widgetChartMetricsMap,
    permissionsWidgetsLineChartInterval,
    permissionsWidgetsLineChartSampling,
    permissionsWidgetsLineChartFilters,
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
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      downloading: false,
    };
  },

  computed: {
    ...mapState(useMetricsStore, {
      getVectorMetricsListByWidgetId: 'getVectorMetricsListByWidgetId',
      getVectorMetricsPendingByWidgetId: 'getVectorMetricsPendingByWidgetId',
      getVectorMetricsMetaByWidgetId: 'getVectorMetricsMetaByWidgetId',
    }),

    vectorMetrics() {
      return this.getVectorMetricsListByWidgetId(this.widget._id);
    },

    vectorMetricsPending() {
      return this.getVectorMetricsPendingByWidgetId(this.widget._id);
    },

    vectorMetricsMeta() {
      return this.getVectorMetricsMetaByWidgetId(this.widget._id);
    },

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

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },
  },

  watch: {
    minAvailableDate() {
      const { from } = this.getIntervalQuery();

      if (this.minAvailableDate && from < this.minAvailableDate) {
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
    ...mapActions(useMetricsStore, {
      createKpiAlarmExport: 'createKpiAlarmExport',
      fetchMetricExport: 'fetchMetricExport',
      fetchVectorMetricsList: 'fetchVectorMetricsList',
    }),

    customQueryCondition(query, oldQuery) {
      return isMetricsQueryChanged(query, oldQuery, this.minAvailableDate);
    },

    updateSampling(sampling) {
      this.updateContentInUserPreference({ sampling });

      this.query = {
        ...this.query,
        sampling,
      };
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'sampling', 'with_history']),
        widget_filters: convertFiltersToQuery(this.query.filter, this.query.lockedFilter),
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
