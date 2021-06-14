<template lang="pug">
  div
    v-layout.position-relative(justify-center)
      horizontal-bar(
        ref="horizontalBar",
        :labels="labels",
        :datasets="datasets",
        :options="options",
        :width="width",
        :height="height"
      )
      div.v-tooltip__content.menuable__content__active(ref="tooltip")
    c-table-pagination(
      :total-items="totalItems",
      :rows-per-page="query.rowsPerPage",
      :rows-per-page-items="rowsPerPageItems",
      :page="query.page",
      @update:page="updatePage",
      @update:rows-per-page="updateRowsPerPage"
    )
</template>

<script>
import { get } from 'lodash';

import { TEST_SUITE_COLORS, TEST_SUITE_STATUSES } from '@/constants';

import HorizontalBar from './horizontal-bar.vue';

/**
 * Local constants
 */
const ICON_NAMES = {
  done: 'done',
  close: 'close',
  arrowUpward: 'arrow_upward',
  arrowDownward: 'arrow_downward',
};

const ICONS_COLORS = {
  [ICON_NAMES.done]: 'primary',
  [ICON_NAMES.close]: 'error',
  [ICON_NAMES.arrowUpward]: 'primary',
  [ICON_NAMES.arrowDownward]: 'error',
};

const FAILED_STATUSES = [
  TEST_SUITE_STATUSES.error,
  TEST_SUITE_STATUSES.skipped,
  TEST_SUITE_STATUSES.failed,
];

const MIN_SKIPPED_TO_SECTIONS_COUNT = 40;

export default {
  components: { HorizontalBar },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    query: {
      type: Object,
      default: () => ({}),
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    historical: {
      type: Boolean,
      default: false,
    },
    width: {
      type: Number,
      required: false,
    },
    height: {
      type: Number,
      required: false,
    },
  },
  computed: {
    rowsPerPageItems() {
      return [5, 10, 20];
    },

    labels() {
      return this.items.map(({ name }) => name);
    },

    /**
     * We've using this value if test case was skipped with 0 time value.
     * Otherwise we will not have possibility to see this test case in the chart.
     *
     * @returns {number}
     */
    minSkippedTo() {
      const firstFrom = get(this.items, [0, 'from'], 0);
      const lastTo = get(this.items, [this.items.length - 1, 'to'], 0);

      return (lastTo - firstFrom) / MIN_SKIPPED_TO_SECTIONS_COUNT;
    },

    datasets() {
      const defaultDatasetParameters = {
        barPercentage: 0.9,
        categoryPercentage: 1,
      };

      const { items, historical, minSkippedTo } = this;

      if (!items.length) {
        return [];
      }

      const datasets = [
        {
          ...defaultDatasetParameters,

          backgroundColor: items.map(({ status }) => TEST_SUITE_COLORS[status]),
          data: items.map(({ from, time, to }) => [from, !time ? from + minSkippedTo : to]),
        },
      ];

      if (historical) {
        datasets.unshift({
          ...defaultDatasetParameters,

          backgroundColor: 'rgba(0, 0, 0, .2)',
          data: items.map(({ from, avg_to: avgTo }) => [from, avgTo]),
        });
      }

      return datasets;
    },

    options() {
      return {
        animation: false,
        responsive: false,
        maintainAspectRatio: false,
        legend: { display: false },
        scales: {
          yAxes: [{
            stacked: true,
          }],
          xAxes: [{
            position: 'top',
            ticks: {
              callback: value => `${value}s`,
            },
            suggestedMin: 0,
          }],
        },
        tooltips: {
          enabled: false,
          custom: this.getTooltip,
        },
      };
    },
  },
  methods: {
    getTooltipContent(tooltipModel) {
      const { dataPoints: [dataPoint] } = tooltipModel;

      const { time = 0, status, message } = this.items[dataPoint.index];

      const timeDiv = `<div>${time.toFixed(3)}s</div>`;
      const statusDiv = `<div>${this.$t(`testSuite.statuses.${status}`)}${message ? `: ${message}` : ''}</div>`;

      return `<div>${timeDiv}${statusDiv}</div>`;
    },

    getHistoricalTooltipIcon(item) {
      const {
        time = 0,
        status,
        avg_time: avgTime = 0,
        avg_status: avgStatus,
      } = item;

      let icon;

      if (status === TEST_SUITE_STATUSES.passed) {
        if (avgStatus === TEST_SUITE_STATUSES.passed) {
          if (time < avgTime) {
            icon = ICON_NAMES.arrowUpward;
          } else if (time > avgTime) {
            icon = ICON_NAMES.arrowDownward;
          }
        } else if (FAILED_STATUSES.includes(avgStatus)) {
          icon = ICON_NAMES.done;
        }
      } else if (avgStatus === TEST_SUITE_STATUSES.passed && FAILED_STATUSES.includes(status)) {
        icon = ICON_NAMES.close;
      }

      if (icon) {
        const iconClass = `v-icon material-icons ${ICONS_COLORS[icon]}--text`;

        return `<i class="${iconClass}" style="font-size: 16px;">${icon}</i>`;
      }

      return '';
    },

    getHistoricalTooltipContent(tooltipModel) {
      const { dataPoints: [dataPoint] } = tooltipModel;
      const item = this.items[dataPoint.index] || {};
      const {
        time = 0,
        status,
        avg_time: avgTime = 0,
        avg_status: avgStatus,
      } = item;

      const fixedTime = time.toFixed(3);
      const fixedAvgTime = avgTime.toFixed(3);

      const icon = `&nbsp;${this.getHistoricalTooltipIcon(item)}`;
      const currentDiv =
        `<div>${this.$t('common.current')}: ${this.$t(`testSuite.statuses.${status}`)} ${fixedTime}s${icon}</div>`;
      const averageDiv =
        `<div>${this.$t('common.average')}: ${this.$t(`testSuite.statuses.${avgStatus}`)} ${fixedAvgTime}s</div>`;

      return `<div>${currentDiv}${averageDiv}</div>`;
    },

    getTooltip(tooltipModel) {
      const { tooltip: tooltipEl, horizontalBar: horizontalBarEl } = this.$refs;

      if (tooltipModel.opacity === 0) {
        tooltipEl.style.opacity = 0;
        return;
      }

      tooltipEl.innerHTML = this.historical
        ? this.getHistoricalTooltipContent(tooltipModel)
        : this.getTooltipContent(tooltipModel);

      const leftOffset = 10;
      const topOffset = 5;

      const position = horizontalBarEl.$refs.canvas.getBoundingClientRect();
      const tooltipPosition = tooltipEl.getBoundingClientRect();
      const top = (tooltipModel.caretY + (tooltipPosition.height / 2)) - topOffset;
      let left = tooltipModel.caretX + leftOffset;

      if (left + tooltipPosition.width > position.width) {
        left -= tooltipPosition.width + (leftOffset * 2);
      }

      tooltipEl.style.opacity = 1;
      tooltipEl.style.left = `${left}px`;
      tooltipEl.style.top = `${top}px`;
      tooltipEl.style.pointerEvents = 'none';
    },

    updateRowsPerPage(rowsPerPage) {
      this.$emit('update:query', { ...this.query, rowsPerPage, page: 1 });
    },

    updatePage(page) {
      this.$emit('update:query', { ...this.query, page });
    },
  },
};
</script>

<style lang="scss" scoped>
.v-tooltip__content {
  position: absolute;
  opacity: 0;
  transition: all .2s linear;
}
</style>
