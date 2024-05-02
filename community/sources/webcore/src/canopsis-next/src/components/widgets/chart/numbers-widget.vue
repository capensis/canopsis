<template>
  <v-layout
    class="py-2"
    column
  >
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreference.filters"
      :widget-filters="widget.filters"
      :locked-filter="query.lockedFilter"
      :filters="query.filter"
      :interval="query.interval"
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
        @export:csv="exportMetricsAsCsv"
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import {
  NUMBERS_CHART_MAX_AUTO_FONT_SIZE,
  NUMBERS_CHART_FONT_SIZE_WIDTH_COEFFICIENT,
  NUMBERS_CHART_DEFAULT_FONT_SIZE,
  NUMBERS_CHART_MIN_AUTO_FONT_SIZE,
} from '@/constants';

import { convertFilterToQuery } from '@/helpers/entities/shared/query';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';
import { permissionsWidgetsNumbersInterval } from '@/mixins/permissions/widgets/chart/numbers/interval';
import { permissionsWidgetsNumbersSampling } from '@/mixins/permissions/widgets/chart/numbers/sampling';
import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';

import KpiWidgetFilters from '../partials/kpi-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import NumbersMetrics from './partials/numbers-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$system'],
  components: {
    KpiWidgetFilters,
    ChartLoader,
    NumbersMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    widgetChartMetricsMap,
    entitiesAggregatedMetricsMixin,
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
    tabId: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      containerWidth: null,
    };
  },
  computed: {
    preparedMetrics() {
      return this.aggregatedMetrics.map((metric) => {
        const parameters = this.widgetMetricsMap[metric.title] ?? {};

        return {
          ...metric,

          label: parameters.label,
        };
      });
    },

    hasMetrics() {
      return !!this.aggregatedMetrics.length;
    },

    valueFontSize() {
      if (this.widget.parameters.font_size) {
        return this.widget.parameters.font_size;
      }

      if (this.containerWidth) {
        const size = Math.round(this.containerWidth / NUMBERS_CHART_FONT_SIZE_WIDTH_COEFFICIENT);

        return Math.max(Math.min(size, NUMBERS_CHART_MAX_AUTO_FONT_SIZE), NUMBERS_CHART_MIN_AUTO_FONT_SIZE);
      }

      return NUMBERS_CHART_DEFAULT_FONT_SIZE;
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
  beforeDestroy() {
    this.resizeObserver.unobserve(this.$el);
    this.resizeObserver.disconnect();
  },
  methods: {
    ...mapMetricsActions({
      createKpiAlarmAggregateExport: 'createKpiAlarmAggregateExport',
      fetchMetricExport: 'fetchMetricExport',
    }),

    setElementWidth() {
      if (this.fontSize) {
        return;
      }

      const { width } = this.$el.getBoundingClientRect();

      this.containerWidth = width;
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'sampling', 'with_history']),
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
