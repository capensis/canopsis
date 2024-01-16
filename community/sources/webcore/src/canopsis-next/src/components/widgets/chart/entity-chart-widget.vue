<template>
  <v-layout
    class="py-2"
    column
  >
    <kpi-widget-filters
      class="mx-3"
      :interval="query.interval"
      :min-interval-date="minAvailableDate"
      :sampling="query.sampling"
      show-interval
      show-sampling
      @update:sampling="updateQueryField('sampling', $event)"
      @update:interval="updateQueryField('interval', $event)"
    />
    <v-layout
      class="pa-3"
      column
    >
      <chart-loader
        v-if="pending"
        :has-metrics="hasMetrics"
      />
      <component
        v-if="hasMetrics"
        v-bind="component"
        :is="component.is"
        @export:png="exportMetricsAsPng"
        @export:csv="exportMetricsAsCsv"
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick } from 'lodash';

import { WIDGET_TYPES } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertWidgetToQuery } from '@/helpers/entities/widget/query';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';

import { localQueryMixin } from '@/mixins/query-local/query';
import { metricsIntervalFilterMixin } from '@/mixins/widget/metrics/interval';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { entitiesVectorMetricsMixin } from '@/mixins/entities/vector-metrics';
import { widgetChartMetricsMap } from '@/mixins/widget/chart/metrics-map';

import KpiWidgetFilters from '@/components/widgets/partials/kpi-widget-filters.vue';
import ChartLoader from '@/components/widgets/chart/partials/chart-loader.vue';
import BarChartMetrics from '@/components/widgets/chart/partials/bar-chart-metrics.vue';
import LineChartMetrics from '@/components/widgets/chart/partials/line-chart-metrics.vue';
import NumbersMetrics from '@/components/widgets/chart/partials/numbers-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$system'],
  components: {
    ChartLoader,
    KpiWidgetFilters,
    BarChartMetrics,
    LineChartMetrics,
    NumbersMetrics,
  },
  mixins: [
    localQueryMixin,
    metricsIntervalFilterMixin,
    entitiesVectorMetricsMixin,
    widgetChartMetricsMap,
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
    entity: {
      type: Object,
      required: true,
    },
    availableMetrics: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      metrics: [],
      minDate: null,
      pending: false,
    };
  },
  computed: {
    hasMetrics() {
      return !!this.metrics.length;
    },

    isVectorChart() {
      return [WIDGET_TYPES.barChart, WIDGET_TYPES.lineChart].includes(this.widget.type);
    },

    minAvailableDate() {
      return this.minDate
        ? convertDateToStartOfDayTimestampByTimezone(this.minDate, this.$system.timezone)
        : null;
    },

    preparedMetrics() {
      return convertMetricsToTimezone(this.metrics, this.$system.timezone).map((metric) => {
        const parameters = this.widgetMetricsMap[metric.title] ?? {};

        return {
          ...metric,

          color: parameters.color,
          label: parameters.label,
        };
      });
    },

    component() {
      const props = {
        metrics: this.preparedMetrics,
        downloading: this.downloading,
        title: this.widget.parameters.chart_title || this.widget.title,
      };

      if (this.widget.type === WIDGET_TYPES.barChart) {
        return {
          is: 'bar-chart-metrics',
          sampling: this.query.sampling,
          stacked: this.widget.parameters.stacked,
          height: 400,
          ...props,
        };
      }

      if (this.widget.type === WIDGET_TYPES.lineChart) {
        return {
          is: 'line-chart-metrics',
          sampling: this.query.sampling,
          height: 400,
          ...props,
        };
      }

      return {
        is: 'numbers-metrics',
        showTrend: this.widget.parameters.show_trend,
        fontSize: this.widget.parameters.font_size,
        ...props,
      };
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

    widget: {
      immediate: true,
      handler() {
        this.query = convertWidgetToQuery(this.widget);
      },
    },
  },
  created() {
    this.setWidgetMetricsMap();
  },
  methods: {
    ...mapMetricsActions({
      createKpiAlarmExport: 'createKpiAlarmExport',
      fetchMetricExport: 'fetchMetricExport',
      fetchEntityAlarmsMetricsWithoutStore: 'fetchEntityAlarmsMetricsWithoutStore',
      fetchEntityAggregateMetricsWithoutStore: 'fetchEntityAggregateMetricsWithoutStore',
    }),

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['sampling', 'with_history']),
        parameters: this.query.parameters.filter(({ metric }) => this.availableMetrics.includes(metric)),
        entity: this.entity._id,
      };
    },

    async fetchList() {
      this.pending = true;

      try {
        const fetchList = this.isVectorChart
          ? this.fetchEntityAlarmsMetricsWithoutStore
          : this.fetchEntityAggregateMetricsWithoutStore;

        const { data, meta: { min_date: minDate } = {} } = await fetchList({ params: this.getQuery() });

        this.metrics = data;
        this.minDate = minDate;

        this.setWidgetMetricsMap();
      } catch (error) {
        console.error(error);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
