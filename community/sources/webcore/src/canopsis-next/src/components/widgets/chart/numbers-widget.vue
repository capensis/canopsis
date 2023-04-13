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
        :metrics="preparedAggregatedMetrics",
        :title="widget.parameters.chart_title",
        :show-trend="widget.parameters.show_trend",
        :font-size="valueFontSize",
        :downloading="downloading",
        @export:csv="exportMetricsAsCsv"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { keyBy, pick } from 'lodash';

import {
  NUMBERS_CHART_MAX_AUTO_FONT_SIZE,
  NUMBERS_CHART_FONT_SIZE_WIDTH_COEFFICIENT,
  NUMBERS_CHART_DEFAULT_FONT_SIZE,
  NUMBERS_CHART_MIN_AUTO_FONT_SIZE,
} from '@/constants';

import { convertFilterToQuery } from '@/helpers/query';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { entitiesAggregatedMetricsMixin } from '@/mixins/entities/aggregated-metrics';
import { permissionsWidgetsNumbersInterval } from '@/mixins/permissions/widgets/chart/numbers/interval';
import { permissionsWidgetsNumbersSampling } from '@/mixins/permissions/widgets/chart/numbers/sampling';
import { permissionsWidgetsNumbersFilters } from '@/mixins/permissions/widgets/chart/numbers/filters';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import NumbersMetrics from './partials/numbers-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

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
      widgetMetricsMap: {},
    };
  },
  computed: {
    preparedAggregatedMetrics() {
      return this.aggregatedMetrics.map((metric) => {
        const parameters = this.widgetMetricsMap[metric.title];

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
        ...pick(this.query, ['parameters', 'sampling']),
        widget_filters: convertFilterToQuery(this.query.filter),
      };
    },

    setWidgetMetricsMap() {
      this.widgetMetricsMap = keyBy(this.widget.parameters?.metrics ?? [], 'metric');
    },

    fetchList() {
      this.fetchAggregatedMetricsList({
        widgetId: this.widget._id,
        trend: this.widget.parameters.show_trend,
        params: this.getQuery(),
      });

      this.setWidgetMetricsMap();
    },
  },
};
</script>
