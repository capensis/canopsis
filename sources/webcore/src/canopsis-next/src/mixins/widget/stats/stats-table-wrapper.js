export default {
  computed: {
    getFormattedValue() {
      const PROPERTIES_FILTERS_MAP = {
        state_rate: value => this.$options.filters.percentage(value),
        ack_time_sla: value => this.$options.filters.percentage(value),
        resolve_time_sla: value => this.$options.filters.percentage(value),
        time_in_state: value => this.$options.filters.duration({ value }),
        mtbf: value => this.$options.filters.duration({ value }),
      };

      return (value, columnValue) => {
        if (PROPERTIES_FILTERS_MAP[columnValue]) {
          return PROPERTIES_FILTERS_MAP[columnValue](value);
        }

        return value;
      };
    },
    trendFormat() {
      return (value) => {
        if (value > 0) {
          return {
            icon: 'trending_up',
            color: 'primary',
          };
        } else if (value < 0) {
          return {
            icon: 'trending_down',
            color: 'error',
          };
        }

        return {
          icon: 'trending_flat',
        };
      };
    },
  },
};
