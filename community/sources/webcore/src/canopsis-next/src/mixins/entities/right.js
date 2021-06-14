import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('right');

export default {
  methods: {
    ...mapActions({
      createRight: 'create',
      removeRight: 'remove',
      fetchRightsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
