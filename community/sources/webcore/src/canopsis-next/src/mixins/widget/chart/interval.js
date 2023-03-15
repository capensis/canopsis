export const widgetIntervalFilterMixin = {
  inject: ['$system'],
  methods: {
    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },
  },
};
