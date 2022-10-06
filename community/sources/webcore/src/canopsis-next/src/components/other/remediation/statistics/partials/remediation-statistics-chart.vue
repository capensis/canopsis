<template lang="pug">
  bar-chart(:datasets="datasets", :options="chartOptions")
</template>

<script>
import { COLORS } from '@/config';
import {
  SAMPLINGS,
  TIME_UNITS,
  REMEDIATION_STATISTICS_BAR_PERCENTAGE,
  REMEDIATION_STATISTICS_CHART_DATA_TYPE,
} from '@/constants';

import { colorToRgba } from '@/helpers/color';
import { getDateLabelBySampling } from '@/helpers/metrics';

import BarChart from '@/components/common/chart/bar-chart.vue';

export default {
  components: { BarChart },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    dataType: {
      type: String,
      default: REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent,
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    responsive: {
      type: Boolean,
      default: false,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },
  computed: {
    isPercentDataType() {
      return this.dataType === REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent;
    },

    maxValueBySampling() {
      return {
        [SAMPLINGS.hour]: 60,
        [SAMPLINGS.day]: 24,
        [SAMPLINGS.week]: 7,
        [SAMPLINGS.month]: 31,
      }[this.sampling];
    },

    maxValueByType() {
      if (this.isPercentDataType) {
        return 100;
      }

      return this.maxValueBySampling;
    },

    samplingUnit() {
      return {
        [SAMPLINGS.hour]: TIME_UNITS.minute,
        [SAMPLINGS.day]: TIME_UNITS.hour,
        [SAMPLINGS.week]: TIME_UNITS.day,
        [SAMPLINGS.month]: TIME_UNITS.day,
      }[this.sampling];
    },

    datasets() {
      if (this.isPercentDataType) {
        const ratio = this.metrics.map(metric => ({
          x: metric.timestamp * 1000,
          y: metric.ratio,
        }));

        return [{
          backgroundColor: colorToRgba(COLORS.metrics.remediationStatisticRatioRemediations),
          barPercentage: REMEDIATION_STATISTICS_BAR_PERCENTAGE,
          label: this.$tc('common.percent', 1),
          order: 1,
          data: ratio,
        }];
      }

      const { assigned, executed } = this.metrics.reduce((acc, metric) => {
        const x = metric.timestamp * 1000;

        acc.assigned.push({ x, y: metric.assigned, assigned: true });
        acc.executed.push({ x, y: metric.executed });

        return acc;
      }, {
        assigned: [],
        executed: [],
      });

      return [{
        backgroundColor: colorToRgba(COLORS.metrics.remediationStatisticAssignedRemediations),
        barPercentage: REMEDIATION_STATISTICS_BAR_PERCENTAGE,
        label: this.$t('remediationStatistic.labels.remediated'),
        order: 2,
        data: assigned,
      }, {
        backgroundColor: colorToRgba(COLORS.metrics.remediationStatisticExecutedRemediations),
        barPercentage: REMEDIATION_STATISTICS_BAR_PERCENTAGE,
        label: this.$t('remediationStatistic.labels.notRemediated'),
        order: 1,
        data: executed,
      }];
    },

    chartOptions() {
      return {
        animation: false,
        responsive: this.responsive,
        scales: {
          x: {
            type: 'time',
            stacked: true,
            ticks: {
              min: this.minDate * 1000,
              max: Date.now(),
              source: 'data',
              font: {
                size: 11,
                family: 'Arial, sans-serif',
              },
              callback: this.getChartTimeTickLabel,
            },
          },
          y: {
            max: this.maxValueByType,
            beginAtZero: true,
            ticks: {
              font: {
                size: 11,
                family: 'Arial, sans-serif',
              },
              callback: this.getChartValueTickLabel,
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'index',
        },
        plugins: {
          background: {
            color: 'white',
          },
          legend: {
            position: 'top',
            align: 'end',
            labels: {
              boxWidth: 20,
              boxHeight: 20,
            },
          },
          tooltip: {
            itemSort: (a, b) => b.dataset.order - a.dataset.order,
            callbacks: {
              title: this.getChartTooltipTitle,
              label: this.getChartTooltipLabel,
            },
          },
        },
      };
    },
  },
  methods: {
    getChartValueTickLabel(value) {
      return `${value}${this.isPercentDataType ? '%' : ''}`;
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

    getChartTooltipLabel(tooltip) {
      const { raw } = tooltip;

      if (this.isPercentDataType) {
        return this.$t('remediationStatistic.tooltips.remediated', { value: `${raw.y}%` });
      }

      return this.$t(`remediationStatistic.tooltips.${raw.assigned ? 'assigned' : 'remediated'}`, { value: raw.y });
    },
  },
};
</script>