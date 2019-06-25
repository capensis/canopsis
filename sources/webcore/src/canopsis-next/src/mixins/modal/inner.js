import { createNamespacedHelpers } from 'vuex';
import { OPTIONS_SANITIZE_TEXT_EDITOR } from '@/constants';

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
    sanitizedValue(value) {
      try {
        return this.$sanitize(value, OPTIONS_SANITIZE_TEXT_EDITOR);
      } catch (err) {
        console.warn(err);

        return '';
      }
    },
  },
};
