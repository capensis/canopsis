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
  },
};
