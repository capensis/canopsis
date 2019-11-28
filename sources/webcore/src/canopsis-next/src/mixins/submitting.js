export default {
  data() {
    return {
      submitting: false,
    };
  },
  created() {
    const sourceSubmit = this.submit;

    if (sourceSubmit) {
      this.submit = async function submit(...args) {
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
};
