export default {
  data() {
    return {
      serverErrorMessage: null,
    };
  },
  computed: {
    hasError() {
      return !!this.serverErrorMessage;
    },

    errorMessage() {
      return this.serverErrorMessage ? this.$t('errors.statsRequestProblem') : null;
    },
  },
};
