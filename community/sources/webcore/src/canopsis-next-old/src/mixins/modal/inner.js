export const modalInnerMixin = {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  computed: {
    config() {
      return this.modal.config;
    },
  },
};
