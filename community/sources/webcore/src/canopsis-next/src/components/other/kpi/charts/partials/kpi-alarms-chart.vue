<template>
  <bar-chart
    :datasets="datasets"
    :options="alarmsChartOptions"
    :width="width"
    :height="height"
    :dark="$system.dark"
  >
    <template #actions="{ chart }">
      <chart-export-actions
        :downloading="downloading"
        :chart="chart"
        class="mt-4"
        v-on="$listeners"
      />
    </template>
  </bar-chart>
</template>

<script>
import { debounce } from 'lodash';

import { KPI_ALARMS_GRAPH_BAR_PERCENTAGE, SAMPLINGS, TIME_UNITS } from '@/constants';

import { getMetricColor } from '@/helpers/entities/metric/color';
import { convertDateToStartOfUnitTimestamp, getNowTimestamp } from '@/helpers/date/date';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

import BarChart from '@/components/common/chart/bar-chart.vue';
import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

export default {
  inject: ['$system'],
  components: { ChartExportActions, BarChart },
  mixins: [chartMetricsOptionsMixin],
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
    width: {
      default: 600,
      type: Number,
    },
    height: {
      default: 440,
      type: Number,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },
  computed: {
    datasets() {
      return this.metrics.map(({ title: metric, data }) => ({
        metric,
        backgroundColor: getMetricColor(metric),
        barPercentage: KPI_ALARMS_GRAPH_BAR_PERCENTAGE,
        yAxisID: this.getMetricYAxisId(metric),
        label: this.$t(`alarm.metrics.${metric}`),
        data: data.map(({ timestamp, value }) => ({
          x: timestamp * 1000,
          y: value,
        })),
      }));
    },

    alarmsChartOptions() {
      return {
        animation: false,
        responsive: this.responsive,
        scales: {
          ...this.xAxes,
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
              onZoom: this.updateChartInterval,
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
