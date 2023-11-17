<template>
  <div>
    <v-layout
      class="position-relative"
      justify-center
    >
      <horizontal-bar
        ref="horizontalBar"
        :labels="labels"
        :datasets="datasets"
        :options="options"
        :width="width"
        :height="height"
      />
      <div
        class="v-tooltip__content menuable__content__active"
        ref="tooltip"
      />
    </v-layout>
    <c-table-pagination
      :total-items="totalItems"
      :items-per-page="query.itemsPerPage"
      :items="itemsPerPageItems"
      :page="query.page"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
    />
  </div>
</template>

<script>
import { get } from 'lodash';

import { TEST_SUITE_COLORS, TEST_SUITE_STATUSES } from '@/constants';

import { colorToRgba } from '@/helpers/color';

const HorizontalBar = () => import(/* webpackChunkName: "Charts" */ './horizontal-bar.vue');

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

const PASSED_STATUSES = [
  TEST_SUITE_STATUSES.passed,
  TEST_SUITE_STATUSES.total,
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
    itemsPerPageItems() {
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

          backgroundColor: colorToRgba('#000', 0.2),
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
        scales: {
          y: {
            stacked: true,
            ticks: {
              font: {
                family: 'Arial, sans-serif',
              },
            },
          },
          x: {
            position: 'top',
            ticks: {
              font: {
                family: 'Arial, sans-serif',
              },
              callback: value => `${value}s`,
            },
            suggestedMin: 0,
          },
        },
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            enabled: false,
            external: this.getTooltip,
          },
        },
      };
    },
  },
  methods: {
    getTooltipContent(tooltipModel) {
      const { dataPoints: [dataPoint] } = tooltipModel;

      const { time = 0, status, message } = this.items[dataPoint.dataIndex];

      const timeDiv = `<div>${time.toFixed(3)}s</div>`;
      const statusDiv = status !== TEST_SUITE_STATUSES.total
        ? `<div>${this.$t(`testSuite.statuses.${status}`)}${message ? `: ${message}` : ''}</div>`
        : '';

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

      if (PASSED_STATUSES.includes(status)) {
        if (PASSED_STATUSES.includes(avgStatus)) {
          if (time < avgTime) {
            icon = ICON_NAMES.arrowUpward;
          } else if (time > avgTime) {
            icon = ICON_NAMES.arrowDownward;
          }
        } else {
          icon = ICON_NAMES.done;
        }
      } else if (avgStatus === TEST_SUITE_STATUSES.passed) {
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
      const item = this.items[dataPoint.dataIndex];
      const {
        time = 0,
        status,
        avg_time: avgTime = 0,
        avg_status: avgStatus,
      } = item;

      const fixedTime = time.toFixed(3);
      const fixedAvgTime = avgTime.toFixed(3);

      const icon = `&nbsp;${this.getHistoricalTooltipIcon(item)}`;

      const statusMessage = status !== TEST_SUITE_STATUSES.total
        ? `${this.$t(`testSuite.statuses.${status}`)} `
        : '';

      const avgStatusMessage = avgStatus !== TEST_SUITE_STATUSES.total
        ? `${this.$t(`testSuite.statuses.${avgStatus}`)} `
        : '';

      const currentDiv = `<div>${this.$t('common.current')}: ${statusMessage}${fixedTime}s${icon}</div>`;
      const averageDiv = `<div>${this.$t('common.average')}: ${avgStatusMessage}${fixedAvgTime}s</div>`;

      return `<div>${currentDiv}${averageDiv}</div>`;
    },

    getTooltip({ tooltip: tooltipModel }) {
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

    updateItemsPerPage(itemsPerPage) {
      this.$emit('update:query', { ...this.query, itemsPerPage, page: 1 });
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
