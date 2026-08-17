<template>
  <c-layout class="entity-chart-widget py-2 position-relative" column>
    <c-progress-overlay :pending="pending" />
    <div class="text-subtitle-1 font-weight-medium px-3">
      {{ title }}
    </div>
    <v-table v-if="metrics.length" density="compact">
      <thead>
        <tr>
          <th>Metric</th>
          <th>{{ $t('common.value') }}</th>
          <th>{{ $t('common.date') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="metric in metrics" :key="`${metric.title}-${metric.timestamp}`">
          <td>{{ metric.label || metric.title }}</td>
          <td>{{ metric.value ?? metric.val ?? '-' }}</td>
          <td>{{ dateToString(metric.timestamp) }}</td>
        </tr>
      </tbody>
    </v-table>
    <div v-else-if="!pending" class="text-medium-emphasis px-3">
      {{ $t('common.noData') }}
    </div>
  </c-layout>
</template>

<script>
import pick from 'lodash/pick';
import { mapActions } from 'pinia';

import { DATETIME_FORMATS, QUICK_RANGES, WIDGET_TYPES } from '@/constants';

import { useMetricsStore } from '@/stores/modules/metrics';

import { convertDateToString } from '@/helpers/date/date';
import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';

import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
  mixins: [entitiesInfoMixin],

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
    const range = QUICK_RANGES[this.widget.parameters.default_time_range] ?? QUICK_RANGES.last30Days;

    return {
      metrics: [],
      pending: false,
      query: {
        interval: {
          from: range.start,
          to: range.stop,
        },
        sampling: this.widget.parameters.default_sampling,
        parameters: this.widget.parameters.metrics ?? [],
        with_history: this.widget.parameters.comparison,
      },
    };
  },

  computed: {
    title() {
      return this.widget.parameters.chart_title || this.widget.title;
    },

    isVectorChart() {
      return [WIDGET_TYPES.barChart, WIDGET_TYPES.lineChart].includes(this.widget.type);
    },
  },

  watch: {
    widget: {
      deep: true,
      handler() {
        this.fetchList();
      },
    },
  },

  mounted() {
    this.fetchList();
  },

  methods: {
    ...mapActions(useMetricsStore, [
      'fetchEntityAlarmsMetricsWithoutStore',
      'fetchEntityAggregateMetricsWithoutStore',
    ]),

    dateToString(value) {
      return value ? convertDateToString(value, DATETIME_FORMATS.long, '-') : '-';
    },

    getQuery() {
      return {
        ...convertMetricIntervalToTimestamp({
          interval: this.query.interval,
          timezone: this.timezone,
        }),
        ...pick(this.query, ['sampling', 'with_history']),
        parameters: this.query.parameters.filter(({ metric }) =>
          this.availableMetrics.includes(metric),
        ),
        entity: this.entity._id,
      };
    },

    async fetchList() {
      this.pending = true;

      try {
        const fetchList = this.isVectorChart
          ? this.fetchEntityAlarmsMetricsWithoutStore
          : this.fetchEntityAggregateMetricsWithoutStore;
        const { data } = await fetchList({ params: this.getQuery() });

        this.metrics = convertMetricsToTimezone(data ?? [], this.timezone).map((metric) => {
          const parameters = this.widget.parameters.metrics?.find(({ metric: name }) => name === metric.title);

          return {
            ...metric,
            color: parameters?.color,
            label: parameters?.label,
          };
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
