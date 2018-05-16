import { createNamespacedHelpers } from 'vuex';

const { mapActions: modalMapActions, mapGetters: modalMapGetters } = createNamespacedHelpers('modal');

export default {
  computed: {
    ...modalMapGetters(['config']),
  },
  methods: {
    ...modalMapActions({
      hideModal: 'hide',
    }),
  },
};
