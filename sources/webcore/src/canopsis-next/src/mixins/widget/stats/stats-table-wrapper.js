import Vue from 'vue';

export default {
  filters: {
    trendColor(value) {
      if (value > 0) {
        return 'primary';
      } else if (value < 0) {
        return 'error';
      }

      return 'black';
    },

    trendIcon(value) {
      if (value > 0) {
        return 'trending_up';
      } else if (value < 0) {
        return 'trending_down';
      }

      return 'trending_flat';
    },

    formatValue(statValue, columnValue) {
      const PROPERTIES_FILTERS_MAP = {
        state_rate: value => Vue.options.filters.percentage(value),
        ack_time_sla: value => Vue.options.filters.percentage(value),
        resolve_time_sla: value => Vue.options.filters.percentage(value),
        time_in_state: value => Vue.options.filters.duration({ value }),
        mtbf: value => Vue.options.filters.duration({ value }),
      };

      if (PROPERTIES_FILTERS_MAP[columnValue]) {
        return PROPERTIES_FILTERS_MAP[columnValue](statValue);
      }

      return statValue;
    },
  },
  computed: {
    hasTrend() {
      return value => value.trend !== undefined && value.trend !== null;
    },
  },
};
