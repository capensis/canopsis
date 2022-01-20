import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehaviorTimespan');

export const entitiesPbehaviorTimespansMixin = {
  methods: {
    ...mapActions({
      fetchTimespansListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
