<template lang="pug">
  v-layout.py-2(column)
    chart-widget-filters.mx-3(
      :widget-id="widget._id",
      :user-filters="userPreference.filters",
      :widget-filters="widget.filters",
      :locked-value="lockedFilter",
      :filters="mainFilter",
      :interval="query.interval",
      :min-interval-date="minAvailableDate",
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
      chart-loader(v-if="vectorMetricsPending", :has-metrics="hasMetrics")
      line-chart-metrics(
        v-if="hasMetrics",
        :metrics="preparedVectorMetrics",
        :title="widget.parameters.chart_title",
        :sampling="query.sampling",
        :downloading="downloading",
        @export:png="exportMetricsAsPng",
        @export:csv="exportMetricsAsCsv"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { keyBy, pick } from 'lodash';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertFilterToQuery } from '@/helpers/query';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetChartExportMixinCreator } from '@/mixins/widget/chart/export';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { entitiesVectorMetricsMixin } from '@/mixins/entities/vector-metrics';
import { permissionsWidgetsLineChartInterval } from '@/mixins/permissions/widgets/chart/line/interval';
import { permissionsWidgetsLineChartSampling } from '@/mixins/permissions/widgets/chart/line/sampling';
import { permissionsWidgetsLineChartFilters } from '@/mixins/permissions/widgets/chart/line/filters';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

import ChartLoader from './partials/chart-loader.vue';
import LineChartMetrics from './partials/line-chart-metrics.vue';

const { mapActions: mapMetricsActions } = createNamespacedHelpers('metrics');

export default {
  inject: ['$system'],
  components: {
    ChartWidgetFilters,
    ChartLoader,
    LineChartMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    entitiesVectorMetricsMixin,
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
      return !!this.vectorMetrics.length;
    },

    minAvailableDate() {
      const { min_date: minDate } = this.vectorMetricsMeta;

      return minDate
        ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
        : null;
    },

    preparedVectorMetrics() {
      return this.vectorMetrics.map(metric => ({
        ...metric,

        color: this.widgetMetricsMap[metric.title].color,
      }));
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
  created() {
    this.setWidgetMetricsMap();
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

    setWidgetMetricsMap() {
      this.widgetMetricsMap = keyBy(this.widget.parameters?.metrics ?? [], 'metric');
    },

    async fetchList() {
      await this.fetchVectorMetricsList({ params: this.getQuery(), widgetId: this.widget._id });

      this.setWidgetMetricsMap();
    },
  },
};
</script>