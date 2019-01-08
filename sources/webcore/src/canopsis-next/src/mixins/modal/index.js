import { createNamespacedHelpers } from 'vuex';

const { mapActions: modalMapActions } = createNamespacedHelpers('modal');

/**
 * @mixin
 */
export default {
  methods: {
    ...modalMapActions({
      showModal: 'show',
      hideModal: 'hide',
    }),
  },
};
