export default {
  computed: {
    hasTrend() {
      return value => value.trend !== undefined && value.trend !== null;
    },
  },
};
