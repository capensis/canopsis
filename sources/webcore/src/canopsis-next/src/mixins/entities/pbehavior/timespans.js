import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehaviorTimespan');

export default {
  methods: {
    ...mapActions({
      fetchTimespansListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
