<template>
  <bar-chart :datasets="datasets" :options="chartOptions">
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

<script setup>
import { computed } from 'vue';

import { COLORS } from '@/config';
import {
  REMEDIATION_STATISTICS_BAR_PERCENTAGE,
  REMEDIATION_STATISTICS_CHART_DATA_TYPE,
  SAMPLINGS,
} from '@/constants';

import { colorToRgba } from '@/helpers/color';
import { getDateLabelBySampling } from '@/helpers/entities/metric/list';

import { useI18n } from '@/hooks/i18n';

import BarChart from '@/components/common/chart/bar-chart.vue';
import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

const props = defineProps({
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
  downloading: {
    type: Boolean,
    default: false,
  },
  minDate: {
    type: Number,
    required: false,
  },
});

defineEmits(['export:csv', 'export:png']);

const { t } = useI18n();

const isPercentDataType = computed(
  () => props.dataType === REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent,
);
const maxValueByType = computed(() => (isPercentDataType.value ? 100 : null));

const datasets = computed(() => {
  if (isPercentDataType.value) {
    const ratio = props.metrics.map(metric => ({
      x: metric.timestamp * 1000,
      y: metric.ratio,
    }));

    return [
      {
        backgroundColor: colorToRgba(COLORS.metrics.remediationStatisticRatioRemediations),
        barPercentage: REMEDIATION_STATISTICS_BAR_PERCENTAGE,
        label: t('common.percent'),
        order: 1,
        data: ratio,
      },
    ];
  }

  const { assigned, executed } = props.metrics.reduce(
    (acc, metric) => {
      const x = metric.timestamp * 1000;

      acc.assigned.push({ x, y: metric.assigned, assigned: true });
      acc.executed.push({ x, y: metric.executed });

      return acc;
    },
    {
      assigned: [],
      executed: [],
    },
  );

  return [
    {
      backgroundColor: colorToRgba(COLORS.metrics.remediationStatisticAssignedRemediations),
      barPercentage: REMEDIATION_STATISTICS_BAR_PERCENTAGE,
      label: t('remediation.statistic.labels.withAssignedRemediations'),
      order: 2,
      data: assigned,
    },
    {
      backgroundColor: colorToRgba(COLORS.metrics.remediationStatisticExecutedRemediations),
      barPercentage: REMEDIATION_STATISTICS_BAR_PERCENTAGE,
      label: t('remediation.statistic.labels.remediated'),
      order: 1,
      data: executed,
    },
  ];
});

const getChartValueTickLabel = value => `${value}${isPercentDataType.value ? '%' : ''}`;
const getChartTimeTickLabel = (_, index, data) => {
  const { value } = data[index] ?? {};

  return getDateLabelBySampling(value, props.sampling).split('\n');
};

const getChartTooltipTitle = (data) => {
  const [dataset] = data;
  const { x: timestamp } = dataset.raw;

  return getDateLabelBySampling(timestamp, props.sampling);
};

const getChartTooltipLabel = ({ raw }) => {
  if (isPercentDataType.value) {
    return t('remediation.statistic.tooltips.remediated', { value: `${raw.y}%` });
  }

  return t(`remediation.statistic.tooltips.${raw.assigned ? 'assigned' : 'remediated'}`, {
    value: raw.y,
  });
};

const chartOptions = computed(() => ({
  animation: false,
  responsive: true,
  scales: {
    x: {
      type: 'time',
      stacked: true,
      ticks: {
        min: props.minDate ? props.minDate * 1000 : undefined,
        max: Date.now(),
        source: 'data',
        font: {
          size: 11,
          family: 'Arial, sans-serif',
        },
        callback: getChartTimeTickLabel,
      },
    },
    y: {
      max: maxValueByType.value,
      beginAtZero: true,
      ticks: {
        font: {
          size: 11,
          family: 'Arial, sans-serif',
        },
        callback: getChartValueTickLabel,
      },
    },
  },
  interaction: {
    intersect: false,
    mode: 'index',
  },
  plugins: {
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
        title: getChartTooltipTitle,
        label: getChartTooltipLabel,
      },
    },
  },
}));
</script>
