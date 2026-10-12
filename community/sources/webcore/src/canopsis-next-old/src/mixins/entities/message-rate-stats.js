import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('messageRateStats');

export const entitiesMessageRateStatsMixin = {
  methods: {
    ...mapActions({
      fetchMessageRateStatsWithoutStore: 'fetchMessageRateStatsWithoutStore',
    }),
  },
};
