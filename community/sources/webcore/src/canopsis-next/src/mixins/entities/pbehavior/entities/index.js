import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehavior/entities');

export default {
  methods: {
    ...mapActions({
      fetchPbehaviorEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
