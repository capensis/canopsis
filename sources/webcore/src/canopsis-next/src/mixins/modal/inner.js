import { createNamespacedHelpers } from 'vuex';

const { mapActions: modalMapActions } = createNamespacedHelpers('modal');

/**
 * @mixin
 */
export default {
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
  methods: {
    ...modalMapActions({
      showModal: 'show',
      hideModalAction: 'hide',
    }),

    hideModal() {
      this.hideModalAction({ id: this.modal.id });
    },
  },
};
