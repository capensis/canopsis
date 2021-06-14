import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('engineRunInfo');

export default {
  methods: {
    ...mapActions({
      fetchEnginesListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
