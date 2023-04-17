<template lang="pug">
  v-layout.py-2(column)
    chart-widget-filters.mx-3(
      :interval="query.interval",
      :min-interval-date="minAvailableDate",
      :sampling="query.sampling",
      show-interval,
      show-sampling,
      @update:sampling="updateQueryField('sampling', $event)",
      @update:interval="updateQueryField('interval', $event)"
    )
    v-layout.pa-3(column)
      chart-loader(v-if="pending", :has-metrics="hasMetrics")
      component(
        v-if="hasMetrics",
        v-bind="component",
        @export:png="exportMetricsAsPng",
        @export:csv="exportMetricsAsCsv"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { pick, keyBy } from 'lodash';

import { WIDGET_TYPES } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertWidgetToQuery } from '@/helpers/query';

import { localQueryMixin } from '@/mixins/query-local/query';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { entitiesVectorMetricsMixin } from '@/mixins/entities/vector-metrics';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';
import ChartLoader from '@/components/widgets/chart/partials/chart-loader.vue';
import BarChartMetrics from '@/components/widgets/chart/partials/bar-chart-metrics.vue';
import LineChartMetrics from '@/components/widgets/chart/partials/line-chart-metrics.vue';
import NumbersMetrics from '@/components/widgets/chart/partials/numbers-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');
const { mapActions: mapAggregatedMetricsActions } = createNamespacedHelpers('aggregatedMetrics');

export default {
  inject: ['$system'],
  components: {
    ChartLoader,
    ChartWidgetFilters,
    BarChartMetrics,
    LineChartMetrics,
    NumbersMetrics,
  },
  mixins: [
    localQueryMixin,
    widgetIntervalFilterMixin,
    entitiesVectorMetricsMixin,
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
    alarm: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      widgetMetricsMap: {},
      metrics: [],
      minAvailableDate: null,
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

    preparedMetrics() {
      return this.metrics.map((metric) => {
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
        title: this.widget.parameters.chart_title,
      };

      if (this.widget.type === WIDGET_TYPES.barChart) {
        return {
          is: 'bar-chart-metrics',
          sampling: this.query.sampling,
          stacked: this.widget.parameters.stacked,
          ...props,
        };
      }

      if (this.widget.type === WIDGET_TYPES.lineChart) {
        return {
          is: 'line-chart-metrics',
          sampling: this.query.sampling,
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
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapMetricsActions({
      createKpiAlarmExport: 'createKpiAlarmExport',
      fetchMetricExport: 'fetchMetricExport',
      fetchAlarmsMetricsWithoutStore: 'fetchAlarmsMetricsWithoutStore',
    }),
    ...mapAggregatedMetricsActions({
      fetchAggregatedMetricsWithoutStore: 'fetchListWithoutStore',
    }),

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'sampling', 'with_history']),
        entity: this.alarm.entity._id,
      };
    },

    setWidgetMetricsMap() {
      this.widgetMetricsMap = keyBy(this.widget.parameters?.metrics ?? [], 'metric');
    },

    async fetchList() {
      this.pending = true;

      try {
        const fetchList = this.isVectorChart
          ? this.fetchAlarmsMetricsWithoutStore
          : this.fetchAggregatedMetricsWithoutStore;

        const { data, meta: { min_date: minDate } = {} } = await fetchList({ params: this.getQuery() });

        this.metrics = data;
        this.minAvailableDate = minDate
          ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
          : null;

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
