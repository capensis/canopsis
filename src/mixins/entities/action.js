import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('action');

export default {
  methods: {
    ...mapActions({
      fetchActionsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
