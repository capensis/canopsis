import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('pattern');

export const entitiesPatternsMixin = {
  computed: {
    ...mapGetters({
      patternsMeta: 'meta',
      patternsPending: 'pending',
      patterns: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchPatternsList: 'fetchList',
      createPattern: 'create',
      updatePattern: 'update',
      removePattern: 'remove',
      bulkRemovePatterns: 'bulkRemove',
      checkPatternsCount: 'checkPatternsCount',
      fetchPatternsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchPatternsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
