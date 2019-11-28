export default (method = 'submit') => ({
  data() {
    return {
      submitting: false,
    };
  },
  created() {
    const sourceSubmit = this[method];

    if (sourceSubmit) {
      this[method] = async function submit(...args) {
        try {
          this.submitting = true;

          await sourceSubmit.apply(this, args);
        } catch (err) {
          this.$popups.error({ text: err.details || this.$t('errors.default') });
        } finally {
          this.submitting = false;
        }
      };
    }
  },
  computed: {
    isDisabled() {
      if (this.errors) {
        return this.errors.any() || this.submitting;
      }

      return this.submitting;
    },
  },
});
