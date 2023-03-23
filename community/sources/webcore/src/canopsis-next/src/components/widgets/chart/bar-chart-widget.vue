<template lang="pug">
  v-layout.px-3.pb-3(column)
    chart-widget-filters(
      :widget-id="widget._id",
      :user-filters="userPreference.filters",
      :widget-filters="widget.filters",
      :locked-value="lockedFilter",
      :filters="mainFilter",
      :interval="query.interval",
      :min-interval-date="minAvailableDate",
      :sampling="query.sampling",
      :show-filter="true",
      :show-interval="true",
      :show-sampling="true",
      :filter-disabled="!true",
      :filter-addable="true",
      :filter-editable="true",
      @update:filters="updateSelectedFilter",
      @update:sampling="updateSampling",
      @update:interval="updateInterval"
    )
    v-layout.ma-4(column)
      div.position-relative
        bar-chart-metrics(
          :metrics="preparedVectorMetrics",
          :sampling="query.sampling",
          :stacked="widget.parameters.stacked"
        )
</template>

<script>
import { omit, keyBy } from 'lodash';
import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetIntervalFilterMixin } from '@/mixins/widget/chart/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { entitiesVectorMetricsMixin } from '@/mixins/entities/vector-metrics';
import { permissionsWidgetsBarChartInterval } from '@/mixins/permissions/widgets/chart/bar/interval';
import { permissionsWidgetsBarChartSampling } from '@/mixins/permissions/widgets/chart/bar/sampling';
import { permissionsWidgetsBarChartFilters } from '@/mixins/permissions/widgets/chart/bar/filters';

import ChartWidgetFilters from '@/components/widgets/chart/partials/chart-widget-filters.vue';

import BarChartMetrics from './partials/bar-chart-metrics.vue';

export default {
  inject: ['$system'],
  components: {
    ChartWidgetFilters,
    BarChartMetrics,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetIntervalFilterMixin,
    widgetSamplingFilterMixin,
    entitiesVectorMetricsMixin,
    permissionsWidgetsBarChartInterval,
    permissionsWidgetsBarChartSampling,
    permissionsWidgetsBarChartFilters,
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
    minAvailableDate() {
      const { min_date: minDate } = this.vectorMetricsMeta;

      return minDate
        ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
        : null;
    },

    widgetMetricsMap() {
      return keyBy(this.widget.parameters?.metrics ?? [], 'metric');
    },

    preparedVectorMetrics() {
      return this.vectorMetrics.map(metric => ({
        ...metric,

        color: this.widgetMetricsMap.color,
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
  methods: {
    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...omit(this.query, ['interval']),
      };
    },

    async fetchList() {
      this.fetchVectorMetricsList({ params: this.getQuery(), widgetId: this.widget._id });
    },
  },
};
</script>
