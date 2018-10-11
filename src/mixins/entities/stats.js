import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('stats');

export default {
  methods: {
    ...mapActions({
      fetchStatsListWithoutStore: 'fetchListWithoutStore',
      fetchStatsEvolutionWithoutStore: 'fetchEvolutionWithoutStore',
    }),
  },
};
