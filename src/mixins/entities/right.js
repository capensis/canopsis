import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('right');

export default {
  methods: {
    ...mapActions({
      createRight: 'create',
      fetchRightsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
