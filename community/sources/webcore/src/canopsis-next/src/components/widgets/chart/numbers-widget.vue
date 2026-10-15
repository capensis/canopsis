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

      <numbers-metrics
        v-if="hasMetrics"
        :metrics="preparedMetrics"
        :title="widget.parameters.chart_title"
        :show-trend="widget.parameters.show_trend"
        :font-size="valueFontSize"
        :downloading="downloading"
        :hide-actions="hideActions"
        @export:csv="exportMetricsAsCsv"
      />
    </c-layout>
  </c-layout>
</template>

<script>
import { pick } from 'lodash';
import { mapActions, mapState } from 'pinia';

import {
  NUMBERS_CHART_DEFAULT_FONT_SIZE,
  NUMBERS_CHART_FONT_SIZE_WIDTH_COEFFICIENT,
  NUMBERS_CHART_MAX_AUTO_FONT_SIZE,
  NUMBERS_CHART_MIN_AUTO_FONT_SIZE,
} from '@/constants';

import { useMetricsStore } from '@/stores/modules/metrics';

import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';
import { convertFiltersToQuery } from '@/helpers/entities/shared/query';

import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';
import { permissionsWidgetsNumbersInterval } from '@/mixins/permissions/widgets/chart/numbers/interval';
import { permissionsWidgetsNumbersSampling } from '@/mixins/permissions/widgets/chart/numbers/sampling';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import KpiWidgetFilters from '@/components/widgets/partials/kpi-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import NumbersMetrics from './partials/numbers-metrics.vue';

export default {
  inject: ['$system'],

  components: {
    ChartLoader,
    KpiWidgetFilters,
    NumbersMetrics,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetPeriodicRefreshMixin,
    widgetChartMetricsMap,
    permissionsWidgetsNumbersInterval,
    permissionsWidgetsNumbersSampling,
    permissionsWidgetsNumbersFilters,
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
      containerWidth: null,
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

          label: parameters.label,
        };
      });
    },

    valueFontSize() {
      if (this.widget.parameters.font_size) {
        return this.widget.parameters.font_size;
      }

      if (this.containerWidth) {
        const size = Math.round(this.containerWidth / NUMBERS_CHART_FONT_SIZE_WIDTH_COEFFICIENT);

        return Math.max(
          Math.min(size, NUMBERS_CHART_MAX_AUTO_FONT_SIZE),
          NUMBERS_CHART_MIN_AUTO_FONT_SIZE,
        );
      }

      return NUMBERS_CHART_DEFAULT_FONT_SIZE;
    },

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },
  },

  created() {
    this.setWidgetMetricsMap();

    this.resizeObserver = new ResizeObserver(this.setElementWidth);
  },

  mounted() {
    this.resizeObserver.observe(this.$el);
    this.setElementWidth();
  },

  beforeUnmount() {
    this.resizeObserver.unobserve(this.$el);
    this.resizeObserver.disconnect();
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

    setElementWidth() {
      if (this.widget.parameters.font_size) {
        return;
      }

      const { width } = this.$el.getBoundingClientRect();

      this.containerWidth = width;
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'sampling', 'with_history']),
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
