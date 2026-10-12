<template>
  <c-layout class="py-2" column>
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreferenceFilters"
      :widget-filters="widget.filters"
      :locked-filters="query.lockedFilter"
      :filters="query.filter"
      :interval="query.interval"
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
        v-if="aggregatedMetricsPending"
        :has-data="hasMetrics"
      />

      <pie-chart-metrics
        v-if="hasMetrics"
        :metrics="preparedMetrics"
        :title="widget.parameters.chart_title"
        :show-mode="widget.parameters.show_mode"
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

import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';
import { convertFiltersToQuery } from '@/helpers/entities/shared/query';

import { permissionsWidgetsPieChartFilters } from '@/mixins/permissions/widgets/chart/pie/filters';
import { permissionsWidgetsPieChartInterval } from '@/mixins/permissions/widgets/chart/pie/interval';
import { permissionsWidgetsPieChartSampling } from '@/mixins/permissions/widgets/chart/pie/sampling';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import KpiWidgetFilters from '@/components/widgets/partials/kpi-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import PieChartMetrics from './partials/pie-chart-metrics.vue';

export default {
  inject: ['$system'],

  components: {
    ChartLoader,
    KpiWidgetFilters,
    PieChartMetrics,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetPeriodicRefreshMixin,
    widgetChartMetricsMap,
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
      getAggregatedMetricsListByWidgetId: 'getListByWidgetId',
      getAggregatedMetricsPendingByWidgetId: 'getPendingByWidgetId',
    }),

    aggregatedMetrics() {
      return this.getAggregatedMetricsListByWidgetId(this.widget._id);
    },

    aggregatedMetricsPending() {
      return this.getAggregatedMetricsPendingByWidgetId(this.widget._id);
    },

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

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },
  },

  methods: {
    ...mapActions(useMetricsStore, {
      createKpiAlarmAggregateExport: 'createKpiAlarmAggregateExport',
      fetchMetricExport: 'fetchMetricExport',
      fetchAggregatedMetrics: 'fetchAggregatedMetrics',
    }),

    customQueryCondition(query, oldQuery) {
      return isMetricsQueryChanged(query, oldQuery);
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
        ...pick(this.query, ['parameters', 'sampling']),
        widget_filters: convertFiltersToQuery(this.query.filter, this.query.lockedFilter),
      };
    },

    async fetchList() {
      await this.fetchAggregatedMetrics({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });

      this.setWidgetMetricsMap();
    },
  },
};
</script>
