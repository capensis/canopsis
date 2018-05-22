import { createNamespacedHelpers } from 'vuex';

const { mapActions: modalMapActions } = createNamespacedHelpers('modal');

export default {
  methods: {
    ...modalMapActions({
      showModal: 'show',
      hideModal: 'hide',
    }),
  },
};
