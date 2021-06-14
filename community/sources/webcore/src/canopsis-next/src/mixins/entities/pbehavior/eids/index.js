import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehavior/eids');

export default {
  methods: {
    ...mapActions({
      fetchPbehaviorEidsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
