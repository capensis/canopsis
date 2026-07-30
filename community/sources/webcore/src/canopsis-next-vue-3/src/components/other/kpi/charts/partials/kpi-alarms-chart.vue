<template>
  <bar-chart :datasets="datasets" :options="alarmsChartOptions" :dark="$system.dark">
    <template #actions="{ chart }">
      <chart-export-actions
        :downloading="downloading"
        :chart="chart"
        class="mt-4"
        @export:csv="$emit('export:csv')"
        @export:png="$emit('export:png', $event)"
      />
    </template>
  </bar-chart>
</template>

<script>
import { debounce } from 'lodash';

import { KPI_ALARMS_GRAPH_BAR_PERCENTAGE, SAMPLINGS, TIME_UNITS } from '@/constants';

import { convertDateToStartOfUnitTimestamp, getNowTimestamp } from '@/helpers/date/date';
import {
  convertDurationToMaxUnitDuration,
  convertDurationToString,
  fromSeconds,
} from '@/helpers/date/duration';
import { getMetricColor } from '@/helpers/entities/metric/color';
import { getDateLabelBySampling } from '@/helpers/entities/metric/list';
import { isRatioMetric, isTimeMetric } from '@/helpers/entities/metric/form';

import BarChart from '@/components/common/chart/bar-chart.vue';
import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

export default {
  inject: ['$system'],

  components: {
    BarChart,
    ChartExportActions,
  },

  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    interval: {
      type: Object,
      default: () => ({
        from: 0,
        to: getNowTimestamp(),
      }),
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    responsive: {
      type: Boolean,
      default: false,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },

  emits: ['export:csv', 'export:png', 'zoom'],

  computed: {
    labelsFont() {
      return {
        size: 11,
        family: 'Arial, sans-serif',
      };
    },

    maxDurationsByMetric() {
      return this.metrics.reduce((acc, { title, data }) => {
        if (isTimeMetric(title)) {
          const values = data.map(({ value }) => value);
          const maxValue = Math.max(...values, 0);

          acc[title] = convertDurationToMaxUnitDuration({
            value: maxValue,
            unit: TIME_UNITS.second,
          });
        }

        return acc;
      }, {});
    },

    datasets() {
      return this.metrics.map(({ title: metric, data }) => ({
        metric,
        backgroundColor: getMetricColor(metric),
        barPercentage: KPI_ALARMS_GRAPH_BAR_PERCENTAGE,
        yAxisID: this.getMetricYAxisId(metric),
        label: this.$t(`alarm.metrics.${metric}`),
        data: data.map(({ timestamp, value }) => ({
          x: timestamp * 1000,
          y: this.convertValueByMetric(metric, value),
        })),
      }));
    },

    yAxes() {
      return this.metrics.reduce((acc, { title: metric }, index) => {
        const yAxisID = this.getMetricYAxisId(metric);

        if (acc[yAxisID]) {
          return acc;
        }

        acc[yAxisID] = {
          type: 'linear',
          display: index === 0,
          position: index % 2 ? 'right' : 'left',
          beginAtZero: true,
          max: isRatioMetric(metric) ? 100 : undefined,
          ticks: {
            font: this.labelsFont,
            callback: value => this.getChartYTick(metric, value),
          },
        };

        return acc;
      }, {});
    },

    alarmsChartOptions() {
      return {
        animation: false,
        responsive: this.responsive,
        scales: {
          x: {
            type: 'time',
            min: this.interval.from * 1000,
            max: this.interval.to * 1000,
            ticks: {
              source: 'data',
              font: this.labelsFont,
              callback: this.getChartTimeTickLabel,
            },
          },
          ...this.yAxes,
        },
        interaction: {
          intersect: false,
          mode: 'x',
        },
        plugins: {
          legend: {
            position: 'right',
            align: 'start',
            maxWidth: 600,
            labels: {
              font: this.labelsFont,
              boxWidth: 15,
              boxHeight: 15,
            },
          },
          tooltip: {
            displayColors: false,
            callbacks: {
              title: this.getChartTooltipTitle,
              label: this.getChartTooltipLabel,
            },
          },
          zoom: {
            limits: {
              x: {
                min: this.minDate * 1000,
                max: convertDateToStartOfUnitTimestamp(getNowTimestamp(), TIME_UNITS.hour) * 1000,
              },
            },
            pan: {
              enabled: true,
              mode: 'x',
              threshold: 3,
              onPanComplete: this.updateChartInterval,
            },
            zoom: {
              mode: 'x',
              wheel: {
                enabled: true,
                speed: 0.1,
              },
              pinch: {
                enabled: true,
              },
              drag: {
                enabled: true,
                modifierKey: 'ctrl',
              },
              onZoomComplete: this.updateChartInterval,
            },
          },
        },
      };
    },
  },

  created() {
    this.debouncedUpdateInterval = debounce(this.updateInterval, 300);
  },

  methods: {
    getMetricYAxisId(metric) {
      if (isTimeMetric(metric)) {
        return 'time';
      }

      if (isRatioMetric(metric)) {
        return 'ratio';
      }

      return 'count';
    },

    convertValueByMetric(metric, value) {
      if (!isTimeMetric(metric)) {
        return value;
      }

      return fromSeconds(value, this.maxDurationsByMetric[metric]?.unit);
    },

    getChartYTick(metric, value) {
      if (isTimeMetric(metric)) {
        return `${value}${this.maxDurationsByMetric[metric]?.unit}`;
      }

      if (isRatioMetric(metric)) {
        return `${value}%`;
      }

      return value;
    },

    getChartTimeTickLabel(_, index, data) {
      const { value } = data[index] ?? {};

      return getDateLabelBySampling(value, this.sampling).split('\n');
    },

    getChartTooltipTitle(data) {
      const [dataset] = data;
      const { x: timestamp } = dataset.raw;

      return getDateLabelBySampling(timestamp, this.sampling);
    },

    getChartTooltipLabel({ raw, dataset }) {
      const { metric } = dataset;
      const value = isTimeMetric(metric)
        ? convertDurationToString(raw.y, undefined, this.maxDurationsByMetric[metric]?.unit)
        : raw.y;

      return this.$t(`kpi.metrics.tooltip.${metric}`, { value });
    },

    updateChartInterval({ chart }) {
      const { min, max } = chart.scales.x;

      this.debouncedUpdateInterval({
        from: Math.floor(min / 1000),
        to: Math.ceil(max / 1000),
      });
    },

    updateInterval(interval) {
      this.$emit('zoom', interval);
    },
  },
};
</script>
